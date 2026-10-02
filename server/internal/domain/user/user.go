// Package user 提供用户领域模型（实体、行为与 Port 接口）。
// 本包不依赖任何基础设施库，仅依赖标准库、domain/shared 与 domain/auth。
package user

import (
	"context"
	"time"

	"server/internal/domain/shared"

	"github.com/gofrs/uuid"
)

// User 用户实体（充血模型：冻结、角色校验等规则收敛于此）。
type User struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	UUID      uuid.UUID       `json:"uuid"`
	Username  string          `json:"username"`
	Password  string          `json:"-"`
	Email     string          `json:"email"`
	Openid    string          `json:"openid"`
	Avatar    string          `json:"avatar"`
	Address   string          `json:"address"`
	Signature string          `json:"signature"`
	RoleID    shared.RoleID   `json:"role_id"`
	Register  shared.Register `json:"register"`
	Freeze    bool            `json:"freeze"`
}

// IsFrozen 判断用户是否被冻结。
func (u *User) IsFrozen() bool {
	return u.Freeze
}

// SetFrozen 设置冻结状态。
func (u *User) SetFrozen(frozen bool) {
	u.Freeze = frozen
}

// IsAdmin 是否为管理员。
func (u *User) IsAdmin() bool {
	return u.RoleID == shared.Admin
}

// IsNormalUser 是否为普通用户。
func (u *User) IsNormalUser() bool {
	return u.RoleID == shared.User
}

// SetPasswordHash 设置密码哈希。
func (u *User) SetPasswordHash(hash string) {
	u.Password = hash
}

// ResetRegisterInfo 重置为邮箱注册用户默认信息（注册流程使用）。
func (u *User) ResetRegisterInfo(uuidVal uuid.UUID, hash, avatar string) {
	u.UUID = uuidVal
	u.Password = hash
	u.Avatar = avatar
	u.RoleID = shared.User
	u.Register = shared.Email
}

// LoginRecord 登录日志实体。
type LoginRecord struct {
	ID          uint      `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UserID      uint      `json:"user_id"`
	User        *User     `json:"user"`
	LoginMethod string    `json:"login_method"`
	IP          string    `json:"ip"`
	Address     string    `json:"address"`
	OS          string    `json:"os"`
	DeviceInfo  string    `json:"device_info"`
	BrowserInfo string    `json:"browser_info"`
	Status      int       `json:"status"`
}

// ListCond 用户分页查询条件（对应原 admin 列表筛选）。
type ListCond struct {
	UUID     *string
	Register *string
	Page     int
	PageSize int
}

// LoginListCond 登录日志分页查询条件。
type LoginListCond struct {
	UUID     *string
	Page     int
	PageSize int
}

// UserRepository 用户仓储端口。
type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByUUID(ctx context.Context, uuid uuid.UUID) (*User, error)
	FindByID(ctx context.Context, id uint) (*User, error)
	Create(ctx context.Context, u *User) error
	Update(ctx context.Context, u *User) error
	UpdateFreeze(ctx context.Context, id uint, frozen bool) (*User, error)
	UpdateInfo(ctx context.Context, id uint, username, address, signature string) error
	UpdateAvatar(ctx context.Context, id uint, avatarURL string) (*User, error)
	CountByDate(ctx context.Context, days int) (map[string]int, error)
	Page(ctx context.Context, cond ListCond) ([]*User, int64, error)
	FindIDByUUID(ctx context.Context, uuidVal uuid.UUID) (uint, error)
}

// LoginRecordRepository 登录日志仓储端口。
type LoginRecordRepository interface {
	Create(ctx context.Context, r *LoginRecord) error
	CountByDate(ctx context.Context, days int) (map[string]int, error)
	Page(ctx context.Context, cond LoginListCond) ([]*LoginRecord, int64, error)
}

// Cache 通用键值缓存端口（Redis 实现；用于天气等聚合数据缓存）。
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Del(ctx context.Context, key string) error
}

// GeoInfo IP 定位信息。
type GeoInfo struct {
	Province string
	City     string
	Adcode   string
}

// WeatherInfo 实时天气信息。
type WeatherInfo struct {
	Province      string
	City          string
	Weather       string
	Temperature   string
	WindDirection string
	WindPower     string
	Humidity      string
}

// GeoProvider 高德定位/天气端口。
type GeoProvider interface {
	LocationByIP(ctx context.Context, ip string) (GeoInfo, error)
	WeatherByAdcode(ctx context.Context, adcode string) (WeatherInfo, error)
}
