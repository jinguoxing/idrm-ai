package menu

import "context"

// Model 菜单数据访问接口
// 支持GORM和SQLx双ORM实现，通过工厂模式自动选择
type Model interface {
	// Insert 插入新的菜单记录
	// 参数data中的Id会被忽略（自增），其他字段保持不变
	// 返回插入后的完整记录（包含Id和时间戳）
	Insert(ctx context.Context, data *Menu) (*Menu, error)

	// FindOne 根据ID查询单个菜单
	// 若不存在则返回 ErrNotFound
	FindOne(ctx context.Context, id int64) (*Menu, error)

	// FindByPermTag 根据权限标识查询菜单
	// 若不存在则返回 nil
	FindByPermTag(ctx context.Context, permTag string) (*Menu, error)

	// Update 更新菜单记录
	// 参数data中必须包含Id字段
	// 只更新非零值字段（使用GORM的Select或Map实现）
	Update(ctx context.Context, data *Menu) error

	// Delete 删除菜单记录（软删除）
	// 实际操作是设置deleted_at字段
	Delete(ctx context.Context, id int64) error

	// FindAll 查询所有菜单（支持状态过滤）
	// 如果status为-1，则查询所有状态的菜单
	FindAll(ctx context.Context, status int) ([]*Menu, error)

	// FindByParentId 根据父级ID查询子菜单列表
	// 如果parentId为0，则查询根菜单
	FindByParentId(ctx context.Context, parentId int64) ([]*Menu, error)

	// CountByParentId 统计指定父级下的子菜单数量
	CountByParentId(ctx context.Context, parentId int64) (int64, error)

	// HasChildren 检查菜单是否有子菜单
	HasChildren(ctx context.Context, parentId int64) (bool, error)

	// 事务支持方法

	// WithTx 创建事务副本
	// 参数tx是事务对象（GORM的*gorm.DB或SQLx的*sqlx.Tx）
	// 返回使用该事务的Model副本
	WithTx(tx interface{}) Model

	// Trans 在事务中执行多个操作
	// 若fn返回error，则回滚事务；否则提交
	Trans(ctx context.Context, fn func(ctx context.Context, model Model) error) error
}
