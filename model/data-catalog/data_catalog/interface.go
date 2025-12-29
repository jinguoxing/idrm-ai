package data_catalog

import "context"

// Model 数据资源目录数据访问接口
// 支持GORM和SQLx双ORM实现，通过工厂模式自动选择
type Model interface {
	// ========== CRUD基础方法 ==========

	// Insert 插入新的目录记录
	// 参数data中的Id会被自动生成，其他字段保持不变
	// 返回插入后的完整记录
	Insert(ctx context.Context, data *DataCatalog) (*DataCatalog, error)

	// FindOne 根据ID查询单个目录
	// 若不存在则返回 ErrNotFound
	FindOne(ctx context.Context, id string) (*DataCatalog, error)

	// Update 更新目录记录
	// 参数data中必须包含Id字段
	// 只更新非零值字段
	Update(ctx context.Context, data *DataCatalog) error

	// Delete 删除目录记录（软删除）
	// 实际操作是设置deleted_at字段
	Delete(ctx context.Context, id string) error

	// ========== 查询方法 ==========

	// FindBySourceDept 根据来源部门查询目录列表
	FindBySourceDept(ctx context.Context, deptId string, status string) ([]*DataCatalog, error)

	// FindByDraftId 根据草稿ID查询目录
	FindByDraftId(ctx context.Context, draftId string) (*DataCatalog, error)

	// FindDrafts 查询草稿列表
	FindDrafts(ctx context.Context, deptId string, status string, page, pageSize int) ([]*DataCatalog, int64, error)

	// CheckNameExists 检查目录名称是否重复
	CheckNameExists(ctx context.Context, deptId string, title string, excludeId string) (bool, error)

	// ========== 事务方法 ==========

	// WithTx 创建事务副本
	// 参数tx是事务对象（GORM的*gorm.DB或SQLx的*sqlx.Tx）
	// 返回使用该事务的Model副本
	WithTx(tx interface{}) Model

	// Trans 在事务中执行多个操作
	// 若fn返回error，则回滚事务；否则提交
	Trans(ctx context.Context, fn func(ctx context.Context, model Model) error) error
}
