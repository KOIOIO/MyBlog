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
	articleapp "server/internal/application/article"
	authapp "server/internal/application/auth"
	userapp "server/internal/application/user"
	userdomain "server/internal/domain/user"
	articlehandler "server/internal/interface/http/handler/article"
	basehandler "server/internal/interface/http/handler/base"
	userhandler "server/internal/interface/http/handler/user"
	"server/internal/interface/http/middleware"
	oldrouter "server/router"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Deps 路由装配依赖。
type Deps struct {
	Config         *config.Config
	Log            *zap.Logger
	Auth           *authapp.AuthService
	User           *userapp.UserService
	Article        *articleapp.Service
	BaseHandler    *basehandler.Handler
	UserHandler    *userhandler.Handler
	ArticleHandler *articlehandler.Handler
	Geo            userdomain.GeoProvider
	Logins         userdomain.LoginRecordRepository
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

	routerGroup := oldrouter.RouterGroupApp

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

	// ---- 未迁移 BC：委托旧 router 注册（保持路由不丢失） ----
	{
		routerGroup.InitCommentRouter(privateGroup, publicGroup, adminGroup)
		routerGroup.InitFeedbackRouter(privateGroup, publicGroup, adminGroup)
		routerGroup.InitForumRouter(privateGroup, publicGroup, adminGroup)
	}
	{
		routerGroup.InitImageRouter(adminGroup)
		routerGroup.InitAdvertisementRouter(adminGroup, publicGroup)
		routerGroup.InitFriendLinkRouter(adminGroup, publicGroup)
		routerGroup.InitWebsiteRouter(adminGroup, publicGroup)
		routerGroup.InitConfigRouter(adminGroup)
	}

	return engine
}
