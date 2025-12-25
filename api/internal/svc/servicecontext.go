package svc

import (
	"database/sql"
	"fmt"
	"strings"

	"idrm/api/internal/config"
	"idrm/model/system/menu"

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
	SqlConn *sql.DB // SQLx连接
	GormDB  *gorm.DB // GORM连接

	// 菜单Model（接口类型，遵循依赖倒置原则）
	MenuModel menu.Model
}

// NewServiceContext 创建ServiceContext实例
// 初始化所有服务依赖
func NewServiceContext(c config.Config) *ServiceContext {
	var sqlConn *sql.DB
	var gormDB *gorm.DB
	var err error

	// 判断数据库类型（根据DSN格式）
	if strings.HasPrefix(c.DataSource, "./") || strings.HasPrefix(c.DataSource, "/") {
		// SQLite数据库
		sqlConn, err = sql.Open("sqlite3", c.DataSource)
		if err != nil {
			panic(fmt.Sprintf("failed to connect to SQLite: %v", err))
		}
		if err := sqlConn.Ping(); err != nil {
			panic(fmt.Sprintf("failed to ping SQLite: %v", err))
		}

		gormDB, err = gorm.Open(sqlite.Open(c.GormDSN), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			panic(fmt.Sprintf("failed to connect to SQLite (GORM): %v", err))
		}

		// 自动创建表结构
		err = gormDB.AutoMigrate(&menu.Menu{})
		if err != nil {
			panic(fmt.Sprintf("failed to migrate database: %v", err))
		}

		// 插入演示数据
		insertDemoData(gormDB)
	} else {
		// MySQL数据库
		sqlConn, err = sql.Open("mysql", c.DataSource)
		if err != nil {
			panic(fmt.Sprintf("failed to connect to database (SQLx): %v", err))
		}
		if err := sqlConn.Ping(); err != nil {
			panic(fmt.Sprintf("failed to ping database (SQLx): %v", err))
		}

		gormDB, err = gorm.Open(mysql.Open(c.GormDSN), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			panic(fmt.Sprintf("failed to connect to database (GORM): %v", err))
		}
	}

	// 使用工厂模式初始化MenuModel
	// 优先使用GORM，如果不可用则降级到SQLx
	menuModel := menu.NewModel(sqlConn, gormDB)

	return &ServiceContext{
		Config:    c,
		SqlConn:   sqlConn,
		GormDB:    gormDB,
		MenuModel: menuModel,
	}
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

	fmt.Println("✅ 演示数据已插入到数据库")
}
