package data_resource

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

func RegisterGormFactory(fn gormFactoryFunc) {
	gormFactory = fn
}

func RegisterSqlxFactory(fn sqlxFactoryFunc) {
	sqlxFactory = fn
}

// ========== 工厂创建函数 ==========

func NewModel(sqlConn *sql.DB, gormDB *gorm.DB) Model {
	if gormDB != nil && gormFactory != nil {
		return gormFactory(gormDB)
	}

	if sqlConn != nil && sqlxFactory != nil {
		return sqlxFactory(sqlConn)
	}

	panic("data_resource: no database connection available")
}

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

// StrategyModel 策略Model实现
type StrategyModel struct {
	gormModel Model
	sqlxModel Model
	selector  *orm.ORMSelector
}

func (m *StrategyModel) Insert(ctx context.Context, data *DataResource) error {
	model := m.selector.Select(m.gormModel, m.sqlxModel)
	if model == nil {
		panic("no available ORM implementation")
	}
	return model.(Model).Insert(ctx, data)
}

func (m *StrategyModel) FindOne(ctx context.Context, id string) (*DataResource, error) {
	model := m.selector.Select(m.gormModel, m.sqlxModel)
	if model == nil {
		panic("no available ORM implementation")
	}
	return model.(Model).FindOne(ctx, id)
}

func (m *StrategyModel) Update(ctx context.Context, data *DataResource) error {
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

func (m *StrategyModel) FindByCatalogId(ctx context.Context, catalogId string) ([]*DataResource, error) {
	model := m.selector.Select(m.gormModel, m.sqlxModel)
	if model == nil {
		panic("no available ORM implementation")
	}
	return model.(Model).FindByCatalogId(ctx, catalogId)
}

func (m *StrategyModel) FindByCatalogIdAndType(ctx context.Context, catalogId string, resourceType ResourceType) ([]*DataResource, error) {
	model := m.selector.Select(m.gormModel, m.sqlxModel)
	if model == nil {
		panic("no available ORM implementation")
	}
	return model.(Model).FindByCatalogIdAndType(ctx, catalogId, resourceType)
}

func (m *StrategyModel) DeleteByCatalogId(ctx context.Context, catalogId string) error {
	model := m.selector.Select(m.gormModel, m.sqlxModel)
	if model == nil {
		panic("no available ORM implementation")
	}
	return model.(Model).DeleteByCatalogId(ctx, catalogId)
}

func (m *StrategyModel) CountByCatalogId(ctx context.Context, catalogId string) (int64, error) {
	model := m.selector.Select(m.gormModel, m.sqlxModel)
	if model == nil {
		panic("no available ORM implementation")
	}
	return model.(Model).CountByCatalogId(ctx, catalogId)
}

func (m *StrategyModel) CountLogicalViews(ctx context.Context, catalogId string) (int64, error) {
	model := m.selector.Select(m.gormModel, m.sqlxModel)
	if model == nil {
		panic("no available ORM implementation")
	}
	return model.(Model).CountLogicalViews(ctx, catalogId)
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
