package data_resource

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

func init() {
	RegisterGormFactory(func(db *gorm.DB) Model {
		return &gormModel{db: db}
	})
}

type gormModel struct {
	db *gorm.DB
}

func (m *gormModel) Insert(ctx context.Context, data *DataResource) error {
	if err := m.db.WithContext(ctx).Create(data).Error; err != nil {
		return fmt.Errorf("failed to insert data resource: %w", err)
	}
	return nil
}

func (m *gormModel) FindOne(ctx context.Context, id string) (*DataResource, error) {
	var resource DataResource

	err := m.db.WithContext(ctx).Where("id = ?", id).First(&resource).Error
	if err == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find data resource by id %s: %w", id, err)
	}

	return &resource, nil
}

func (m *gormModel) Update(ctx context.Context, data *DataResource) error {
	err := m.db.WithContext(ctx).Model(&DataResource{}).
		Where("id = ?", data.Id).
		Updates(data).Error

	if err != nil {
		return fmt.Errorf("failed to update data resource %s: %w", data.Id, err)
	}

	return nil
}

func (m *gormModel) Delete(ctx context.Context, id string) error {
	result := m.db.WithContext(ctx).Delete(&DataResource{}, id)

	if result.Error != nil {
		return fmt.Errorf("failed to delete data resource %s: %w", id, result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (m *gormModel) FindByCatalogId(ctx context.Context, catalogId string) ([]*DataResource, error) {
	var resources []*DataResource

	err := m.db.WithContext(ctx).
		Where("catalog_id = ?", catalogId).
		Order("sort_order ASC, id ASC").
		Find(&resources).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find resources by catalog id %s: %w", catalogId, err)
	}

	return resources, nil
}

func (m *gormModel) FindByCatalogIdAndType(ctx context.Context, catalogId string, resourceType ResourceType) ([]*DataResource, error) {
	var resources []*DataResource

	err := m.db.WithContext(ctx).
		Where("catalog_id = ? AND resource_type = ?", catalogId, resourceType).
		Order("sort_order ASC, id ASC").
		Find(&resources).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find resources by catalog id %s and type %s: %w", catalogId, resourceType, err)
	}

	return resources, nil
}

func (m *gormModel) DeleteByCatalogId(ctx context.Context, catalogId string) error {
	result := m.db.WithContext(ctx).
		Where("catalog_id = ?", catalogId).
		Delete(&DataResource{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete resources by catalog id %s: %w", catalogId, result.Error)
	}

	return nil
}

func (m *gormModel) CountByCatalogId(ctx context.Context, catalogId string) (int64, error) {
	var count int64

	err := m.db.WithContext(ctx).
		Model(&DataResource{}).
		Where("catalog_id = ?", catalogId).
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to count resources by catalog id %s: %w", catalogId, err)
	}

	return count, nil
}

func (m *gormModel) CountLogicalViews(ctx context.Context, catalogId string) (int64, error) {
	var count int64

	err := m.db.WithContext(ctx).
		Model(&DataResource{}).
		Where("catalog_id = ? AND resource_type = ?", catalogId, ResourceTypeLogicalView).
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to count logical views by catalog id %s: %w", catalogId, err)
	}

	return count, nil
}

func (m *gormModel) WithTx(tx interface{}) Model {
	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		panic("invalid transaction type, expected *gorm.DB")
	}

	return &gormModel{db: gormTx}
}

func (m *gormModel) Trans(ctx context.Context, fn func(ctx context.Context, model Model) error) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		modelTx := &gormModel{db: tx}
		return fn(ctx, modelTx)
	})
}
