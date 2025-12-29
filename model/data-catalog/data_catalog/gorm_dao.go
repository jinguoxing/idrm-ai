package data_catalog

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

func init() {
	// 注册GORM工厂函数
	RegisterGormFactory(func(db *gorm.DB) Model {
		return &gormDao{db: db}
	})
}

// gormDao GORM实现的数据访问对象
type gormDao struct {
	db *gorm.DB
}

// Insert 插入新的目录记录
func (d *gormDao) Insert(ctx context.Context, data *DataCatalog) (*DataCatalog, error) {
	if err := d.db.WithContext(ctx).Create(data).Error; err != nil {
		return nil, fmt.Errorf("failed to insert data catalog: %w", err)
	}
	return data, nil
}

// FindOne 根据ID查询单个目录
func (d *gormDao) FindOne(ctx context.Context, id string) (*DataCatalog, error) {
	var catalog DataCatalog
	err := d.db.WithContext(ctx).
		Where("id = ?", id).
		First(&catalog).Error

	if err == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find data catalog by id %s: %w", id, err)
	}

	return &catalog, nil
}

// Update 更新目录记录
func (d *gormDao) Update(ctx context.Context, data *DataCatalog) error {
	result := d.db.WithContext(ctx).
		Model(&DataCatalog{}).
		Where("id = ?", data.Id).
		Updates(data)

	if result.Error != nil {
		return fmt.Errorf("failed to update data catalog %s: %w", data.Id, result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// Delete 删除目录记录（软删除）
func (d *gormDao) Delete(ctx context.Context, id string) error {
	result := d.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&DataCatalog{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete data catalog %s: %w", id, result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// FindBySourceDept 根据来源部门查询目录列表
func (d *gormDao) FindBySourceDept(ctx context.Context, deptId string, status string) ([]*DataCatalog, error) {
	var catalogs []*DataCatalog

	query := d.db.WithContext(ctx).
		Where("source_department_id = ?", deptId)

	if status != "" {
		query = query.Where("publish_status = ?", status)
	}

	if err := query.
		Order("created_at DESC").
		Find(&catalogs).Error; err != nil {
		return nil, fmt.Errorf("failed to find catalogs by dept %s: %w", deptId, err)
	}

	return catalogs, nil
}

// FindByDraftId 根据草稿ID查询目录
func (d *gormDao) FindByDraftId(ctx context.Context, draftId string) (*DataCatalog, error) {
	var catalog DataCatalog
	err := d.db.WithContext(ctx).
		Where("id = ?", draftId).
		First(&catalog).Error

	if err == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find catalog by draft_id %s: %w", draftId, err)
	}

	return &catalog, nil
}

// FindDrafts 查询草稿列表
func (d *gormDao) FindDrafts(ctx context.Context, deptId string, status string, page, pageSize int) ([]*DataCatalog, int64, error) {
	var catalogs []*DataCatalog
	var total int64

	query := d.db.WithContext(ctx).Model(&DataCatalog{})

	if deptId != "" {
		query = query.Where("source_department_id = ?", deptId)
	}

	if status != "" {
		query = query.Where("publish_status = ?", status)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count drafts: %w", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&catalogs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to find drafts: %w", err)
	}

	return catalogs, total, nil
}

// CheckNameExists 检查目录名称是否重复
func (d *gormDao) CheckNameExists(ctx context.Context, deptId string, title string, excludeId string) (bool, error) {
	var count int64

	query := d.db.WithContext(ctx).
		Model(&DataCatalog{}).
		Where("source_department_id = ? AND title = ?", deptId, title)

	if excludeId != "" {
		query = query.Where("id != ?", excludeId)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check name exists: %w", err)
	}

	return count > 0, nil
}

// ========== 事务方法实现 ==========

// WithTx 创建事务副本
func (d *gormDao) WithTx(tx interface{}) Model {
	// 类型断言获取GORM事务对象
	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		panic("invalid transaction type, expected *gorm.DB")
	}

	return &gormDao{db: gormTx}
}

// Trans 在事务中执行多个操作
func (d *gormDao) Trans(ctx context.Context, fn func(ctx context.Context, model Model) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(ctx, &gormDao{db: tx})
	})
}
