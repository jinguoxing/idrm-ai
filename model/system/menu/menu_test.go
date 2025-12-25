package menu

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB 创建测试数据库
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// 自动创建测试表
	err = db.AutoMigrate(&Menu{})
	assert.NoError(t, err)

	return db
}

// setupTestData 创建测试数据
func setupTestData(t *testing.T, db *gorm.DB) {
	menus := []Menu{
		{
			ParentId:  0,
			Name:      "系统管理",
			SortOrder: 100,
			Status:    1,
		},
		{
			ParentId:  1,
			Name:      "菜单管理",
			SortOrder: 1,
			Status:    1,
		},
		{
			ParentId:  1,
			Name:      "用户管理",
			SortOrder: 2,
			Status:    1,
		},
	}

	for _, menu := range menus {
		err := db.Create(&menu).Error
		assert.NoError(t, err)
	}
}

// TestMenuDao_Insert 测试插入菜单
func TestMenuDao_Insert(t *testing.T) {
	db := setupTestDB(t)
	dao := &MenuDao{db: db}
	ctx := context.Background()

	// 插入根菜单
	menuData := &Menu{
		ParentId:  0,
		Name:      "系统管理",
		RoutePath: "/system",
		SortOrder: 100,
		Status:    1,
	}

	result, err := dao.Insert(ctx, menuData)
	assert.NoError(t, err)
	assert.NotZero(t, result.Id)
	assert.Equal(t, "系统管理", result.Name)
	assert.Equal(t, int64(0), result.ParentId)

	// 插入子菜单
	childMenu := &Menu{
		ParentId:  result.Id,
		Name:      "菜单管理",
		RoutePath: "/system/menu",
		SortOrder: 1,
		Status:    1,
	}

	childResult, err := dao.Insert(ctx, childMenu)
	assert.NoError(t, err)
	assert.NotZero(t, childResult.Id)
	assert.Equal(t, result.Id, childResult.ParentId)
}

// TestMenuDao_FindOne 测试根据ID查询菜单
func TestMenuDao_FindOne(t *testing.T) {
	db := setupTestDB(t)
	setupTestData(db)
	dao := &MenuDao{db: db}
	ctx := context.Background()

	// 查询存在的菜单
	foundMenu, err := dao.FindOne(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), foundMenu.Id)
	assert.Equal(t, "系统管理", foundMenu.Name)

	// 查询不存在的菜单
	_, err = dao.FindOne(ctx, 999)
	assert.Error(t, err)
	assert.Equal(t, ErrNotFound, err)
}

// TestMenuDao_Update 测试更新菜单
func TestMenuDao_Update(t *testing.T) {
	db := setupTestDB(t)
	setupTestData(db)
	dao := &MenuDao{db: db}
	ctx := context.Background()

	// 更新菜单
	updateData := &Menu{
		Id:       1,
		Name:     "系统管理（已更新）",
		SortOrder: 200,
		Status:   1,
	}

	err := dao.Update(ctx, updateData)
	assert.NoError(t, err)

	// 验证更新结果
	updatedMenu, _ := dao.FindOne(ctx, 1)
	assert.Equal(t, "系统管理（已更新）", updatedMenu.Name)
	assert.Equal(t, 200, updatedMenu.SortOrder)

	// 更新不存在的菜单
	nonExistent := &Menu{
		Id:   999,
		Name: "不存在的菜单",
	}
	err = dao.Update(ctx, nonExistent)
	assert.Error(t, err)
	assert.Equal(t, ErrNotFound, err)
}

// TestMenuDao_Delete 测试删除菜单（软删除）
func TestMenuDao_Delete(t *testing.T) {
	db := setupTestDB(t)
	setupTestData(db)
	dao := &MenuDao{db: db}
	ctx := context.Background()

	// 删除存在的菜单
	err := dao.Delete(ctx, 1)
	assert.NoError(t, err)

	// 验证软删除（查询应该失败）
	_, err = dao.FindOne(ctx, 1)
	assert.Error(t, err)
	assert.Equal(t, ErrNotFound, err)

	// 删除不存在的菜单
	err = dao.Delete(ctx, 999)
	assert.Error(t, err)
	assert.Equal(t, ErrNotFound, err)
}

// TestMenuDao_FindByParentId 测试根据父级ID查询子菜单
func TestMenuDao_FindByParentId(t *testing.T) {
	db := setupTestDB(t)
	setupTestData(db)
	dao := &MenuDao{db: db}
	ctx := context.Background()

	// 查询根菜单（parentId = 0）
	rootMenus, err := dao.FindByParentId(ctx, 0)
	assert.NoError(t, err)
	assert.Len(t, rootMenus, 1) // 只有"系统管理"是根菜单

	// 查询"系统管理"的子菜单
	childMenus, err := dao.FindByParentId(ctx, 1)
	assert.NoError(t, err)
	assert.Len(t, childMenus, 2) // "菜单管理"和"用户管理"

	// 验证排序
	assert.Equal(t, "菜单管理", childMenus[0].Name)
	assert.Equal(t, "用户管理", childMenus[1].Name)
}

