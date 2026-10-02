// Package es 提供文章 Elasticsearch 适配器（实现 domain/article.EsArticleStore）。
package es

import (
	"context"
	"encoding/json"
	"strconv"

	"server/internal/common/page"
	"server/internal/domain/article"
	"server/model/other"
	"server/model/request"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/bulk"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/update"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/refresh"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/scriptlanguage"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
)

const articleIndex = "article_index"

// ArticleStore ES 文章存储。
type ArticleStore struct {
	es *elasticsearch.TypedClient
}

// NewArticleStore 构造 ES 文章存储。
func NewArticleStore(es *elasticsearch.TypedClient) *ArticleStore {
	return &ArticleStore{es: es}
}

// Index 索引文章（refresh=true）。
func (s *ArticleStore) Index(ctx context.Context, a *article.Article) error {
	_, err := s.es.Index(articleIndex).Request(a).Refresh(refresh.True).Do(ctx)
	return err
}

// Update 更新文章字段（refresh=true）。
func (s *ArticleStore) Update(ctx context.Context, id string, doc any) error {
	bytes, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	_, err = s.es.Update(articleIndex, id).Request(&update.Request{Doc: bytes}).Refresh(refresh.True).Do(ctx)
	return err
}

// Get 按 ID 获取文章。
func (s *ArticleStore) Get(ctx context.Context, id string) (*article.Article, error) {
	res, err := s.es.Get(articleIndex, id).Do(ctx)
	if err != nil {
		return nil, err
	}
	if !res.Found {
		return nil, article.ErrDocumentNotFound
	}
	var a article.Article
	if err := json.Unmarshal(res.Source_, &a); err != nil {
		return nil, err
	}
	return &a, nil
}

// DeleteByIDs 批量删除（refresh=true）。
func (s *ArticleStore) DeleteByIDs(ctx context.Context, ids []string) error {
	var request bulk.Request
	for _, id := range ids {
		request = append(request, types.OperationContainer{Delete: &types.DeleteOperation{Id_: &id}})
	}
	_, err := s.es.Bulk().Request(&request).Index(articleIndex).Refresh(refresh.True).Do(ctx)
	return err
}

// Exists 检查标题是否已存在（按 keyword 字段精确匹配）。
func (s *ArticleStore) Exists(ctx context.Context, title string) (bool, error) {
	req := &search.Request{
		Query: &types.Query{
			Match: map[string]types.MatchQuery{"keyword": {Query: title}},
		},
	}
	res, err := s.es.Search().Index(articleIndex).Request(req).Size(1).Do(ctx)
	if err != nil {
		return false, err
	}
	return res.Hits.Total.Value > 0, nil
}

// AddLikes 通过脚本增减收藏数。
func (s *ArticleStore) AddLikes(ctx context.Context, id string, delta int) error {
	source := "ctx._source.likes += " + strconv.Itoa(delta)
	script := types.Script{Source: &source, Lang: &scriptlanguage.Painless}
	_, err := s.es.Update(articleIndex, id).Script(&script).Do(ctx)
	return err
}

// AddViews 通过脚本累加浏览量。
func (s *ArticleStore) AddViews(ctx context.Context, id string, num int) error {
	source := "ctx._source.views += " + strconv.Itoa(num)
	script := types.Script{Source: &source, Lang: &scriptlanguage.Painless}
	_, err := s.es.Update(articleIndex, id).Script(&script).Do(ctx)
	return err
}

// Search 按规格搜索。
func (s *ArticleStore) Search(ctx context.Context, spec article.SearchSpec) (article.SearchResult, error) {
	req := &search.Request{Query: &types.Query{}}
	boolQuery := &types.BoolQuery{}

	if spec.ListMode {
		buildListQuery(boolQuery, spec)
	} else {
		buildSearchQuery(boolQuery, spec)
	}

	if boolQuery.Must != nil || boolQuery.Filter != nil || boolQuery.Should != nil {
		req.Query.Bool = boolQuery
	} else {
		req.Query.MatchAll = &types.MatchAllQuery{}
	}

	// 排序：置顶文章固定最前
	topOrder := sortorder.Desc
	req.Sort = []types.SortCombinations{
		types.SortOptions{SortOptions: map[string]types.FieldSort{"is_top": {Order: &topOrder}}},
	}
	if spec.ListMode {
		req.Sort = append(req.Sort, types.SortOptions{
			SortOptions: map[string]types.FieldSort{"created_at": {Order: &sortorder.Desc}},
		})
	} else if spec.Sort != "" {
		var sortField string
		switch spec.Sort {
		case "time":
			sortField = "created_at"
		case "view":
			sortField = "views"
		case "comment":
			sortField = "comments"
		case "like":
			sortField = "likes"
		default:
			sortField = "created_at"
		}
		order := sortorder.Desc
		if spec.Order == "asc" {
			order = sortorder.Asc
		}
		req.Sort = append(req.Sort, types.SortOptions{
			SortOptions: map[string]types.FieldSort{sortField: {Order: &order}},
		})
	} else {
		req.Sort = append(req.Sort, types.SortOptions{
			SortOptions: map[string]types.FieldSort{"created_at": {Order: &sortorder.Desc}},
		})
	}

	var sourceIncludes []string
	if !spec.ListMode {
		sourceIncludes = []string{"created_at", "cover", "title", "abstract", "category", "tags", "views", "comments", "likes", "is_top"}
	}

	option := other.EsOption{
		PageInfo: request.PageInfo{
			Page:     spec.Page,
			PageSize: spec.PageSize,
		},
		Index:          articleIndex,
		Request:        req,
		SourceIncludes: sourceIncludes,
	}
	res, total, err := page.EsPagination(ctx, s.es, option)
	if err != nil {
		return article.SearchResult{}, err
	}
	return article.SearchResult{Hits: res, Total: total}, nil
}

// buildSearchQuery 前台搜索：关键词多字段匹配 + 标签 Must + 类别 Filter。
func buildSearchQuery(boolQuery *types.BoolQuery, spec article.SearchSpec) {
	if spec.Query != "" {
		boolQuery.Should = []types.Query{
			{Match: map[string]types.MatchQuery{"title": {Query: spec.Query}}},
			{Match: map[string]types.MatchQuery{"keyword": {Query: spec.Query}}},
			{Match: map[string]types.MatchQuery{"abstract": {Query: spec.Query}}},
			{Match: map[string]types.MatchQuery{"content": {Query: spec.Query}}},
		}
	}
	if spec.Tag != "" {
		boolQuery.Must = []types.Query{
			{Match: map[string]types.MatchQuery{"tags": {Query: spec.Tag}}},
		}
	}
	if spec.Category != "" {
		boolQuery.Filter = []types.Query{
			{Term: map[string]types.TermQuery{"category": {Value: spec.Category}}},
		}
	}
}

// buildListQuery 后台列表：标题/简介模糊匹配 + 类别 Filter。
func buildListQuery(boolQuery *types.BoolQuery, spec article.SearchSpec) {
	if spec.Title != nil {
		boolQuery.Must = append(boolQuery.Must, types.Query{Match: map[string]types.MatchQuery{"title": {Query: *spec.Title}}})
	}
	if spec.Abstract != nil {
		boolQuery.Must = append(boolQuery.Must, types.Query{Match: map[string]types.MatchQuery{"abstract": {Query: *spec.Abstract}}})
	}
	if spec.Category != "" {
		boolQuery.Filter = []types.Query{
			{Term: map[string]types.TermQuery{"category": {Value: spec.Category}}},
		}
	}
}

// 编译期断言。
var _ article.EsArticleStore = (*ArticleStore)(nil)
