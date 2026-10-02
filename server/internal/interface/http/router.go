// Package http 提供 Gin 路由装配（替代 initialize/router.go）。
//
// 绞杀者迁移说明：
//   - Phase 0：仅建立骨架，按 BC 注册空路由组，未挂载业务，亦未接入 main。
//   - Phase 1 起：本包成为主路由入口，已迁移 BC 使用 internal/handler 下的新 Handler，
//     未迁移 BC 委托旧 router/ 文件注册，保证路由不丢失。
//   - Phase 5：删除旧 router/ 后，本包为唯一路由定义处。
package http

import (
	"net/http"

	"server/config"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// NewRouter 构建 Gin Engine（骨架：中间件 + 三组路由 + 空 BC 子组）。
func NewRouter(cfg *config.Config) *gin.Engine {
	gin.SetMode(cfg.System.Env)
	router := gin.Default()
	router.Use(ginLoggerPlaceholder(), ginRecoveryPlaceholder(true))

	// 会话中间件（与旧实现一致）
	var store = cookie.NewStore([]byte(cfg.System.SessionsSecret))
	router.Use(sessions.Sessions("session", store))

	// 静态文件
	router.StaticFS(cfg.Upload.Path, http.Dir(cfg.Upload.Path))

	publicGroup := router.Group(cfg.System.RouterPrefix)
	privateGroup := router.Group(cfg.System.RouterPrefix)
	adminGroup := router.Group(cfg.System.RouterPrefix)

	// TODO(Phase 1): 中间件改为注入式后在此挂载 JWTAuth/AdminAuth。
	_ = privateGroup
	_ = adminGroup

	// TODO(各 Phase)：按 BC 迁移路由。
	// 已迁移 BC 在此注册 internal/handler 下的新 Handler；
	// 未迁移 BC 调用旧 router 的 Init 函数保持路由不丢失。
	_ = publicGroup
	_ = cfg

	return router
}

// ginLoggerPlaceholder 骨架占位：Phase 5 迁移 middleware/logger.go 后替换为真实实现。
func ginLoggerPlaceholder() gin.HandlerFunc {
	return func(c *gin.Context) { c.Next() }
}

// ginRecoveryPlaceholder 骨架占位：Phase 5 迁移 middleware/logger.go 后替换为真实实现。
func ginRecoveryPlaceholder(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) { c.Next() }
}
