// WebsiteRepo 网站仓储：图片端口（首页背景/类别）+ 页脚链接仓储。
package mysql

import (
	"context"

	"server/global"
	"server/internal/domain/website"
	"server/model/appTypes"
	"server/model/database"

	"gorm.io/gorm"
)

// WebsiteRepo 网站仓储。
type WebsiteRepo struct {
	db *gorm.DB
}

// NewWebsiteRepo 构造网站仓储。
func NewWebsiteRepo(db *gorm.DB) *WebsiteRepo {
	return &WebsiteRepo{db: db}
}

var (
	_ website.WebsiteImagePort     = (*WebsiteRepo)(nil)
	_ website.FooterLinkRepository = (*WebsiteRepo)(nil)
)

// CarouselURLs 首页背景图片 URL 列表。
func (r *WebsiteRepo) CarouselURLs(ctx context.Context) ([]string, error) {
	var urls []string
	err := r.db.WithContext(ctx).Model(&database.Image{}).
		Where("category = ?", appTypes.Carousel).Pluck("url", &urls).Error
	return urls, err
}

// ChangeCategory 修改图片类别（非事务）。
func (r *WebsiteRepo) ChangeCategory(ctx context.Context, urls []string, category string) error {
	return ChangeImagesCategory(r.db.WithContext(ctx), urls, appTypes.ToCategory(category))
}

// InitCategory 重置图片类别（非事务）。
func (r *WebsiteRepo) InitCategory(ctx context.Context, urls []string) error {
	return InitImagesCategory(r.db.WithContext(ctx), urls)
}

// List 页脚链接全量。
func (r *WebsiteRepo) List(ctx context.Context) ([]*website.FooterLink, error) {
	var links []database.FriendLink
	if err := r.db.WithContext(ctx).Find(&links).Error; err != nil {
		return nil, err
	}
	out := make([]*website.FooterLink, 0, len(links))
	for i := range links {
		out = append(out, toDomainFriendLink(&links[i]))
	}
	return out, nil
}

// Save 保存（创建或更新）页脚链接。
func (r *WebsiteRepo) Save(ctx context.Context, link *website.FooterLink) error {
	return r.db.WithContext(ctx).Save(&database.FriendLink{
		MODEL:       global.MODEL{ID: link.ID},
		Logo:        link.Logo,
		Link:        link.Link,
		Name:        link.Name,
		Description: link.Description,
	}).Error
}

// Delete 删除页脚链接。
func (r *WebsiteRepo) Delete(ctx context.Context, link *website.FooterLink) error {
	return r.db.WithContext(ctx).Delete(&database.FriendLink{
		MODEL:       global.MODEL{ID: link.ID},
		Logo:        link.Logo,
		Link:        link.Link,
		Name:        link.Name,
		Description: link.Description,
	}).Error
}
