# API设计指南

> **文档版本**: v1.0 (大纲版)  
> **最后更新**: 2025-12-24  
> **状态**: 📝 待完善

---

## RESTful 原则

### 资源导向

```
GET    /api/v1/catalog/categories        # 列表
POST   /api/v1/catalog/categories        # 创建
GET    /api/v1/catalog/categories/:id    # 详情
PUT    /api/v1/catalog/categories/:id    # 更新
DELETE /api/v1/catalog/categories/:id    # 删除
```

### HTTP动词

- GET: 查询
- POST: 创建
- PUT: 更新（全量）
- PATCH: 更新（部分）
- DELETE: 删除

---

## 请求格式

### URL参数

```
GET /api/v1/categories?status=1&page=1&pageSize=20
```

### Body (JSON)

```json
{
  "name": "类别名称",
  "code": "CODE001",
  "status": 1
}
```

---

## 统一响应格式

### 成功响应

```json
{
  "code": 0,
  "msg": "success",
  "data": {},
  "requestId": "req-123456"
}
```

### 分页响应

```json
{
  "code": 0,
  "data": {
    "items": [],
    "total": 100,
    "page": 1,
    "pageSize": 20
  }
}
```

### 错误响应

```json
{
  "code": 1001,
  "msg": "参数错误",
  "requestId": "req-123456"
}
```

---

## 错误码设计

```
0: 成功
1001: 参数错误
1002: 资源不存在
1003: 权限不足
1004: 操作冲突
5000: 服务器错误
```

---

## 版本管理

使用URL路径: `/api/v1/`, `/api/v2/`

---

## 📌 待补充内容

- [ ] 完整错误码列表
- [ ] 认证授权方案
- [ ] 限流策略
- [ ] API文档生成

---

**参考**: [分层架构](./layered-architecture.md) | [Constitution](../constitution.md)
