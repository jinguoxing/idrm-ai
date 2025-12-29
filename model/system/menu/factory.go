package menu

import (
	"database/sql"
	"sync"

	"idrm/pkg/config"
	"idrm/pkg/orm"

	"gorm.io/gorm"
)

// 工厂函数类型
type gormFactoryFunc func(db *gorm.DB) Model
type sqlxFactoryFunc func(db *sql.DB) Model

var (
	gormFactory gormFactoryFunc
	sqlxFactory sqlxFactoryFunc

	once sync.Once
)

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

	panic("menu: no database connection available, both GORM and SQLx are nil")
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

// MustNewModelWithStrategy 使用策略创建Model实例，若失败则panic
func MustNewModelWithStrategy(sqlConn *sql.DB, gormDB *gorm.DB, strategy config.ORMStrategy, metrics *orm.ORMMetrics) Model {
	model, err := func() (Model, error) {
		defer func() {
			if r := recover(); r != nil {
				// panic已被转换为error返回
			}
		}()
		return NewModelWithStrategy(sqlConn, gormDB, strategy, metrics), nil
	}()

	if err != nil {
		panic(err)
	}

	return model
}
