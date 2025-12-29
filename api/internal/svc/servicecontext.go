package svc

import (
	"database/sql"
	"fmt"
	"strings"

	"idrm/api/internal/config"
	"idrm/model/data-catalog/data_catalog"
	"idrm/model/data-catalog/data_catalog_category"
	"idrm/model/data-catalog/data_catalog_column"
	"idrm/model/data-catalog/data_resource"
	"idrm/model/system/menu"
	"idrm/pkg/orm"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ServiceContext 服务上下文
// 包含所有Handler需要的服务依赖
type ServiceContext struct {
	Config config.Config

	// 数据库连接
	SqlConn *sql.DB  // SQLx连接
	GormDB  *gorm.DB // GORM连接

	// ORM指标
	ORMMetrics *orm.ORMMetrics

	// 菜单Model（接口类型，遵循依赖倒置原则）
	MenuModel menu.Model

	// 数据资源目录Models
	DataCatalogModel       data_catalog.Model
	DataCatalogColumnModel data_catalog_column.Model
	DataCatalogCategoryModel data_catalog_category.Model
	DataResourceModel      data_resource.Model
}

// NewServiceContext 创建ServiceContext实例
// 初始化所有服务依赖
func NewServiceContext(c config.Config) *ServiceContext {
	// 验证和迁移配置
	if err := c.Validate(); err != nil {
		panic(fmt.Sprintf("invalid config: %v", err))
	}

	var sqlConn *sql.DB
	var gormDB *gorm.DB

	// 初始化ORM指标
	ormMetrics := orm.NewORMMetrics()

	// 获取数据源配置
	sqlxDataSource := c.Database.SQLx.DataSource
	gormDSN := c.Database.GORM.DSN

	// 判断数据库类型（根据DSN格式）
	if strings.HasPrefix(sqlxDataSource, "./") || strings.HasPrefix(sqlxDataSource, "/") || 
	   strings.HasPrefix(gormDSN, "./") || strings.HasPrefix(gormDSN, "/") {
		// SQLite数据库
		sqlConn, gormDB = initSQLite(sqlxDataSource, gormDSN, c.Database.Pool)
	} else {
		// MySQL数据库
		sqlConn, gormDB = initMySQL(sqlxDataSource, gormDSN, c.Database.Pool, c.Database.GORM)
	}

	// 使用策略模式初始化MenuModel
	var menuModel menu.Model
	if c.ORM.EnableMetrics {
		// 使用带指标的策略模型
		menuModel = menu.NewModelWithStrategy(sqlConn, gormDB, c.ORM.Strategy, ormMetrics)
	} else {
		// 使用传统工厂模式（向后兼容）
		menuModel = menu.NewModel(sqlConn, gormDB)
	}

	// 初始化数据资源目录Models
	var dataCatalogModel data_catalog.Model
	var dataCatalogColumnModel data_catalog_column.Model
	var dataCatalogCategoryModel data_catalog_category.Model
	var dataResourceModel data_resource.Model

	if c.ORM.EnableMetrics {
		// 使用带指标的策略模型
		dataCatalogModel = data_catalog.NewModelWithStrategy(sqlConn, gormDB, c.ORM.Strategy, ormMetrics)
		dataCatalogColumnModel = data_catalog_column.NewModelWithStrategy(sqlConn, gormDB, c.ORM.Strategy, ormMetrics)
		dataCatalogCategoryModel = data_catalog_category.NewModelWithStrategy(sqlConn, gormDB, c.ORM.Strategy, ormMetrics)
		dataResourceModel = data_resource.NewModelWithStrategy(sqlConn, gormDB, c.ORM.Strategy, ormMetrics)
	} else {
		// 使用传统工厂模式（向后兼容）
		dataCatalogModel = data_catalog.NewModel(sqlConn, gormDB)
		dataCatalogColumnModel = data_catalog_column.NewModel(sqlConn, gormDB)
		dataCatalogCategoryModel = data_catalog_category.NewModel(sqlConn, gormDB)
		dataResourceModel = data_resource.NewModel(sqlConn, gormDB)
	}

	logx.Infow("ServiceContext initialized",
		logx.Field("ormStrategy", c.ORM.Strategy),
		logx.Field("enableMetrics", c.ORM.EnableMetrics),
		logx.Field("maxOpenConns", c.Database.Pool.MaxOpenConns),
		logx.Field("maxIdleConns", c.Database.Pool.MaxIdleConns),
	)

	return &ServiceContext{
		Config:                  c,
		SqlConn:                 sqlConn,
		GormDB:                  gormDB,
		ORMMetrics:              ormMetrics,
		MenuModel:               menuModel,
		DataCatalogModel:        dataCatalogModel,
		DataCatalogColumnModel:  dataCatalogColumnModel,
		DataCatalogCategoryModel: dataCatalogCategoryModel,
		DataResourceModel:       dataResourceModel,
	}
}

// initSQLite 初始化SQLite数据库
func initSQLite(sqlxDataSource, gormDSN string, poolConfig config.DBPoolConfig) (*sql.DB, *gorm.DB) {
	var sqlConn *sql.DB
	var gormDB *gorm.DB
	var err error

	// 初始化SQLx连接
	if sqlxDataSource != "" {
		sqlConn, err = sql.Open("sqlite3", sqlxDataSource)
		if err != nil {
			logx.Errorf("failed to connect to SQLite (SQLx): %v", err)
		} else {
			if err := sqlConn.Ping(); err != nil {
				logx.Errorf("failed to ping SQLite (SQLx): %v", err)
				sqlConn = nil
			} else {
				configureConnectionPool(sqlConn, poolConfig)
				logx.Info("SQLite (SQLx) connected successfully")
			}
		}
	}

	// 初始化GORM连接
	if gormDSN != "" {
		gormDB, err = gorm.Open(sqlite.Open(gormDSN), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			logx.Errorf("failed to connect to SQLite (GORM): %v", err)
		} else {
			// 配置GORM连接池
			if sqlDB, err := gormDB.DB(); err == nil {
				configureConnectionPool(sqlDB, poolConfig)
			}

			// 自动创建表结构
			err = gormDB.AutoMigrate(&menu.Menu{})
			if err != nil {
				logx.Errorf("failed to migrate database: %v", err)
			} else {
				logx.Info("Database migration completed")
			}

			// 插入演示数据
			insertDemoData(gormDB)
			logx.Info("SQLite (GORM) connected successfully")
		}
	}

	// 至少需要一个可用的连接
	if sqlConn == nil && gormDB == nil {
		panic("no database connection available")
	}

	return sqlConn, gormDB
}

// initMySQL 初始化MySQL数据库
func initMySQL(sqlxDataSource, gormDSN string, poolConfig config.DBPoolConfig, gormConfig config.GORMConfig) (*sql.DB, *gorm.DB) {
	var sqlConn *sql.DB
	var gormDB *gorm.DB
	var err error

	// 初始化SQLx连接
	if sqlxDataSource != "" {
		sqlConn, err = sql.Open("mysql", sqlxDataSource)
		if err != nil {
			logx.Errorf("failed to connect to MySQL (SQLx): %v", err)
		} else {
			if err := sqlConn.Ping(); err != nil {
				logx.Errorf("failed to ping MySQL (SQLx): %v", err)
				sqlConn = nil
			} else {
				configureConnectionPool(sqlConn, poolConfig)
				logx.Info("MySQL (SQLx) connected successfully")
			}
		}
	}

	// 初始化GORM连接
	if gormDSN != "" {
		// 配置GORM日志级别
		var logLevel logger.LogLevel
		switch strings.ToLower(gormConfig.LogLevel) {
		case "silent":
			logLevel = logger.Silent
		case "error":
			logLevel = logger.Error
		case "warn":
			logLevel = logger.Warn
		case "info":
			logLevel = logger.Info
		default:
			logLevel = logger.Info
		}

		gormDB, err = gorm.Open(mysql.Open(gormDSN), &gorm.Config{
			Logger: logger.Default.LogMode(logLevel),
		})
		if err != nil {
			logx.Errorf("failed to connect to MySQL (GORM): %v", err)
		} else {
			// 配置GORM连接池
			if sqlDB, err := gormDB.DB(); err == nil {
				configureConnectionPool(sqlDB, poolConfig)
			}

			// 自动迁移（如果启用）
			if gormConfig.AutoMigrate {
				err = gormDB.AutoMigrate(&menu.Menu{})
				if err != nil {
					logx.Errorf("failed to migrate database: %v", err)
				} else {
					logx.Info("Database migration completed")
				}
			}

			logx.Info("MySQL (GORM) connected successfully")
		}
	}

	// 至少需要一个可用的连接
	if sqlConn == nil && gormDB == nil {
		panic("no database connection available")
	}

	return sqlConn, gormDB
}

// configureConnectionPool 配置数据库连接池
func configureConnectionPool(db *sql.DB, config config.DBPoolConfig) {
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	logx.Infow("Database connection pool configured",
		logx.Field("maxOpenConns", config.MaxOpenConns),
		logx.Field("maxIdleConns", config.MaxIdleConns),
		logx.Field("connMaxLifetime", config.ConnMaxLifetime),
		logx.Field("connMaxIdleTime", config.ConnMaxIdleTime),
	)
}

// insertDemoData 插入演示数据
func insertDemoData(db *gorm.DB) {
	// 检查是否已有数据
	var count int64
	db.Table("sys_menus").Count(&count)
	if count > 0 {
		return // 已有数据，不重复插入
	}

	// 插入演示菜单数据
	demoMenus := []menu.Menu{
		{
			ParentId:  0,
			Name:      "系统管理",
			RoutePath: "/system",
			Icon:      "setting",
			SortOrder: 100,
			PermTag:   "system",
			Status:    1,
		},
		{
			ParentId:      1,
			Name:          "菜单管理",
			RoutePath:     "/system/menu",
			ComponentPath: "system/menu/index",
			Icon:          "menu",
			SortOrder:     1,
			PermTag:       "system:menu",
			Status:        1,
		},
		{
			ParentId:      1,
			Name:          "用户管理",
			RoutePath:     "/system/user",
			ComponentPath: "system/user/index",
			Icon:          "user",
			SortOrder:     2,
			PermTag:       "system:user",
			Status:        1,
		},
		{
			ParentId:  0,
			Name:      "数据资源",
			RoutePath: "/resource",
			Icon:      "database",
			SortOrder: 200,
			PermTag:   "resource",
			Status:    1,
		},
		{
			ParentId:      4,
			Name:          "资源目录",
			RoutePath:     "/resource/catalog",
			ComponentPath: "resource/catalog/index",
			Icon:          "folder",
			SortOrder:     1,
			PermTag:       "resource:catalog",
			Status:        1,
		},
	}

	for _, m := range demoMenus {
		db.Create(&m)
	}

	logx.Info("Demo data inserted successfully")
}
