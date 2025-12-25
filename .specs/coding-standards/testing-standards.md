# 测试规范

> **文档版本**: v1.0 (大纲版)  
> **最后更新**: 2025-12-24  
> **状态**: 📝 待完善

---

## 覆盖率要求

- 核心业务逻辑: >80%
- 工具函数: >90%
- Handler (可选): >60%

---

## 命名规范

```go
func TestCreateCategory(t *testing.T) {}
func TestCategoryModel_Insert(t *testing.T) {}
func TestCategoryModel_Insert_DuplicateCode(t *testing.T) {}
```

---

## 表驱动测试

```go
func TestValidate(t *testing.T) {
    tests := []struct {
        name    string
        input   *Category
        wantErr bool
    }{
        {
            name:    "valid category",
            input:   &Category{Name: "test", Code: "T001"},
            wantErr: false,
        },
        {
            name:    "empty name",
            input:   &Category{Code: "T001"},
            wantErr: true,
        },
        {
            name:    "empty code",
            input:   &Category{Name: "test"},
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validate(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

---

## Mock使用

### 接口Mock

```go
type MockCategoryModel struct {
    mock.Mock
}

func (m *MockCategoryModel) Insert(ctx context.Context, data *Category) (*Category, error) {
    args := m.Called(ctx, data)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*Category), args.Error(1)
}

// 使用
func TestCreateCategoryLogic(t *testing.T) {
    mockModel := new(MockCategoryModel)
    mockModel.On("Insert", mock.Anything, mock.Anything).
        Return(&Category{Id: 1}, nil)
    
    // 测试逻辑
}
```

---

## 集成测试

### 测试数据库

```go
func setupTestDB() *gorm.DB {
    db, _ := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
    db.AutoMigrate(&Category{})
    return db
}

func TestCategoryModel_Integration(t *testing.T) {
    db := setupTestDB()
    defer db.Close()
    
    model := NewCategoryModel(db)
    // 测试...
}
```

---

## 测试命令

```bash
# 运行所有测试
go test ./...

# 带覆盖率
go test ./... -cover

# 生成覆盖率报告
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# 详细输出
go test ./... -v

# 指定包
go test ./model/resource_catalog/category/...

# 运行特定测试
go test -run TestCreateCategory
```

---

## 📌 待补充内容

- [ ] Benchmark测试
- [ ] 测试数据管理
- [ ] 并发测试
- [ ] 性能测试
- [ ] 完整示例

---

**参考**: [Go风格指南](./go-style-guide.md) | [Constitution](../constitution.md)
