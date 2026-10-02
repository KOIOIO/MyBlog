package bootstrap

import (
	"server/internal/application/auth"
	"server/internal/application/user"
	"server/internal/common/email"
	"server/internal/common/jwt"
	"server/internal/infrastructure/geo"
	"server/internal/infrastructure/mysql"
	"server/internal/infrastructure/redis"
	ihttp "server/internal/interface/http"
	basehandler "server/internal/interface/http/handler/base"
	userhandler "server/internal/interface/http/handler/user"

	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

// BuildApp 手工依赖组装：基础设施 → 应用服务 → Handler → 路由依赖。
// 全部依赖经构造注入，不读取 global（global 仅作为旧代码兼容层）。
func BuildApp(infra *Infra) *ihttp.Deps {
	cfg := infra.Config
	log := infra.Log

	// ---- auth BC ----
	jwtUtil, err := jwt.New(
		cfg.Jwt.AccessTokenSecret,
		cfg.Jwt.RefreshTokenSecret,
		cfg.Jwt.AccessTokenExpiryTime,
		cfg.Jwt.RefreshTokenExpiryTime,
		cfg.Jwt.Issuer,
	)
	if err != nil {
		log.Fatal("Failed to init JWT:", zap.Error(err))
	}
	refreshTTL, err := jwt.ParseDuration(cfg.Jwt.RefreshTokenExpiryTime)
	if err != nil {
		log.Fatal("Failed to parse refresh token expiry:", zap.Error(err))
	}

	sessionStore := redis.NewSessionStore(&infra.Redis, refreshTTL)
	blacklistStore := mysql.NewBlacklistStore(infra.DB, infra.BlackCache)

	// ---- user BC ----
	userRepo := mysql.NewUserRepository(infra.DB)
	loginRepo := mysql.NewLoginRecordRepository(infra.DB)
	cache := redis.NewCache(&infra.Redis)
	geoClient := geo.NewClient(cfg.Gaode.Key, log)

	authApp := auth.NewAuthService(jwtUtil, sessionStore, blacklistStore, userRepo, refreshTTL)
	userApp := user.NewUserService(userRepo, loginRepo, cache, geoClient, authApp, cfg)

	// ---- interface ----
	captchaStore := base64Captcha.DefaultMemStore
	emailSender := email.New(cfg.Email)
	baseHandler := basehandler.NewHandler(cfg, captchaStore, emailSender, log)
	userHandler := userhandler.NewHandler(userApp, authApp, captchaStore, cfg, log)

	return &ihttp.Deps{
		Config:      cfg,
		Log:         log,
		Auth:        authApp,
		User:        userApp,
		BaseHandler: baseHandler,
		UserHandler: userHandler,
		Geo:         geoClient,
		Logins:      loginRepo,
	}
}
