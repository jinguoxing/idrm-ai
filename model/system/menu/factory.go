package menu

import (
	"database/sql"
	"sync"

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
