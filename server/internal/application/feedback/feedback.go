// Package feedback 提供反馈应用服务。
package feedback

import (
	"context"

	"server/internal/domain/feedback"

	"github.com/gofrs/uuid"
)

// Service 反馈应用服务。
type Service struct {
	fb feedback.FeedbackRepository
}

// NewService 构造反馈应用服务。
func NewService(fb feedback.FeedbackRepository) *Service {
	return &Service{fb: fb}
}

// Newest 最新反馈。
func (s *Service) Newest(ctx context.Context) ([]*feedback.Feedback, error) {
	return s.fb.Newest(ctx, 5)
}

// Create 创建反馈。
func (s *Service) Create(ctx context.Context, f *feedback.Feedback) error {
	return s.fb.Create(ctx, f)
}

// Info 用户反馈。
func (s *Service) Info(ctx context.Context, userUUID uuid.UUID) ([]*feedback.Feedback, error) {
	return s.fb.Info(ctx, userUUID)
}

// Delete 删除反馈。
func (s *Service) Delete(ctx context.Context, ids []uint) error {
	return s.fb.Delete(ctx, ids)
}

// Reply 回复反馈。
func (s *Service) Reply(ctx context.Context, id uint, reply string) error {
	return s.fb.Reply(ctx, id, reply)
}

// List 反馈列表。
func (s *Service) List(ctx context.Context, page, pageSize int) ([]*feedback.Feedback, int64, error) {
	return s.fb.List(ctx, page, pageSize)
}
