# 分层架构详解

> **文档版本**: v1.0  
> **最后更新**: 2025-12-24  
> **适用范围**: IDRM全栈应用

---

## 📖 概述

IDRM采用经典的**四层架构**模式，严格分离各层职责，确保系统的可维护性、可测试性和可扩展性。

### 为什么需要分层？

- ✅ **职责清晰** - 每层专注自己的职责
- ✅ **易于测试** - 可以独立测试每一层
- ✅ **便于维护** - 修改一层不影响其他层
- ✅ **团队协作** - 不同团队可以并行开发不同层

---

## 🏗️ 架构全景图

```
┌────────────────────────────────────────────────────────────┐
│                   Presentation Layer                        │
│                     (表现层/接口层)                           │
│  ┌────────────────────────────────────────────────────┐    │
│  │              HTTP Handler (Go-Zero)                 │    │
│  │  • 接收HTTP请求                                      │    │
│  │  • 参数解析和基础验证                                 │    │
│  │  • 调用Logic层                                       │    │
│  │  • 统一响应格式化                                     │    │
│  └────────────────────────────────────────────────────┘    │
└────────────────────────────────────────────────────────────┘
                           ↓ 依赖
┌────────────────────────────────────────────────────────────┐
│                    Business Layer                           │
│                     (业务逻辑层)                             │
│  ┌────────────────────────────────────────────────────┐    │
│  │                   Logic                             │    │
│  │  • 业务规则验证                                      │    │
│  │  • 复杂业务逻辑处理                                   │    │
│  │  • 流程编排                                          │    │
│  │  • 数据转换                                          │    │
│  └────────────────────────────────────────────────────┘    │
└────────────────────────────────────────────────────────────┘
                           ↓ 依赖
┌────────────────────────────────────────────────────────────┐
│                   Persistence Layer                         │
│                      (持久化层)                              │
│  ┌────────────────────────────────────────────────────┐    │
│  │                    Model                            │    │
│  │  • 数据访问接口(Interface)                           │    │
│  │  • ORM封装(GORM/SQLx)                               │    │
│  │  • 事务管理                                          │    │
│  │  • 查询构建                                          │    │
│  └────────────────────────────────────────────────────┘    │
└────────────────────────────────────────────────────────────┘
                           ↓ 依赖
┌────────────────────────────────────────────────────────────┐
│                 Infrastructure Layer                        │
│                    (基础设施层)                              │
│  ┌────────────────────────────────────────────────────┐    │
│  │   Database / Cache / Message Queue                 │    │
│  │  • MySQL 8.0                                        │    │
│  │  • Redis 7.0                                        │    │
│  │  • Kafka 3.0                                        │    │
│  └────────────────────────────────────────────────────┘    │
└────────────────────────────────────────────────────────────┘
```

---

## 📂 各层详细说明

### 1. Handler层 (表现层)

#### 位置
```
api/internal/handler/
├── resource_catalog/
│   └── category/
│       ├── createcategoryhandler.go
│       ├── getcategoryhandler.go
│       └── listcategoryhandler.go
└── data_view/
    └── ...
```

#### 职责

**✅ 应该做的事情**：
1. 接收HTTP请求
2. 解析请求参数（URL参数、Query参数、Body）
3. 基础参数验证（类型、必填项）
4. 调用对应的Logic
5. 格式化响应（成功/失败）
6. 设置HTTP状态码和Header

**❌ 不应该做的事情**：
1. ❌ 实现业务逻辑
2. ❌ 直接访问数据库
3. ❌ 直接调用Model层
4. ❌ 包含复杂的数据处理

#### 代码示例

```go
// ✅ 好的Handler实现
package category

import (
    "net/http"
    
    "idrm/api/internal/logic/resource_catalog/category"
    "idrm/api/internal/svc"
    "idrm/api/internal/types"
    "idrm/pkg/response"
    
    "github.com/zeromicro/go-zero/rest/httpx"
)

type CreateCategoryHandler struct {
    svcCtx *svc.ServiceContext
}

func NewCreateCategoryHandler(svcCtx *svc.ServiceContext) *CreateCategoryHandler {
    return &CreateCategoryHandler{
        svcCtx: svcCtx,
    }
}

func (h *CreateCategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
    // 1. 解析请求参数
    var req types.CreateCategoryReq
    if err := httpx.Parse(r, &req); err != nil {
        response.Error(w, err)
        return
    }
    
    // 2. 调用Logic层处理业务
    l := category.NewCreateCategoryLogic(r.Context(), h.svcCtx)
    resp, err := l.CreateCategory(&req)
    
    // 3. 返回响应
    if err != nil {
        response.Error(w, err)
    } else {
        response.Success(w, resp)
    }
}
```

