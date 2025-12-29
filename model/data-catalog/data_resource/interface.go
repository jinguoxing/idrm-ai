package data_resource

import "context"

// Model 数据资源挂接数据访问接口
type Model interface {
	// Insert 插入资源记录
	Insert(ctx context.Context, data *DataResource) error

	// FindOne 根据ID查询单个资源
	FindOne(ctx context.Context, id string) (*DataResource, error)

	// Update 更新资源记录
	Update(ctx context.Context, data *DataResource) error

	// Delete 删除资源记录
	Delete(ctx context.Context, id string) error

	// ========== 查询方法 ==========

	// FindByCatalogId 根据目录ID查询资源列表
	FindByCatalogId(ctx context.Context, catalogId string) ([]*DataResource, error)

	// FindByCatalogIdAndType 根据目录ID和资源类型查询
	FindByCatalogIdAndType(ctx context.Context, catalogId string, resourceType ResourceType) ([]*DataResource, error)

	// DeleteByCatalogId 根据目录ID删除所有资源
	DeleteByCatalogId(ctx context.Context, catalogId string) error

	// CountByCatalogId 统计目录下的资源数量
	CountByCatalogId(ctx context.Context, catalogId string) (int64, error)

	// ========== 业务方法 ==========

	// CountLogicalViews 统计目录下的逻辑视图数量（最多只能有一个）
	CountLogicalViews(ctx context.Context, catalogId string) (int64, error)

	// ========== 事务方法 ==========

	// WithTx 创建事务副本
	WithTx(tx interface{}) Model

	// Trans 在事务中执行多个操作
	Trans(ctx context.Context, fn func(ctx context.Context, model Model) error) error
}
