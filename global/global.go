package global

import (
	"sync"

	"github.com/feel-easy/hole-server/config"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	DBList map[string]*gorm.DB
	REDIS  *redis.Client
	CONFIG config.Server
	VIPER  *viper.Viper
	LOG    *zap.Logger

	Concurrency_Control = &singleflight.Group{}

	lock sync.RWMutex
)