```go
// ❌ 错误的Handler实现
func (h *CreateCategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
    var req types.CreateCategoryReq
    httpx.Parse(r, &req)
    
    // ❌ 错误1：包含业务逻辑
    if req.Code == "" {
        req.Code = generateCode(req.Name)
    }
    
    // ❌ 错误2：直接访问Model
    category := &category.Category{
        Name: req.Name,
        Code: req.Code,
    }
    result, _ := h.svcCtx.CategoryModel.Insert(r.Context(), category)
    
    // ❌ 错误3：忽略错误处理
    response.Success(w, result)
}
```

---

### 2. Logic层 (业务逻辑层)

#### 位置
```
api/internal/logic/
├── resource_catalog/
│   └── category/
│       ├── createcategorylogic.go
│       ├── getcategorylogic.go
│       └── listcategorylogic.go
└── data_view/
    └── ...
```

#### 职责

**✅ 应该做的事情**：
1. 实现业务规则
2. 数据验证（业务级别）
3. 流程编排（调用多个Model、组合数据）
4. 数据格式转换（types ↔ model）
5. 调用Model层进行数据操作
6. 错误处理和日志记录

**❌ 不应该做的事情**：
1. ❌ 直接访问数据库
2. ❌ 包含HTTP相关代码
3. ❌ 直接操作Request/Response
4. ❌ 硬编码配置信息

#### 代码示例

```go
// ✅ 好的Logic实现
package category

import (
    "context"
    "fmt"
    
    "idrm/api/internal/svc"
    "idrm/api/internal/types"
    "idrm/model/resource_catalog/category"
    
    "github.com/zeromicro/go-zero/core/logx"
)

type CreateCategoryLogic struct {
    logx.Logger
    ctx    context.Context
    svcCtx *svc.ServiceContext
}

func NewCreateCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCategoryLogic {
    return &CreateCategoryLogic{
        Logger: logx.WithContext(ctx),
        ctx:    ctx,
        svcCtx: svcCtx,
    }
}

func (l *CreateCategoryLogic) CreateCategory(req *types.CreateCategoryReq) (*types.CreateCategoryResp, error) {
    // 1. 业务验证
    if err := l.validateCategory(req); err != nil {
        return nil, err
    }
    
    // 2. 检查code是否已存在
    existing, err := l.svcCtx.CategoryModel.FindByCode(l.ctx, req.Code)
    if err == nil && existing != nil {
        return nil, fmt.Errorf("category code %s already exists", req.Code)
    }
    
    // 3. 数据转换 (types -> model)
    categoryData := &category.Category{
        Name:   req.Name,
        Code:   req.Code,
        Status: req.Status,
    }
    
    // 4. 调用Model层
    result, err := l.svcCtx.CategoryModel.Insert(l.ctx, categoryData)
    if err != nil {
        l.Errorf("failed to create category: %v", err)
        return nil, fmt.Errorf("failed to create category: %w", err)
    }
    
    // 5. 记录日志
    l.Infow("category created successfully",
        logx.Field("categoryId", result.Id),
        logx.Field("categoryCode", result.Code),
    )
    
    // 6. 返回结果 (model -> types)
    return &types.CreateCategoryResp{
        Id:   result.Id,
        Name: result.Name,
        Code: result.Code,
    }, nil
}

// 业务验证方法
func (l *CreateCategoryLogic) validateCategory(req *types.CreateCategoryReq) error {
    if len(req.Name) > 100 {
        return fmt.Errorf("category name too long (max 100 characters)")
    }
    
    if len(req.Code) < 2 || len(req.Code) > 50 {
        return fmt.Errorf("category code length must be between 2 and 50")
    }
    
    if req.Status != 0 && req.Status != 1 {
        return fmt.Errorf("invalid status value")
    }
    
    return nil
}
```

