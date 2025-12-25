package config

// ConsumerConfig 消费者服务配置
type ConsumerConfig struct {
	Name string
	Mode string

	// 多数据库配置
	DataSources DataSourcesConfig

	// Redis配置
	Redis RedisConfig

	// Kafka消费者配置
	Kafka KafkaConsumerConfig

	// 日志配置
	Log LogConfig

	// 重试配置
	Retry RetryConfig
}

// KafkaConsumerConfig Kafka消费者配置
type KafkaConsumerConfig struct {
	Brokers   []string
	Consumers map[string]ConsumerItemConfig
}

// ConsumerItemConfig 单个消费者配置
type ConsumerItemConfig struct {
	Group          string
	Topics         []string
	Workers        int
	AutoCommit     bool
	CommitInterval int
}

// RetryConfig 重试配置
type RetryConfig struct {
	MaxRetries      int
	InitialInterval int
	MaxInterval     int
	Multiplier      int
}
