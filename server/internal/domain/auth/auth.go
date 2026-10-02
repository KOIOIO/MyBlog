// Package auth 提供认证与 JWT 领域模型（实体、Port、领域错误）。
// 本包不依赖任何基础设施库，仅依赖标准库与 domain/shared。
package auth

import (
	"context"
	"errors"
	"time"

	"server/internal/domain/shared"

	"github.com/gofrs/uuid"
)

// JwtBlacklist JWT 黑名单实体。
type JwtBlacklist struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	Jwt       string
}

// TokenSubject 签发 Token 所需的用户主体信息。
type TokenSubject struct {
	UserID uint
	UUID   uuid.UUID
	RoleID shared.RoleID
}

// TokenPair 签发的令牌对。
type TokenPair struct {
	Access          string
	Refresh         string
	AccessExpiresAt int64 // 毫秒（响应 access_token_expires_at 使用）
	RefreshMaxAge   int64 // 秒（refresh token cookie maxAge 使用）
}

// AccessClaims Access Token 携带的声明。
type AccessClaims struct {
	UserID    uint
	UUID      uuid.UUID
	RoleID    shared.RoleID
	ExpiresAt int64 // Unix 秒
}

// RefreshClaims Refresh Token 携带的声明。
type RefreshClaims struct {
	UserID    uint
	ExpiresAt int64 // Unix 秒
}

// AuthResult 中间件鉴权结果；NewAccessToken 非空表示完成了一次自动续期。
type AuthResult struct {
	Claims         AccessClaims
	NewAccessToken string
	NewAccessAt    int64 // Unix 秒
}

// 鉴权流程领域错误（中间件据此映射 HTTP 错误文案，文案必须与重构前一致）。
var (
	ErrRefreshBlacklisted = errors.New("Account logged in from another location or token is invalid")
	ErrRefreshInvalid     = errors.New("Refresh token expired or invalid")
	ErrUserNotFound       = errors.New("The user does not exist")
	ErrTokenCreate        = errors.New("Faild to create new access token")
	ErrAccessInvalid      = errors.New("Invalid access token")
	// ErrSessionNotFound Redis 会话不存在（登录流程区分“首次登录/已登录”）。
	ErrSessionNotFound = errors.New("session not found")
)

// SessionStore Redis 会话存储（uuid → refresh token）。
type SessionStore interface {
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
}

// BlacklistStore JWT 黑名单存储（MySQL 持久化 + 本地缓存加速）。
type BlacklistStore interface {
	Add(ctx context.Context, jwt string) error
	Contains(ctx context.Context, jwt string) bool
	LoadAll(ctx context.Context) error
}

// AuthPort 认证领域端口（由 application/auth 实现）。
type AuthPort interface {
	GenerateToken(ctx context.Context, sub TokenSubject, useMultipoint bool) (TokenPair, error)
	ParseAccess(ctx context.Context, token string) (AccessClaims, error)
	ParseRefresh(ctx context.Context, token string) (RefreshClaims, error)
	Blacklist(ctx context.Context, jwt string) error
	IsBlacklisted(ctx context.Context, jwt string) bool
	SetSession(ctx context.Context, uuid string, refresh string) error
	GetSession(ctx context.Context, uuid string) (string, error)
	DelSession(ctx context.Context, uuid string) error
	Authenticate(ctx context.Context, accessToken, refreshToken string) (*AuthResult, error)
	LoadAll(ctx context.Context) error
}
