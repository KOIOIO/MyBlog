package middleware

import (
	"github.com/gin-gonic/gin"
	"server/model/appTypes"
	"server/model/response"
	"server/utills"
)

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID := utills.GetRoleID(c)
		if roleID != appTypes.Admin {
			response.Forbidden("Access defined,Admin privileges are required", c)
			c.Abort()
			return
		}
		c.Next()
	}
}
