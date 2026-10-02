package middleware

import (
	"context"

	userdomain "server/internal/domain/user"

	"github.com/gin-gonic/gin"
	"github.com/ua-parser/uap-go/uaparser"
	"go.uber.org/zap"
)

// LoginRecord 登录日志中间件（行为与原 middleware/login_record.go 一致，异步写入）。
func LoginRecord(geo userdomain.GeoProvider, logins userdomain.LoginRecordRepository, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 异步记录日志
		go func() {
			var userID uint
			var address string
			ip := c.ClientIP()
			loginMethod := c.DefaultQuery("flag", "email") // 若未传递flag参数，则默认为"email"
			userAgent := c.Request.UserAgent()

			if value, exists := c.Get("user_id"); exists {
				if id, ok := value.(uint); ok {
					userID = id
				}
			}

			address = getAddressFromIP(ip, geo)
			os, deviceInfo, browserInfo := parseUserAgent(userAgent)

			rec := &userdomain.LoginRecord{
				UserID:      userID,
				LoginMethod: loginMethod,
				IP:          ip,
				Address:     address,
				OS:          os,
				DeviceInfo:  deviceInfo,
				BrowserInfo: browserInfo,
				Status:      c.Writer.Status(),
			}
			if err := logins.Create(context.Background(), rec); err != nil {
				log.Error("Failed to record login", zap.Error(err))
			}
		}()
	}
}

// getAddressFromIP 获取 IP 对应的地理位置信息。
func getAddressFromIP(ip string, geo userdomain.GeoProvider) string {
	res, err := geo.LocationByIP(context.Background(), ip)
	if err != nil || res.Province == "" {
		return "未知"
	}
	if res.City != "" && res.Province != res.City {
		return res.Province + "-" + res.City
	}
	return res.Province
}

// parseUserAgent 解析 User-Agent，提取操作系统、设备信息和浏览器信息。
func parseUserAgent(userAgent string) (os, deviceInfo, browserInfo string) {
	os = userAgent
	deviceInfo = userAgent
	browserInfo = userAgent

	parser := uaparser.NewFromSaved()
	cli := parser.Parse(userAgent)
	os = cli.Os.Family
	deviceInfo = cli.Device.Family
	browserInfo = cli.UserAgent.Family

	return
}
