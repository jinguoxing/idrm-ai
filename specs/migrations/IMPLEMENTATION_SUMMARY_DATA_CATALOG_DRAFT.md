# 数据资源目录暂存功能 - 实施总结

## 概述

本文档记录了从原项目（project_bundle/data-catalog）迁移 `SaveDataCatalogDraft` 功能到新 IDRM 项目（基于 Go-Zero 框架）的实施过程。

**实施日期**: 2024-12-27
**功能范围**: 数据资源目录暂存功能（创建暂存、变更暂存、草稿更新）
**框架**: Go-Zero v1.9+

---

## 一、Specification 阶段

### 1.1 规格文档

- **位置**: `specs/migrations/specification-data-catalog-draft.md`
- **内容**: 完整的业务需求、数据库表结构、错误码定义

### 1.2 核心功能

1. **创建暂存**: 首次创建全新的目录草稿（CatalogID 为空）
2. **变更暂存**: 对已发布目录进行编辑，创建或更新草稿副本
3. **草稿更新**: 对现有草稿进行内容更新

---

## 二、Plan 阶段

### 2.1 技术架构

```
┌─────────────┐
│   Handler   │  ← HTTP 请求处理
└──────┬──────┘
       │
┌──────▼──────┐
│    Logic    │  ← 业务逻辑（三种场景处理）
└──────┬──────┘
       │
┌──────▼──────┐
│    Model    │  ← 数据访问（双 ORM 支持）
└──────┬──────┘
       │
   ┌───▼────┐
   │Database│
   └────────┘
```

### 2.2 API 定义

**位置**: `api/doc/api/data-catalog/data-catalog.api`

```go
@server(
    prefix: /api/v1/data-catalog
    group: datacatalog
)
service idrm-api {
    @doc "暂存数据资源目录"
    @handler saveDraft
    post /draft (SaveDraftReq) returns (SaveDraftResp)

    @doc "获取草稿详情"
    @handler getDraft
    get /draft/:id (GetDraftReq) returns (GetDraftResp)

    @doc "草稿列表"
    @handler listDrafts
    get /drafts (ListDraftsReq) returns (ListDraftsResp)

    @doc "删除草稿"
    @handler deleteDraft
    delete /draft/:id (DeleteDraftReq) returns (DeleteDraftResp)
}
```

---

## 三、Generate 阶段

### 3.1 文件结构

```
idrm-ai/
├── api/
│   ├── doc/api/
│   │   └── data-catalog/
│   │       └── data-catalog.api      # API 定义
│   ├── internal/
│   │   ├── handler/datacatalog/      # HTTP 处理器
│   │   ├── logic/datacatalog/        # 业务逻辑
│   │   ├── types/types.go            # 请求/响应类型
│   │   └── svc/servicecontext.go     # 服务上下文
│   └── routes.go                     # 路由注册
├── model/data-catalog/
│   ├── data_catalog/                 # 主目录表
│   │   ├── interface.go
│   │   ├── types.go
│   │   ├── vars.go
│   │   ├── factory.go
│   │   ├── gorm_dao.go
│   │   └── sqlx_model.go
│   ├── data_catalog_column/          # 信息项表
│   ├── data_catalog_category/        # 分类关联表
│   └── data_resource/                # 挂接资源表
├── pkg/
│   ├── errorx/errorx.go              # 错误码扩展
│   └── orm/                          # ORM 选择器
└── deploy/migrations/
    └── 20251226_data_catalog_draft.sql # 数据库迁移
```

### 3.2 核心实现

#### Model 层 - 双 ORM 支持

```go
// 工厂模式
func NewModel(sqlConn *sql.DB, gormDB *gorm.DB) Model {
    if gormDB != nil && gormFactory != nil {
        return gormFactory(gormDB)
    }
    if sqlConn != nil && sqlxFactory != nil {
        return sqlxFactory(sqlConn)
    }
    panic("no database connection available")
}

// 策略模式（带指标）
func NewModelWithStrategy(sqlConn *sql.DB, gormDB *gorm.DB,
    strategy config.ORMStrategy, metrics *orm.ORMMetrics) Model {
    selector := orm.NewORMSelector(strategy, sqlConn, gormDB, metrics)
    return &StrategyModel{
        gormModel: gormFactory(gormDB),
        sqlxModel: sqlxModel(sqlConn),
        selector:  selector,
    }
}
```

