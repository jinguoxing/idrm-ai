package data_catalog_column

import "context"

// Model 数据资源目录信息项数据访问接口
// 支持GORM和SQLx双ORM实现，通过工厂模式自动选择
type Model interface {
	// ========== CRUD基础方法 ==========

	// Insert 批量插入信息项记录
	Insert(ctx context.Context, data []*DataCatalogColumn) error

	// FindOne 根据ID查询单个信息项
	FindOne(ctx context.Context, id string) (*DataCatalogColumn, error)

	// Update 更新信息项记录
	Update(ctx context.Context, data *DataCatalogColumn) error

	// Delete 删除信息项记录
	Delete(ctx context.Context, id string) error

	// ========== 查询方法 ==========

	// FindByCatalogId 根据目录ID查询信息项列表
	FindByCatalogId(ctx context.Context, catalogId string) ([]*DataCatalogColumn, error)

	// DeleteByCatalogId 根据目录ID删除所有信息项
	DeleteByCatalogId(ctx context.Context, catalogId string) error

	// BatchInsert 批量插入信息项
	BatchInsert(ctx context.Context, columns []*DataCatalogColumn) error

	// ========== 事务方法 ==========

	// WithTx 创建事务副本
	WithTx(tx interface{}) Model

	// Trans 在事务中执行多个操作
	Trans(ctx context.Context, fn func(ctx context.Context, model Model) error) error
}
