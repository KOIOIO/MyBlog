// Package advertisement_test 提供广告 Handler 冒烟测试。
package advertisement_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	advapp "server/internal/application/advertisement"
	advdomain "server/internal/domain/advertisement"
	"server/internal/interface/http/handler/advertisement"
	"server/internal/interface/http/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type stubRepo struct {
	info []*advdomain.Advertisement
}

func (s *stubRepo) Info(context.Context) ([]*advdomain.Advertisement, int64, error) {
	return s.info, int64(len(s.info)), nil
}
func (s *stubRepo) Create(context.Context, *advdomain.Advertisement) error { return nil }
func (s *stubRepo) Delete(context.Context, []uint) error                   { return nil }
func (s *stubRepo) Update(context.Context, *advdomain.Advertisement) error { return nil }
func (s *stubRepo) List(context.Context, advdomain.ListCond) ([]*advdomain.Advertisement, int64, error) {
	return nil, 0, nil
}

func TestSmokeAdvertisementInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := zap.NewNop()
	svc := advapp.NewService(&stubRepo{info: []*advdomain.Advertisement{{ID: 1, Title: "t", AdImage: "/a.png"}}})
	h := advertisement.NewHandler(svc, log)
	engine := gin.New()
	engine.Use(middleware.GinLogger(log), middleware.GinRecovery(log, true))
	engine.GET("/api/advertisement/info", h.Info)

	w := httptest.NewRecorder()
	engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/advertisement/info", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("status: %d body=%s", w.Code, w.Body.String())
	}
	if got := w.Body.String(); got != `{"code":0,"data":{"list":[{"id":1,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","ad_image":"/a.png","link":"","title":"t","content":""}],"total":1},"msg":"success"}` {
		t.Fatalf("body: %s", got)
	}
}