```go
// ❌ 错误的Logic实现
func (l *CreateCategoryLogic) CreateCategory(req *types.CreateCategoryReq) (*types.CreateCategoryResp, error) {
    // ❌ 错误1：直接使用SQL
    _, err := l.svcCtx.DB.Exec("INSERT INTO categories (name, code) VALUES (?, ?)", 
        req.Name, req.Code)
    
    // ❌ 错误2：包含HTTP相关逻辑
    if err != nil {
        http.Error(w, "failed", 500) // Logic不应该知道HTTP
    }
    
    return nil, nil
}
```

---

### 3. Model层 (持久化层)

#### 位置
```
model/resource_catalog/category/
├── interface.go      # 接口定义
├── types.go          # 数据结构
├── vars.go           # 常量/错误
├── factory.go        # ORM工厂
├── gorm_dao.go       # GORM实现
└── sqlx_model.go     # SQLx实现
```

#### 职责

**✅ 应该做的事情**：
1. 定义数据访问接口
2. 实现CRUD操作
3. 事务管理
4. 复杂查询封装
5. 数据库错误处理

**❌ 不应该做的事情**：
1. ❌ 实现业务逻辑
2. ❌ 了解上层业务概念
3. ❌ 直接返回数据库错误
4. ❌ 包含HTTP或其他非持久化相关的代码

#### 代码示例

```go
// ✅ 好的Model接口定义
package category

import "context"

type Model interface {
    // CRUD操作
    Insert(ctx context.Context, data *Category) (*Category, error)
    FindOne(ctx context.Context, id int64) (*Category, error)
    FindByCode(ctx context.Context, code string) (*Category, error)
    Update(ctx context.Context, data *Category) error
    Delete(ctx context.Context, id int64) error
    
    // 列表查询
    FindAll(ctx context.Context) ([]*Category, error)
    List(ctx context.Context, page, pageSize int) ([]*Category, int64, error)
    
    // 事务支持
    WithTx(tx interface{}) Model
    Trans(ctx context.Context, fn func(ctx context.Context, model Model) error) error
}
```

```go
// ✅ 好的Model实现 (GORM示例)
package category

import (
    "context"
    "gorm.io/gorm"
)

type CategoryDao struct {
    db *gorm.DB
}

// Insert 插入新记录
func (d *CategoryDao) Insert(ctx context.Context, data *Category) (*Category, error) {
    if err := d.db.WithContext(ctx).Create(data).Error; err != nil {
        return nil, err
    }
    return data, nil
}

// FindOne 根据ID查询
func (d *CategoryDao) FindOne(ctx context.Context, id int64) (*Category, error) {
    var category Category
    err := d.db.WithContext(ctx).
        Where("id = ?", id).
        First(&category).Error
    
    if err == gorm.ErrRecordNotFound {
        return nil, ErrNotFound
    }
    if err != nil {
        return nil, err
    }
    
    return &category, nil
}

// Trans 事务处理
func (d *CategoryDao) Trans(ctx context.Context, fn func(ctx context.Context, model Model) error) error {
    return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        txModel := &CategoryDao{db: tx}
        return fn(ctx, txModel)
    })
}
```

```go
// ❌ 错误的Model实现
func (d *CategoryDao) Insert(ctx context.Context, data *Category) (*Category, error) {
    // ❌ 错误1：包含业务逻辑
    if data.Status == 0 {
        data.Status = 1 // 这是业务规则，不应该在Model层
    }
    
    // ❌ 错误2：直接panic
    if err := d.db.Create(data).Error; err != nil {
        panic(err) // 应该返回error
    }
    
    // ❌ 错误3：返回数据库特定错误
    return data, nil // 应该转换为业务错误
}
```

---

## 🔗 层间依赖规则

### 依赖方向（单向）

```
Handler → Logic → Model → Database
  ↓        ↓        ↓
不可反向依赖！
```

**规则**：
1. ✅ 上层可以依赖下层
2. ❌ 下层**绝对不可以**依赖上层
3. ✅ 同层之间可以调用（避免循环依赖）
4. ✅ 通过接口解耦

### 跨层通信

#### ✅ 正确的做法：使用接口

```go
// ServiceContext中使用接口
type ServiceContext struct {
    CategoryModel category.Model  // 接口类型
}

// Logic层使用接口
func (l *Logic) DoSomething() error {
    // 通过接口调用，不关心具体实现
    result, err := l.svcCtx.CategoryModel.Insert(l.ctx, data)
    return err
}
```

#### ❌ 错误的做法：直接依赖具体实现

```go
// ❌ 不好
type ServiceContext struct {
    CategoryDao *gorm.CategoryDao  // 具体类型
}
```

