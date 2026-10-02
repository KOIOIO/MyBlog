// Package bootstrap 负责基础设施初始化与手工依赖组装（替代 initialize + global 的写入点）。
package bootstrap

import (
	"log"
	"os"

	"server/config"
	"server/internal/common/jwt"
	"server/internal/common/logger"
	iconfig "server/internal/config"

	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-redis/redis"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Infra 持有全部基础设施句柄，供后续各层构造注入使用。
type Infra struct {
	Config     *config.Config
	Log        *zap.Logger
	DB         *gorm.DB
	Redis      redis.Client
	ESClient   *elasticsearch.TypedClient
	BlackCache local_cache.Cache
}

// InitConfig 加载配置（构造注入起点，不再经过 global.Config）。
func InitConfig() *config.Config {
	cfg, err := iconfig.Load("")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	return cfg
}

// InitLogger 初始化日志。
func InitLogger(cfg *config.Config) *zap.Logger {
	return logger.Init(&cfg.Zap)
}

// InitInfra 实例化 DB/Redis/ES/本地缓存并返回句柄。
func InitInfra(cfg *config.Config, log *zap.Logger) *Infra {
	db := initGorm(cfg, log)
	rdb := connectRedis(cfg, log)

	es := connectEs(cfg, log)

	bc := initBlackCache(cfg, log)

	return &Infra{
		Config:     cfg,
		Log:        log,
		DB:         db,
		Redis:      rdb,
		ESClient:   es,
		BlackCache: bc,
	}
}

// initGorm 初始化并返回使用 MySQL 配置的 GORM 连接。
func initGorm(cfg *config.Config, log *zap.Logger) *gorm.DB {
	mysqlCfg := cfg.Mysql

	db, err := gorm.Open(mysqlDriver(mysqlCfg), &gorm.Config{
		Logger: mysqlLogger(mysqlCfg),
	})
	if err != nil {
		log.Error("Failed to connect to MySQL:", zap.Error(err))
		os.Exit(1)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(mysqlCfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(mysqlCfg.MaxOpenConns)
	return db
}

// connectRedis 初始化并返回 Redis 客户端。
func connectRedis(cfg *config.Config, log *zap.Logger) redis.Client {
	redisCfg := cfg.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Address,
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	})

	if _, err := client.Ping().Result(); err != nil {
		log.Error("Failed to connect to Redis:", zap.Error(err))
		os.Exit(1)
	}
	return *client
}

// connectEs 初始化并返回 Elasticsearch TypedClient。
func connectEs(cfg *config.Config, log *zap.Logger) *elasticsearch.TypedClient {
	esCfg := cfg.ES
	ecfg := elasticsearch.Config{
		Addresses: []string{esCfg.URL},
		Username:  esCfg.Username,
		Password:  esCfg.Password,
	}
	if esCfg.IsConsolePrint {
		ecfg.Logger = &elastictransport.ColorLogger{
			Output:             os.Stdout,
			EnableRequestBody:  true,
			EnableResponseBody: true,
		}
	}

	client, err := elasticsearch.NewTypedClient(ecfg)
	if err != nil {
		log.Error("Failed to connect to Elasticsearch", zap.Error(err))
		os.Exit(1)
	}
	return client
}

// initBlackCache 配置本地缓存（刷新令牌过期时间），用于 JWT 黑名单。
func initBlackCache(cfg *config.Config, log *zap.Logger) local_cache.Cache {
	refreshTokenExpiry, err := jwt.ParseDuration(cfg.Jwt.RefreshTokenExpiryTime)
	if err != nil {
		log.Error("Failed to parse refresh token expiry time configuration:", zap.Error(err))
		os.Exit(1)
	}
	if _, err := jwt.ParseDuration(cfg.Jwt.AccessTokenExpiryTime); err != nil {
		log.Error("Failed to parse access token expiry time configuration:", zap.Error(err))
		os.Exit(1)
	}
	return local_cache.NewCache(
		local_cache.SetDefaultExpire(refreshTokenExpiry),
	)
}
