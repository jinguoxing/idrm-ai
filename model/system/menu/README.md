# 菜单管理模块 (Menu Management Module)

> **版本**: v1.0
> **更新日期**: 2024-12-24
> **状态**: ✅ 已实现

---

## 📖 模块概述

菜单管理模块提供完整的菜单CRUD操作，支持多级树形结构，用于后台管理系统的导航栏渲染和权限控制。

### 核心功能

- ✅ **增删改查** - 完整的菜单CRUD操作
- ✅ **树形结构** - 支持多级菜单嵌套
- ✅ **双ORM模式** - 同时支持GORM和SQLx
- ✅ **事务支持** - 完整的事务管理
- ✅ **Redis缓存** - 菜单树查询结果缓存（待实现）
- ✅ **单元测试** - 完整的测试覆盖

---

## 🏗️ 架构设计

### 分层架构

```
┌─────────────────────────────────────┐
│      API Layer (Go-Zero)            │
│  Handler → Logic → Model             │
└─────────────────────────────────────┘
             ↓
┌─────────────────────────────────────┐
│      Model Layer (Dual ORM)         │
│  Interface → Factory → GORM/SQLx     │
└─────────────────────────────────────┘
             ↓
┌─────────────────────────────────────┐
│      Infrastructure Layer           │
│  MySQL / Redis                       │
└─────────────────────────────────────┘
```

### 目录结构

```
model/system/menu/
├── interface.go      # Model接口定义
├── types.go          # Menu数据结构
├── vars.go           # 常量和错误定义
├── factory.go        # ORM工厂函数
├── gorm_dao.go       # GORM实现
├── sqlx_model.go     # SQLx实现
├── menu_test.go      # 单元测试
└── README.md         # 本文档

api/internal/logic/sys/menu/
├── createmenulogic.go      # 创建菜单
├── updatemenulogic.go      # 更新菜单
├── deletemenulogic.go      # 删除菜单
├── getmenulogic.go         # 获取菜单详情
├── getmenutreelogic.go     # 获取菜单树
└── listmenuslogic.go       # 获取菜单列表

deploy/migrations/
└── 20241224_create_sys_menu.sql  # 数据库迁移脚本
```

---

## 📊 数据模型

### sys_menu 表结构

| 字段 | 类型 | 说明 | 约束 |
|-----|------|-----|------|
| id | BIGINT | 菜单ID | 主键，自增 |
| parent_id | BIGINT | 父级菜单ID | 0表示根菜单 |
| name | VARCHAR(50) | 菜单名称 | 必填 |
| route_path | VARCHAR(200) | 路由路径 | 可选 |
| component_path | VARCHAR(200) | 组件路径 | 可选 |
| icon | VARCHAR(100) | 菜单图标 | 可选 |
| sort_order | INT | 排序号 | 默认0 |
| perm_tag | VARCHAR(100) | 权限标识 | 可选 |
| status | TINYINT | 状态 | 1-启用，0-禁用 |
| created_at | DATETIME | 创建时间 | 自动生成 |
| updated_at | DATETIME | 更新时间 | 自动更新 |
| deleted_at | DATETIME | 删除时间 | 软删除 |

### 索引

- `uk_id` - 主键唯一索引
- `idx_parent_id` - 父级ID索引
- `idx_perm_tag` - 权限标识索引
- `idx_status` - 状态索引
- `idx_sort_order` - 排序索引
- `idx_deleted_at` - 软删除索引

---

## 🔌 API接口

### 基础路径
```
/api/v1/sys/menu
```

### 接口列表

| 方法 | 路径 | 说明 | 权限标识 |
|-----|------|-----|---------|
| POST | /menu | 创建菜单 | - |
| PUT | /menu/:id | 更新菜单 | - |
| DELETE | /menu/:id | 删除菜单 | - |
| GET | /menu/:id | 获取菜单详情 | - |
| GET | /menu/tree | 获取菜单树 | - |
| GET | /menus | 获取菜单列表 | - |

### 请求/响应示例

#### 1. 创建菜单

**请求**:
```json
POST /api/v1/sys/menu
{
  "parentId": 0,
  "name": "系统管理",
  "routePath": "/system",
  "componentPath": "Layout",
  "icon": "setting",
  "sortOrder": 100,
  "permTag": "system",
  "status": 1
}
```

**响应**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "name": "系统管理"
  },
  "requestId": "req-123456"
}
```

#### 2. 获取菜单树

**请求**:
```
GET /api/v1/sys/menu/tree?status=1
```

**响应**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "tree": [
      {
        "id": 1,
        "parentId": 0,
        "name": "系统管理",
        "routePath": "/system",
        "icon": "setting",
        "sortOrder": 100,
        "status": 1,
        "children": [
          {
            "id": 2,
            "parentId": 1,
            "name": "菜单管理",
            "routePath": "/system/menu",
            "sortOrder": 1,
            "status": 1,
            "children": []
          }
        ]
      }
    ]
  }
}
```

