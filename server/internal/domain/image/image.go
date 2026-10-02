// Package image 提供图片领域模型与仓储端口。
package image

import (
	"context"
	"time"
)

// Image 图片实体（JSON 字段与 database.Image 一致，category/storage 为中文串）。
type Image struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Category  string    `json:"category"`
	Storage   string    `json:"storage"`
}

// ListCond 图片列表查询条件。
type ListCond struct {
	Name     *string
	Category *string
	Storage  *string
	Page     int
	PageSize int
}

// ImageRepository 图片仓储端口。
type ImageRepository interface {
	// Create 记录上传图片（类别初始化为"未使用"，存储按配置）。
	Create(ctx context.Context, name, url, storage string) error
	// Delete 按 ID 查找并删除图片记录，返回被删记录（供存储侧删除文件）。
	Delete(ctx context.Context, ids []uint) ([]Image, error)
	// List 图片分页列表（name/category/storage 过滤）。
	List(ctx context.Context, cond ListCond) ([]Image, int64, error)
}
