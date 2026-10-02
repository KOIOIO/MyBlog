// Package comment_test 提供评论 Handler 的 gin 冒烟测试。
package comment_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	commentapp "server/internal/application/comment"
	commentdomain "server/internal/domain/comment"
	"server/internal/domain/shared"
	"server/internal/interface/http/handler/comment"
	"server/internal/interface/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

type stubRepo struct {
	tree    []*commentdomain.Comment
	newest  []*commentdomain.Comment
	list    []*commentdomain.Comment
	total   int64
	created []*commentdomain.Comment
	err     error
}

func (s *stubRepo) ByArticle(_ context.Context, _ string) ([]*commentdomain.Comment, error) {
	return s.tree, s.err
}
func (s *stubRepo) Newest(_ context.Context, _ int) ([]*commentdomain.Comment, error) {
	return s.newest, s.err
}
func (s *stubRepo) Create(_ context.Context, c *commentdomain.Comment) error {
	s.created = append(s.created, c)
	return s.err
}
func (s *stubRepo) DeleteTree(_ context.Context, _ []uint, _ uuid.UUID, _ shared.RoleID) error {
	return s.err
}
func (s *stubRepo) ByUser(_ context.Context, _ uuid.UUID) ([]*commentdomain.Comment, error) {
	return s.tree, s.err
}
func (s *stubRepo) List(_ context.Context, _ commentdomain.ListCond) ([]*commentdomain.Comment, int64, error) {
	return s.list, s.total, s.err
}

func newEnv(t *testing.T) (*gin.Engine, *stubRepo) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	log := zap.NewNop()
	stub := &stubRepo{
		tree: []*commentdomain.Comment{{
			ID:        1,
			ArticleID: "15",
			Children:  []*commentdomain.Comment{{ID: 2, ArticleID: "15"}},
		}},
		newest: []*commentdomain.Comment{{ID: 3}},
	}
	svc := commentapp.NewService(stub, log)
	h := comment.NewHandler(svc, log)

	engine := gin.New()
	engine.Use(middleware.GinLogger(log), middleware.GinRecovery(log, true))
	const prefix = "api"
	publicGroup := engine.Group(prefix)
	commentPublicRouter := publicGroup.Group("comment")
	commentPublicRouter.GET(":article_id", h.InfoByArticleID)
	commentPublicRouter.GET("new", h.New)
	return engine, stub
}

func TestSmokeCommentInfoByArticleID(t *testing.T) {
	engine, _ := newEnv(t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/comment/15", nil)
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status want 200, got %d body=%s", w.Code, w.Body.String())
	}
	body := w.Body.String()
	if body != `{"code":0,"data":[{"id":1,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","article_id":"15","p_id":null,"children":[{"id":2,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","article_id":"15","p_id":null,"children":null,"user_uuid":"00000000-0000-0000-0000-000000000000","user":{"id":0,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","uuid":"00000000-0000-0000-0000-000000000000","username":"","email":"","openid":"","avatar":"","address":"","signature":"","role_id":0,"register":"邮箱","freeze":false},"content":""}],"user_uuid":"00000000-0000-0000-0000-000000000000","user":{"id":0,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","uuid":"00000000-0000-0000-0000-000000000000","username":"","email":"","openid":"","avatar":"","address":"","signature":"","role_id":0,"register":"邮箱","freeze":false},"content":""}],"msg":"success"}` {
		t.Fatalf("unexpected body:\n%s", body)
	}
}

func TestSmokeCommentNew(t *testing.T) {
	engine, _ := newEnv(t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/comment/new", nil)
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status want 200, got %d", w.Code)
	}
	if got := w.Body.String(); got != `{"code":0,"data":[{"id":3,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","article_id":"","p_id":null,"children":null,"user_uuid":"00000000-0000-0000-0000-000000000000","user":{"id":0,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","uuid":"00000000-0000-0000-0000-000000000000","username":"","email":"","openid":"","avatar":"","address":"","signature":"","role_id":0,"register":"邮箱","freeze":false},"content":""}],"msg":"success"}` {
		t.Fatalf("unexpected body:\n%s", got)
	}
}
