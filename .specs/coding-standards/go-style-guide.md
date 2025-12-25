# Go语言风格指南

> **文档版本**: v1.0 (大纲版)  
> **最后更新**: 2025-12-24  
> **状态**: 📝 待完善

---

## 代码组织

### 包结构

```go
package mypackage

// 1. 常量
const (
    MaxRetry = 3
)

// 2. 变量
var (
    ErrNotFound = errors.New("not found")
)

// 3. 类型定义
type MyStruct struct {
    Field string
}

// 4. 构造函数
func NewMyStruct() *MyStruct {
    return &MyStruct{}
}

// 5. 公开方法
func (m *MyStruct) PublicMethod() {}

// 6. 私有方法
func (m *MyStruct) privateMethod() {}

// 7. init函数
func init() {}
```

### 导入分组

```go
import (
    // 标准库
    "context"
    "fmt"
    
    // 第三方库
    "github.com/zeromicro/go-zero/core/logx"
    "gorm.io/gorm"
    
    // 项目内部
    "idrm/model/resource_catalog/category"
    "idrm/pkg/response"
)
```

---

## 命名规范

- **文件**: 全小写下划线 `create_category_logic.go`
- **包**: 全小写简短 `category`, `middleware`
- **类型**: 大驼峰 `CategoryModel`, `UserInfo`
- **函数**: 公开大驼峰，私有小驼峰
- **常量**: 全大写下划线 `MAX_RETRY_COUNT`

---

## 代码风格

### 变量声明

```go
// ✅ 好
var name string
count := 0
items := make([]string, 0, 10)

// ❌ 不好
var name string = ""  // 不要显式零值
var count int = 0
```

### 错误处理

```go
// ✅ 好
result, err := doSomething()
if err != nil {
    return fmt.Errorf("failed: %w", err)
}

// ❌ 不好
result, _ := doSomething()  // 永不忽略error
```

### 函数长度

- 理想: <20行
- 警戒: <50行
- 超过50行应考虑拆分

---

## 注释规范

```go
// CreateCategory 创建新的类别记录
// 参数name不能为空，code必须唯一
// 返回创建后的类别ID和error
func CreateCategory(name, code string) (int64, error) {
    // 实现...
}
```

---

## 最佳实践

### 使用context

```go
// ✅ 好
func doWork(ctx context.Context) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
        // do work
    }
}
```

### 并发控制

```go
var wg sync.WaitGroup
for _, item := range items {
    wg.Add(1)
    go func(i Item) {
        defer wg.Done()
        process(i)
    }(item)
}
wg.Wait()
```

---

## 📌 待补充内容

- [ ] 并发编程详解
- [ ] Interface设计原则
- [ ] 性能优化建议
- [ ] 完整代码示例

---

**参考**: [命名规范](./naming-conventions.md) | [错误处理](./error-handling.md)
