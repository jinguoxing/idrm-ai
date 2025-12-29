package data_catalog_category

import "context"

// Model 数据资源目录分类关联数据访问接口
type Model interface {
	// Insert 批量插入分类关联记录
	Insert(ctx context.Context, data []*DataCatalogCategory) error

	// FindByCatalogId 根据目录ID查询分类关联列表
	FindByCatalogId(ctx context.Context, catalogId string) ([]*DataCatalogCategory, error)

	// DeleteByCatalogId 根据目录ID删除所有分类关联
	DeleteByCatalogId(ctx context.Context, catalogId string) error

	// BatchDelete 批量删除分类关联
	BatchDelete(ctx context.Context, ids []string) error

	// ========== 事务方法 ==========

	// WithTx 创建事务副本
	WithTx(tx interface{}) Model

	// Trans 在事务中执行多个操作
	Trans(ctx context.Context, fn func(ctx context.Context, model Model) error) error
}
