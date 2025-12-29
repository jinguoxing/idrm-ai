package data_catalog_column

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

func init() {
	// 注册GORM工厂函数
	RegisterGormFactory(func(db *gorm.DB) Model {
		return &gormModel{db: db}
	})
}

// gormModel GORM实现的数据访问对象
type gormModel struct {
	db *gorm.DB
}

// Insert 批量插入信息项记录
func (m *gormModel) Insert(ctx context.Context, data []*DataCatalogColumn) error {
	if len(data) == 0 {
		return nil
	}

	if err := m.db.WithContext(ctx).Create(data).Error; err != nil {
		return fmt.Errorf("failed to insert data catalog columns: %w", err)
	}

	return nil
}

// FindOne 根据ID查询单个信息项
func (m *gormModel) FindOne(ctx context.Context, id string) (*DataCatalogColumn, error) {
	var column DataCatalogColumn

	err := m.db.WithContext(ctx).Where("id = ?", id).First(&column).Error
	if err == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find data catalog column by id %s: %w", id, err)
	}

	return &column, nil
}

// Update 更新信息项记录
func (m *gormModel) Update(ctx context.Context, data *DataCatalogColumn) error {
	err := m.db.WithContext(ctx).Model(&DataCatalogColumn{}).
		Where("id = ?", data.Id).
		Updates(data).Error

	if err != nil {
		return fmt.Errorf("failed to update data catalog column %s: %w", data.Id, err)
	}

	return nil
}

// Delete 删除信息项记录
func (m *gormModel) Delete(ctx context.Context, id string) error {
	result := m.db.WithContext(ctx).Delete(&DataCatalogColumn{}, id)

	if result.Error != nil {
		return fmt.Errorf("failed to delete data catalog column %s: %w", id, result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// FindByCatalogId 根据目录ID查询信息项列表
func (m *gormModel) FindByCatalogId(ctx context.Context, catalogId string) ([]*DataCatalogColumn, error) {
	var columns []*DataCatalogColumn

	err := m.db.WithContext(ctx).
		Where("catalog_id = ?", catalogId).
		Order("sort_order ASC, id ASC").
		Find(&columns).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find columns by catalog id %s: %w", catalogId, err)
	}

	return columns, nil
}

// DeleteByCatalogId 根据目录ID删除所有信息项
func (m *gormModel) DeleteByCatalogId(ctx context.Context, catalogId string) error {
	result := m.db.WithContext(ctx).
		Where("catalog_id = ?", catalogId).
		Delete(&DataCatalogColumn{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete columns by catalog id %s: %w", catalogId, result.Error)
	}

	return nil
}

// BatchInsert 批量插入信息项
func (m *gormModel) BatchInsert(ctx context.Context, columns []*DataCatalogColumn) error {
	if len(columns) == 0 {
		return nil
	}

	// 使用批量插入提高性能
	if err := m.db.WithContext(ctx).CreateInBatches(columns, 100).Error; err != nil {
		return fmt.Errorf("failed to batch insert data catalog columns: %w", err)
	}

	return nil
}

// ========== 事务方法实现 ==========

// WithTx 创建事务副本
func (m *gormModel) WithTx(tx interface{}) Model {
	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		panic("invalid transaction type, expected *gorm.DB")
	}

	return &gormModel{db: gormTx}
}

// Trans 在事务中执行多个操作
func (m *gormModel) Trans(ctx context.Context, fn func(ctx context.Context, model Model) error) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		modelTx := &gormModel{db: tx}
		return fn(ctx, modelTx)
	})
}
