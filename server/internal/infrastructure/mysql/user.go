// Package mysql 提供 MySQL/GORM 仓储实现（构造注入，不依赖 global）。
package mysql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"server/global"
	"server/internal/common/errs"
	"server/internal/common/page"
	"server/internal/domain/shared"
	"server/internal/domain/user"
	"server/model/appTypes"
	"server/model/database"
	"server/model/other"
	"server/model/request"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// UserRepository GORM 用户仓储实现。
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 构造用户仓储。
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// toUserModel 领域实体 → GORM 模型。
func toUserModel(u *user.User) *database.User {
	return &database.User{
		MODEL:     modelOf(u.ID, u.CreatedAt, u.UpdatedAt),
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

// fromUserModel GORM 模型 → 领域实体。
func fromUserModel(m *database.User) *user.User {
	return &user.User{
		ID:        m.ID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		UUID:      m.UUID,
		Username:  m.Username,
		Password:  m.Password,
		Email:     m.Email,
		Openid:    m.Openid,
		Avatar:    m.Avatar,
		Address:   m.Address,
		Signature: m.Signature,
		RoleID:    shared.RoleID(m.RoleID),
		Register:  shared.Register(m.Register),
		Freeze:    m.Freeze,
	}
}

func fromUserModels(ms []database.User) []*user.User {
	out := make([]*user.User, 0, len(ms))
	for i := range ms {
		out = append(out, fromUserModel(&ms[i]))
	}
	return out
}

// FindByEmail 按邮箱查询用户；不存在返回 errs.ErrNotFound。
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var m database.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}
	return fromUserModel(&m), nil
}

// FindByUUID 按 UUID 查询用户（仅取卡片展示字段，与原 Select 一致）。
func (r *UserRepository) FindByUUID(ctx context.Context, uuidVal uuid.UUID) (*user.User, error) {
	var m database.User
	if err := r.db.WithContext(ctx).Where("uuid = ?", uuidVal).Select("uuid", "username", "avatar", "address", "signature").First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}
	return fromUserModel(&m), nil
}

// FindByID 按主键查询用户。
func (r *UserRepository) FindByID(ctx context.Context, id uint) (*user.User, error) {
	var m database.User
	if err := r.db.WithContext(ctx).Take(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}
	return fromUserModel(&m), nil
}

// Create 创建用户。
func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	return r.db.WithContext(ctx).Create(toUserModel(u)).Error
}

// Update 全量更新用户（对应旧 Save 语义）。
func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	return r.db.WithContext(ctx).Save(toUserModel(u)).Error
}

// UpdateFreeze 更新冻结状态并返回更新前用户；用户不存在时返回错误（保持旧链式 Take+Update 语义）。
func (r *UserRepository) UpdateFreeze(ctx context.Context, id uint, frozen bool) (*user.User, error) {
	var m database.User
	if err := r.db.WithContext(ctx).Take(&m, id).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Model(&m).Update("freeze", frozen).Error; err != nil {
		return nil, err
	}
	return fromUserModel(&m), nil
}

// UpdateInfo 更新个人信息（对应旧 Updates(struct) 非零字段更新语义）。
func (r *UserRepository) UpdateInfo(ctx context.Context, id uint, username, address, signature string) error {
	var m database.User
	if err := r.db.WithContext(ctx).Take(&m, id).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).Model(&m).Updates(database.User{Username: username, Address: address, Signature: signature}).Error
}

// UpdateAvatar 更新头像并返回更新前用户（供旧头像清理）。
func (r *UserRepository) UpdateAvatar(ctx context.Context, id uint, avatarURL string) (*user.User, error) {
	var m database.User
	if err := r.db.WithContext(ctx).Take(&m, id).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Model(&m).Update("avatar", avatarURL).Error; err != nil {
		return nil, err
	}
	return fromUserModel(&m), nil
}

// CountByDate 按注册日期统计用户数（对应 utils.FetchDateCounts 语义）。
func (r *UserRepository) CountByDate(ctx context.Context, days int) (map[string]int, error) {
	return fetchDateCounts(r.db.WithContext(ctx).Model(&database.User{}), days)
}

// Page 用户分页列表。
func (r *UserRepository) Page(ctx context.Context, cond user.ListCond) ([]*user.User, int64, error) {
	db := r.db.WithContext(ctx)
	if cond.UUID != nil {
		db = db.Where("uuid = ?", *cond.UUID)
	}
	if cond.Register != nil {
		db = db.Where("register = ?", *cond.Register)
	}

	option := other.MySQLOption{
		PageInfo: pageInfoOf(cond.Page, cond.PageSize),
		Where:    db,
	}
	list, total, err := page.MySQLPagination(r.db.WithContext(ctx), &database.User{}, option)
	if err != nil {
		return nil, 0, err
	}
	return fromUserModels(list), total, nil
}

