package menu

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

func init() {
	// 注册GORM工厂函数
	RegisterGormFactory(func(db *gorm.DB) Model {
		return &MenuDao{db: db}
	})
}

// MenuDao GORM实现的菜单数据访问对象
type MenuDao struct {
	db *gorm.DB
}

// Insert 插入新的菜单记录
func (d *MenuDao) Insert(ctx context.Context, data *Menu) (*Menu, error) {
	if err := d.db.WithContext(ctx).Create(data).Error; err != nil {
		return nil, fmt.Errorf("failed to insert menu: %w", err)
	}
	return data, nil
}

// FindOne 根据ID查询单个菜单
func (d *MenuDao) FindOne(ctx context.Context, id int64) (*Menu, error) {
	var menu Menu
	err := d.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&menu).Error

	if err == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find menu by id %d: %w", id, err)
	}

	return &menu, nil
}

// FindByPermTag 根据权限标识查询菜单
func (d *MenuDao) FindByPermTag(ctx context.Context, permTag string) (*Menu, error) {
	var menu Menu
	err := d.db.WithContext(ctx).
		Where("perm_tag = ? AND deleted_at IS NULL", permTag).
		First(&menu).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find menu by perm_tag %s: %w", permTag, err)
	}

	return &menu, nil
}

// Update 更新菜单记录
func (d *MenuDao) Update(ctx context.Context, data *Menu) error {
	result := d.db.WithContext(ctx).
		Model(&Menu{}).
		Where("id = ? AND deleted_at IS NULL", data.Id).
		Updates(data)

	if result.Error != nil {
		return fmt.Errorf("failed to update menu %d: %w", data.Id, result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// Delete 删除菜单记录（软删除）
func (d *MenuDao) Delete(ctx context.Context, id int64) error {
	result := d.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&Menu{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete menu %d: %w", id, result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// FindAll 查询所有菜单（支持状态过滤）
func (d *MenuDao) FindAll(ctx context.Context, status int) ([]*Menu, error) {
	var menus []*Menu

	query := d.db.WithContext(ctx).
		Where("deleted_at IS NULL")

	if status != -1 {
		query = query.Where("status = ?", status)
	}

	if err := query.
		Order("sort_order ASC, id ASC").
		Find(&menus).Error; err != nil {
		return nil, fmt.Errorf("failed to find all menus: %w", err)
	}

	return menus, nil
}

// FindByParentId 根据父级ID查询子菜单列表
func (d *MenuDao) FindByParentId(ctx context.Context, parentId int64) ([]*Menu, error) {
	var menus []*Menu

	err := d.db.WithContext(ctx).
		Where("parent_id = ? AND deleted_at IS NULL", parentId).
		Order("sort_order ASC, id ASC").
		Find(&menus).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find menus by parent_id %d: %w", parentId, err)
	}

	return menus, nil
}

// CountByParentId 统计指定父级下的子菜单数量
func (d *MenuDao) CountByParentId(ctx context.Context, parentId int64) (int64, error) {
	var count int64

	err := d.db.WithContext(ctx).
		Model(&Menu{}).
		Where("parent_id = ? AND deleted_at IS NULL", parentId).
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to count menus by parent_id %d: %w", parentId, err)
	}

	return count, nil
}

// HasChildren 检查菜单是否有子菜单
func (d *MenuDao) HasChildren(ctx context.Context, parentId int64) (bool, error) {
	count, err := d.CountByParentId(ctx, parentId)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// WithTx 创建事务副本
func (d *MenuDao) WithTx(tx interface{}) Model {
	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		panic(fmt.Sprintf("invalid transaction type: expected *gorm.DB, got %T", tx))
	}
	return &MenuDao{db: gormTx}
}

// Trans 在事务中执行多个操作
func (d *MenuDao) Trans(ctx context.Context, fn func(ctx context.Context, model Model) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txModel := &MenuDao{db: tx}
		return fn(ctx, txModel)
	})
}
