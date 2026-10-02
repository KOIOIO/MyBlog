// Package user 提供 user 路由 Handler（迁移自 api/user.go，依赖构造注入）。
package user

import (
	"server/config"
	authapp "server/internal/application/auth"
	userapp "server/internal/application/user"
	"server/internal/domain/auth"
	"server/internal/interface/http/middleware"
	"server/model/request"
	"server/model/response"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

// Handler user 接口处理器。
type Handler struct {
	users *userapp.UserService
	auth  *authapp.AuthService
	store base64Captcha.Store
	cfg   *config.Config
	log   *zap.Logger
}

// NewHandler 构造 user 处理器。
func NewHandler(users *userapp.UserService, authService *authapp.AuthService, store base64Captcha.Store, cfg *config.Config, log *zap.Logger) *Handler {
	return &Handler{users: users, auth: authService, store: store, cfg: cfg, log: log}
}

// Register 注册（行为与原 api/user.go 一致）。
func (h *Handler) Register(c *gin.Context) {
	var req request.Register
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	session := sessions.Default(c)
	// 两次邮箱一致性判断
	savedEmail := session.Get("email")
	if savedEmail == nil || savedEmail.(string) != req.Email {
		response.FailWithMessage("This email doesn't match the email to be verified", c)
		return
	}
	// 验证码一致性判断
	savedCode := session.Get("verification_code")
	if savedCode == nil || savedCode.(string) != req.VerificationCode {
		response.FailWithMessage("Invalid verification code", c)
		return
	}
	// 过期判断
	savedTime, ok := session.Get("expire_time").(int64)
	if !ok || savedTime < timeNow().Unix() {
		response.FailWithMessage("The verification code has expired, please resend it", c)
		return
	}

	u, err := h.users.Register(c.Request.Context(), req.Username, req.Password, req.Email)
	if err != nil {
		h.log.Error("Failed to register user:", zap.Error(err))
		response.FailWithMessage("Failed to register user", c)
		return
	}
	h.TokenNext(c, u)
}

// Login 登录接口（按 flag 分派，行为与原 api/user.go 一致）。
func (h *Handler) Login(c *gin.Context) {
	h.EmailLogin(c)
}

// EmailLogin 邮箱登录（行为与原 api/user.go 一致）。
func (h *Handler) EmailLogin(c *gin.Context) {
	var req request.Login
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if h.store.Verify(req.CaptchaID, req.Captcha, true) {
		u, err := h.users.EmailLogin(c.Request.Context(), req.Email, req.Password)
		if err != nil {
			h.log.Error("Failed to login:", zap.Error(err))
			response.FailWithMessage("Failed to login", c)
			return
		}
		h.TokenNext(c, u)
		return
	}
	response.FailWithMessage("Incorrect verification code", c)
}

// TokenNext 登录成功后生成 token（行为与原 api/user.go 一致）。
func (h *Handler) TokenNext(c *gin.Context, u *userEntity) {
	// 检查用户是否被冻结
	if u.IsFrozen() {
		response.FailWithMessage("The user is frozen, contact the administrator", c)
		return
	}

	pair, err := h.auth.GenerateToken(c.Request.Context(), auth.TokenSubject{
		UserID: u.ID,
		UUID:   u.UUID,
		RoleID: u.RoleID,
	}, h.cfg.System.UseMultipoint)
	if err != nil {
		h.logErrorForToken(err)
		response.FailWithMessage(tokenErrMessage(err), c)
		return
	}

	middleware.SetRefreshToken(c, pair.Refresh, int(pair.RefreshMaxAge))
	c.Set("user_id", u.ID)
	response.OkWithDetailed(response.Login{
		User:                 toDBUser(u),
		AccessToken:          pair.Access,
		AccessTokenExpiresAt: pair.AccessExpiresAt,
	}, "Successful login", c)
}

// ForgotPassword 找回密码（行为与原 api/user.go 一致）。
func (h *Handler) ForgotPassword(c *gin.Context) {
	var req request.ForgotPassword
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	session := sessions.Default(c)
	savedEmail := session.Get("email")
	if savedEmail == nil || savedEmail.(string) != req.Email {
		response.FailWithMessage("This email doesn't match the email to be verified", c)
		return
	}
	savedCode := session.Get("verification_code")
	if savedCode == nil || savedCode.(string) != req.VerificationCode {
		response.FailWithMessage("Invalid verification code", c)
		return
	}
	savedTime, ok := session.Get("expire_time").(int64)
	if !ok || savedTime < timeNow().Unix() {
		response.FailWithMessage("The verification code has expired, please resend it", c)
		return
	}

	err = h.users.ForgotPassword(c.Request.Context(), req.Email, req.NewPassword)
	if err != nil {
		h.log.Error("Failed to retrieve the password:", zap.Error(err))
		response.FailWithMessage("Failed to retrieve the password", c)
		return
	}
	response.OkWithMessage("Successfully retrieved", c)
}

// UserCard 获取用户卡片信息（行为与原 api/user.go 一致）。
func (h *Handler) UserCard(c *gin.Context) {
	var req request.UserCard
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	u, err := h.users.Card(c.Request.Context(), req.UUID)
	if err != nil {
		h.log.Error("Failed to get card:", zap.Error(err))
		response.FailWithMessage("Failed to get card", c)
		return
	}
	response.OkWithData(response.UserCard{
		UUID:      u.UUID,
		Username:  u.Username,
		Avatar:    u.Avatar,
		Address:   u.Address,
		Signature: u.Signature,
	}, c)
}

// Logout 登出（行为与原 api/user.go + service/user.go 一致）。
func (h *Handler) Logout(c *gin.Context) {
	h.logout(c)
	response.OkWithMessage("Successful logout", c)
}

// logout 登出流程：清 cookie + 删会话 + 刷新令牌入黑名单。
func (h *Handler) logout(c *gin.Context) {
	refreshToken := middleware.GetRefreshToken(c)
	middleware.ClearRefreshToken(c)
	_ = h.users.Logout(c.Request.Context(), middleware.GetUUID(c), refreshToken)
}

// UserResetPassword 修改密码（行为与原 api/user.go 一致）。
func (h *Handler) UserResetPassword(c *gin.Context) {
	var req request.UserResetPassword
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.UserID = middleware.GetUserID(c)
	err = h.users.ResetPassword(c.Request.Context(), req.UserID, req.Password, req.NewPassword)
	if err != nil {
		h.log.Error("Failed to modify:", zap.Error(err))
		response.FailWithMessage("Failed to modify, orginal password does not match the current account", c)
		return
	}
	response.OkWithMessage("Successfully changed password, please log in again", c)
	h.logout(c)
}

// UserInfo 获取个人信息（行为与原 api/user.go 一致）。
func (h *Handler) UserInfo(c *gin.Context) {
	u, err := h.users.Info(c.Request.Context(), middleware.GetUserID(c))
	if err != nil {
		h.log.Error("Failed to get user information:", zap.Error(err))
		response.FailWithMessage("Failed to get user information", c)
		return
	}
	response.OkWithData(u, c)
}

// UploadAvatar 上传头像（行为与原 api/user.go 一致）。
func (h *Handler) UploadAvatar(c *gin.Context) {
	_, header, err := c.Request.FormFile("avatar")
	if err != nil {
		h.log.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	url, err := h.users.UploadAvatar(c.Request.Context(), middleware.GetUserID(c), header)
	if err != nil {
		h.log.Error("Failed to upload avatar:", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.ImageUpload{
		Url:     url,
		OssType: h.cfg.System.OssType,
	}, "Successfully uploaded avatar", c)
}

// UserChangeInfo 修改个人信息（行为与原 api/user.go 一致）。
func (h *Handler) UserChangeInfo(c *gin.Context) {
	var req request.UserChangeInfo
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = h.users.ChangeInfo(c.Request.Context(), middleware.GetUserID(c), req.Username, req.Address, req.Signature)
	if err != nil {
		h.log.Error("Failed to change user information:", zap.Error(err))
		response.FailWithMessage("Failed to change user information", c)
		return
	}
	response.OkWithMessage("Successfully changed user information", c)
}

// UserWeather 获取天气（行为与原 api/user.go 一致）。
func (h *Handler) UserWeather(c *gin.Context) {
	ip := c.ClientIP()
	weather, err := h.users.Weather(c.Request.Context(), ip)
	if err != nil {
		h.log.Error("Failed to get user weather", zap.Error(err))
		response.FailWithMessage("Failed to get user weather", c)
		return
	}
	response.OkWithData(weather, c)
}

// UserChart 获取用户图表数据（行为与原 api/user.go 一致）。
func (h *Handler) UserChart(c *gin.Context) {
	var req request.UserChart
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	data, err := h.users.Chart(c.Request.Context(), req.Date)
	if err != nil {
		h.log.Error("Failed to get user chart:", zap.Error(err))
		response.FailWithMessage("Failed to user chart", c)
		return
	}
	response.OkWithData(response.UserChart{
		DateList:     data.DateList,
		LoginData:    data.LoginData,
		RegisterData: data.RegisterData,
	}, c)
}

// UserList 获取用户列表（行为与原 api/user.go 一致）。
func (h *Handler) UserList(c *gin.Context) {
	var pageInfo request.UserList
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := h.users.List(c.Request.Context(), userListCond(pageInfo))
	if err != nil {
		h.log.Error("Failed to get user list:", zap.Error(err))
		response.FailWithMessage("Failed to get user list", c)
		return
	}
	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)
}

// UserFreeze 冻结用户（行为与原 api/user.go 一致）。
func (h *Handler) UserFreeze(c *gin.Context) {
	var req request.UserOperation
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = h.users.Freeze(c.Request.Context(), req.ID)
	if err != nil {
		h.log.Error("Failed to freeze user:", zap.Error(err))
		response.FailWithMessage("Failed to freeze user", c)
		return
	}
	response.OkWithMessage("Successfully freeze user", c)
}

// UserUnfreeze 解冻用户（行为与原 api/user.go 一致）。
func (h *Handler) UserUnfreeze(c *gin.Context) {
	var req request.UserOperation
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = h.users.Unfreeze(c.Request.Context(), req.ID)
	if err != nil {
		h.log.Error("Failed to unfreeze user:", zap.Error(err))
		response.FailWithMessage("Failed to unfreeze user", c)
		return
	}
	response.OkWithMessage("Successfully unfreeze user", c)
}

// UserLoginList 获取登录日志列表（行为与原 api/user.go 一致）。
func (h *Handler) UserLoginList(c *gin.Context) {
	var pageInfo request.UserLoginList
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := h.users.LoginList(c.Request.Context(), userLoginListCond(pageInfo))
	if err != nil {
		h.log.Error("Failed to get user login list:", zap.Error(err))
		response.FailWithMessage("Failed to get user login list", c)
		return
	}
	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)
}
