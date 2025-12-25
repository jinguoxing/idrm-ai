package config

// JobConfig 定时任务服务配置
type JobConfig struct {
	Name string
	Mode string

	// 多数据库配置
	DataSources DataSourcesConfig

	// Redis配置
	Redis RedisConfig

	// Kafka配置
	Kafka KafkaProducerConfig

	// 定时任务配置
	Jobs JobsConfig

	// 日志配置
	Log LogConfig
}

// KafkaProducerConfig Kafka生产者配置
type KafkaProducerConfig struct {
	Brokers []string
}

// JobsConfig 定时任务配置
type JobsConfig struct {
	SyncData   JobItemConfig
	Statistics JobItemConfig
	Cleanup    CleanupJobConfig
}

// JobItemConfig 单个任务配置
type JobItemConfig struct {
	Cron    string
	Enabled bool
}

// CleanupJobConfig 清理任务配置
type CleanupJobConfig struct {
	Cron          string
	Enabled       bool
	RetentionDays int
}

// LogConfig 日志配置
type LogConfig struct {
	ServiceName string
	Mode        string
	Encoding    string
	Level       string
	Path        string
}
