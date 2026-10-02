// Package base 提供 base 路由 Handler（验证码 / 邮箱验证码发送）。
package base

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"server/config"
	"server/internal/common/email"
	"server/internal/model/request"
	"server/internal/model/response"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

// Handler base 接口处理器。
type Handler struct {
	cfg    *config.Config
	store  base64Captcha.Store
	sender *email.Sender
	log    *zap.Logger
}

// NewHandler 构造 base 处理器。
func NewHandler(cfg *config.Config, store base64Captcha.Store, sender *email.Sender, log *zap.Logger) *Handler {
	return &Handler{cfg: cfg, store: store, sender: sender, log: log}
}

// Captcha 生成数字验证码（行为与原 api/base.go 一致）。
func (h *Handler) Captcha(c *gin.Context) {
	driver := base64Captcha.NewDriverDigit(
		h.cfg.Captcha.Height,
		h.cfg.Captcha.Width,
		h.cfg.Captcha.Length,
		h.cfg.Captcha.MaxSkew,
		h.cfg.Captcha.DotCount,
	)

	captcha := base64Captcha.NewCaptcha(driver, h.store)
	id, b64s, _, err := captcha.Generate()
	if err != nil {
		h.log.Error("Failed to generate captcha:", zap.Error(err))
		response.FailWithMessage("Failed to generate captcha", c)
		return
	}
	response.OkWithData(response.Captcha{
		CaptchaID: id,
		PicPath:   b64s,
	}, c)
}

// SendEmailVerificationCode 发送邮箱验证码（行为与原 api/base.go + service/base.go 一致）。
func (h *Handler) SendEmailVerificationCode(c *gin.Context) {
	var req request.SendEmailVerificationCode
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if h.store.Verify(req.CaptchaID, req.Captcha, true) {
		err = h.sendVerificationCode(c, req.Email)
		if err != nil {
			h.log.Error("Failed to send email:", zap.Error(err))
			response.FailWithMessage("Failed to send email", c)
			return
		}
		response.OkWithMessage("Successfully sent email", c)
		return
	}
	response.FailWithMessage("Incorrect verification code", c)
}

// sendVerificationCode 生成验证码、写入会话并发送邮件（行为与原 service/base.go 一致）。
func (h *Handler) sendVerificationCode(c *gin.Context, to string) error {
	verificationCode := generateVerificationCode(6)
	expireTime := time.Now().Add(5 * time.Minute).Unix()

	session := sessions.Default(c)
	session.Set("verification_code", verificationCode)
	session.Set("email", to)
	session.Set("expire_time", expireTime)
	_ = session.Save()

	subject := "您的邮箱验证码"
	body := `亲爱的用户[` + to + `]，<br/>
<br/>
感谢您注册` + h.cfg.Website.Name + `的个人博客！为了确保您的邮箱安全，请使用以下验证码进行验证：<br/>
<br/>
验证码：[<font color="blue"><u>` + verificationCode + `</u></font>]<br/>
该验证码在 5 分钟内有效，请尽快使用。<br/>
<br/>
如果您没有请求此验证码，请忽略此邮件。
<br/>
如有任何疑问，请联系我们的支持团队：<br/>
邮箱：` + h.cfg.Email.From + `<br/>
<br/>
祝好，<br/>` +
		h.cfg.Website.Title + `<br/>
<br/>`

	return h.sender.Send(to, subject, body)
}

// generateVerificationCode 生成指定长度的随机验证码（与 utils.GenerateVerificationCode 一致）。
func generateVerificationCode(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%0*d", length, r.Intn(int(math.Pow10(length))))
}
