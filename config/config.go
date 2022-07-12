package config

type Server struct {
	JWT    JWT    `json:"jwt" yaml:"jwt"`
	Redis  Redis  `json:"redis" yaml:"redis"`
	Email  Email  `json:"email" yaml:"email"`
	System System `json:"system",yaml:"system"`
	Zap    Zap    `json:"zap",yaml:"zap"`
	Mysql  Mysql  `json:"mysql" yaml:"mysql"`
}
