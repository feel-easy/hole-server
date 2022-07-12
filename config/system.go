package config

type System struct {
	Env           string `json:"env" yaml:"env"` // 环境值
	Host          string `json:"host" yaml:"host"`
	Port          string `json:"port" yaml:"port"`                     // 端口值
	DbType        string `json:"db-type" yaml:"db-type"`               // 数据库类型:mysql(默认)|sqlite|sqlserver|postgresql
	OssType       string `json:"oss-type" yaml:"oss-type"`             // Oss类型
	UseMultipoint bool   `json:"use-multipoint" yaml:"use-multipoint"` // 多点登录拦截
	UseRedis      bool   `json:"use-redis" yaml:"use-redis"`           // 使用redis
	LimitCountIP  int    `json:"iplimit-count" yaml:"iplimit-count"`
	LimitTimeIP   int    `json:"iplimit-time" yaml:"iplimit-time"`
}
