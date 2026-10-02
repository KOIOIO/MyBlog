package bootstrap

import (
	"context"

	"server/internal/application/article"
	"server/internal/application/auth"
	"server/internal/application/comment"
	"server/internal/application/forum"
	"server/internal/application/user"
	"server/internal/common/email"
	"server/internal/common/jwt"
	"server/internal/infrastructure/es"
	"server/internal/infrastructure/geo"
	"server/internal/infrastructure/mysql"
	"server/internal/infrastructure/redis"
	ihttp "server/internal/interface/http"
	articlehandler "server/internal/interface/http/handler/article"
	basehandler "server/internal/interface/http/handler/base"
	commenthandler "server/internal/interface/http/handler/comment"
	forumhandler "server/internal/interface/http/handler/forum"
	userhandler "server/internal/interface/http/handler/user"

	"github.com/mojocn/base64Captcha"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// cronZapLogger 将 cron 日志转发到 zap。
type cronZapLogger struct {
	log *zap.Logger
}

func (l *cronZapLogger) Info(msg string, keysAndValues ...interface{}) {
	l.log.Info(msg, zap.Any("keysAndValues", keysAndValues))
}

func (l *cronZapLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	l.log.Error(msg, zap.Error(err), zap.Any("keysAndValues", keysAndValues))
}

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

	// ---- article BC ----
	articleRepo := mysql.NewArticleRepository(infra.DB)
	esArticleStore := es.NewArticleStore(infra.ESClient)
	viewCounter := redis.NewViewCounter(&infra.Redis)
	articleApp := article.NewService(articleRepo, esArticleStore, viewCounter)

	// ---- comment BC ----
	commentRepo := mysql.NewCommentRepo(infra.DB)
	commentApp := comment.NewService(commentRepo, log)

	// ---- forum BC ----
	forumRepo := mysql.NewForumRepo(infra.DB)
	forumApp := forum.NewService(forumRepo, cfg, log)

	// ---- interface ----
	captchaStore := base64Captcha.DefaultMemStore
	emailSender := email.New(cfg.Email)
	baseHandler := basehandler.NewHandler(cfg, captchaStore, emailSender, log)
	userHandler := userhandler.NewHandler(userApp, authApp, captchaStore, cfg, log)
	articleHandler := articlehandler.NewHandler(articleApp, log)
	commentHandler := commenthandler.NewHandler(commentApp, log)
	forumHandler := forumhandler.NewHandler(forumApp, cfg, log)

	return &ihttp.Deps{
		Config:         cfg,
		Log:            log,
		Auth:           authApp,
		User:           userApp,
		Article:        articleApp,
		Comment:        commentApp,
		Forum:          forumApp,
		BaseHandler:    baseHandler,
		UserHandler:    userHandler,
		ArticleHandler: articleHandler,
		CommentHandler: commentHandler,
		ForumHandler:   forumHandler,
		Geo:            geoClient,
		Logins:         loginRepo,
	}
}

// InitCron 注册已迁移 BC 的定时任务（浏览量同步）；未迁移任务仍由 initialize.InitCron 注册。
func InitCron(deps *ihttp.Deps, log *zap.Logger) {
	c := cron.New(cron.WithLogger(&cronZapLogger{log}))
	if _, err := c.AddFunc("@hourly", func() {
		if err := deps.Article.SyncViews(context.Background()); err != nil {
			log.Error("Failed to update article views:", zap.Error(err))
		}
	}); err != nil {
		log.Error("Error scheduling article views cron:", zap.Error(err))
		return
	}
	c.Start()
}
