package data_catalog_category

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

func (m *gormModel) Insert(ctx context.Context, data []*DataCatalogCategory) error {
	if len(data) == 0 {
		return nil
	}

	if err := m.db.WithContext(ctx).Create(data).Error; err != nil {
		return fmt.Errorf("failed to insert data catalog categories: %w", err)
	}

	return nil
}

func (m *gormModel) FindByCatalogId(ctx context.Context, catalogId string) ([]*DataCatalogCategory, error) {
	var categories []*DataCatalogCategory

	err := m.db.WithContext(ctx).
		Where("catalog_id = ?", catalogId).
		Find(&categories).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find categories by catalog id %s: %w", catalogId, err)
	}

	return categories, nil
}

func (m *gormModel) DeleteByCatalogId(ctx context.Context, catalogId string) error {
	err := m.db.WithContext(ctx).
		Where("catalog_id = ?", catalogId).
		Delete(&DataCatalogCategory{}).Error

	if err != nil {
		return fmt.Errorf("failed to delete categories by catalog id %s: %w", catalogId, err)
	}

	return nil
}

func (m *gormModel) BatchDelete(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	err := m.db.WithContext(ctx).
		Where("id IN ?", ids).
		Delete(&DataCatalogCategory{}).Error

	if err != nil {
		return fmt.Errorf("failed to batch delete categories: %w", err)
	}

	return nil
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
