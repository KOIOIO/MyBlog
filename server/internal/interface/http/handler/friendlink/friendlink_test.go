// Package friendlink_test 提供友链 Handler 冒烟测试。
package friendlink_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	flapp "server/internal/application/friendlink"
	fldomain "server/internal/domain/friendlink"
	"server/internal/interface/http/handler/friendlink"
	"server/internal/interface/http/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type stubRepo struct {
	info []*fldomain.FriendLink
}

func (s *stubRepo) Info(context.Context) ([]*fldomain.FriendLink, int64, error) {
	return s.info, int64(len(s.info)), nil
}
func (s *stubRepo) Create(context.Context, *fldomain.FriendLink) error { return nil }
func (s *stubRepo) Delete(context.Context, []uint) error               { return nil }
func (s *stubRepo) Update(context.Context, *fldomain.FriendLink) error { return nil }
func (s *stubRepo) List(context.Context, fldomain.ListCond) ([]*fldomain.FriendLink, int64, error) {
	return nil, 0, nil
}

func TestSmokeFriendLinkInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := zap.NewNop()
	svc := flapp.NewService(&stubRepo{info: []*fldomain.FriendLink{{ID: 1, Name: "n", Logo: "/l.png"}}})
	h := friendlink.NewHandler(svc, log)
	engine := gin.New()
	engine.Use(middleware.GinLogger(log), middleware.GinRecovery(log, true))
	engine.GET("/api/friendLink/info", h.Info)

	w := httptest.NewRecorder()
	engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/friendLink/info", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("status: %d body=%s", w.Code, w.Body.String())
	}
	if got := w.Body.String(); got != `{"code":0,"data":{"list":[{"id":1,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","logo":"/l.png","link":"","name":"n","description":""}],"total":1},"msg":"success"}` {
		t.Fatalf("body: %s", got)
	}
}