// TestMenuDao_HasChildren 测试检查是否有子菜单
func TestMenuDao_HasChildren(t *testing.T) {
	db := setupTestDB(t)
	setupTestData(db)
	dao := &MenuDao{db: db}
	ctx := context.Background()

	// "系统管理"有子菜单
	hasChildren, err := dao.HasChildren(ctx, 1)
	assert.NoError(t, err)
	assert.True(t, hasChildren)

	// "菜单管理"没有子菜单
	hasChildren, err = dao.HasChildren(ctx, 2)
	assert.NoError(t, err)
	assert.False(t, hasChildren)
}

// TestMenuDao_CountByParentId 测试统计子菜单数量
func TestMenuDao_CountByParentId(t *testing.T) {
	db := setupTestDB(t)
	setupTestData(db)
	dao := &MenuDao{db: db}
	ctx := context.Background()

	// 统计"系统管理"的子菜单数量
	count, err := dao.CountByParentId(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)

	// 统计"菜单管理"的子菜单数量
	count, err = dao.CountByParentId(ctx, 2)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

// TestMenuDao_FindAll 测试查询所有菜单
func TestMenuDao_FindAll(t *testing.T) {
	db := setupTestDB(t)
	setupTestData(db)
	dao := &MenuDao{db: db}
	ctx := context.Background()

	// 查询所有菜单（不过滤状态）
	allMenus, err := dao.FindAll(ctx, -1)
	assert.NoError(t, err)
	assert.Len(t, allMenus, 3)

	// 只查询启用的菜单
	enabledMenus, err := dao.FindAll(ctx, 1)
	assert.NoError(t, err)
	assert.Len(t, enabledMenus, 3)
}

// TestMenuDao_Trans 测试事务
func TestMenuDao_Trans(t *testing.T) {
	db := setupTestDB(t)
	dao := &MenuDao{db: db}
	ctx := context.Background()

	// 测试事务提交
	err := dao.Trans(ctx, func(ctx context.Context, model Model) error {
		_, err := model.Insert(ctx, &Menu{
			ParentId:  0,
			Name:      "事务测试菜单1",
			SortOrder: 1,
			Status:    1,
		})
		assert.NoError(t, err)

		_, err = model.Insert(ctx, &Menu{
			ParentId:  0,
			Name:      "事务测试菜单2",
			SortOrder: 2,
			Status:    1,
		})
		return err
	})
	assert.NoError(t, err)

	// 验证数据已插入
	allMenus, _ := dao.FindAll(ctx, -1)
	assert.Len(t, allMenus, 2)

	// 测试事务回滚
	err = dao.Trans(ctx, func(ctx context.Context, model Model) error {
		_, err := model.Insert(ctx, &Menu{
			ParentId:  0,
			Name:      "应该回滚的菜单",
			SortOrder: 3,
			Status:    1,
		})
		assert.NoError(t, err)

		// 故意返回错误以触发回滚
		return assert.AnError
	})
	assert.Error(t, err)

	// 验证数据未插入（应该还是2条）
	allMenus, _ = dao.FindAll(ctx, -1)
	assert.Len(t, allMenus, 2)
}

// TestMenuDao_WithTx 测试创建事务副本
func TestMenuDao_WithTx(t *testing.T) {
	db := setupTestDB(t)
	dao := &MenuDao{db: db}
	ctx := context.Background()

	// 开始事务
	err := db.Transaction(func(tx *gorm.DB) error {
		txModel := dao.WithTx(tx)

		// 在事务中插入数据
		_, err := txModel.Insert(ctx, &Menu{
			ParentId:  0,
			Name:      "事务副本测试",
			SortOrder: 1,
			Status:    1,
		})
		assert.NoError(t, err)

		// 在事务中查询
		menus, err := txModel.FindByParentId(ctx, 0)
		assert.NoError(t, err)
		assert.Len(t, menus, 1)

		return nil
	})
	assert.NoError(t, err)
}

// BenchmarkMenuDao_Insert 性能测试：插入操作
func BenchmarkMenuDao_Insert(b *testing.B) {
	db := setupTestDB(&testing.T{})
	dao := &MenuDao{db: db}
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := dao.Insert(ctx, &Menu{
			ParentId:  0,
			Name:      "性能测试菜单",
			SortOrder: i,
			Status:    1,
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}
