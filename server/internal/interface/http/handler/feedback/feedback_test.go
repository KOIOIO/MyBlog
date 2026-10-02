// Package feedback_test 提供反馈 Handler 冒烟测试。
package feedback_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	fbapp "server/internal/application/feedback"
	fbdomain "server/internal/domain/feedback"
	"server/internal/interface/http/handler/feedback"
	"server/internal/interface/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

type stubRepo struct {
	newest []*fbdomain.Feedback
}

func (s *stubRepo) Newest(context.Context, int) ([]*fbdomain.Feedback, error) { return s.newest, nil }
func (s *stubRepo) Create(context.Context, *fbdomain.Feedback) error          { return nil }
func (s *stubRepo) Info(context.Context, uuid.UUID) ([]*fbdomain.Feedback, error) {
	return s.newest, nil
}
func (s *stubRepo) Delete(context.Context, []uint) error      { return nil }
func (s *stubRepo) Reply(context.Context, uint, string) error { return nil }
func (s *stubRepo) List(context.Context, int, int) ([]*fbdomain.Feedback, int64, error) {
	return nil, 0, nil
}

func TestSmokeFeedbackNew(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := zap.NewNop()
	u := uuid.FromStringOrNil("fbd5364d-bb15-11f1-b230-16b7a2303b52")
	svc := fbapp.NewService(&stubRepo{newest: []*fbdomain.Feedback{{ID: 1, UserUUID: u, Content: "c"}}})
	h := feedback.NewHandler(svc, log)
	engine := gin.New()
	engine.Use(middleware.GinLogger(log), middleware.GinRecovery(log, true))
	engine.GET("/api/feedback/new", h.New)

	w := httptest.NewRecorder()
	engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/feedback/new", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("status: %d body=%s", w.Code, w.Body.String())
	}
	if got := w.Body.String(); got != `{"code":0,"data":[{"id":1,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","user_uuid":"fbd5364d-bb15-11f1-b230-16b7a2303b52","content":"c","reply":""}],"msg":"success"}` {
		t.Fatalf("body: %s", got)
	}
}
