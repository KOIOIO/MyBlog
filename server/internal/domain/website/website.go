// Package website 提供网站配置聚合领域模型与端口。
// 消费图片、友链、热搜、日历端口，聚合后供 HTTP 层使用。
package website

import (
	"context"

	"server/internal/domain/friendlink"
)

// HotItem 热搜条目。
type HotItem struct {
	Index       int    `json:"index"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Popularity  string `json:"popularity"`
	URL         string `json:"url"`
}

// HotSearchData 热搜数据。
type HotSearchData struct {
	Source     string    `json:"source"`
	UpdateTime string    `json:"update_time"`
	HotList    []HotItem `json:"hot_list"`
}

// Calendar 农历日历信息。
type Calendar struct {
	Date         string `json:"date"`
	LunarDate    string `json:"lunar_date"`
	Ganzhi       string `json:"ganzhi"`
	Zodiac       string `json:"zodiac"`
	DayOfYear    string `json:"day_of_year"`
	SolarTerm    string `json:"solar_term"`
	Auspicious   string `json:"auspicious"`
	Inauspicious string `json:"inauspicious"`
}

// FooterLink 页脚链接实体（复用友链域实体，JSON 形状一致）。
type FooterLink = friendlink.FriendLink

// WebsiteImagePort 网站图片端口（首页背景/类别维护）。
type WebsiteImagePort interface {
	// CarouselURLs 首页背景图片 URL 列表。
	CarouselURLs(ctx context.Context) ([]string, error)
	// ChangeCategory 修改图片类别（非事务）。
	ChangeCategory(ctx context.Context, urls []string, category string) error
	// InitCategory 重置图片类别（非事务）。
	InitCategory(ctx context.Context, urls []string) error
}

// FooterLinkRepository 页脚链接仓储端口。
type FooterLinkRepository interface {
	List(ctx context.Context) ([]*FooterLink, error)
	Save(ctx context.Context, link *FooterLink) error
	Delete(ctx context.Context, link *FooterLink) error
}

// HotSearchProvider 热搜提供端口（含缓存语义）。
type HotSearchProvider interface {
	GetHotSearchData(ctx context.Context, source string) (HotSearchData, error)
	// WarmAll 预热全部平台热搜（定时任务）。
	WarmAll(ctx context.Context) error
}

// CalendarProvider 日历提供端口（含缓存语义）。
type CalendarProvider interface {
	GetCalendarByDate(ctx context.Context, dateStr string) (Calendar, error)
}
