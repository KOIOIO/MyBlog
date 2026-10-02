// Package advertisement 提供广告领域模型与仓储端口。
package advertisement

import (
	"context"
	"time"
)

// Advertisement 广告实体（JSON 与 database.Advertisement 一致）。
type Advertisement struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	AdImage   string    `json:"ad_image"`
	Link      string    `json:"link"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
}

// ListCond 广告列表查询条件。
type ListCond struct {
	Title    *string
	Content  *string
	Page     int
	PageSize int
}

// AdvertisementRepository 广告仓储端口。
type AdvertisementRepository interface {
	// Info 广告信息（全量 + 总数）。
	Info(ctx context.Context) ([]*Advertisement, int64, error)
	// Create 创建广告（事务：图片类别改为"广告" + 创建）。
	Create(ctx context.Context, ad *Advertisement) error
	// Delete 删除广告（事务：图片类别重置 + 删除）。
	Delete(ctx context.Context, ids []uint) error
	// Update 更新广告（仅 link/title/content）。
	Update(ctx context.Context, ad *Advertisement) error
	// List 广告分页列表。
	List(ctx context.Context, cond ListCond) ([]*Advertisement, int64, error)
}
