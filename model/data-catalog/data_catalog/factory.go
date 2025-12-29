package data_catalog

import (
	"context"
	"database/sql"

	"idrm/pkg/config"
	"idrm/pkg/orm"

	"gorm.io/gorm"
)

// ========== 工厂函数类型 ==========

type gormFactoryFunc func(db *gorm.DB) Model
type sqlxFactoryFunc func(db *sql.DB) Model

var (
	gormFactory gormFactoryFunc
	sqlxFactory sqlxFactoryFunc
)

// ========== 工厂注册函数 ==========

// RegisterGormFactory 注册GORM工厂函数
// 在GORM实现的init()中调用
func RegisterGormFactory(fn gormFactoryFunc) {
	gormFactory = fn
}

// RegisterSqlxFactory 注册SQLx工厂函数
// 在SQLx实现的init()中调用
func RegisterSqlxFactory(fn sqlxFactoryFunc) {
	sqlxFactory = fn
}

// ========== 工厂创建函数 ==========

// NewModel 创建Model实例（工厂函数）
// 优先使用GORM，如果GORM不可用则降级到SQLx
// 若两者都不可用，则panic
func NewModel(sqlConn *sql.DB, gormDB *gorm.DB) Model {
	// 优先使用GORM
	if gormDB != nil && gormFactory != nil {
		return gormFactory(gormDB)
	}

	// 降级到SQLx
	if sqlConn != nil && sqlxFactory != nil {
		return sqlxFactory(sqlConn)
	}

	panic("data_catalog: no database connection available, both GORM and SQLx are nil")
}

// NewModelWithStrategy 使用策略创建Model实例
// 根据配置的ORM策略选择最优实现
func NewModelWithStrategy(sqlConn *sql.DB, gormDB *gorm.DB, strategy config.ORMStrategy, metrics *orm.ORMMetrics) Model {
	selector := orm.NewORMSelector(strategy, sqlConn, gormDB, metrics)

	return &StrategyModel{
		gormModel: func() Model {
		if gormDB != nil && gormFactory != nil {
			return gormFactory(gormDB)
		}
		return nil
	}(),
		sqlxModel: func() Model {
		if sqlConn != nil && sqlxFactory != nil {
			return sqlxFactory(sqlConn)
		}
		return nil
	}(),
		selector: selector,
	}
}

// MustNewModel 创建Model实例，若失败则panic
// 用于初始化阶段，确保数据库连接可用
func MustNewModel(sqlConn *sql.DB, gormDB *gorm.DB) Model {
	model, err := func() (Model, error) {
		defer func() {
			if r := recover(); r != nil {
				// panic已被转换为error返回
			}
		}()
		return NewModel(sqlConn, gormDB), nil
	}()

	if err != nil {
		panic(err)
	}

	return model
}

// StrategyModel 策略Model实现
type StrategyModel struct {
	gormModel Model
	sqlxModel Model
	selector  *orm.ORMSelector
}

// 实现 Model 接口的方法，根据策略选择使用 GORM 或 SQLx
func (m *StrategyModel) Insert(ctx context.Context, data *DataCatalog) (*DataCatalog, error) {
	model := m.selector.Select(m.gormModel, m.sqlxModel)
	if model == nil {
		panic("no available ORM implementation")
	}
	return model.(Model).Insert(ctx, data)
}

func (m *StrategyModel) FindOne(ctx context.Context, id string) (*DataCatalog, error) {
	model := m.selector.Select(m.gormModel, m.sqlxModel)
	if model == nil {
		panic("no available ORM implementation")
	}
	return model.(Model).FindOne(ctx, id)
}

func (m *StrategyModel) Update(ctx context.Context, data *DataCatalog) error {
	model := m.selector.Select(m.gormModel, m.sqlxModel)
	if model == nil {
		panic("no available ORM implementation")
	}
	return model.(Model).Update(ctx, data)
}

func (m *StrategyModel) Delete(ctx context.Context, id string) error {
	model := m.selector.Select(m.gormModel, m.sqlxModel)
	if model == nil {
		panic("no available ORM implementation")
	}
	return model.(Model).Delete(ctx, id)
}

func (m *StrategyModel) FindBySourceDept(ctx context.Context, deptId string, status string) ([]*DataCatalog, error) {
	model := m.selector.Select(m.gormModel, m.sqlxModel)
	if model == nil {
		panic("no available ORM implementation")
	}
	return model.(Model).FindBySourceDept(ctx, deptId, status)
}

func (m *StrategyModel) FindByDraftId(ctx context.Context, draftId string) (*DataCatalog, error) {
	model := m.selector.Select(m.gormModel, m.sqlxModel)
	if model == nil {
		panic("no available ORM implementation")
	}
	return model.(Model).FindByDraftId(ctx, draftId)
}

func (m *StrategyModel) FindDrafts(ctx context.Context, deptId string, status string, page, pageSize int) ([]*DataCatalog, int64, error) {
	model := m.selector.Select(m.gormModel, m.sqlxModel)
	if model == nil {
		panic("no available ORM implementation")
	}
	return model.(Model).FindDrafts(ctx, deptId, status, page, pageSize)
}

func (m *StrategyModel) CheckNameExists(ctx context.Context, deptId string, title string, excludeId string) (bool, error) {
	model := m.selector.Select(m.gormModel, m.sqlxModel)
	if model == nil {
		panic("no available ORM implementation")
	}
	return model.(Model).CheckNameExists(ctx, deptId, title, excludeId)
}

func (m *StrategyModel) WithTx(tx interface{}) Model {
	model := m.selector.Select(m.gormModel, m.sqlxModel)
	if model == nil {
		panic("no available ORM implementation")
	}
	return model.(Model).WithTx(tx)
}

func (m *StrategyModel) Trans(ctx context.Context, fn func(ctx context.Context, model Model) error) error {
	model := m.selector.Select(m.gormModel, m.sqlxModel)
	if model == nil {
		panic("no available ORM implementation")
	}
	return model.(Model).Trans(ctx, fn)
}
