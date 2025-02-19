package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"server/global"
)

type BaseApi struct {
}

var store = base64Captcha.DefaultMemStore

func (baseApi *BaseApi) Captcha(c *gin.Context) {
	driver := base64Captcha.NewDriverDigit(
		global.Config.Captcha.Height,
		global.Config.Captcha.Width,
		global.Config.Captcha.Length,
		global.Config.Captcha.MaxSkew,
		global.Config.Captcha.DotCount,
	)
	captcha := base64Captcha.NewCaptcha(driver, store)

	id, b64s, _, err := captcha.Generate()
	if err != nil {
		global.Log.Error("验证码获取失败!", zap.Error(err))

	}
}
