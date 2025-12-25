package db

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// Config 数据库配置
type Config struct {
	// 基础连接配置
	Host     string `json:",default=127.0.0.1"`
	Port     int    `json:",default=3306"`
	Database string
	Username string
	Password string
	Charset  string `json:",default=utf8mb4"`

	// 连接池配置
	MaxIdleConns    int `json:",default=10"`   // 最大空闲连接数
	MaxOpenConns    int `json:",default=100"`  // 最大打开连接数
	ConnMaxLifetime int `json:",default=3600"` // 连接最大生存时间(秒)
	ConnMaxIdleTime int `json:",default=600"`  // 连接最大空闲时间(秒)

	// 日志配置
	LogLevel          string `json:",default=warn"` // silent/error/warn/info
	SlowThreshold     int    `json:",default=200"`  // 慢查询阈值(毫秒)
	SkipDefaultTxn    bool   `json:",default=true"` // 跳过默认事务
	PrepareStmt       bool   `json:",default=true"` // 预编译语句
	SingularTable     bool   `json:",default=true"` // 使用单数表名
	DisableForeignKey bool   `json:",default=true"` // 禁用外键约束
}

// InitGorm 初始化 GORM 连接
func InitGorm(c Config) (*gorm.DB, error) {
	// 1. 构建 DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		c.Charset,
	)

	// 2. 配置 GORM
	gormConfig := &gorm.Config{
		// 日志配置
		Logger: logger.Default.LogMode(getLogLevel(c.LogLevel)),

		// 命名策略
		NamingStrategy: schema.NamingStrategy{
			SingularTable: c.SingularTable, // 使用单数表名
		},

		// 性能优化
		SkipDefaultTransaction: c.SkipDefaultTxn, // 跳过默认事务
		PrepareStmt:            c.PrepareStmt,    // 预编译语句

		// 外键约束
		DisableForeignKeyConstraintWhenMigrating: c.DisableForeignKey,
	}

	// 3. 打开连接
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// 4. 获取底层 sql.DB 并配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// 连接池配置
	sqlDB.SetMaxIdleConns(c.MaxIdleConns)
	sqlDB.SetMaxOpenConns(c.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(c.ConnMaxLifetime) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(c.ConnMaxIdleTime) * time.Second)

	return db, nil
}

// getLogLevel 获取日志级别
func getLogLevel(level string) logger.LogLevel {
	switch level {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Warn
	}
}
