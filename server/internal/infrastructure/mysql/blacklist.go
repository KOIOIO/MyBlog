// Package mysql 提供 MySQL/GORM 仓储实现（构造注入，不依赖 global）。
package mysql

import (
	"context"

	"server/internal/domain/auth"
	"server/internal/model/database"

	"github.com/songzhibin97/gkit/cache/local_cache"
	"gorm.io/gorm"
)

// BlacklistStore JWT 黑名单仓储（MySQL 持久化 + 本地缓存）。
type BlacklistStore struct {
	db    *gorm.DB
	cache local_cache.Cache
}

// NewBlacklistStore 构造黑名单仓储。
func NewBlacklistStore(db *gorm.DB, cache local_cache.Cache) *BlacklistStore {
	return &BlacklistStore{db: db, cache: cache}
}

// Add 写入黑名单（DB 持久化 + 本地缓存）。
func (s *BlacklistStore) Add(ctx context.Context, jwtStr string) error {
	if err := s.db.WithContext(ctx).Create(&database.JwtBlacklist{Jwt: jwtStr}).Error; err != nil {
		return err
	}
	s.cache.SetDefault(jwtStr, struct{}{})
	return nil
}

// Contains 判断是否在黑名单（本地缓存查询）。
func (s *BlacklistStore) Contains(ctx context.Context, jwtStr string) bool {
	_, ok := s.cache.Get(jwtStr)
	return ok
}

// LoadAll 从数据库加载全部黑名单 JWT 到本地缓存。
func (s *BlacklistStore) LoadAll(ctx context.Context) error {
	var data []string
	if err := s.db.WithContext(ctx).Model(&database.JwtBlacklist{}).Pluck("jwt", &data).Error; err != nil {
		return err
	}
	for i := range data {
		s.cache.SetDefault(data[i], struct{}{})
	}
	return nil
}

// 接口编译期校验。
var _ auth.BlacklistStore = (*BlacklistStore)(nil)
