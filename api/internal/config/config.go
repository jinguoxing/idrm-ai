package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

// Config 配置结构
type Config struct {
	rest.RestConf // Go-Zero REST配置（包含Name, Host, Port等）

	// 数据库配置
	DataSource string `json:",default=,env=DATASOURCE"` // SQLx数据源
	GormDSN    string `json:",default=,env=GORM_DSN"`   // GORM数据源

	// Redis配置
	RedisHost string `json:",default=127.0.0.1:6379,env=REDIS_HOST"` // Redis主机
	RedisPass string `json:",default=,env=REDIS_PASS"`               // Redis密码
	RedisType string `json:",default=node,env=REDIS_TYPE"`          // Redis类型

	// RPC配置（如果需要）
	RPCServer zrpc.RpcClientConf `json:",optional"` // RPC服务配置
}
