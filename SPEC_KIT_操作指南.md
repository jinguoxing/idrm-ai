# IDRM Spec-Kit 操作指南

> **版本**: v1.0
> **更新日期**: 2024-12-26
> **适用范围**: IDRM项目全栈开发
> **基于**: 菜单管理模块实施经验

---

## 📖 文档概述

本文档详细说明了如何使用**Spec-Kit开发流程**在IDRM项目中实现新功能。这是基于AI Native理念的现代化开发方法，强调"契约高于一切"，通过Spec文件驱动整个开发流程。

### 核心理念

**AI Native 哲学**：
1. **契约高于一切** - Spec文件是"单一真理来源"
2. **AI优先，人类审核** - 代码生成由AI完成，人类负责决策和审查
3. **意图显式化** - 自然语言描述 → 技术规约 → 代码生成

### 四阶段工作流

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Specify   │ => │    Plan     │ => │   Generate  │ => │   Verify    │
│  功能规范    │    │  技术方案    │    │   代码生成   │    │   测试验证   │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
      ↓                  ↓                  ↓                  ↓
   需求澄清          架构设计          增量交付          质量保证
```

---

## 🎯 实施路径总览

### 完整文件结构

```
idrm-ai/
├── specs/                           # Spec文件目录
│   ├── specification-{module}.md    # 功能规格说明
│   └── IMPLEMENTATION_SUMMARY.md    # 实施总结
│
├── api/                             # API层
│   ├── {module}.api                 # API定义文件
│   ├── internal/
│   │   ├── config/                  # 配置
│   │   │   └── config.go
│   │   ├── handler/                 # Handler层（goctl生成）
│   │   │   └── sys/{module}/
│   │   ├── logic/                   # Logic层（AI生成）
│   │   │   └── sys/{module}/
│   │   ├── svc/                     # ServiceContext
│   │   │   └── servicecontext.go
│   │   └── types/                   # 类型定义
│   │       └── types.go
│   ├── etc/
│   │   └── {module}.yaml            # 配置文件
│   └── {module}.go                  # 主入口文件
│
├── model/                           # Model层
│   └── {domain}/{module}/
│       ├── interface.go             # Model接口
│       ├── types.go                 # 数据结构
│       ├── vars.go                  # 常量和错误
│       ├── factory.go               # ORM工厂
│       ├── gorm_dao.go              # GORM实现
│       ├── sqlx_model.go            # SQLx实现
│       ├── {module}_test.go         # 单元测试
│       └── README.md
│
├── pkg/errorx/                      # 错误处理
│   └── errorx.go                    # 统一错误定义
│
└── deploy/migrations/               # 数据库迁移
    └── YYYYMMDD_create_{table}.sql
```

---

## 📋 Stage 1: Specify（规范阶段）

### 目标
将业务需求转化为清晰的功能规格说明。

### 输入
- 业务需求文档
- 用户故事
- 产品原型

### 输出
- `specs/specification-{module}.md`

### 操作步骤

#### 1.1 创建Spec文件

```bash
# 在specs目录下创建功能规格文件
cd specs
touch specification-{module}.md
```

#### 1.2 编写Spec内容

使用以下模板：

```markdown
# Specification: {模块名称}

## 1. 业务目标
{描述业务目标和价值}

## 2. 核心功能
- **功能点1**: 描述
- **功能点2**: 描述

## 3. 技术约束 (基于项目宪法)
- **存储**: 数据库类型
- **缓存**: 缓存策略
- **接口路径**: API路径规范
- **数据结构**: 数据模型要求

## 4. 异常处理
- `错误码: 错误名称` (错误描述)
```

#### 1.3 Spec文件示例（菜单管理模块）

```markdown
# Specification: 菜单管理模块

## 1. 业务目标
实现后台管理系统的菜单维护，支持多级树形结构，用于前端动态渲染导航栏和权限控制。

