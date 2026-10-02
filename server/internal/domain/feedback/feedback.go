// Package feedback 提供反馈领域模型与仓储端口。
package feedback

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

// Feedback 反馈实体（JSON 与 database.Feedback 一致）。
type Feedback struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserUUID  uuid.UUID `json:"user_uuid"`
	Content   string    `json:"content"`
	Reply     string    `json:"reply"`
}

// FeedbackRepository 反馈仓储端口。
type FeedbackRepository interface {
	// Newest 最新 N 条反馈。
	Newest(ctx context.Context, limit int) ([]*Feedback, error)
	// Create 创建反馈。
	Create(ctx context.Context, f *Feedback) error
	// Info 用户反馈（按时间倒序）。
	Info(ctx context.Context, userUUID uuid.UUID) ([]*Feedback, error)
	// Delete 批量删除反馈。
	Delete(ctx context.Context, ids []uint) error
	// Reply 回复反馈。
	Reply(ctx context.Context, id uint, reply string) error
	// List 反馈分页列表。
	List(ctx context.Context, page, pageSize int) ([]*Feedback, int64, error)
}
