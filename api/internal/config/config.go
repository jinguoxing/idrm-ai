package config

import (
	"time"

	"idrm/pkg/config"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

// Config 配置结构
type Config struct {
	rest.RestConf // Go-Zero REST配置（包含Name, Host, Port等）

	// 新的数据库配置结构
	Database DatabaseConfig `json:"database"`

	// ORM配置
	ORM config.ORMConfig `json:"orm"`

	// 缓存配置
	Cache CacheConfig `json:"cache"`

	// RPC配置（如果需要）
	RPCServer zrpc.RpcClientConf `json:",optional"` // RPC服务配置

	// 向后兼容的配置字段（已废弃，但保留以避免破坏现有配置）
	DataSource string `json:",default=,env=DATASOURCE,optional"` // 已废弃：使用Database.SQLx.DataSource
	GormDSN    string `json:",default=,env=GORM_DSN,optional"`   // 已废弃：使用Database.GORM.DSN
	RedisHost  string `json:",default=127.0.0.1:6379,env=REDIS_HOST,optional"` // 已废弃：使用Cache.RedisHost
	RedisPass  string `json:",default=,env=REDIS_PASS,optional"`               // 已废弃：使用Cache.RedisPass
	RedisType  string `json:",default=node,env=REDIS_TYPE,optional"`           // 已废弃：使用Cache.RedisType
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	SQLx SQLxConfig   `json:"sqlx"`
	GORM GORMConfig   `json:"gorm"`
	Pool DBPoolConfig `json:"pool"`
}

// SQLxConfig SQLx配置
type SQLxConfig struct {
	DataSource string `json:"dataSource,default=,env=DATASOURCE"` // SQLx数据源
	Driver     string `json:"driver,default=mysql"`               // 数据库驱动
}

// GORMConfig GORM配置
type GORMConfig struct {
	DSN         string `json:"dsn,default=,env=GORM_DSN"`     // GORM数据源
	LogLevel    string `json:"logLevel,default=info"`         // 日志级别
	AutoMigrate bool   `json:"autoMigrate,default=true"`      // 自动迁移
}

// DBPoolConfig 数据库连接池配置
type DBPoolConfig struct {
	MaxOpenConns    int           `json:"maxOpenConns,default=100"`    // 最大打开连接数
	MaxIdleConns    int           `json:"maxIdleConns,default=10"`     // 最大空闲连接数
	ConnMaxLifetime time.Duration `json:"connMaxLifetime,default=1h"`  // 连接最大生存时间
	ConnMaxIdleTime time.Duration `json:"connMaxIdleTime,default=10m"` // 连接最大空闲时间
}

// CacheConfig 缓存配置
type CacheConfig struct {
	RedisHost string `json:"redisHost,default=127.0.0.1:6379,env=REDIS_HOST"` // Redis主机
	RedisPass string `json:"redisPass,default=,env=REDIS_PASS"`               // Redis密码
	RedisType string `json:"redisType,default=node,env=REDIS_TYPE"`           // Redis类型
	RedisDB   int    `json:"redisDB,default=0"`                               // Redis数据库
}

// MigrateOldConfig 迁移旧配置到新结构
// 确保向后兼容性
func (c *Config) MigrateOldConfig() {
	// 如果使用了旧的配置字段，迁移到新结构
	if c.DataSource != "" && c.Database.SQLx.DataSource == "" {
		c.Database.SQLx.DataSource = c.DataSource
	}
	if c.GormDSN != "" && c.Database.GORM.DSN == "" {
		c.Database.GORM.DSN = c.GormDSN
	}
	if c.RedisHost != "" && c.Cache.RedisHost == "127.0.0.1:6379" {
		c.Cache.RedisHost = c.RedisHost
	}
	if c.RedisPass != "" && c.Cache.RedisPass == "" {
		c.Cache.RedisPass = c.RedisPass
	}
	if c.RedisType != "" && c.Cache.RedisType == "node" {
		c.Cache.RedisType = c.RedisType
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	// 迁移旧配置
	c.MigrateOldConfig()

	// 验证ORM策略
	switch c.ORM.Strategy {
	case config.ORMStrategyGORMFirst, config.ORMStrategySQLxFirst, config.ORMStrategyGORMOnly, config.ORMStrategySQLxOnly, config.ORMStrategyAuto:
		// 有效策略
	default:
		c.ORM.Strategy = config.ORMStrategyGORMFirst // 使用默认策略
	}

	// 验证连接池配置
	if c.Database.Pool.MaxOpenConns <= 0 {
		c.Database.Pool.MaxOpenConns = 100
	}
	if c.Database.Pool.MaxIdleConns <= 0 {
		c.Database.Pool.MaxIdleConns = 10
	}
	if c.Database.Pool.ConnMaxLifetime <= 0 {
		c.Database.Pool.ConnMaxLifetime = time.Hour
	}
	if c.Database.Pool.ConnMaxIdleTime <= 0 {
		c.Database.Pool.ConnMaxIdleTime = 10 * time.Minute
	}

	return nil
}