#### Logic 层 - 三种场景处理

```go
func (l *SaveDraftLogic) SaveDraft(req *types.SaveDraftReq) (*types.SaveDraftResp, error) {
    if req.CatalogId == "" {
        // 场景1: 创建新草稿
        return l.createNewDraft(req)
    }

    existingCatalog, _ := l.svcCtx.DataCatalogModel.FindOne(l.ctx, req.CatalogId)
    if existingCatalog.PublishStatus == data_catalog.StatusPublished {
        // 场景2: 已发布目录的变更暂存
        return l.createChangeDraft(req, existingCatalog)
    } else {
        // 场景3: 更新现有草稿
        return l.updateExistingDraft(req, existingCatalog)
    }
}
```

---

## 四、Verify 阶段

### 4.1 编译状态

- ✅ Model 层编译通过
- ✅ Logic 层编译通过（datacatalog 模块）
- ✅ Handler 层编译通过（datacatalog 模块）
- ⚠️ Menu 模块需要类型更新（与草稿功能无关）

### 4.2 已解决的问题

1. **Trans 方法签名**: 添加缺失的返回 error
2. **SQLx 接口兼容性**: 创建 sqlxDBer 和 sqlxTxer 接口
3. **StrategyModel 类型断言**: 添加 model.(Model) 类型断言
4. **常量重复声明**: 移除 vars.go 中重复的状态常量
5. **TableName 字段冲突**: 重命名为 TblName
6. **Import 路径修复**: 修正 goctl 生成的 `idrm/api/api/` 路径
7. **类型映射**: 匹配 goctl 生成的类型与手动创建的 Logic

### 4.3 待测试功能

1. 创建新草稿 API
2. 获取草稿详情 API
3. 草稿列表查询 API
4. 删除草稿 API
5. 变更暂存（已发布目录）
6. 四表事务完整性

---

## 五、数据库结构

### 5.1 主表: t_data_catalog

```sql
CREATE TABLE t_data_catalog (
    id VARCHAR(20) PRIMARY KEY,
    title VARCHAR(500) NOT NULL,
    code VARCHAR(50) NOT NULL,
    description VARCHAR(1000),
    source_dept_id VARCHAR(36) NOT NULL,
    publish_status VARCHAR(20) DEFAULT 'unpublished',
    online_status VARCHAR(20) DEFAULT 'notline',
    draft_id BIGINT DEFAULT 0,
    shared_type TINYINT NOT NULL,
    open_type TINYINT NOT NULL,
    -- ... 更多字段
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP
);
```

### 5.2 关联表

- `t_data_catalog_column`: 信息项（字段定义）
- `t_data_catalog_category`: 分类关联（多对多）
- `t_data_resource`: 挂接资源（表、视图、API、文件）

---

## 六、错误码定义

### 6.1 扩展错误码（pkg/errorx/errorx.go）

```go
const (
    ErrCodeDataCatalogNotFound           = 1201 // 数据资源目录不存在
    ErrCodeDataCatalogDepartmentNotFound = 1202 // 关联部门不存在
    ErrCodeDataCatalogColumnNotFound     = 1203 // 信息项不存在
    ErrCodeCatalogNameRepeat             = 1204 // 目录名称重复
    ErrCodeDataResourceNotExist          = 1205 // 数据资源已挂载或不存在
    ErrCodeDataResourceTypeNotSupport    = 1206 // 资源类型不支持
)
```

---

## 七、下一步计划

1. **单元测试**: 为 Logic 层编写测试用例
2. **集成测试**: 测试完整的 API 流程
3. **性能测试**: 验证双 ORM 性能指标
4. **文档完善**: API 文档和使用指南
5. **Menu 模块**: 修复 menu logic 的类型不匹配问题

---

## 八、参考资料

- **项目宪章**: `constitution.md`
- **SPEC-Kit 指南**: `SPEC_KIT_操作指南.md`
- **原项目代码**: `project_bundle/data-catalog/adapter/controller/data_catalog/v1/data_catalog.go`
- **数据库迁移**: `deploy/migrations/20251226_data_catalog_draft.sql`

---

**签名**: Claude AI Assistant
**审核**: 待审核
