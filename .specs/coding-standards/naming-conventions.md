# 命名规范

> **文档版本**: v1.0 (大纲版)  
> **最后更新**: 2025-12-24  
> **状态**: 📝 待完善

---

## 文件命名

### Handler/Logic文件

```
{action}{resource}handler.go
{action}{resource}logic.go

示例:
createcategoryhandler.go
getcategorylogic.go
listcategorylogic.go
```

### Model文件

```
gorm_dao.go      # GORM实现
sqlx_model.go    # SQLx实现
interface.go     # 接口定义
types.go         # 数据结构
vars.go          # 常量错误
factory.go       # ORM工厂
```

---

## 包命名

### 规则

- **全小写**，不使用下划线或驼峰
- **简短有意义**: `category`, `middleware`, `validator`
- **避免**: `util`, `common`, `base` (太泛化)

### 示例

```go
package category      // ✅ 好
package middleware    // ✅ 好

package utils         // ❌ 太泛化
package myPackage     // ❌ 不要驼峰
```

---

## 类型命名

### 结构体

```go
type CategoryModel struct {}  // 大驼峰PascalCase
type UserInfo struct {}
type HTTPClient struct {}     // 缩写全大写
```

### 接口

```go
type Model interface {}       // 简洁名称
type Repository interface {}
type Reader interface {}
```

---

## 函数命名

### 公开函数

```go
func NewCategoryModel() Model {}  // 大驼峰
func CreateCategory() error {}
func GetUserByID(id int64) {}
```

### 私有函数

```go
func validateInput() error {}     // 小驼峰
func buildQuery() string {}
func parseResponse() {}
```

---

## 变量命名

### 局部变量

```go
category := &Category{}  // 小驼峰
userId := 123
httpClient := &http.Client{}
```

### 全局变量/常量

```go
var ErrNotFound = errors.New("not found")
var DefaultTimeout = 30 * time.Second

const MaxRetry = 3
const API_VERSION = "v1"
```

---

## 数据库命名

### 表名

```
categories      # 复数，下划线
user_profiles
data_views
```

### 字段名

```
id
category_name
created_at
updated_at
```

---

## API路径命名

```
/api/v1/catalog/categories       # 复数，小写
/api/v1/catalog/categories/:id
/api/v1/users/:userId/profiles
```

---

## 📌 待补充内容

- [ ] 更多命名示例
- [ ] 反模式说明
- [ ] 领域特定命名
- [ ] 缩写使用规则

---

**参考**: [Go风格指南](./go-style-guide.md) | [Constitution](../constitution.md)
