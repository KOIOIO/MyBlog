// Package auth 提供认证用例编排（实现 domain/auth.AuthPort）。
package auth

import (
	"context"
	"errors"
	"time"

	jwtutil "server/internal/common/jwt"
	"server/internal/domain/auth"
	"server/internal/domain/shared"
	"server/internal/domain/user"
	"server/model/appTypes"
	"server/model/request"
)

// 登录流程错误（handler 据此映射固定文案，文案必须与重构前一致）。
var (
	ErrCreateAccess   = errors.New("failed to create access token")
	ErrCreateRefresh  = errors.New("failed to create refresh token")
	ErrSetLoginStatus = errors.New("failed to set login status")
	ErrInvalidateJWT  = errors.New("failed to invalidate jwt")
)

// AuthService 认证用例。
type AuthService struct {
	jwt        *jwtutil.JWT
	sessions   auth.SessionStore
	blacklist  auth.BlacklistStore
	users      user.UserRepository
	refreshTTL time.Duration
}

// NewAuthService 构造认证用例。
func NewAuthService(jwtUtil *jwtutil.JWT, sessions auth.SessionStore, blacklist auth.BlacklistStore, users user.UserRepository, refreshTTL time.Duration) *AuthService {
	return &AuthService{jwt: jwtUtil, sessions: sessions, blacklist: blacklist, users: users, refreshTTL: refreshTTL}
}

// GenerateToken 签发令牌对（含多地点登录拦截逻辑）。
func (s *AuthService) GenerateToken(ctx context.Context, sub auth.TokenSubject, useMultipoint bool) (auth.TokenPair, error) {
	base := request.BaseClaims{UserID: sub.UserID, UUID: sub.UUID, RoleID: appTypes.RoleID(sub.RoleID)}

	accessClaims := s.jwt.CreateAccessClaims(base)
	accessToken, err := s.jwt.CreateAccessToken(accessClaims)
	if err != nil {
		return auth.TokenPair{}, ErrCreateAccess
	}
	refreshClaims := s.jwt.CreateRefreshClaims(base)
	refreshToken, err := s.jwt.CreateRefreshToken(refreshClaims)
	if err != nil {
		return auth.TokenPair{}, ErrCreateRefresh
	}

	pair := auth.TokenPair{
		Access:          accessToken,
		Refresh:         refreshToken,
		AccessExpiresAt: accessClaims.ExpiresAt.Unix() * 1000,
		RefreshMaxAge:   refreshClaims.ExpiresAt.Unix() - time.Now().Unix(),
	}

	if !useMultipoint {
		return pair, nil
	}

	key := sub.UUID.String()
	if jwtStr, err := s.sessions.Get(ctx, key); errors.Is(err, auth.ErrSessionNotFound) {
		if err := s.sessions.Set(ctx, key, refreshToken, s.refreshTTL); err != nil {
			return auth.TokenPair{}, ErrSetLoginStatus
		}
	} else if err != nil {
		return auth.TokenPair{}, ErrSetLoginStatus
	} else {
		if err := s.blacklist.Add(ctx, jwtStr); err != nil {
			return auth.TokenPair{}, ErrInvalidateJWT
		}
		if err := s.sessions.Set(ctx, key, refreshToken, s.refreshTTL); err != nil {
			return auth.TokenPair{}, ErrSetLoginStatus
		}
	}
	return pair, nil
}

// ParseAccess 解析 Access Token。
func (s *AuthService) ParseAccess(ctx context.Context, token string) (auth.AccessClaims, error) {
	claims, err := s.jwt.ParseAccessToken(token)
	if err != nil {
		return auth.AccessClaims{}, err
	}
	return toDomainAccessClaims(claims), nil
}

// ParseRefresh 解析 Refresh Token。
func (s *AuthService) ParseRefresh(ctx context.Context, token string) (auth.RefreshClaims, error) {
	claims, err := s.jwt.ParseRefreshToken(token)
	if err != nil {
		return auth.RefreshClaims{}, err
	}
	return auth.RefreshClaims{UserID: claims.UserID, ExpiresAt: claims.ExpiresAt.Unix()}, nil
}

// Blacklist 将令牌加入黑名单。
func (s *AuthService) Blacklist(ctx context.Context, jwtStr string) error {
	return s.blacklist.Add(ctx, jwtStr)
}

// IsBlacklisted 判断令牌是否在黑名单。
func (s *AuthService) IsBlacklisted(ctx context.Context, jwtStr string) bool {
	return s.blacklist.Contains(ctx, jwtStr)
}

// SetSession 保存会话。
func (s *AuthService) SetSession(ctx context.Context, uuidVal, refresh string) error {
	return s.sessions.Set(ctx, uuidVal, refresh, s.refreshTTL)
}

// GetSession 读取会话。
func (s *AuthService) GetSession(ctx context.Context, uuidVal string) (string, error) {
	return s.sessions.Get(ctx, uuidVal)
}

// DelSession 删除会话。
func (s *AuthService) DelSession(ctx context.Context, uuidVal string) error {
	return s.sessions.Del(ctx, uuidVal)
}

// Authenticate 中间件鉴权流程：黑名单检查 → access 解析 → 失败时走 refresh 续期。
func (s *AuthService) Authenticate(ctx context.Context, accessToken, refreshToken string) (*auth.AuthResult, error) {
	if s.blacklist.Contains(ctx, refreshToken) {
		return nil, auth.ErrRefreshBlacklisted
	}

	claims, err := s.jwt.ParseAccessToken(accessToken)
	if err == nil {
		return &auth.AuthResult{Claims: toDomainAccessClaims(claims)}, nil
	}

	if accessToken == "" || errors.Is(err, jwtutil.TokenExpired) {
		refreshClaims, err := s.jwt.ParseRefreshToken(refreshToken)
		if err != nil {
			return nil, auth.ErrRefreshInvalid
		}

		u, err := s.users.FindByID(ctx, refreshClaims.UserID)
		if err != nil {
			return nil, auth.ErrUserNotFound
		}

		base := request.BaseClaims{UserID: u.ID, UUID: u.UUID, RoleID: appTypes.RoleID(u.RoleID)}
		newAccessClaims := s.jwt.CreateAccessClaims(base)
		newAccessToken, err := s.jwt.CreateAccessToken(newAccessClaims)
		if err != nil {
			return nil, auth.ErrTokenCreate
		}

		return &auth.AuthResult{
			Claims:         toDomainAccessClaims(&newAccessClaims),
			NewAccessToken: newAccessToken,
			NewAccessAt:    newAccessClaims.ExpiresAt.Unix(),
		}, nil
	}

	return nil, auth.ErrAccessInvalid
}

// LoadAll 从数据库加载全部黑名单到本地缓存。
func (s *AuthService) LoadAll(ctx context.Context) error {
	return s.blacklist.LoadAll(ctx)
}

// toDomainAccessClaims request 声明 → 领域声明。
func toDomainAccessClaims(c *request.JwtCustomClaims) auth.AccessClaims {
	return auth.AccessClaims{
		UserID:    c.UserID,
		UUID:      c.UUID,
		RoleID:    shared.RoleID(c.RoleID),
		ExpiresAt: c.ExpiresAt.Unix(),
	}
}

// 接口编译期校验。
var _ auth.AuthPort = (*AuthService)(nil)
