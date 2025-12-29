package data_catalog_category

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

	panic("data_catalog_category: no database connection available")
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

func (m *StrategyModel) Insert(ctx context.Context, data []*DataCatalogCategory) error {
	model := m.selector.Select(m.gormModel, m.sqlxModel)
	if model == nil {
		panic("no available ORM implementation")
	}
	return model.(Model).Insert(ctx, data)
}

func (m *StrategyModel) FindByCatalogId(ctx context.Context, catalogId string) ([]*DataCatalogCategory, error) {
	model := m.selector.Select(m.gormModel, m.sqlxModel)
	if model == nil {
		panic("no available ORM implementation")
	}
	return model.(Model).FindByCatalogId(ctx, catalogId)
}

func (m *StrategyModel) DeleteByCatalogId(ctx context.Context, catalogId string) error {
	model := m.selector.Select(m.gormModel, m.sqlxModel)
	if model == nil {
		panic("no available ORM implementation")
	}
	return model.(Model).DeleteByCatalogId(ctx, catalogId)
}

func (m *StrategyModel) BatchDelete(ctx context.Context, ids []string) error {
	model := m.selector.Select(m.gormModel, m.sqlxModel)
	if model == nil {
		panic("no available ORM implementation")
	}
	return model.(Model).BatchDelete(ctx, ids)
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
