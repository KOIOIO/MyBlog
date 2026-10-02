package database

// JwtBlacklist JWT 黑名单表
type JwtBlacklist struct {
	MODEL
	Jwt string `json:"jwt" gorm:"type:text"` // Jwt
}