---

## ⚙️ 使用方法

### 1. 数据库初始化

执行数据库迁移脚本：
```bash
mysql -u root -p idrm_db < deploy/migrations/20241224_create_sys_menu.sql
```

### 2. 生成Go-Zero代码

```bash
# 生成Handler和Type定义
goctl api go -api api/sys/menu.api -dir .
```

### 3. 配置ServiceContext

在 `api/internal/svc/servicecontext.go` 中添加：
```go
type ServiceContext struct {
    Config config.Config

    // 菜单Model（接口类型）
    MenuModel menu.Model
}

func NewServiceContext(c config.Config) *ServiceContext {
    return &ServiceContext{
        Config: c,

        // 初始化MenuModel（工厂模式）
        MenuModel: menu.NewModel(c.SqlConn, c.GormDB),
    }
}
```

### 4. 运行服务

```bash
# 开发模式
go run api/menu.go -f api/etc/menu.yaml

# 生产构建
go build -o menu-api api/menu.go
./menu-api -f api/etc/menu.yaml
```

### 5. 测试接口

```bash
# 创建菜单
curl -X POST http://localhost:8080/api/v1/sys/menu \
  -H "Content-Type: application/json" \
  -d '{"parentId":0,"name":"系统管理","sortOrder":100,"status":1}'

# 获取菜单树
curl http://localhost:8080/api/v1/sys/menu/tree

# 获取菜单详情
curl http://localhost:8080/api/v1/sys/menu/1
```

---

## 🧪 测试

### 运行单元测试

```bash
# 运行所有测试
go test ./model/system/menu/... -v

# 运行特定测试
go test ./model/system/menu/... -v -run TestMenuDao_Insert

# 查看测试覆盖率
go test ./model/system/menu/... -cover

# 生成覆盖率报告
go test ./model/system/menu/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### 性能测试

```bash
# 运行性能测试
go test ./model/system/menu/... -bench=. -benchmem
```

---

## 🔍 业务规则

### 1. 创建菜单

- ✅ 菜单名称不能为空，最长50个字符
- ✅ 父级菜单必须存在（parentId为0时表示根菜单）
- ✅ 排序号默认为0
- ✅ 状态值必须为0或1

### 2. 更新菜单

- ✅ 菜单必须存在
- ✅ 不能将自己设置为父级菜单（防止循环引用）
- ✅ 父级菜单必须存在

### 3. 删除菜单

- ✅ 菜单必须存在
- ✅ 如果存在子菜单，则不允许删除（防止孤立节点）
- ✅ 软删除（设置deleted_at字段）

### 4. 菜单树

- ✅ 支持状态过滤（status参数）
- ✅ 递归构建树形结构
- ✅ 按sort_order和id排序

---

## ❌ 错误处理

### 错误码定义

| 错误码 | 说明 | 处理方式 |
|-------|------|---------|
| 1101 | 父级菜单不存在 | 检查parentId参数 |
| 1102 | 存在子菜单无法删除 | 先删除子菜单或调整层级 |
| 30001 | 菜单不存在 | 检查菜单ID |
| 20000 | 参数错误 | 检查请求参数 |

### 错误响应示例

```json
{
  "code": 1101,
  "msg": "父级菜单不存在",
  "requestId": "req-123456"
}
```

---

## 🚀 性能优化

### 1. Redis缓存（待实现）

菜单树查询结果缓存，减少数据库压力：

```go
// 缓存key
menu:tree:all          // 所有菜单
menu:tree:status:1     // 启用的菜单

// TTL
1小时
```

### 2. 索引优化

已创建的索引：
- `idx_parent_id` - 加速父级查询
- `idx_status` - 加速状态过滤
- `idx_sort_order` - 加速排序

### 3. 查询优化

- ✅ 使用批量查询减少数据库往返
- ✅ 软删除使用索引
- ✅ 树形结构在内存中构建

---

## 📝 注意事项

### 1. 循环引用检测

更新菜单时，需要检查是否会形成循环引用：
```
菜单A (id=1) → 菜单B (id=2) → 菜单A (id=1) ❌
```

### 2. 软删除

删除操作是软删除，数据不会真正删除，只是设置deleted_at字段。

### 3. 排序规则

菜单排序规则：
1. 按 `sort_order` 升序
2. 如果 `sort_order` 相同，按 `id` 升序

### 4. 权限标识

`perm_tag` 字段用于权限控制，建议格式：
```
模块:功能:操作
例: system:menu:view
```

---

## 🔗 相关文档

- [项目宪章](../constitution.md) - 项目开发规范
- [分层架构](../.specs/architecture/layered-architecture.md) - 架构设计
- [双ORM模式](../.specs/architecture/dual-orm-pattern.md) - ORM使用指南
- [API设计指南](../.specs/architecture/api-design-guide.md) - RESTful规范

---

## 📞 联系方式

如有问题，请联系IDRM开发团队。

---

**本模块遵循IDRM项目宪章的所有规范和标准。**
