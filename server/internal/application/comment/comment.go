// Package comment 提供评论应用服务（用例编排，依赖领域端口）。
package comment

import (
	"context"

	"server/internal/domain/comment"
	"server/internal/domain/shared"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

// Service 评论应用服务。
type Service struct {
	comments comment.CommentRepository
	log      *zap.Logger
}

// NewService 构造评论应用服务。
func NewService(comments comment.CommentRepository, log *zap.Logger) *Service {
	return &Service{comments: comments, log: log}
}

// InfoByArticleID 根据文章 ID 获取评论信息。
func (s *Service) InfoByArticleID(ctx context.Context, articleID string) ([]*comment.Comment, error) {
	return s.comments.ByArticle(ctx, articleID)
}

// Newest 获取最新评论（5 条）。
func (s *Service) Newest(ctx context.Context) ([]*comment.Comment, error) {
	return s.comments.Newest(ctx, 5)
}

// Create 创建评论。
func (s *Service) Create(ctx context.Context, c *comment.Comment) error {
	return s.comments.Create(ctx, c)
}

// Delete 删除评论（含权限校验与级联删除）。
func (s *Service) Delete(ctx context.Context, ids []uint, userUUID uuid.UUID, roleID shared.RoleID) error {
	return s.comments.DeleteTree(ctx, ids, userUUID, roleID)
}

// Info 获取用户评论（含去重规则）。
func (s *Service) Info(ctx context.Context, userUUID uuid.UUID) ([]*comment.Comment, error) {
	raw, err := s.comments.ByUser(ctx, userUUID)
	if err != nil {
		return nil, err
	}
	return comment.DedupByRootUser(raw), nil
}

// List 获取评论管理列表。
func (s *Service) List(ctx context.Context, cond comment.ListCond) ([]*comment.Comment, int64, error) {
	return s.comments.List(ctx, cond)
}
