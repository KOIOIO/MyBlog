// FeedbackRepo 反馈仓储（GORM）。
package mysql

import (
	"context"

	"server/internal/common/page"
	"server/internal/domain/feedback"
	"server/model/database"
	"server/model/other"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// FeedbackRepo 反馈仓储。
type FeedbackRepo struct {
	db *gorm.DB
}

// NewFeedbackRepo 构造反馈仓储。
func NewFeedbackRepo(db *gorm.DB) *FeedbackRepo {
	return &FeedbackRepo{db: db}
}

var _ feedback.FeedbackRepository = (*FeedbackRepo)(nil)

func toDomainFeedback(m *database.Feedback) *feedback.Feedback {
	return &feedback.Feedback{
		ID:        m.ID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		UserUUID:  m.UserUUID,
		Content:   m.Content,
		Reply:     m.Reply,
	}
}

// Newest 最新 N 条反馈。
func (r *FeedbackRepo) Newest(ctx context.Context, limit int) ([]*feedback.Feedback, error) {
	var list []database.Feedback
	if err := r.db.WithContext(ctx).Order("id desc").Limit(limit).Find(&list).Error; err != nil {
		return nil, err
	}
	out := make([]*feedback.Feedback, 0, len(list))
	for i := range list {
		out = append(out, toDomainFeedback(&list[i]))
	}
	return out, nil
}

// Create 创建反馈。
func (r *FeedbackRepo) Create(ctx context.Context, f *feedback.Feedback) error {
	return r.db.WithContext(ctx).Create(&database.Feedback{
		UserUUID: f.UserUUID,
		Content:  f.Content,
	}).Error
}

// Info 用户反馈（按时间倒序）。
func (r *FeedbackRepo) Info(ctx context.Context, userUUID uuid.UUID) ([]*feedback.Feedback, error) {
	var list []database.Feedback
	if err := r.db.WithContext(ctx).Model(&database.Feedback{}).
		Order("id desc").Where("user_uuid = ?", userUUID).Find(&list).Error; err != nil {
		return nil, err
	}
	out := make([]*feedback.Feedback, 0, len(list))
	for i := range list {
		out = append(out, toDomainFeedback(&list[i]))
	}
	return out, nil
}

// Delete 批量删除反馈。
func (r *FeedbackRepo) Delete(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Delete(&database.Feedback{}, ids).Error
}

// Reply 回复反馈。
func (r *FeedbackRepo) Reply(ctx context.Context, id uint, reply string) error {
	return r.db.WithContext(ctx).Take(&database.Feedback{}, id).Update("reply", reply).Error
}

// List 反馈分页列表。
func (r *FeedbackRepo) List(ctx context.Context, pageNum, pageSize int) ([]*feedback.Feedback, int64, error) {
	option := other.MySQLOption{
		PageInfo: pageInfoOf(pageNum, pageSize),
	}
	list, total, err := page.MySQLPagination(r.db.WithContext(ctx), &database.Feedback{}, option)
	if err != nil {
		return nil, 0, err
	}
	out := make([]*feedback.Feedback, 0, len(list))
	for i := range list {
		out = append(out, toDomainFeedback(&list[i]))
	}
	return out, total, nil
}
