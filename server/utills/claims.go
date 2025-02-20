package utills

import (
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"net"
	"server/global"
	"server/model/appTypes"
	"server/model/request"
)

func SetRefreshToken(c *gin.Context, token string, maxAge int) {
	host, _, err := net.SplitHostPort(c.Request.Host)
	if err != nil {
		host = c.Request.Host
	}
	setCookie(c, "x-refresh-token", token, host, maxAge)

}

func ClearRefreshToken(c *gin.Context) {
	host, _, err := net.SplitHostPort(c.Request.Host)
	if err != nil {
		host = c.Request.Host
	}
	setCookie(c, "x-refresh-token", "", host, -1)
}
func setCookie(c *gin.Context, name, value string, host string, maxAge int) {
	if net.ParseIP(host) != nil {
		c.SetCookie(name, value, maxAge, "/", "", false, true)
	} else {
		c.SetCookie(name, value, maxAge, "/", host, false, true)
	}
}

func GetAccessToken(c *gin.Context) string {
	token := c.Request.Header.Get("x-access-token")
	return token
}

func GetRefreshToken(c *gin.Context) string {
	token, _ := c.Cookie("x-refresh-token")
	return token
}

// GetClaims 从Gin的Context中解析并获取JWT的Claims
func GetClaims(c *gin.Context) (*request.JwtCustomClaims, error) {
	// 获取请求头中的Access Token
	token := GetAccessToken(c)
	// 创建JWT实例
	j := NewJWT()
	// 解析Access Token
	claims, err := j.ParseAccessToken(token)
	if err != nil {
		// 如果解析失败，记录错误日志
		global.Log.Error("Failed to retrieve JWT parsing information from Gin's Context. Please check if the request header contains 'x-access-token' and if the claims structure is correct.", zap.Error(err))
	}
	return claims, err
}

// GetRefreshClaims 从Gin的Context中解析并获取Refresh Token的Claims
func GetRefreshClaims(c *gin.Context) (*request.JwtCustomRefreshClaims, error) {
	// 获取Refresh Token
	token := GetRefreshToken(c)
	// 创建JWT实例
	j := NewJWT()
	// 解析Refresh Token
	claims, err := j.ParseRefreshToken(token)
	if err != nil {
		// 如果解析失败，记录错误日志
		global.Log.Error("Failed to retrieve JWT parsing information from Gin's Context. Please check if the request header contains 'x-refresh-token' and if the claims structure is correct.", zap.Error(err))
	}
	return claims, err
}

func GetUserInfo(c *gin.Context) *request.JwtCustomClaims {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return nil
		} else {
			return cl
		}
	} else {
		waitUse := claims.(*request.JwtCustomClaims)
		return waitUse
	}
}

func GetUserID(c *gin.Context) uint {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return 0
		} else {
			return cl.UserID
		}
	} else {
		waitUse := claims.(*request.JwtCustomClaims)
		return waitUse.UserID
	}
}

// GetUUID 从Gin的Context中获取JWT解析出来的用户UUID
func GetUUID(c *gin.Context) uuid.UUID {
	// 首先尝试从Context中获取"claims"
	if claims, exists := c.Get("claims"); !exists {
		// 如果不存在，则重新解析Access Token
		if cl, err := GetClaims(c); err != nil {
			// 如果解析失败，返回一个空UUID
			return uuid.UUID{}
		} else {
			// 返回解析出来的UUID
			return cl.UUID
		}
	} else {
		// 如果已存在claims，则直接返回UUID
		waitUse := claims.(*request.JwtCustomClaims)
		return waitUse.UUID
	}
}

// GetRoleID 从Gin的Context中获取JWT解析出来的用户角色ID
func GetRoleID(c *gin.Context) appTypes.RoleID {
	// 首先尝试从Context中获取"claims"
	if claims, exists := c.Get("claims"); !exists {
		// 如果不存在，则重新解析Access Token
		if cl, err := GetClaims(c); err != nil {
			// 如果解析失败，返回0
			return 0
		} else {
			// 返回解析出来的角色ID
			return cl.RoleID
		}
	} else {
		// 如果已存在claims，则直接返回角色ID
		waitUse := claims.(*request.JwtCustomClaims)
		return waitUse.RoleID
	}
}
