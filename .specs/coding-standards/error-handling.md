# 错误处理规范

> **文档版本**: v1.0 (大纲版)  
> **最后更新**: 2025-12-24  
> **状态**: 📝 待完善

---

## 错误定义

### 位置: vars.go

```go
var (
    ErrNotFound          = errors.New("category not found")
    ErrCodeAlreadyExists = errors.New("code already exists")
    ErrInvalidStatus     = errors.New("invalid status")
)
```

---

## 错误封装

### 使用 %w 保留错误链

```go
// ✅ 好: 保留错误链
return fmt.Errorf("failed to create category: %w", err)

// ❌ 不好: 丢失错误链
return fmt.Errorf("failed to create category: %v", err)
```

---

## 错误检查

### 使用 errors.Is 和 errors.As

```go
// 检查特定错误
if errors.Is(err, ErrNotFound) {
    return nil  // 忽略"未找到"错误
}

// 获取错误类型
var validationErr *ValidationError
if errors.As(err, &validationErr) {
    // 处理验证错误
    log.Printf("validation failed: %s", validationErr.Field)
}
```

---

## 自定义错误

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// 使用
return &ValidationError{
    Field:   "name",
    Message: "name is required",
}
```

---

## 各层错误处理

### Model层

```go
// 转换数据库错误为业务错误
func (d *Dao) FindOne(ctx context.Context, id int64) (*Category, error) {
    var category Category
    err := d.db.First(&category, id).Error
    
    if err == gorm.ErrRecordNotFound {
        return nil, ErrNotFound  // 转换为业务错误
    }
    return &category, err
}
```

### Logic层

```go
// 添加业务上下文
func (l *Logic) CreateCategory(req *Req) error {
    result, err := l.svcCtx.CategoryModel.Insert(l.ctx, data)
    if err != nil {
        return fmt.Errorf("failed to create category %s: %w", req.Name, err)
    }
    return nil
}
```

### Handler层

```go
// 统一响应格式
func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
    resp, err := logic.CreateCategory(&req)
    if err != nil {
        response.Error(w, err)  // 统一错误响应
        return
    }
    response.Success(w, resp)
}
```

---

## 日志记录

```go
if err != nil {
    l.Errorf("operation failed: %v", err)
    return fmt.Errorf("operation failed: %w", err)
}
```

---

## 📌 待补充内容

- [ ] 错误码设计方案
- [ ] 多语言错误消息
- [ ] 错误监控集成
- [ ] 完整错误处理示例

---

**参考**: [Go风格指南](./go-style-guide.md) | [Constitution](../constitution.md)
