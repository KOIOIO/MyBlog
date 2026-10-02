package middleware

import (
	"errors"
	"strconv"

	authapp "server/internal/application/auth"
	authdomain "server/internal/domain/auth"
	"server/internal/model/appTypes"
	"server/internal/model/response"

	"github.com/gin-gonic/gin"
)

// JWTAuth 基于注入 AuthService 的 JWT 鉴权中间件（行为与原 middleware/jwt.go 一致）。
func JWTAuth(authService *authapp.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := GetAccessToken(c)
		refreshToken := GetRefreshToken(c)

		result, err := authService.Authenticate(c.Request.Context(), accessToken, refreshToken)
		if err != nil {
			switch {
			case errors.Is(err, authdomain.ErrRefreshBlacklisted):
				ClearRefreshToken(c)
				response.NoAuth("Account logged in from another location or token is invalid", c)
			case errors.Is(err, authdomain.ErrRefreshInvalid):
				ClearRefreshToken(c)
				response.NoAuth("Refresh token expired or invalid", c)
			case errors.Is(err, authdomain.ErrUserNotFound):
				ClearRefreshToken(c)
				response.NoAuth("The user does not exist", c)
			case errors.Is(err, authdomain.ErrTokenCreate):
				ClearRefreshToken(c)
				response.NoAuth("Faild to create new access token", c)
			default:
				ClearRefreshToken(c)
				response.NoAuth("Invalid access token", c)
			}
			c.Abort()
			return
		}

		if result.NewAccessToken != "" {
			c.Header("new-access-token", result.NewAccessToken)
			c.Header("new-access-expires-at", strconv.FormatInt(result.NewAccessAt, 10))
		}
		c.Set("claims", toRequestClaims(result.Claims))
		c.Next()
	}
}

// AdminAuth 管理员权限中间件（行为与原 middleware/admin.go 一致）。
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if GetRoleID(c) != appTypes.Admin {
			response.Forbidden("Access denied. Adimin privileges are required", c)
			c.Abort()
			return
		}
		c.Next()
	}
}
