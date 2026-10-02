// Package website 提供网站应用服务（配置聚合）。
package website

import (
	"context"

	"server/config"
	"server/internal/domain/website"
)

// Service 网站应用服务。
type Service struct {
	images  website.WebsiteImagePort
	footers website.FooterLinkRepository
	news    website.HotSearchProvider
	cal     website.CalendarProvider
	cfg     *config.Config
}

// NewService 构造网站应用服务。
func NewService(
	images website.WebsiteImagePort,
	footers website.FooterLinkRepository,
	news website.HotSearchProvider,
	cal website.CalendarProvider,
	cfg *config.Config,
) *Service {
	return &Service{images: images, footers: footers, news: news, cal: cal, cfg: cfg}
}

// Carousel 首页背景图片列表。
func (s *Service) Carousel(ctx context.Context) ([]string, error) {
	return s.images.CarouselURLs(ctx)
}

// News 获取热搜数据。
func (s *Service) News(ctx context.Context, sourceStr string) (website.HotSearchData, error) {
	return s.news.GetHotSearchData(ctx, sourceStr)
}

// WarmHotSearch 预热全部平台热搜（定时任务）。
func (s *Service) WarmHotSearch(ctx context.Context) error {
	return s.news.WarmAll(ctx)
}

// Calendar 获取日历信息。
func (s *Service) Calendar(ctx context.Context, dateStr string) (website.Calendar, error) {
	return s.cal.GetCalendarByDate(ctx, dateStr)
}

// FooterLink 页脚链接全量。
func (s *Service) FooterLink(ctx context.Context) ([]*website.FooterLink, error) {
	return s.footers.List(ctx)
}

// AddCarousel 添加首页背景。
func (s *Service) AddCarousel(ctx context.Context, url string) error {
	return s.images.ChangeCategory(ctx, []string{url}, "背景")
}

// CancelCarousel 移除首页背景。
func (s *Service) CancelCarousel(ctx context.Context, url string) error {
	return s.images.InitCategory(ctx, []string{url})
}

// CreateFooterLink 创建页脚链接。
func (s *Service) CreateFooterLink(ctx context.Context, link *website.FooterLink) error {
	return s.footers.Save(ctx, link)
}

// DeleteFooterLink 删除页脚链接。
func (s *Service) DeleteFooterLink(ctx context.Context, link *website.FooterLink) error {
	return s.footers.Delete(ctx, link)
}
