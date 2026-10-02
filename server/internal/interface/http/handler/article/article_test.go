package article_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	articleapp "server/internal/application/article"
	articledomain "server/internal/domain/article"
	articlehandler "server/internal/interface/http/handler/article"
	"server/internal/interface/http/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ---- 手写 stub（不引入 mock 框架） ----

type fakeEsStore struct {
	mu   sync.Mutex
	docs map[string]*articledomain.Article
}

func (s *fakeEsStore) Index(ctx context.Context, a *articledomain.Article) error { return nil }
func (s *fakeEsStore) Update(ctx context.Context, id string, doc any) error {
	return nil
}
func (s *fakeEsStore) Get(ctx context.Context, id string) (*articledomain.Article, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if a, ok := s.docs[id]; ok {
		cp := *a
		return &cp, nil
	}
	return nil, articledomain.ErrDocumentNotFound
}
func (s *fakeEsStore) DeleteByIDs(ctx context.Context, ids []string) error { return nil }
func (s *fakeEsStore) Exists(ctx context.Context, title string) (bool, error) {
	return false, nil
}
func (s *fakeEsStore) AddLikes(ctx context.Context, id string, delta int) error { return nil }
func (s *fakeEsStore) AddViews(ctx context.Context, id string, num int) error   { return nil }
func (s *fakeEsStore) Search(ctx context.Context, spec articledomain.SearchSpec) (articledomain.SearchResult, error) {
	return articledomain.SearchResult{Hits: []any{}, Total: 0}, nil
}

type stubRepo struct{}

func (r *stubRepo) Categories(ctx context.Context) ([]articledomain.ArticleCategory, error) {
	return []articledomain.ArticleCategory{{Category: "技术", Number: 2}}, nil
}
func (r *stubRepo) Tags(ctx context.Context) ([]articledomain.BlogTag, error) {
	return []articledomain.BlogTag{{Tag: "Gin", Group: "tech", Number: 1}}, nil
}
func (r *stubRepo) CheckTagsExist(ctx context.Context, tags []string) error { return nil }
func (r *stubRepo) LikeTx(ctx context.Context, userID uint, articleID string) (int, error) {
	return 1, nil
}
func (r *stubRepo) IsLike(ctx context.Context, userID uint, articleID string) (bool, error) {
	return false, nil
}
func (r *stubRepo) LikesList(ctx context.Context, userID uint, page, pageSize int) ([]articledomain.ArticleLike, int64, error) {
	return nil, 0, nil
}
func (r *stubRepo) CreateWithCounts(ctx context.Context, a *articledomain.Article, cover string, illustrations []string) error {
	return nil
}
func (r *stubRepo) UpdateWithCounts(ctx context.Context, a *articledomain.Article, old *articledomain.Article, newCover string, newIllustrations []string, addedIllustrations, removedIllustrations []string) error {
	return nil
}
func (r *stubRepo) DeleteWithCounts(ctx context.Context, a *articledomain.Article, cover string, illustrations []string) error {
	return nil
}

type stubViews struct{}

func (v *stubViews) Set(ctx context.Context, id string) error { return nil }
func (v *stubViews) GetInfo(ctx context.Context) map[string]int {
	return map[string]int{}
}
func (v *stubViews) Clear(ctx context.Context) {}

func newSmokeEnv(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	log := zap.NewNop()
	esStore := &fakeEsStore{docs: map[string]*articledomain.Article{
		"4": {CreatedAt: "2026-10-01 00:00:00", Title: "Gin 框架中间件机制深度解析", Category: "技术",
			Tags: []string{"Gin"}, Abstract: "摘要", Content: "内容", Views: 10, Comments: 1, Likes: 2},
	}}
	svc := articleapp.NewService(&stubRepo{}, esStore, &stubViews{})
	h := articlehandler.NewHandler(svc, log)

	engine := gin.New()
	engine.Use(middleware.GinLogger(log), middleware.GinRecovery(log, true))
	const prefix = "api"
	publicGroup := engine.Group(prefix)
	articlePublicRouter := publicGroup.Group("article")

	articlePublicRouter.GET(":id", h.InfoByID)
	articlePublicRouter.GET("search", h.Search)
	articlePublicRouter.GET("category", h.Category)
	articlePublicRouter.GET("tags", h.Tags)
	return engine
}

func doGet(engine *gin.Engine, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	return w
}

func TestSmokeArticleInfoByID(t *testing.T) {
	engine := newSmokeEnv(t)
	w := doGet(engine, "/api/article/4")
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), `"title":"Gin 框架中间件机制深度解析"`) {
		t.Fatalf("want article title in body, got %s", w.Body.String())
	}
}

func TestSmokeArticleInfoNotFound(t *testing.T) {
	engine := newSmokeEnv(t)
	w := doGet(engine, "/api/article/999")
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), `"msg":"Failed to get article information"`) {
		t.Fatalf("want failure message, got %s", w.Body.String())
	}
}

func TestSmokeArticleSearch(t *testing.T) {
	engine := newSmokeEnv(t)
	w := doGet(engine, "/api/article/search?query=Gin&order=desc")
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), `"total":0`) {
		t.Fatalf("want empty search result, got %s", w.Body.String())
	}
}

func TestSmokeArticleCategoryAndTags(t *testing.T) {
	engine := newSmokeEnv(t)
	w := doGet(engine, "/api/article/category")
	if !strings.Contains(w.Body.String(), `"category":"技术"`) {
		t.Fatalf("want category data, got %s", w.Body.String())
	}
	w = doGet(engine, "/api/article/tags")
	if !strings.Contains(w.Body.String(), `"tag":"Gin"`) {
		t.Fatalf("want tag data, got %s", w.Body.String())
	}
}
