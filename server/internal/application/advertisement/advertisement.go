// Package advertisement 提供广告应用服务。
package advertisement

import (
	"context"

	"server/internal/domain/advertisement"
)

// Service 广告应用服务。
type Service struct {
	ads advertisement.AdvertisementRepository
}

// NewService 构造广告应用服务。
func NewService(ads advertisement.AdvertisementRepository) *Service {
	return &Service{ads: ads}
}

// Info 广告信息。
func (s *Service) Info(ctx context.Context) ([]*advertisement.Advertisement, int64, error) {
	return s.ads.Info(ctx)
}

// Create 创建广告。
func (s *Service) Create(ctx context.Context, ad *advertisement.Advertisement) error {
	return s.ads.Create(ctx, ad)
}

// Delete 删除广告。
func (s *Service) Delete(ctx context.Context, ids []uint) error {
	return s.ads.Delete(ctx, ids)
}

// Update 更新广告。
func (s *Service) Update(ctx context.Context, ad *advertisement.Advertisement) error {
	return s.ads.Update(ctx, ad)
}

// List 广告列表。
func (s *Service) List(ctx context.Context, cond advertisement.ListCond) ([]*advertisement.Advertisement, int64, error) {
	return s.ads.List(ctx, cond)
}
