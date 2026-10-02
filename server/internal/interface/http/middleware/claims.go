// Package middleware 提供注入式 Gin 中间件（JWT/Admin/登录记录/日志）。
package middleware

import (
	"net"
	"time"

	"server/internal/domain/auth"
	"server/internal/model/appTypes"
	"server/internal/model/request"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
)

// SetRefreshToken 设置 Refresh Token 的 cookie。
func SetRefreshToken(c *gin.Context, token string, maxAge int) {
	host, _, err := net.SplitHostPort(c.Request.Host)
	if err != nil {
		host = c.Request.Host
	}
	setCookie(c, "x-refresh-token", token, maxAge, host)
}

// ClearRefreshToken 清除 Refresh Token 的 cookie。
func ClearRefreshToken(c *gin.Context) {
	host, _, err := net.SplitHostPort(c.Request.Host)
	if err != nil {
		host = c.Request.Host
	}
	setCookie(c, "x-refresh-token", "", -1, host)
}

// setCookie 设置指定名称和值的 cookie（与 utils/claims.go 行为一致）。
func setCookie(c *gin.Context, name, value string, maxAge int, host string) {
	if net.ParseIP(host) != nil {
		c.SetCookie(name, value, maxAge, "/", "", false, true)
	} else {
		c.SetCookie(name, value, maxAge, "/", host, false, true)
	}
}

// GetAccessToken 从请求头获取 Access Token。
func GetAccessToken(c *gin.Context) string {
	return c.Request.Header.Get("x-access-token")
}

// GetRefreshToken 从 cookie 获取 Refresh Token。
func GetRefreshToken(c *gin.Context) string {
	token, _ := c.Cookie("x-refresh-token")
	return token
}

// getClaims 从 Gin 上下文读取 JWT 声明（JWTAuth 中间件保证已设置）。
func getClaims(c *gin.Context) *request.JwtCustomClaims {
	if claims, exists := c.Get("claims"); exists {
		if waitUse, ok := claims.(*request.JwtCustomClaims); ok {
			return waitUse
		}
	}
	return nil
}

// GetUserInfo 获取 JWT 解析出来的用户信息。
func GetUserInfo(c *gin.Context) *request.JwtCustomClaims {
	return getClaims(c)
}

// GetUserID 获取用户 ID。
func GetUserID(c *gin.Context) uint {
	if claims := getClaims(c); claims != nil {
		return claims.UserID
	}
	return 0
}

// GetUUID 获取用户 UUID。
func GetUUID(c *gin.Context) uuid.UUID {
	if claims := getClaims(c); claims != nil {
		return claims.UUID
	}
	return uuid.UUID{}
}

// GetRoleID 获取用户角色 ID。
func GetRoleID(c *gin.Context) appTypes.RoleID {
	if claims := getClaims(c); claims != nil {
		return claims.RoleID
	}
	return 0
}

// toRequestClaims 领域声明 → Gin 上下文兼容声明（保留 Audience/ExpiresAt 结构）。
func toRequestClaims(c auth.AccessClaims) *request.JwtCustomClaims {
	return &request.JwtCustomClaims{
		BaseClaims: request.BaseClaims{
			UserID: c.UserID,
			UUID:   c.UUID,
			RoleID: appTypes.RoleID(c.RoleID),
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"TAP"},
			ExpiresAt: jwt.NewNumericDate(time.Unix(c.ExpiresAt, 0)),
		},
	}
}
