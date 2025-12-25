package config

import "github.com/zeromicro/go-zero/rest"

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
