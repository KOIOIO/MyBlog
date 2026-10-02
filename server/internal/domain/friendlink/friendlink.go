// Package friendlink 提供友链领域模型与仓储端口。
package friendlink

import (
	"context"
	"time"
)

// FriendLink 友链实体（JSON 与 database.FriendLink 一致）。
type FriendLink struct {
	ID          uint      `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Logo        string    `json:"logo"`
	Link        string    `json:"link"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

// ListCond 友链列表查询条件。
type ListCond struct {
	Name        *string
	Description *string
	Page        int
	PageSize    int
}

// FriendLinkRepository 友链仓储端口。
type FriendLinkRepository interface {
	// Info 友链信息（全量 + 总数）。
	Info(ctx context.Context) ([]*FriendLink, int64, error)
	// Create 创建友链（事务：图片类别改为"友链" + 创建）。
	Create(ctx context.Context, link *FriendLink) error
	// Delete 删除友链（事务：图片类别重置 + 删除）。
	Delete(ctx context.Context, ids []uint) error
	// Update 更新友链（仅 link/name/description）。
	Update(ctx context.Context, link *FriendLink) error
	// List 友链分页列表。
	List(ctx context.Context, cond ListCond) ([]*FriendLink, int64, error)
}