// FindIDByUUID 按 UUID 查询用户 ID（保留旧 Pluck 语义）。
func (r *UserRepository) FindIDByUUID(ctx context.Context, uuidVal uuid.UUID) (uint, error) {
	var userID uint
	err := r.db.WithContext(ctx).Model(&database.User{}).Where("uuid = ?", uuidVal).Pluck("id", &userID).Error
	if err != nil {
		return 0, err
	}
	return userID, nil
}

// ---------------------------------------------------------------------------
// 登录日志仓储
// ---------------------------------------------------------------------------

// LoginRecordRepository GORM 登录日志仓储实现。
type LoginRecordRepository struct {
	db *gorm.DB
}

// NewLoginRecordRepository 构造登录日志仓储。
func NewLoginRecordRepository(db *gorm.DB) *LoginRecordRepository {
	return &LoginRecordRepository{db: db}
}

// toLoginModel 领域实体 → GORM 模型。
func toLoginModel(rec *user.LoginRecord) *database.Login {
	return &database.Login{
		MODEL:       modelOf(rec.ID, rec.CreatedAt, time.Time{}),
		UserID:      rec.UserID,
		LoginMethod: rec.LoginMethod,
		IP:          rec.IP,
		Address:     rec.Address,
		OS:          rec.OS,
		DeviceInfo:  rec.DeviceInfo,
		BrowserInfo: rec.BrowserInfo,
		Status:      rec.Status,
	}
}

func fromLoginModels(ms []database.Login) []*user.LoginRecord {
	out := make([]*user.LoginRecord, 0, len(ms))
	for i := range ms {
		rec := &user.LoginRecord{
			ID:          ms[i].ID,
			CreatedAt:   ms[i].CreatedAt,
			UserID:      ms[i].UserID,
			LoginMethod: ms[i].LoginMethod,
			IP:          ms[i].IP,
			Address:     ms[i].Address,
			OS:          ms[i].OS,
			DeviceInfo:  ms[i].DeviceInfo,
			BrowserInfo: ms[i].BrowserInfo,
			Status:      ms[i].Status,
		}
		if ms[i].User.ID != 0 {
			rec.User = fromUserModel(&ms[i].User)
		}
		out = append(out, rec)
	}
	return out
}

// Create 创建登录记录。
func (r *LoginRecordRepository) Create(ctx context.Context, rec *user.LoginRecord) error {
	return r.db.WithContext(ctx).Create(toLoginModel(rec)).Error
}

// CountByDate 按登录日期统计（对应 utils.FetchDateCounts 语义）。
func (r *LoginRecordRepository) CountByDate(ctx context.Context, days int) (map[string]int, error) {
	return fetchDateCounts(r.db.WithContext(ctx).Model(&database.Login{}), days)
}

// Page 登录日志分页（预加载用户，保持旧行为：uuid 筛选失败时返回空列表）。
func (r *LoginRecordRepository) Page(ctx context.Context, cond user.LoginListCond) ([]*user.LoginRecord, int64, error) {
	db := r.db.WithContext(ctx)
	if cond.UUID != nil {
		var userID uint
		if err := db.Model(&database.User{}).Where("uuid = ?", *cond.UUID).Pluck("id", &userID).Error; err != nil {
			return nil, 0, nil
		}
		db = db.Where("user_id = ?", userID)
	}

	option := other.MySQLOption{
		PageInfo: pageInfoOf(cond.Page, cond.PageSize),
		Where:    db,
		Preload:  []string{"User"},
	}
	list, total, err := page.MySQLPagination(r.db.WithContext(ctx), &database.Login{}, option)
	if err != nil {
		return nil, 0, err
	}
	return fromLoginModels(list), total, nil
}

// ---------------------------------------------------------------------------
// 包内共享小工具
// ---------------------------------------------------------------------------

// modelOf 由领域字段构造 GORM 基础模型。
func modelOf(id uint, createdAt, updatedAt time.Time) global.MODEL {
	return global.MODEL{ID: id, CreatedAt: createdAt, UpdatedAt: updatedAt}
}

// pageInfoOf 构造分页参数。
func pageInfoOf(pageNum, pageSize int) request.PageInfo {
	return request.PageInfo{Page: pageNum, PageSize: pageSize}
}

// fetchDateCounts 按日期分组统计数量（SQL 与 utils.FetchDateCounts 一致）。
func fetchDateCounts(db *gorm.DB, days int) (map[string]int, error) {
	var dateCounts []struct {
		Date  string `json:"date"`
		Count int    `json:"count"`
	}
	db.Where(fmt.Sprintf("date_sub(curdate(), interval %d day) <= created_at", days)).
		Select("date_format(created_at, '%Y-%m-%d') as date", "count(id) as count").
		Group("date").Scan(&dateCounts)

	dateCountMap := make(map[string]int)
	for _, c := range dateCounts {
		dateCountMap[c.Date] = c.Count
	}
	return dateCountMap, nil
}
