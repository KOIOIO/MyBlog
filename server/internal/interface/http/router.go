// Package http 提供 Gin 路由装配（替代 initialize/router.go）。
//
// 绞杀者迁移说明：
//   - 已迁移 BC（base、user）在本文件注册新 Handler；
//   - 未迁移 BC 委托旧 router/ 的 Init 函数注册，保证路由不丢失；
//   - 中间件（JWT/Admin/日志/登录记录）已迁移为注入式（internal/interface/http/middleware）。
package http

import (
	"net/http"

	"server/config"
	advertisementapp "server/internal/application/advertisement"
	articleapp "server/internal/application/article"
	authapp "server/internal/application/auth"
	commentapp "server/internal/application/comment"
	configapp "server/internal/application/config"
	feedbackapp "server/internal/application/feedback"
	forumapp "server/internal/application/forum"
	friendlinkapp "server/internal/application/friendlink"
	imageapp "server/internal/application/image"
	userapp "server/internal/application/user"
	websiteapp "server/internal/application/website"
	userdomain "server/internal/domain/user"
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
	"server/internal/interface/http/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Deps 路由装配依赖。
type Deps struct {
	Config               *config.Config
	Log                  *zap.Logger
	Auth                 *authapp.AuthService
	User                 *userapp.UserService
	Article              *articleapp.Service
	Comment              *commentapp.Service
	Forum                *forumapp.Service
	Image                *imageapp.Service
	Website              *websiteapp.Service
	Advertisement        *advertisementapp.Service
	FriendLink           *friendlinkapp.Service
	Feedback             *feedbackapp.Service
	ConfigSvc            *configapp.Service
	BaseHandler          *basehandler.Handler
	UserHandler          *userhandler.Handler
	ArticleHandler       *articlehandler.Handler
	CommentHandler       *commenthandler.Handler
	ForumHandler         *forumhandler.Handler
	ImageHandler         *imagehandler.Handler
	WebsiteHandler       *websitehandler.Handler
	AdvertisementHandler *advertisementhandler.Handler
	AgentHandler         *agenthandler.Handler
	FriendLinkHandler    *friendlinkhandler.Handler
	FeedbackHandler      *feedbackhandler.Handler
	ConfigHandler        *confighandler.Handler
	Geo                  userdomain.GeoProvider
	Logins               userdomain.LoginRecordRepository
}

// NewRouter 构建 Gin Engine（中间件链 + 三组路由 + 已迁移/未迁移 BC 路由）。
func NewRouter(deps *Deps) *gin.Engine {
	cfg := deps.Config
	gin.SetMode(cfg.System.Env)
	engine := gin.Default()
	engine.Use(middleware.GinLogger(deps.Log), middleware.GinRecovery(deps.Log, true))

	var store = cookie.NewStore([]byte(cfg.System.SessionsSecret))
	engine.Use(sessions.Sessions("session", store))

	engine.StaticFS(cfg.Upload.Path, http.Dir(cfg.Upload.Path))

	publicGroup := engine.Group(cfg.System.RouterPrefix)
	privateGroup := engine.Group(cfg.System.RouterPrefix)
	privateGroup.Use(middleware.JWTAuth(deps.Auth))
	adminGroup := engine.Group(cfg.System.RouterPrefix)
	adminGroup.Use(middleware.JWTAuth(deps.Auth)).Use(middleware.AdminAuth())

	// ---- 已迁移 BC：base ----
	{
		baseRouter := publicGroup.Group("base")
		baseRouter.POST("captcha", deps.BaseHandler.Captcha)
		baseRouter.POST("sendEmailVerificationCode", deps.BaseHandler.SendEmailVerificationCode)
	}

	// ---- 已迁移 BC：user ----
	{
		userRouter := privateGroup.Group("user")
		userPublicRouter := publicGroup.Group("user")
		userLoginRouter := publicGroup.Group("user").Use(middleware.LoginRecord(deps.Geo, deps.Logins, deps.Log))
		userAdminRouter := adminGroup.Group("user")
		userHandler := deps.UserHandler
		{
			userRouter.POST("logout", userHandler.Logout)
			userRouter.PUT("resetPassword", userHandler.UserResetPassword)
			userRouter.GET("info", userHandler.UserInfo)
			userRouter.PUT("changeInfo", userHandler.UserChangeInfo)
			userRouter.POST("avatar", userHandler.UploadAvatar)
			userRouter.GET("weather", userHandler.UserWeather)
			userRouter.GET("chart", userHandler.UserChart)
		}
		{
			userPublicRouter.POST("forgotPassword", userHandler.ForgotPassword)
			userPublicRouter.GET("card", userHandler.UserCard)
		}
		{
			userLoginRouter.POST("register", userHandler.Register)
			userLoginRouter.POST("login", userHandler.Login)
		}
		{
			userAdminRouter.GET("list", userHandler.UserList)
			userAdminRouter.PUT("freeze", userHandler.UserFreeze)
			userAdminRouter.PUT("unfreeze", userHandler.UserUnfreeze)
			userAdminRouter.GET("loginList", userHandler.UserLoginList)
		}
	}

	// ---- 已迁移 BC：article ----
	{
		articleRouter := privateGroup.Group("article")
		articlePublicRouter := publicGroup.Group("article")
		articleAdminRouter := adminGroup.Group("article")
		articleHandler := deps.ArticleHandler
		{
			articleRouter.POST("like", articleHandler.Like)
			articleRouter.GET("isLike", articleHandler.IsLike)
			articleRouter.GET("likesList", articleHandler.LikesList)
		}
		{
			articlePublicRouter.GET(":id", articleHandler.InfoByID)
			articlePublicRouter.GET("search", articleHandler.Search)
			articlePublicRouter.GET("category", articleHandler.Category)
			articlePublicRouter.GET("tags", articleHandler.Tags)
		}
		{
			articleAdminRouter.POST("create", articleHandler.Create)
			articleAdminRouter.DELETE("delete", articleHandler.Delete)
			articleAdminRouter.PUT("update", articleHandler.Update)
			articleAdminRouter.PUT("setTop", articleHandler.SetTop)
			articleAdminRouter.GET("list", articleHandler.List)
		}
	}

	// ---- 已迁移 BC：comment ----
	{
		commentRouter := privateGroup.Group("comment")
		commentPublicRouter := publicGroup.Group("comment")
		commentAdminRouter := adminGroup.Group("comment")
		commentHandler := deps.CommentHandler
		{
			commentRouter.POST("create", commentHandler.Create)
			commentRouter.DELETE("delete", commentHandler.Delete)
			commentRouter.GET("info", commentHandler.Info)
		}
		{
			commentPublicRouter.GET(":article_id", commentHandler.InfoByArticleID)
			commentPublicRouter.GET("new", commentHandler.New)
		}
		{
			commentAdminRouter.GET("list", commentHandler.List)
		}
	}

	// ---- 已迁移 BC：forum ----
	{
		forumRouter := privateGroup.Group("forum")
		forumPublicRouter := publicGroup.Group("forum")
		forumHandler := deps.ForumHandler
		{
			forumPublicRouter.GET("list", forumHandler.List)
			forumPublicRouter.GET("detail", forumHandler.Detail)
			forumPublicRouter.GET("tags", forumHandler.Tags)
		}
		{
			forumRouter.POST("publish", forumHandler.Publish)
			forumRouter.POST("upload", forumHandler.Upload)
			forumRouter.POST("like", forumHandler.Like)
			forumRouter.POST("comment", forumHandler.Comment)
			forumRouter.GET("manageList", forumHandler.ManageList)
			forumRouter.DELETE("delete", forumHandler.Delete)
			forumRouter.GET("manageComments", forumHandler.ManageComments)
			forumRouter.DELETE("comment", forumHandler.CommentDelete)
		}
	}

	// ---- 已迁移 BC：agent ----
	{
		agentRouter := privateGroup.Group("agent")
		agentHandler := deps.AgentHandler
		{
			agentRouter.POST("conversations", agentHandler.Create)
			agentRouter.GET("conversations", agentHandler.List)
			agentRouter.GET("conversations/:id/messages", agentHandler.Messages)
			agentRouter.DELETE("conversations/:id", agentHandler.Delete)
			agentRouter.POST("chat", agentHandler.Chat)
		}
	}

	// ---- 已迁移 BC：image ----
	{
		imageRouter := adminGroup.Group("image")
		imageHandler := deps.ImageHandler
		{
			imageRouter.POST("upload", imageHandler.Upload)
			imageRouter.DELETE("delete", imageHandler.Delete)
			imageRouter.GET("list", imageHandler.List)
		}
	}

	// ---- 已迁移 BC：website ----
	{
		websiteRouter := privateGroup.Group("website")
		websitePublicRouter := publicGroup.Group("website")
		websiteHandler := deps.WebsiteHandler
		{
			websiteRouter.POST("addCarousel", websiteHandler.AddCarousel)
			websiteRouter.PUT("cancelCarousel", websiteHandler.CancelCarousel)
			websiteRouter.POST("createFooterLink", websiteHandler.CreateFooterLink)
			websiteRouter.DELETE("deleteFooterLink", websiteHandler.DeleteFooterLink)
		}
		{
			websitePublicRouter.GET("logo", websiteHandler.Logo)
			websitePublicRouter.GET("title", websiteHandler.Title)
			websitePublicRouter.GET("info", websiteHandler.Info)
			websitePublicRouter.GET("carousel", websiteHandler.Carousel)
			websitePublicRouter.GET("news", websiteHandler.News)
			websitePublicRouter.GET("calendar", websiteHandler.Calendar)
			websitePublicRouter.GET("footerLink", websiteHandler.FooterLink)
		}
	}

	// ---- 已迁移 BC：advertisement ----
	{
		advertisementRouter := privateGroup.Group("advertisement")
		advertisementPublicRouter := publicGroup.Group("advertisement")
		advertisementHandler := deps.AdvertisementHandler
		{
			advertisementRouter.POST("create", advertisementHandler.Create)
			advertisementRouter.DELETE("delete", advertisementHandler.Delete)
			advertisementRouter.PUT("update", advertisementHandler.Update)
			advertisementRouter.GET("list", advertisementHandler.List)
		}
		{
			advertisementPublicRouter.GET("info", advertisementHandler.Info)
		}
	}

	// ---- 已迁移 BC：friendLink ----
	{
		friendLinkRouter := privateGroup.Group("friendLink")
		friendLinkPublicRouter := publicGroup.Group("friendLink")
		friendLinkHandler := deps.FriendLinkHandler
		{
			friendLinkRouter.POST("create", friendLinkHandler.Create)
			friendLinkRouter.DELETE("delete", friendLinkHandler.Delete)
			friendLinkRouter.PUT("update", friendLinkHandler.Update)
			friendLinkRouter.GET("list", friendLinkHandler.List)
		}
		{
			friendLinkPublicRouter.GET("info", friendLinkHandler.Info)
		}
	}

	// ---- 已迁移 BC：feedback ----
	{
		feedbackRouter := privateGroup.Group("feedback")
		feedbackPublicRouter := publicGroup.Group("feedback")
		feedbackAdminRouter := adminGroup.Group("feedback")
		feedbackHandler := deps.FeedbackHandler
		{
			feedbackRouter.POST("create", feedbackHandler.Create)
			feedbackRouter.GET("info", feedbackHandler.Info)
		}
		{
			feedbackPublicRouter.GET("new", feedbackHandler.New)
		}
		{
			feedbackAdminRouter.DELETE("delete", feedbackHandler.Delete)
			feedbackAdminRouter.PUT("reply", feedbackHandler.Reply)
			feedbackAdminRouter.GET("list", feedbackHandler.List)
		}
	}

	// ---- 已迁移 BC：config ----
	{
		configRouter := privateGroup.Group("config")
		configHandler := deps.ConfigHandler
		{
			configRouter.GET("website", configHandler.GetWebsite)
			configRouter.PUT("website", configHandler.UpdateWebsite)
			configRouter.GET("system", configHandler.GetSystem)
			configRouter.PUT("system", configHandler.UpdateSystem)
			configRouter.GET("email", configHandler.GetEmail)
			configRouter.PUT("email", configHandler.UpdateEmail)
			configRouter.GET("qiniu", configHandler.GetQiniu)
			configRouter.PUT("qiniu", configHandler.UpdateQiniu)
			configRouter.GET("jwt", configHandler.GetJwt)
			configRouter.PUT("jwt", configHandler.UpdateJwt)
			configRouter.GET("gaode", configHandler.GetGaode)
			configRouter.PUT("gaode", configHandler.UpdateGaode)
		}
	}

	return engine
}
