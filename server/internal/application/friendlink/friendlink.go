// Package friendlink 提供友链应用服务。
package friendlink

import (
	"context"

	"server/internal/domain/friendlink"
)

// Service 友链应用服务。
type Service struct {
	links friendlink.FriendLinkRepository
}

// NewService 构造友链应用服务。
func NewService(links friendlink.FriendLinkRepository) *Service {
	return &Service{links: links}
}

// Info 友链信息。
func (s *Service) Info(ctx context.Context) ([]*friendlink.FriendLink, int64, error) {
	return s.links.Info(ctx)
}

// Create 创建友链。
func (s *Service) Create(ctx context.Context, link *friendlink.FriendLink) error {
	return s.links.Create(ctx, link)
}

// Delete 删除友链。
func (s *Service) Delete(ctx context.Context, ids []uint) error {
	return s.links.Delete(ctx, ids)
}

// Update 更新友链。
func (s *Service) Update(ctx context.Context, link *friendlink.FriendLink) error {
	return s.links.Update(ctx, link)
}

// List 友链列表。
func (s *Service) List(ctx context.Context, cond friendlink.ListCond) ([]*friendlink.FriendLink, int64, error) {
	return s.links.List(ctx, cond)
}
