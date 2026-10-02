// Package forum_test 提供论坛 Handler 的 gin 冒烟测试。
package forum_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"server/config"
	forumapp "server/internal/application/forum"
	forumdomain "server/internal/domain/forum"
	"server/internal/domain/shared"
	"server/internal/interface/http/handler/forum"
	"server/internal/interface/http/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type stubRepo struct {
	tags    []forumdomain.Tag
	list    []*forumdomain.ForumPost
	total   int64
	detail  *forumdomain.ForumDetail
	publish []*forumdomain.ForumPost
	err     error
}

func (s *stubRepo) Tags(_ context.Context) ([]forumdomain.Tag, error) { return s.tags, s.err }
func (s *stubRepo) Publish(_ context.Context, p *forumdomain.ForumPost) (uint, error) {
	s.publish = append(s.publish, p)
	return 9, s.err
}
func (s *stubRepo) List(_ context.Context, _ forumdomain.ListCond) ([]*forumdomain.ForumPost, int64, error) {
	return s.list, s.total, s.err
}
func (s *stubRepo) Detail(_ context.Context, id uint) (*forumdomain.ForumPost, []*forumdomain.ForumComment, error) {
	if s.detail == nil {
		return nil, nil, nil
	}
	return &s.detail.ForumPost, s.detail.Comments, s.err
}
func (s *stubRepo) Like(_ context.Context, _, _ uint) (bool, int, error) { return false, 0, s.err }
func (s *stubRepo) Comment(_ context.Context, _ *forumdomain.ForumComment) error {
	return s.err
}
func (s *stubRepo) ManageList(_ context.Context, _ forumdomain.ManageListCond) ([]*forumdomain.ForumPost, int64, error) {
	return s.list, s.total, s.err
}
func (s *stubRepo) DeletePosts(_ context.Context, _ []uint, _ uint, _ shared.RoleID) error {
	return s.err
}
func (s *stubRepo) ManageComments(_ context.Context, _ forumdomain.ManageCommentCond) ([]*forumdomain.ManageComment, int64, error) {
	return nil, 0, s.err
}
func (s *stubRepo) DeleteComments(_ context.Context, _ []uint, _ uint, _ shared.RoleID) error {
	return s.err
}

func newEnv(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	log := zap.NewNop()
	stub := &stubRepo{
		tags:  []forumdomain.Tag{{Tag: "Go", Group: "tech", Number: 3}},
		list:  []*forumdomain.ForumPost{{ID: 3, Title: "帖子"}},
		total: 1,
		detail: &forumdomain.ForumDetail{
			ForumPost: forumdomain.ForumPost{ID: 3, Title: "帖子", Tags: []string{"Go"}, Images: []string{}, LikeCount: 2},
			Comments:  []*forumdomain.ForumComment{{ID: 1, PostID: 3, Content: "评论"}},
		},
	}
	svc := forumapp.NewService(stub, &config.Config{}, log)
	h := forum.NewHandler(svc, &config.Config{}, log)

	engine := gin.New()
	engine.Use(middleware.GinLogger(log), middleware.GinRecovery(log, true))
	const prefix = "api"
	publicGroup := engine.Group(prefix)
	forumPublicRouter := publicGroup.Group("forum")
	forumPublicRouter.GET("list", h.List)
	forumPublicRouter.GET("detail", h.Detail)
	forumPublicRouter.GET("tags", h.Tags)
	return engine
}

func TestSmokeForumList(t *testing.T) {
	engine := newEnv(t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/forum/list", nil)
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status want 200, got %d body=%s", w.Code, w.Body.String())
	}
	if got := w.Body.String(); got != `{"code":0,"data":{"list":[{"id":3,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","user_id":0,"user":{"id":0,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","uuid":"00000000-0000-0000-0000-000000000000","username":"","email":"","openid":"","avatar":"","address":"","signature":"","role_id":0,"register":"邮箱","freeze":false},"title":"帖子","content":"","category":"","tags":null,"images":null,"like_count":0,"comment_count":0,"view_count":0}],"total":1},"msg":"success"}` {
		t.Fatalf("unexpected body:\n%s", got)
	}
}

func TestSmokeForumDetail(t *testing.T) {
	engine := newEnv(t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/forum/detail?id=3", nil)
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status want 200, got %d body=%s", w.Code, w.Body.String())
	}
	if got := w.Body.String(); got != `{"code":0,"data":{"id":3,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","user_id":0,"user":{"id":0,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","uuid":"00000000-0000-0000-0000-000000000000","username":"","email":"","openid":"","avatar":"","address":"","signature":"","role_id":0,"register":"邮箱","freeze":false},"title":"帖子","content":"","category":"","tags":["Go"],"images":[],"like_count":2,"comment_count":0,"view_count":0,"comments":[{"id":1,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","post_id":3,"parent_id":0,"user_id":0,"user":{"id":0,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","uuid":"00000000-0000-0000-0000-000000000000","username":"","email":"","openid":"","avatar":"","address":"","signature":"","role_id":0,"register":"邮箱","freeze":false},"content":"评论","children":null}]},"msg":"success"}` {
		t.Fatalf("unexpected body:\n%s", got)
	}
}

func TestSmokeForumTags(t *testing.T) {
	engine := newEnv(t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/forum/tags", nil)
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status want 200, got %d body=%s", w.Code, w.Body.String())
	}
	if got := w.Body.String(); got != `{"code":0,"data":[{"tag":"Go","group":"tech","number":3}],"msg":"success"}` {
		t.Fatalf("unexpected body:\n%s", got)
	}
}
