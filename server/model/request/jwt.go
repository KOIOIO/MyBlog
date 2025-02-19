package request

import (
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	"server/model/appTypes"
)

type JwtCustomClaims struct {
	BaseClaims
	jwt.RegisteredClaims
}

type JwtCustomRefreshClaims struct {
	UserID uint
	jwt.RegisteredClaims
}

type BaseClaims struct {
	UserID uint
	UUID   uuid.UUID
	RoleID appTypes.RoleID
}
