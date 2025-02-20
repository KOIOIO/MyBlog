package service

import (
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"server/global"
	"server/model/database"
	"server/utills"
)

type JWTService struct {
}

func (jwtService *JWTService) SetRedisJWT(jwt string, uuid uuid.UUID) error {
	dr, err := utills.ParseDuration(global.Config.Jwt.RefreshTokenExpiryTime)
	if err != nil {
		return err
	}
	return global.Redis.Set(uuid.String(), jwt, dr).Err()
}

func (jwtService *JWTService) GetRedisJWT(uuid uuid.UUID) (string, error) {
	return global.Redis.Get(uuid.String()).Result()
}

func (jwtService *JWTService) JoinInBlacklist(jwtList database.JwtBlacklist) error {
	if err := global.DB.Create(&jwtList).Error; err != nil {
		return err
	}
	global.BlackCache.SetDefault(jwtList.Jwt, struct{}{})
	return nil
}

func (jwtService *JWTService) IsInBlacklist(jwt string) bool {
	_, isIn := global.BlackCache.Get(jwt)
	return isIn
}

func LoadAll() {
	var data []string
	if err := global.DB.Model(&database.JwtBlacklist{}).Pluck("jwt", &data).Error; err != nil {
		global.Log.Error("Failed to load jet black list from the database", zap.Error(err))
		return
	}
	for i := 0; i < len(data); i++ {
		global.BlackCache.SetDefault(data[i], struct{}{})
	}

}
