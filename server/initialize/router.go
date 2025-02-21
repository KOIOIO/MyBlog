package initialize

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"server/global"
	"server/middleware"
	"server/router"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	// 设置gin模式
	gin.SetMode(global.Config.System.Env)
	Router := gin.Default()
	//使用日志记录中间键
	Router.Use(middleware.GinLogger(), middleware.GinRecovery(true))
	//使用gin会话路由
	var store = cookie.NewStore([]byte(global.Config.System.SessionsSecret))
	Router.Use(sessions.Sessions("session", store))

	routerGroup := router.RouterGroupApp

	publicGroup := Router.Group(global.Config.System.RouterPrefix)
	privateGroup := Router.Group((global.Config.System.RouterPrefix))
	privateGroup.Use(middleware.JWTAuth())
	adminGroup := Router.Group((global.Config.System.RouterPrefix))
	adminGroup.Use(middleware.JWTAuth()).Use(middleware.AdminAuth())
	{
		routerGroup.InitBaseRouter(publicGroup)
	}
	{
		routerGroup.InitUserRouter(privateGroup, publicGroup, adminGroup)
	}
	return Router
}
