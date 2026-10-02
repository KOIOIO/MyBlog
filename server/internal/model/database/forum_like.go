package database

import "time"

// ForumLike 论坛点赞表（防重复：post_id + user_id 唯一）
type ForumLike struct {
	ID        uint      `json:"id" gorm:"primarykey"` // 主键 ID
	CreatedAt time.Time `json:"created_at"`           // 创建时间
	UpdatedAt time.Time `json:"updated_at"`           // 更新时间
	PostID    uint      `json:"post_id"`              // 帖子 ID
	UserID    uint      `json:"user_id"`              // 点赞用户 ID
}
