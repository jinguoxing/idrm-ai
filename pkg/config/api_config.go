package config

import (
	"time"

	"github.com/zeromicro/go-zero/rest"
)

// ORMStrategy ORM选择策略
type ORMStrategy string

const (
	ORMStrategyGORMFirst ORMStrategy = "gorm_first" // GORM优先（默认）
	ORMStrategySQLxFirst ORMStrategy = "sqlx_first" // SQLx优先
	ORMStrategyGORMOnly  ORMStrategy = "gorm_only"  // 仅GORM
	ORMStrategySQLxOnly  ORMStrategy = "sqlx_only"  // 仅SQLx
	ORMStrategyAuto      ORMStrategy = "auto"       // 智能选择
)

// Config API服务配置
type Config struct {
	rest.RestConf
	
	// 多数据库配置
	DataSources DataSourcesConfig
	
	// Redis配置
	Redis RedisConfig
	
	// Kafka配置
	Kafka KafkaConfig
	
	// JWT配置
	Auth AuthConfig
	
	// CORS配置
	Cors CorsConfig

	// ORM配置（新增）
	ORM ORMConfig
}

// DataSourcesConfig 多数据库配置
type DataSourcesConfig struct {
	DataView          DatabaseConfig
	DataUnderstanding DatabaseConfig
	ResourceCatalog   DatabaseConfig
}

// DatabaseConfig 单个数据库配置
type DatabaseConfig struct {
	Driver          string
	Source          string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host string
	Type string
	Pass string
	Db   int
}

// KafkaConfig Kafka配置
type KafkaConfig struct {
	Brokers []string
	Group   string
	Topics  []string
}

// AuthConfig JWT配置
type AuthConfig struct {
	AccessSecret string
	AccessExpire int64
}

// CorsConfig CORS配置
type CorsConfig struct {
	AllowOrigins []string
	AllowMethods []string
	AllowHeaders []string
}

// ORMConfig ORM配置
type ORMConfig struct {
	Strategy      ORMStrategy   `json:"strategy,default=gorm_first"`  // ORM选择策略
	EnableMetrics bool          `json:"enableMetrics,default=true"`   // 启用性能指标
	EnableCache   bool          `json:"enableCache,default=false"`    // 启用缓存（暂未实现）
	CacheTTL      time.Duration `json:"cacheTTL,default=5m"`          // 缓存TTL
}

// DBPoolConfig 数据库连接池配置
type DBPoolConfig struct {
	MaxOpenConns    int           `json:"maxOpenConns,default=100"`    // 最大打开连接数
	MaxIdleConns    int           `json:"maxIdleConns,default=10"`     // 最大空闲连接数
	ConnMaxLifetime time.Duration `json:"connMaxLifetime,default=1h"`  // 连接最大生存时间
	ConnMaxIdleTime time.Duration `json:"connMaxIdleTime,default=10m"` // 连接最大空闲时间
}