---

## 📊 数据流转示例

### 创建Category的完整流程

```
用户请求
   ↓
[HTTP POST /api/v1/catalog/categories]
   ↓
┌─────────────────────────────────────┐
│  Handler Layer                      │
│  1. 解析JSON请求体                   │
│  2. 校验基础参数（必填项、类型）       │
│  3. 创建Logic实例                    │
│  4. 调用Logic.CreateCategory()      │
└─────────────────────────────────────┘
   ↓ types.CreateCategoryReq
┌─────────────────────────────────────┐
│  Logic Layer                        │
│  1. 业务验证（长度、格式等）          │
│  2. 检查code唯一性                   │
│  3. 数据转换 (Req -> Model)          │
│  4. 调用Model.Insert()              │
│  5. 数据转换 (Model -> Resp)         │
│  6. 记录日志                        │
└─────────────────────────────────────┘
   ↓ category.Category
┌─────────────────────────────────────┐
│  Model Layer                        │
│  1. 执行SQL INSERT                  │
│  2. 返回插入后的记录（含ID）          │
│  3. 错误转换                        │
└─────────────────────────────────────┘
   ↓
[MySQL Database]
```

---

## ⚠️ 常见错误和反模式

### 错误1：跨层直接访问

```go
// ❌ 错误：Handler直接访问Model
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
    // 跳过Logic层
    result, _ := h.svcCtx.CategoryModel.Insert(ctx, data)
}
```

**正确做法**：
```go
// ✅ 正确：通过Logic层
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
    l := logic.NewCreateCategoryLogic(r.Context(), h.svcCtx)
    result, err := l.CreateCategory(&req)
}
```

### 错误2：Logic层包含SQL

```go
// ❌ 错误：Logic层直接写SQL
func (l *Logic) CreateCategory(req *types.Req) error {
    _, err := l.svcCtx.DB.Exec("INSERT INTO categories ...")
    return err
}
```

**正确做法**：
```go
// ✅ 正确：调用Model层
func (l *Logic) CreateCategory(req *types.Req) error {
    result, err := l.svcCtx.CategoryModel.Insert(l.ctx, data)
    return err
}
```

### 错误3：Model层包含业务逻辑

```go
// ❌ 错误：Model层验证业务规则
func (d *Dao) Insert(ctx context.Context, data *Category) error {
    if data.Name == "admin" {
        return errors.New("reserved name") // 这是业务规则
    }
    return d.db.Create(data).Error
}
```

**正确做法**：
```go
// ✅ 正确：业务验证在Logic层
func (l *Logic) CreateCategory(req *types.Req) error {
    // 业务规则在Logic层检查
    if req.Name == "admin" {
        return errors.New("reserved name")
    }
    return l.svcCtx.CategoryModel.Insert(l.ctx, data)
}
```

---

## 🎯 最佳实践

### 1. 保持层的纯净

每一层只做自己的事情，不要越界。

### 2. 使用接口解耦

所有跨层调用都通过接口，便于测试和替换实现。

### 3. 明确数据转换边界

- Handler → Logic: `types.Req`
- Logic → Model: `model.Entity`
- Model → Logic: `model.Entity`
- Logic → Handler: `types.Resp`

### 4. 集中错误处理

每层处理自己的错误，向上传递时添加上下文信息。

### 5. 日志记录

- Handler: 记录请求/响应
- Logic: 记录业务操作
- Model: 记录数据库操作（调试时）

---

## 🤔 FAQ

**Q: 简单的CRUD也要这样分层吗？**  
A: 是的。即使现在很简单，未来可能会变复杂。保持一致的架构有助于长期维护。

**Q: Logic层可以调用多个Model吗？**  
A: 可以。流程编排是Logic层的职责之一。

**Q: Model层可以调用其他Model吗？**  
A: 不建议。如果需要组合多个数据源，应该在Logic层完成。

**Q: 可以在Handler中直接返回Model数据吗？**  
A: 不推荐。应该在Logic层转换为`types.Resp`后返回，避免暴露内部数据结构。

---

## 📚 相关文档

- [Constitution](../constitution.md) - 项目宪章
- [双ORM模式](./dual-orm-pattern.md) - Model层双ORM实现
- [API设计指南](./api-design-guide.md) - RESTful API规范

---

**本文档定义了IDRM项目的核心架构模式，所有开发活动必须遵循此架构。**
