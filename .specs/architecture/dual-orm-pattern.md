# 双ORM模式详解

> **文档版本**: v1.0 (大纲版)  
> **最后更新**: 2025-12-24  
> **状态**: 📝 待完善

---

## 概述

IDRM Model层采用双ORM设计：同时支持**GORM**和**SQLx**，通过工厂模式实现灵活切换。

### 为什么需要双ORM？

- **GORM**: 功能丰富，适合复杂查询和关联
- **SQLx**: 轻量高效，适合简单CRUD和性能敏感场景
- **灵活切换**: 根据场景选择最合适的ORM

---

## 目录结构

```
model/resource_catalog/category/
├── interface.go      # 统一接口
├── types.go          # 共享数据结构
├── factory.go        # ORM工厂
├── gorm_dao.go       # GORM实现
└── sqlx_model.go     # SQLx实现
```

---

## 核心设计

### 1. 统一接口

```go
type Model interface {
    Insert(ctx context.Context, data *T) (*T, error)
    FindOne(ctx context.Context, id int64) (*T, error)
    Update(ctx context.Context, data *T) error
    Delete(ctx context.Context, id int64) error
    WithTx(tx interface{}) Model
    Trans(ctx context.Context, fn func(...) error) error
}
```

### 2. 工厂模式

```go
func NewModel(sqlConn *sql.DB, gormDB *gorm.DB) Model {
    if gormDB != nil && gormFactory != nil {
        return gormFactory(gormDB)  // 优先GORM
    }
    if sqlConn != nil && sqlxFactory != nil {
        return sqlxFactory(sqlConn)  // 降级SQLx
    }
    panic("no database connection available")
}
```

### 3. 自动注册

每个实现在init()中注册自己的工厂函数。

```go
func init() {
    RegisterGormFactory(newGormDao)
}
```

---

## GORM vs SQLx 对比

| 特性 | GORM | SQLx |
|-----|------|------|
| 学习曲线 | 中等 | 低 |
| 功能丰富度 | 高 | 中 |
| 性能 | 中 | 高 |
| 类型安全 | 强 | 中 |
| 适用场景 | 复杂查询、关联 | 简单CRUD、性能优先 |

---

## 使用示例

```go
// 业务层无需关心底层ORM
func (l *Logic) CreateCategory(req *types.Req) error {
    // 自动使用可用的ORM（GORM优先）
    result, err := l.svcCtx.CategoryModel.Insert(l.ctx, data)
    return err
}
```

---

## 事务处理

```go
// 统一的事务接口
err := model.Trans(ctx, func(ctx context.Context, m Model) error {
    // 在事务中执行多个操作
    _, err := m.Insert(ctx, data1)
    if err != nil {
        return err
    }
    _, err = m.Insert(ctx, data2)
    return err
})
```

---

## 📌 待补充内容

- [ ] GORM实现详解
- [ ] SQLx实现详解
- [ ] 事务处理对比
- [ ] 性能测试数据
- [ ] 迁移指南

---

**参考**: [分层架构](./layered-architecture.md) | [Constitution](../constitution.md)
