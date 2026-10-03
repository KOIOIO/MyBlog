// Package article 提供文章用例编排（ES 读写 + MySQL 计数联动 + 浏览量计数）。
package article

import (
	"context"
	"errors"
	"time"

	"server/internal/domain/article"
)

// LikeHit 收藏列表条目（保持原响应 JSON 结构：_id / _source）。
type LikeHit struct {
	Id_     string          `json:"_id"`
	Source_ article.Article `json:"_source"`
}

// Service 文章用例。
type Service struct {
	repo  article.ArticleRepository
	es    article.EsArticleStore
	views article.ViewCounter
}

// NewService 构造文章用例。
func NewService(repo article.ArticleRepository, esStore article.EsArticleStore, viewCounter article.ViewCounter) *Service {
	return &Service{repo: repo, es: esStore, views: viewCounter}
}

// InfoByID 获取文章详情（异步累加浏览量）。
func (s *Service) InfoByID(ctx context.Context, id string) (*article.Article, error) {
	go func() {
		_ = s.views.Set(context.Background(), id)
	}()
	return s.es.Get(ctx, id)
}

// Search 前台搜索。
func (s *Service) Search(ctx context.Context, spec article.SearchSpec) (any, int64, error) {
	res, err := s.es.Search(ctx, spec)
	if err != nil {
		return nil, 0, err
	}
	return res.Hits, res.Total, nil
}

// Categories 全部类别及数量。
func (s *Service) Categories(ctx context.Context) ([]article.ArticleCategory, error) {
	return s.repo.Categories(ctx)
}

// Tags 全部固定标签。
func (s *Service) Tags(ctx context.Context) ([]article.BlogTag, error) {
	return s.repo.Tags(ctx)
}

// Like 收藏/取消收藏（DB 事务 + ES 计数脚本）。
func (s *Service) Like(ctx context.Context, userID uint, articleID string) error {
	delta, err := s.repo.LikeTx(ctx, userID, articleID)
	if err != nil {
		return err
	}
	return s.es.AddLikes(ctx, articleID, delta)
}

// IsLike 收藏状态。
func (s *Service) IsLike(ctx context.Context, userID uint, articleID string) (bool, error) {
	return s.repo.IsLike(ctx, userID, articleID)
}

// LikesList 收藏列表（逐篇回查 ES，清空 UpdatedAt/Keyword/Content）。
func (s *Service) LikesList(ctx context.Context, userID uint, pageNo, pageSize int) (any, int64, error) {
	likes, total, err := s.repo.LikesList(ctx, userID, pageNo, pageSize)
	if err != nil {
		return nil, 0, err
	}
	list := make([]LikeHit, 0, len(likes))
	for _, like := range likes {
		articleDoc, err := s.es.Get(ctx, like.ArticleID)
		if err != nil {
			return nil, 0, err
		}
		articleDoc.UpdatedAt = ""
		articleDoc.Keyword = ""
		articleDoc.Content = ""
		list = append(list, LikeHit{Id_: like.ArticleID, Source_: *articleDoc})
	}
	return list, total, nil
}

// Create 创建文章（标题唯一校验 + 分类/标签校验 + 计数联动 + ES 索引）。
func (s *Service) Create(ctx context.Context, a *article.Article) error {
	exists, err := s.es.Exists(ctx, a.Title)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("the article already exists")
	}
	if err := article.ValidateCategory(a.Category); err != nil {
		return err
	}
	if err := s.repo.CheckTagsExist(ctx, a.Tags); err != nil {
		return err
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	a.CreatedAt = now
	a.UpdatedAt = now
	a.Keyword = a.Title

	illustrations := article.IllustrationsOf(a.Content)
	if err := s.repo.CreateWithCounts(ctx, a, a.Cover, illustrations); err != nil {
		return err
	}
	return s.es.Index(ctx, a)
}

// Delete 删除文章（计数联动 + ES 批量删除）。
func (s *Service) Delete(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	for _, id := range ids {
		a, err := s.es.Get(ctx, id)
		if err != nil {
			return err
		}
		illustrations := article.IllustrationsOf(a.Content)
		if err := s.repo.DeleteWithCounts(ctx, a, a.Cover, illustrations); err != nil {
			return err
		}
	}
	return s.es.DeleteByIDs(ctx, ids)
}

// Update 更新文章（分类/标签校验 + 封面与插图差异联动 + ES 更新）。
func (s *Service) Update(ctx context.Context, id string, a *article.Article) error {
	if err := article.ValidateCategory(a.Category); err != nil {
		return err
	}
	if err := s.repo.CheckTagsExist(ctx, a.Tags); err != nil {
		return err
	}

	a.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	a.Keyword = a.Title

	old, err := s.es.Get(ctx, id)
	if err != nil {
		return err
	}
	// 编辑表单不携带创建时间，继承 ES 旧文档值，避免以空串覆盖 created_at 触发 ES date 解析失败
	a.CreatedAt = old.CreatedAt

	oldIllustrations := article.IllustrationsOf(old.Content)
	newIllustrations := article.IllustrationsOf(a.Content)
	added, removed := article.DiffTags(oldIllustrations, newIllustrations)

	if err := s.repo.UpdateWithCounts(ctx, a, old, a.Cover, newIllustrations, added, removed); err != nil {
		return err
	}
	return s.es.Update(ctx, id, a)
}

// SetTop 置顶（仅更新 is_top 字段）。
func (s *Service) SetTop(ctx context.Context, id string, isTop bool) error {
	topValue := 0
	if isTop {
		topValue = 1
	}
	return s.es.Update(ctx, id, map[string]any{"is_top": topValue})
}

// List 后台文章列表。
func (s *Service) List(ctx context.Context, spec article.SearchSpec) (any, int64, error) {
	spec.ListMode = true
	res, err := s.es.Search(ctx, spec)
	if err != nil {
		return nil, 0, err
	}
	return res.Hits, res.Total, nil
}

// SyncViews 将 Redis 浏览量增量同步到 ES 并清空
// （修复原 task 仅同步第一条且从不清空的缺陷；HTTP 行为不受影响）。
func (s *Service) SyncViews(ctx context.Context) error {
	info := s.views.GetInfo(ctx)
	for id, num := range info {
		if num == 0 {
			continue
		}
		if err := s.es.AddViews(ctx, id, num); err != nil {
			return err
		}
	}
	s.views.Clear(ctx)
	return nil
}
