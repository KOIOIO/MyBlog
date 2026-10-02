// Package jwt 提供 JWT 签发与解析能力（构造注入，不依赖 global）。
package jwt

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"server/internal/model/request"

	"github.com/golang-jwt/jwt/v4"
)

// 通用 Token 校验错误。
var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenInvalid     = errors.New("couldn't handle this token")
)

// JWT 持有签名密钥与过期配置。
type JWT struct {
	AccessTokenSecret  []byte
	RefreshTokenSecret []byte
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	Issuer             string
}

// New 根据 JWT 配置构造实例；过期时间解析失败时返回错误。
func New(accessSecret, refreshSecret, accessExpiry, refreshExpiry, issuer string) (*JWT, error) {
	access, err := ParseDuration(accessExpiry)
	if err != nil {
		return nil, fmt.Errorf("parse access token expiry: %w", err)
	}
	refresh, err := ParseDuration(refreshExpiry)
	if err != nil {
		return nil, fmt.Errorf("parse refresh token expiry: %w", err)
	}
	return &JWT{
		AccessTokenSecret:  []byte(accessSecret),
		RefreshTokenSecret: []byte(refreshSecret),
		AccessTokenExpiry:  access,
		RefreshTokenExpiry: refresh,
		Issuer:             issuer,
	}, nil
}

// CreateAccessClaims 创建 Access Token 的 Claims，包含基本信息和过期时间等。
func (j *JWT) CreateAccessClaims(baseClaims request.BaseClaims) request.JwtCustomClaims {
	claims := request.JwtCustomClaims{
		BaseClaims: baseClaims,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"TAP"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.AccessTokenExpiry)),
			Issuer:    j.Issuer,
		},
	}
	return claims
}

// CreateAccessToken 通过 Claims 生成 Access Token。
func (j *JWT) CreateAccessToken(claims request.JwtCustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.AccessTokenSecret)
}

// CreateRefreshClaims 创建 Refresh Token 的 Claims。
func (j *JWT) CreateRefreshClaims(baseClaims request.BaseClaims) request.JwtCustomRefreshClaims {
	claims := request.JwtCustomRefreshClaims{
		UserID: baseClaims.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"TAP"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.RefreshTokenExpiry)),
			Issuer:    j.Issuer,
		},
	}
	return claims
}

// CreateRefreshToken 通过 Claims 生成 Refresh Token。
func (j *JWT) CreateRefreshToken(claims request.JwtCustomRefreshClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.RefreshTokenSecret)
}

// ParseAccessToken 解析 Access Token，验证签名并返回 Claims。
func (j *JWT) ParseAccessToken(tokenString string) (*request.JwtCustomClaims, error) {
	claims, err := j.parseToken(tokenString, &request.JwtCustomClaims{}, j.AccessTokenSecret)
	if err != nil {
		return nil, err
	}
	if customClaims, ok := claims.(*request.JwtCustomClaims); ok {
		return customClaims, nil
	}
	return nil, TokenInvalid
}

// ParseRefreshToken 解析 Refresh Token，验证签名并返回 Claims。
func (j *JWT) ParseRefreshToken(tokenString string) (*request.JwtCustomRefreshClaims, error) {
	claims, err := j.parseToken(tokenString, &request.JwtCustomRefreshClaims{}, j.RefreshTokenSecret)
	if err != nil {
		return nil, err
	}
	if refreshClaims, ok := claims.(*request.JwtCustomRefreshClaims); ok {
		return refreshClaims, nil
	}
	return nil, TokenInvalid
}

// parseToken 通用的 Token 解析方法。
func (j *JWT) parseToken(tokenString string, claims jwt.Claims, secretKey interface{}) (interface{}, error) {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			switch {
			case ve.Errors&jwt.ValidationErrorMalformed != 0:
				return nil, TokenMalformed
			case ve.Errors&jwt.ValidationErrorExpired != 0:
				return nil, TokenExpired
			case ve.Errors&jwt.ValidationErrorNotValidYet != 0:
				return nil, TokenNotValidYet
			default:
				return nil, TokenInvalid
			}
		}
		return nil, TokenInvalid
	}
	if token.Valid {
		return token.Claims, nil
	}
	return nil, TokenInvalid
}

// ParseDuration 解析持续时间字符串为 time.Duration，支持 d/h/m/s 单位，例如 "1d2h30m"。
func ParseDuration(d string) (time.Duration, error) {
	d = strings.TrimSpace(d)
	if len(d) == 0 {
		return 0, fmt.Errorf("empty duration string")
	}

	unitPattern := map[string]time.Duration{
		"d": time.Hour * 24,
		"h": time.Hour,
		"m": time.Minute,
		"s": time.Second,
	}

	var totalDuration time.Duration
	for _, unit := range []string{"d", "h", "m", "s"} {
		for strings.Contains(d, unit) {
			unitIndex := strings.Index(d, unit)
			part := d[:unitIndex]
			if part == "" {
				part = "0"
			}
			val, err := strconv.Atoi(part)
			if err != nil {
				return 0, fmt.Errorf("invalid duration part: %v", err)
			}
			totalDuration += time.Duration(val) * unitPattern[unit]
			d = d[unitIndex+len(unit):]
		}
	}
	if len(d) > 0 {
		return 0, fmt.Errorf("unrecognized duration format")
	}
	return totalDuration, nil
}
