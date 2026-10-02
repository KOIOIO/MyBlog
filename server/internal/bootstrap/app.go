package bootstrap

import (
	"context"
	"time"

	"server/internal/application/advertisement"
	agentapp "server/internal/application/agent"
	"server/internal/application/article"
	"server/internal/application/auth"
	"server/internal/application/comment"
	"server/internal/application/config"
	"server/internal/application/feedback"
	"server/internal/application/forum"
	"server/internal/application/friendlink"
	"server/internal/application/image"
	"server/internal/application/user"
	"server/internal/application/website"
	"server/internal/common/email"
	"server/internal/common/jwt"
	"server/internal/infrastructure/calendar"
	"server/internal/infrastructure/configfile"
	"server/internal/infrastructure/es"
	"server/internal/infrastructure/geo"
	"server/internal/infrastructure/hotsearch"
	"server/internal/infrastructure/llm"
	"server/internal/infrastructure/mysql"
	"server/internal/infrastructure/redis"
	"server/internal/infrastructure/storage"
	ihttp "server/internal/interface/http"
	advertisementhandler "server/internal/interface/http/handler/advertisement"
	agenthandler "server/internal/interface/http/handler/agent"
	articlehandler "server/internal/interface/http/handler/article"
	basehandler "server/internal/interface/http/handler/base"
	commenthandler "server/internal/interface/http/handler/comment"
	confighandler "server/internal/interface/http/handler/config"
	feedbackhandler "server/internal/interface/http/handler/feedback"
	forumhandler "server/internal/interface/http/handler/forum"
	friendlinkhandler "server/internal/interface/http/handler/friendlink"
	imagehandler "server/internal/interface/http/handler/image"
	userhandler "server/internal/interface/http/handler/user"
	websitehandler "server/internal/interface/http/handler/website"

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
	commentRepo := mysql.NewCommentRepo(infra.DB, infra.ESClient)
	commentApp := comment.NewService(commentRepo, log)

	// ---- forum BC ----
	forumRepo := mysql.NewForumRepo(infra.DB)
	forumApp := forum.NewService(forumRepo, cfg, log)

	// ---- image / storage BC ----
	imageRepo := mysql.NewImageRepo(infra.DB)
	fileStorage := storage.New(cfg.System.OssType, cfg)
	imageApp := image.NewService(imageRepo, fileStorage, cfg.System.Storage().String())

	// ---- website BC（消费图片/友链/热搜/日历端口）----
	websiteRepo := mysql.NewWebsiteRepo(infra.DB)
	hotSearchProvider := hotsearch.NewProvider(&infra.Redis)
	calendarProvider := calendar.NewProvider(&infra.Redis)
	websiteApp := website.NewService(websiteRepo, websiteRepo, hotSearchProvider, calendarProvider, cfg)

	// ---- advertisement / friendLink / feedback BC ----
	advertisementRepo := mysql.NewAdvertisementRepo(infra.DB)
	advertisementApp := advertisement.NewService(advertisementRepo)
	friendLinkRepo := mysql.NewFriendLinkRepo(infra.DB)
	friendLinkApp := friendlink.NewService(friendLinkRepo)
	feedbackRepo := mysql.NewFeedbackRepo(infra.DB)
	feedbackApp := feedback.NewService(feedbackRepo)

	// ---- config BC（配置指针 + YAML 持久化）----
	configStore := configfile.NewStore(cfg, "config.yaml")
	configApp := config.NewService(cfg, configStore, websiteRepo)

	// ---- interface ----
	captchaStore := base64Captcha.DefaultMemStore
	emailSender := email.New(cfg.Email)
	baseHandler := basehandler.NewHandler(cfg, captchaStore, emailSender, log)
	userHandler := userhandler.NewHandler(userApp, authApp, captchaStore, cfg, log)
	articleHandler := articlehandler.NewHandler(articleApp, log)
	commentHandler := commenthandler.NewHandler(commentApp, log)
	forumHandler := forumhandler.NewHandler(forumApp, cfg, log)
	imageHandler := imagehandler.NewHandler(imageApp, cfg, log)
	websiteHandler := websitehandler.NewHandler(websiteApp, cfg, log)
	advertisementHandler := advertisementhandler.NewHandler(advertisementApp, log)
	friendLinkHandler := friendlinkhandler.NewHandler(friendLinkApp, log)
	feedbackHandler := feedbackhandler.NewHandler(feedbackApp, log)
	configHandler := confighandler.NewHandler(configApp, cfg, log)

	// ---- agent BC（AI 助手）----
	agentConvs := mysql.NewAgentConversationRepo(infra.DB)
	agentMsgs := mysql.NewAgentMessageRepo(infra.DB)
	agentMems := mysql.NewAgentMemoryRepo(infra.DB)
	agentArticles := es.NewAgentArticleReader(esArticleStore)
	modelProvider := llm.NewDashScopeProvider(&cfg.LLM, log)
	agentApp := agentapp.NewService(cfg, log, agentConvs, agentMsgs, agentMems, agentArticles, modelProvider)
	agentHandler := agenthandler.NewHandler(agentApp, log)

	return &ihttp.Deps{
		Config:               cfg,
		Log:                  log,
		Auth:                 authApp,
		User:                 userApp,
		Article:              articleApp,
		Comment:              commentApp,
		Forum:                forumApp,
		Image:                imageApp,
		Website:              websiteApp,
		Advertisement:        advertisementApp,
		FriendLink:           friendLinkApp,
		Feedback:             feedbackApp,
		ConfigSvc:            configApp,
		BaseHandler:          baseHandler,
		UserHandler:          userHandler,
		ArticleHandler:       articleHandler,
		CommentHandler:       commentHandler,
		ForumHandler:         forumHandler,
		ImageHandler:         imageHandler,
		WebsiteHandler:       websiteHandler,
		AdvertisementHandler: advertisementHandler,
		AgentHandler:         agentHandler,
		FriendLinkHandler:    friendLinkHandler,
		FeedbackHandler:      feedbackHandler,
		ConfigHandler:        configHandler,
		Geo:                  geoClient,
		Logins:               loginRepo,
	}
}

// InitCron 注册已迁移 BC 的定时任务（文章浏览量同步 / 热搜预热 / 日历预热）。
func InitCron(deps *ihttp.Deps, log *zap.Logger) {
	c := cron.New(cron.WithLogger(&cronZapLogger{log}))
	if _, err := c.AddFunc("@hourly", func() {
		if err := deps.Article.SyncViews(context.Background()); err != nil {
			log.Error("Failed to update article views:", zap.Error(err))
		}
		if err := deps.Website.WarmHotSearch(context.Background()); err != nil {
			log.Error("Failed to get hot list:", zap.Error(err))
		}
	}); err != nil {
		log.Error("Error scheduling hourly cron:", zap.Error(err))
		return
	}
	if _, err := c.AddFunc("@daily", func() {
		if _, err := deps.Website.Calendar(context.Background(), time.Now().Format("2006/0102")); err != nil {
			log.Error("Failed to get calendar:", zap.Error(err))
		}
	}); err != nil {
		log.Error("Error scheduling daily cron:", zap.Error(err))
		return
	}
	c.Start()
}
