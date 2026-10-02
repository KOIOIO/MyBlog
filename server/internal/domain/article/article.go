// Package article 提供文章领域模型（实体、规则、端口）。
// 本包不依赖任何基础设施库，仅依赖标准库。
package article

import (
	"context"
	"errors"
	"regexp"
)

// Article 文章实体（对应 ES 文档结构，JSON 与 model/elasticsearch.Article 等价）。
type Article struct {
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	Cover     string   `json:"cover"`
	Title     string   `json:"title"`
	Keyword   string   `json:"keyword"`
	Category  string   `json:"category"`
	Tags      []string `json:"tags"`
	Abstract  string   `json:"abstract"`
	Content   string   `json:"content"`
	Views     int      `json:"views"`
	Comments  int      `json:"comments"`
	Likes     int      `json:"likes"`
}

// ArticleCategory 文章类别及数量（对应 database.ArticleCategory）。
type ArticleCategory struct {
	Category string `json:"category"`
	Number   int    `json:"number"`
}

// BlogTag 固定标签库标签（对应 database.BlogTag）。
type BlogTag struct {
	Tag    string `json:"tag"`
	Group  string `json:"group"`
	Number int    `json:"number"`
}

// ArticleLike 文章收藏记录（对应 database.ArticleLike）。
type ArticleLike struct {
	ID        uint
	ArticleID string
	UserID    uint
}

// 领域规则错误（文案必须与重构前一致）。
var (
	ErrInvalidCategory = errors.New("分类只能是 技术 或 生活")
	ErrTagNotExist     = errors.New("标签不存在")
)

// ValidateCategory 校验文章分类（技术/生活）。
func ValidateCategory(category string) error {
	if category != "技术" && category != "生活" {
		return ErrInvalidCategory
	}
	return nil
}

var illustrationRegex = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)

// IllustrationsOf 提取 Markdown 内容中的插图 URL（与 utils.FindIllustrations 一致）。
func IllustrationsOf(text string) []string {
	matches := illustrationRegex.FindAllStringSubmatch(text, -1)
	var out []string
	for _, match := range matches {
		if len(match) > 2 {
			out = append(out, match[2])
		}
	}
	return out
}

// DiffTags 计算标签/URL 集合差异：新增与移除（与 utils.DiffArrays 一致）。
func DiffTags(oldArray, newArray []string) (added, removed []string) {
	oldMap := make(map[string]struct{})
	for _, item := range oldArray {
		oldMap[item] = struct{}{}
	}
	for _, item := range newArray {
		if _, exists := oldMap[item]; !exists {
			added = append(added, item)
		}
	}
	newMap := make(map[string]struct{})
	for _, item := range newArray {
		newMap[item] = struct{}{}
	}
	for _, item := range oldArray {
		if _, exists := newMap[item]; !exists {
			removed = append(removed, item)
		}
	}
	return
}

// SearchSpec 文章搜索条件（前台搜索与后台列表共用，infrastructure 负责映射为 ES 查询）。
type SearchSpec struct {
	Query    string
	Category string
	Tag      string
	Sort     string
	Order    string
	Page     int
	PageSize int
	// ListMode 为 true 时走后台列表逻辑（标题/简介模糊 + 置顶+创建时间排序）；
	// 否则走前台搜索逻辑（关键词多字段匹配 + 自定义排序）。
	ListMode bool
	Title    *string
	Abstract *string
}

// SearchResult 搜索结果；Hits 保持 ES hit 原始结构（_id/_source），JSON 序列化与 types.Hit 等价。
type SearchResult struct {
	Hits  any
	Total int64
}

// ArticleRepository 文章 MySQL 仓储端口（分类/标签计数与图片类别更新的事务联动封装在实现内部）。
type ArticleRepository interface {
	// Categories 获取全部文章类别及数量（保持原排序：无）。
	Categories(ctx context.Context) ([]ArticleCategory, error)
	// Tags 获取全部固定标签（按 `group` ASC, number DESC, tag ASC 排序）。
	Tags(ctx context.Context) ([]BlogTag, error)
	// CheckTagsExist 校验标签必须全部存在于固定标签库。
	CheckTagsExist(ctx context.Context, tags []string) error
	// LikeTx 收藏/取消收藏（DB 事务），返回 delta（+1 收藏 / -1 取消）。
	LikeTx(ctx context.Context, userID uint, articleID string) (int, error)
	// IsLike 用户是否已收藏。
	IsLike(ctx context.Context, userID uint, articleID string) (bool, error)
	// LikesList 用户收藏的分页列表。
	LikesList(ctx context.Context, userID uint, page, pageSize int) ([]ArticleLike, int64, error)
	// CreateWithCounts 创建文章时联动分类/标签计数与图片类别更新（DB 事务）。
	CreateWithCounts(ctx context.Context, a *Article, cover string, illustrations []string) error
	// UpdateWithCounts 更新文章时联动分类/标签计数与图片类别更新（DB 事务）。
	UpdateWithCounts(ctx context.Context, a *Article, old *Article, newCover string, newIllustrations []string, addedIllustrations, removedIllustrations []string) error
	// DeleteWithCounts 删除文章时联动分类/标签计数与图片类别更新（DB 事务）。
	DeleteWithCounts(ctx context.Context, a *Article, cover string, illustrations []string) error
}

// EsArticleStore 文章 Elasticsearch 端口。
type EsArticleStore interface {
	// Index 索引文章（refresh=true）。
	Index(ctx context.Context, a *Article) error
	// Update 更新文章字段（refresh=true）。
	Update(ctx context.Context, id string, doc any) error
	// Get 按 ID 获取文章；不存在时返回 ErrDocumentNotFound。
	Get(ctx context.Context, id string) (*Article, error)
	// DeleteByIDs 批量删除（refresh=true）。
	DeleteByIDs(ctx context.Context, ids []string) error
	// Exists 检查标题是否已存在。
	Exists(ctx context.Context, title string) (bool, error)
	// AddLikes 通过脚本增减收藏数（ctx._source.likes += delta）。
	AddLikes(ctx context.Context, id string, delta int) error
	// AddViews 通过脚本累加浏览量（ctx._source.views += num）。
	AddViews(ctx context.Context, id string, num int) error
	// Search 按规格搜索，返回原始 hit 列表。
	Search(ctx context.Context, spec SearchSpec) (SearchResult, error)
}

// ErrDocumentNotFound ES 文档不存在。
var ErrDocumentNotFound = errors.New("document not found")

// ViewCounter 浏览量计数端口（Redis hash 实现）。
type ViewCounter interface {
	// Set 浏览量 +1。
	Set(ctx context.Context, id string) error
	// GetInfo 取出全部计数。
	GetInfo(ctx context.Context) map[string]int
	// Clear 清空计数。
	Clear(ctx context.Context)
}
