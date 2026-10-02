package user

import (
	"errors"
	"time"

	authapp "server/internal/application/auth"
	userdomain "server/internal/domain/user"
	"server/internal/model/appTypes"
	"server/internal/model/database"
	"server/internal/model/request"

	"go.uber.org/zap"
)

// userEntity 领域用户实体别名（TokenNext 使用其行为方法）。
type userEntity = userdomain.User

// timeNow 当前时间（便于测试替换）。
var timeNow = time.Now

// toDBUser 领域用户 → GORM 模型（用于 response.Login 兼容序列化）。
func toDBUser(u *userdomain.User) database.User {
	return database.User{
		MODEL:     database.MODEL{ID: u.ID, CreatedAt: u.CreatedAt, UpdatedAt: u.UpdatedAt},
		UUID:      u.UUID,
		Username:  u.Username,
		Password:  u.Password,
		Email:     u.Email,
		Openid:    u.Openid,
		Avatar:    u.Avatar,
		Address:   u.Address,
		Signature: u.Signature,
		RoleID:    appTypes.RoleID(u.RoleID),
		Register:  appTypes.Register(u.Register),
		Freeze:    u.Freeze,
	}
}

// logErrorForToken 按错误类型记录与旧实现一致的日志。
func (h *Handler) logErrorForToken(err error) {
	switch {
	case errors.Is(err, authapp.ErrCreateAccess):
		h.log.Error("Failed to get accessToken:", zap.Error(err))
	case errors.Is(err, authapp.ErrCreateRefresh):
		h.log.Error("Failed to get refreshToken:", zap.Error(err))
	case errors.Is(err, authapp.ErrInvalidateJWT):
		h.log.Error("Failed to invalidate jwt:", zap.Error(err))
	default:
		h.log.Error("Failed to set login status:", zap.Error(err))
	}
}

// tokenErrMessage 按错误类型映射与旧实现一致的响应文案。
func tokenErrMessage(err error) string {
	switch {
	case errors.Is(err, authapp.ErrCreateAccess):
		return "Failed to get accessToken"
	case errors.Is(err, authapp.ErrCreateRefresh):
		return "Failed to get refreshToken"
	case errors.Is(err, authapp.ErrInvalidateJWT):
		return "Failed to invalidate jwt"
	default:
		return "Failed to set login status"
	}
}

// userListCond request.UserList → 领域查询条件。
func userListCond(req request.UserList) userdomain.ListCond {
	return userdomain.ListCond{
		UUID:     req.UUID,
		Register: req.Register,
		Page:     req.Page,
		PageSize: req.PageSize,
	}
}

// userLoginListCond request.UserLoginList → 领域查询条件。
func userLoginListCond(req request.UserLoginList) userdomain.LoginListCond {
	return userdomain.LoginListCond{
		UUID:     req.UUID,
		Page:     req.Page,
		PageSize: req.PageSize,
	}
}