## 2. 核心功能
- **基础字段**: 菜单名称、父级ID、路由路径、组件路径、图标、排序、权限标识（PermTag）。
- **菜单树查询**: 一次性获取完整树形结构。
- **增删改查**:
    - 创建/编辑时需校验父级 ID 是否存在。
    - 删除时，若存在子菜单，禁止删除。

## 3. 技术约束
- **存储**: 使用 MySQL
- **缓存**: 菜单树查询结果需通过 Redis 缓存
- **接口路径**: `/v1/sys/menu/*`
- **数据结构**: 树形转换逻辑需封装在 Logic 层

## 4. 异常处理
- `1101: ParentMenuNotFound` (父级菜单不存在)
- `1102: MenuHasSubEntries` (存在子菜单无法直接删除)
```

### 注意事项
- ✅ **只描述做什么，不描述怎么做**
- ✅ **避免包含具体技术实现细节**
- ✅ **明确业务规则和约束条件**
- ❌ **不包含代码片段**

---

## 🏗️ Stage 2: Plan（规划阶段）

### 目标
基于Spec文件设计技术方案，明确实施路径。

### 输入
- `specs/specification-{module}.md`
- 项目架构文档（`.specs/architecture/`）
- 编码规范（`.specs/coding-standards/`）

### 输出
- 技术设计文档（可选，建议记录）
- 文件清单
- 数据模型设计

### 操作步骤

#### 2.1 设计数据模型

```sql
CREATE TABLE {table_name} (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    -- 字段定义
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME,

    INDEX idx_{field} ({field})
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

#### 2.2 设计API端点

遵循RESTful规范：

```
POST   /api/v1/{resource}         - 创建
GET    /api/v1/{resource}/:id     - 详情
PUT    /api/v1/{resource}/:id     - 更新
DELETE /api/v1/{resource}/:id     - 删除
GET    /api/v1/{resource}s        - 列表
```

#### 2.3 确定文件清单

| 类型 | 路径 | 说明 |
|-----|------|------|
| API定义 | `api/{module}.api` | API接口定义 |
| Model层 | `model/{domain}/{module}/` | 数据访问层（7个文件） |
| Logic层 | `api/internal/logic/.../` | 业务逻辑层（6个文件） |
| Handler层 | `api/internal/handler/.../` | HTTP处理层（6个文件） |
| 配置 | `api/internal/svc/` | 服务上下文 |
| 数据库 | `deploy/migrations/` | 迁移脚本 |

#### 2.4 确定错误码

在 `pkg/errorx/errorx.go` 中添加模块专用错误码：

```go
const (
    ErrCode{Module}Error1 = 1101  // 错误描述
    ErrCode{Module}Error2 = 1102  // 错误描述
)
```

### 注意事项
- ✅ **遵循四层架构原则**
- ✅ **使用接口抽象，依赖倒置**
- ✅ **支持双ORM模式**
- ✅ **明确事务和缓存策略**

---

## 🔧 Stage 3: Generate（生成阶段）

### 目标
根据技术方案生成完整可运行的代码。

### 工作分配
- **AI任务**：生成代码实现
- **人类任务**：执行CLI工具（goctl）

### 操作步骤

#### 3.1 创建API定义文件

```bash
# 创建API定义
touch api/{module}.api
```

API定义模板：

```go
syntax = "v1"

info(
    title: "{模块名称}接口"
    desc: "{模块描述}"
    author: "IDRM Team"
    version: "v1.0"
)

import "base.api"

@server(
    prefix: /api/v1/{group}
    group: {resource}
)
service idrm-api {
    @doc "{操作}资源"
    @handler {action}{Resource}
    post /resource (CreateReq) returns (CreateResp)

    @doc "{操作}资源"
    @handler {action}{Resource}
    put /resource/:id (UpdateReq) returns (UpdateResp)

    // ... 更多端点
}
```

#### 3.2 添加错误码

编辑 `pkg/errorx/errorx.go`：

```go
// 在const块中添加
const (
    // {模块}错误 (1100-1199)
    ErrCode{Module}Error1 = 1101
    ErrCode{Module}Error2 = 1102
)

// 在errMsgMap中添加
var errMsgMap = map[int]string{
    // ... 现有错误
    ErrCode{Module}Error1: "{错误描述}",
    ErrCode{Module}Error2: "{错误描述}",
}
```

#### 3.3 创建Model层（AI生成）

Model层包含7个文件：

```
model/{domain}/{module}/
├── interface.go      # Model接口定义
├── types.go          # 数据结构
├── vars.go           # 常量和错误定义
├── factory.go        # ORM工厂函数
├── gorm_dao.go       # GORM实现
├── sqlx_model.go     # SQLx实现
└── {module}_test.go  # 单元测试
```

**关键点**：
- ✅ 统一接口定义
- ✅ 工厂模式自动选择ORM
- ✅ 支持事务
- ✅ 完整的错误处理

#### 3.4 生成Go-Zero代码（人类执行）

```bash
# 生成Handler和Types
goctl api go -api api/{module}.api -dir .
```

这将生成：
- `api/internal/handler/` - HTTP处理器
- `api/internal/types/` - 类型定义
- 其他必要文件

#### 3.5 创建Logic层（AI生成）

Logic层每个端点一个文件：

```
api/internal/logic/{group}/{resource}/
├── create{resource}logic.go
├── update{resource}logic.go
├── delete{resource}logic.go
├── get{resource}logic.go
└── list{resource}slogic.go
```

**Logic层模板**：

```go
package {resource}

import (
    "context"
    "fmt"

    "idrm/api/internal/svc"
    "idrm/api/internal/types"
    "idrm/model/{domain}/{resource}"

    "github.com/zeromicro/go-zero/core/logx"
)

type {Action}{Resource}Logic struct {
    logx.Logger
    ctx    context.Context
    svcCtx *svc.ServiceContext
}

func New{Action}{Resource}Logic(ctx context.Context, svcCtx *svc.ServiceContext) *{Action}{Resource}Logic {
    return &{Action}{Resource}Logic{
        Logger: logx.WithContext(ctx),
        ctx:    ctx,
        svcCtx: svcCtx,
    }
}

func (l *{Action}{Resource}Logic) {Action}{Resource}(req *types.{Action}{Resource}Req) (*types.{Action}{Resource}Resp, error) {
    // 1. 参数验证
    if err := l.validate(req); err != nil {
        return nil, err
    }

    // 2. 业务逻辑处理
    // ...

    // 3. 数据转换 (types -> model)
    data := &{resource}.{Resource}{
        // ...
    }

    // 4. 调用Model层
    result, err := l.svcCtx.{Resource}Model.Insert(l.ctx, data)
    if err != nil {
        l.Errorf("failed to create: %v", err)
        return nil, fmt.Errorf("failed to create: %w", err)
    }

    // 5. 返回结果
    return &types.{Action}{Resource}Resp{
        Id: result.Id,
    }, nil
}

func (l *{Action}{Resource}Logic) validate(req *types.{Action}{Resource}Req) error {
    // 验证逻辑
    return nil
}
```

#### 3.6 配置ServiceContext

编辑 `api/internal/svc/servicecontext.go`：

```go
type ServiceContext struct {
    Config config.Config

    // 数据库连接
    SqlConn *sql.DB
    GormDB  *gorm.DB

    // Model（接口类型）
    {Resource}Model {resource}.Model
}

func NewServiceContext(c config.Config) *ServiceContext {
    // 初始化数据库连接
    // ...

    // 初始化Model（使用工厂）
    {resource}Model := {resource}.NewModel(sqlConn, gormDB)

    return &ServiceContext{
        Config:         c,
        SqlConn:        sqlConn,
        GormDB:         gormDB,
        {Resource}Model: {resource}Model,
    }
}
```

#### 3.7 创建数据库迁移脚本

```bash
# 创建迁移文件
touch deploy/migrations/YYYYMMDD_create_{table}.sql
```

#### 3.8 创建配置文件

```bash
# 创建配置文件
touch api/etc/{module}.yaml
```

配置模板：

```yaml
Name: {module}-api
Host: 0.0.0.0
Port: 8080

Timeout: 30000

# 数据库配置
DataSource: {数据库连接}
GormDSN: {数据库连接}

# 日志配置
Log:
  ServiceName: {module}-api
  Mode: console
  Level: info
```

#### 3.9 创建主入口文件

```bash
touch api/{module}.go
```

主入口模板：

```go
package main

import (
    "flag"

    "idrm/api/internal/config"
    menuhandler "idrm/api/internal/handler/sys/menu"
    "idrm/api/internal/svc"

    "github.com/zeromicro/go-zero/core/conf"
    "github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "api/etc/{module}.yaml", "配置文件路径")

func main() {
    flag.Parse()

    var c config.Config
    conf.MustLoad(*configFile, &c)

    server := rest.MustNewServer(c.RestConf)
    defer server.Stop()

    serverCtx := svc.NewServiceContext(c)

    // 注册路由
    registerRoutes(server, serverCtx)

    server.Start()
}

func registerRoutes(server *rest.Server, serverCtx *svc.ServiceContext) {
    // 注册路由
}
```

### 注意事项
- ✅ **所有生成代码必须有完整注释**
- ✅ **遵循命名规范**
- ✅ **包含错误处理和日志**
- ✅ **Logic层不可直接访问数据库**

---

## ✅ Stage 4: Verify（验证阶段）

### 目标
验证生成的代码质量、功能和性能。

### 操作步骤

#### 4.1 初始化Go模块

```bash
# 初始化模块
go mod init {module-name}

# 整理依赖
go mod tidy
```

#### 4.2 编译检查

```bash
# 编译项目
go build ./...

# 或使用Makefile
make build
```

#### 4.3 初始化数据库

```bash
# 方式1：使用MySQL
mysql -u root -p < deploy/migrations/YYYYMMDD_create_{table}.sql

# 方式2：使用SQLite（开发环境）
# 数据库会在ServiceContext初始化时自动创建
```

#### 4.4 运行单元测试

```bash
# 运行所有测试
go test ./... -v

# 运行Model层测试
go test ./model/{domain}/{module}/... -v -cover

# 查看覆盖率
go test ./model/{domain}/{module}/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

#### 4.5 启动服务

```bash
# 方式1：直接运行
go run api/{module}.go -f api/etc/{module}.yaml

# 方式2：编译后运行
go build -o bin/{module}-api api/{module}.go
./bin/{module}-api -f api/etc/{module}.yaml

# 方式3：后台运行
nohup ./bin/{module}-api -f api/etc/{module}.yaml > /tmp/{module}-api.log 2>&1 &
```

#### 4.6 测试API

```bash
# 测试创建
curl -X POST http://localhost:8080/api/v1/{resource} \
  -H "Content-Type: application/json" \
  -d '{"field":"value"}'

# 测试查询
curl http://localhost:8080/api/v1/{resource}/1

# 测试列表
curl http://localhost:8080/api/v1/{resource}s

# 测试更新
curl -X PUT http://localhost:8080/api/v1/{resource}/1 \
  -H "Content-Type: application/json" \
  -d '{"field":"newvalue"}'

# 测试删除
curl -X DELETE http://localhost:8080/api/v1/{resource}/1
```

#### 4.7 代码质量检查

```bash
# 格式化代码
gofmt -s -w .
goimports -w .

# 代码检查
golangci-lint run

# 静态分析
go vet ./...
```

### 验证清单

- [ ] 编译通过
- [ ] 单元测试通过（覆盖率>80%）
- [ ] 服务正常启动
- [ ] API端点可访问
- [ ] 数据正常读写
- [ ] 错误处理正确
- [ ] 日志记录完整
- [ ] 代码格式规范

---

## 🎯 实战案例：菜单管理模块

### 完整实施路径

```bash
# ========== Stage 1: Specify ==========
# 1. 创建Spec文件
cd specs
cat > specification-menu.md << 'EOF'
# Specification: 菜单管理模块
## 1. 业务目标
实现后台管理系统的菜单维护...
## 2. 核心功能
...
EOF

# ========== Stage 2: Plan ==========
# 设计已完成，进入生成阶段

# ========== Stage 3: Generate ==========

# 1. 创建API定义
cat > api/sys/menu.api << 'EOF'
syntax = "v1"
info(title: "菜单管理接口", ...)
...
EOF

# 2. 添加错误码
# 编辑 pkg/errorx/errorx.go
# 添加: ErrCodeParentMenuNotFound = 1101

# 3. 创建Model层（7个文件）
mkdir -p model/system/menu
# AI生成:
# - interface.go
# - types.go
# - vars.go
# - factory.go
# - gorm_dao.go
# - sqlx_model.go
# - menu_test.go

# 4. 生成Go-Zero代码（人类执行）
goctl api go -api api/sys/menu.api -dir .

# 5. 创建Logic层（6个文件）
# AI生成所有logic文件

# 6. 配置ServiceContext
# 编辑 api/internal/svc/servicecontext.go

# 7. 创建主入口
# AI生成 api/menu.go

# 8. 创建配置文件
# AI生成 api/etc/menu.yaml

# ========== Stage 4: Verify ==========

# 1. 初始化依赖
go mod init idrm
go mod tidy

# 2. 编译
go build -o bin/menu-api api/menu.go

# 3. 运行测试
go test ./model/system/menu/... -v -cover

# 4. 启动服务
./bin/menu-api -f api/etc/menu.yaml

# 5. 测试API
curl http://localhost:8080/api/v1/sys/menu/tree
```

### 实施成果

- ✅ **25+ 文件**生成
- ✅ **~4,500行代码**
- ✅ **6个API端点**
- ✅ **>85%测试覆盖率**
- ✅ **完整文档**

---

## 🛠️ 常用命令

### Makefile命令

```bash
# 查看所有命令
make help

# 初始化数据库
make migrate

# 编译项目
make build

# 运行服务
make run

# 运行测试
make test

# 代码检查
make check

# Docker操作
make docker-build
make docker-up
```

### Git操作

```bash
# 添加文件
git add .

# 提交
git commit -m "feat({module}): implement {module} module

- Add API definition in api/{module}.api
- Implement Model layer with dual ORM support
- Implement Logic layer for all CRUD operations
- Add unit tests with >85% coverage
- Add database migration script

Closes #{issue-number}"

# 推送
git push origin feature/{module}
```

---

## 📊 质量标准

### 代码质量

- ✅ 编译通过（`go build ./...`）
- ✅ 测试通过（`go test ./...`）
- ✅ Lint无错误（`golangci-lint run`）
- ✅ 格式规范（`gofmt`, `goimports`）
- ✅ 完整注释（所有公开API）

### 架构标准

- ✅ 遵循四层架构
- ✅ 依赖倒置（面向接口）
- ✅ 双ORM支持
- ✅ 事务支持
- ✅ 错误处理完整

### 文档标准

- ✅ Spec文件完整
- ✅ 代码注释完整
- ✅ README完整
- ✅ API文档完整

---

## ⚠️ 常见问题和解决方案

### 问题1：goctl生成失败

**原因**：API定义文件语法错误

**解决**：
```bash
# 检查API文件语法
goctl api validate -api api/{module}.api

# 查看详细错误
goctl api go -api api/{module}.api -dir . --verbose
```

### 问题2：数据库连接失败

**原因**：配置文件中的数据库连接字符串不正确

**解决**：
```bash
# 检查配置文件
cat api/etc/{module}.yaml

# 测试数据库连接
mysql -u root -p -h 127.0.0.1

# SQLite版本确保文件路径正确
ls -la ./idrm.db
```

### 问题3：端口被占用

**原因**：8080端口已被其他进程使用

**解决**：
```bash
# 查找占用端口的进程
lsof -ti:8080

# 杀死进程
lsof -ti:8080 | xargs kill -9

# 或修改配置文件使用其他端口
```

### 问题4：依赖包下载失败

**原因**：网络问题或Go版本不兼容

**解决**：
```bash
# 设置Go代理（中国大陆）
export GOPROXY=https://goproxy.cn,direct

# 清理缓存
go clean -modcache

# 重新下载
go mod tidy
```

---

## 📚 参考资源

### 项目文档

- [项目宪章](constitution.md) - 必读
- [分层架构](.specs/architecture/layered-architecture.md) - 架构规范
- [双ORM模式](.specs/architecture/dual-orm-pattern.md) - ORM使用
- [API设计指南](.specs/architecture/api-design-guide.md) - RESTful规范

### 外部资源

- [Go-Zero文档](https://go-zero.dev/)
- [GORM文档](https://gorm.io/)
- [Go规范](https://go.dev/doc/effective_go)

---

## 🎓 最佳实践

### 1. Spec编写

- ✅ **简洁明了** - 一页纸原则
- ✅ **业务语言** - 避免技术术语
- ✅ **明确约束** - 清晰列出限制条件

### 2. 代码生成

- ✅ **分批生成** - 先Model，再Logic，最后Handler
- ✅ **增量提交** - 每完成一层就提交一次
- ✅ **持续测试** - 边生成边测试

### 3. 质量保证

- ✅ **自测优先** - AI生成后立即自测
- ✅ **代码审查** - 人类必须审查所有AI生成的代码
- ✅ **文档同步** - 代码和文档同步更新

---

## 🔄 持续改进

### 版本历史

- **v1.0** (2024-12-26) - 基于菜单管理模块实施经验创建

### 反馈渠道

如有问题或改进建议，请：
1. 提交Issue到项目仓库
2. 更新本文档
3. 团队讨论和评审

---

## 📝 附录

### A. 完整检查清单

#### Specify阶段
- [ ] Spec文件已创建
- [ ] 业务目标清晰
- [ ] 核心功能明确
- [ ] 技术约束完整
- [ ] 异常处理定义

#### Plan阶段
- [ ] 数据模型设计完成
- [ ] API端点设计完成
- [ ] 文件清单明确
- [ ] 错误码已定义
- [ ] 技术方案已评审

#### Generate阶段
- [ ] API定义文件已创建
- [ ] Model层已实现（7个文件）
- [ ] goctl代码已生成
- [ ] Logic层已实现
- [ ] ServiceContext已配置
- [ ] 主入口文件已创建
- [ ] 配置文件已创建
- [ ] 数据库脚本已创建

#### Verify阶段
- [ ] Go模块已初始化
- [ ] 依赖已下载
- [ ] 编译通过
- [ ] 单元测试通过
- [ ] 服务正常启动
- [ ] API测试通过
- [ ] 代码质量检查通过
- [ ] 文档完整

### B. 快速参考卡

```bash
# === 开发流程快速命令 ===

# 1. 开始新功能
cd specs && touch specification-{module}.md

# 2. 生成代码
# - API定义: api/{module}.api
# - Model层: model/{domain}/{module}/
# - goctl: goctl api go -api api/{module}.api -dir .

# 3. 编译测试
go mod tidy
go build ./...
go test ./...

# 4. 运行服务
go run api/{module}.go -f api/etc/{module}.yaml

# 5. 质量检查
make check
```

---

**本文档是IDRM项目Spec-Kit开发流程的完整指南，所有开发活动应遵循此文档的规范和流程。**

**维护者**: IDRM开发团队
**最后更新**: 2024-12-26
