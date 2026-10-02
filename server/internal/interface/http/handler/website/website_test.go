// Package website_test 提供网站 Handler 的 gin 冒烟测试。
package website_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"server/config"
	websiteapp "server/internal/application/website"
	websitedomain "server/internal/domain/website"
	websitehandler "server/internal/interface/http/handler/website"
	"server/internal/interface/http/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type stubImages struct {
	urls []string
}

func (s *stubImages) CarouselURLs(context.Context) ([]string, error)         { return s.urls, nil }
func (s *stubImages) ChangeCategory(context.Context, []string, string) error { return nil }
func (s *stubImages) InitCategory(context.Context, []string) error           { return nil }

type stubFooters struct {
	list []*websitedomain.FooterLink
}

func (s *stubFooters) List(context.Context) ([]*websitedomain.FooterLink, error) { return s.list, nil }
func (s *stubFooters) Save(context.Context, *websitedomain.FooterLink) error     { return nil }
func (s *stubFooters) Delete(context.Context, *websitedomain.FooterLink) error   { return nil }

type stubNews struct {
	data websitedomain.HotSearchData
}

func (s *stubNews) GetHotSearchData(_ context.Context, source string) (websitedomain.HotSearchData, error) {
	s.data.Source = source
	return s.data, nil
}
func (s *stubNews) WarmAll(context.Context) error { return nil }

type stubCal struct {
	data websitedomain.Calendar
}

func (s *stubCal) GetCalendarByDate(_ context.Context, dateStr string) (websitedomain.Calendar, error) {
	s.data.Date = dateStr
	return s.data, nil
}

func newEnv(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	log := zap.NewNop()
	cfg := &config.Config{}
	cfg.Website.Logo = "/uploads/image/logo.png"
	cfg.Website.Title = "Folio"

	svc := websiteapp.NewService(
		&stubImages{urls: []string{"/uploads/image/01.jpg"}},
		&stubFooters{list: []*websitedomain.FooterLink{{ID: 1, Name: "n", Link: "https://x"}}},
		&stubNews{data: websitedomain.HotSearchData{Source: "baidu", UpdateTime: "t", HotList: []websitedomain.HotItem{{Index: 1, Title: "热搜"}}}},
		&stubCal{data: websitedomain.Calendar{Date: "2026年10月2日 五", LunarDate: "八月廿二"}},
		cfg,
	)
	h := websitehandler.NewHandler(svc, cfg, log)

	engine := gin.New()
	engine.Use(middleware.GinLogger(log), middleware.GinRecovery(log, true))
	const prefix = "api"
	publicGroup := engine.Group(prefix)
	wPublic := publicGroup.Group("website")
	wPublic.GET("logo", h.Logo)
	wPublic.GET("title", h.Title)
	wPublic.GET("info", h.Info)
	wPublic.GET("carousel", h.Carousel)
	wPublic.GET("news", h.News)
	wPublic.GET("calendar", h.Calendar)
	wPublic.GET("footerLink", h.FooterLink)
	return engine
}

func TestSmokeWebsiteLogo301(t *testing.T) {
	engine := newEnv(t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/website/logo", nil)
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusMovedPermanently {
		t.Fatalf("status want 301, got %d", w.Code)
	}
	if loc := w.Header().Get("Location"); loc != "/uploads/image/logo.png" {
		t.Fatalf("location want /uploads/image/logo.png, got %s", loc)
	}
}

func TestSmokeWebsiteTitleAndInfo(t *testing.T) {
	engine := newEnv(t)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/website/title", nil))
	if got := w.Body.String(); got != `{"title":"Folio"}` {
		t.Fatalf("title body: %s", got)
	}
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/website/info", nil))
	if got := w.Body.String(); got != `{"code":0,"data":{"logo":"/uploads/image/logo.png","full_logo":"","title":"Folio","slogan":"","slogan_en":"","description":"","version":"","created_at":"","icp_filing":"","public_security_filing":"","bilibili_url":"","gitee_url":"","github_url":"","name":"","job":"","address":"","email":"","qq_image":"","wechat_image":""},"msg":"success"}` {
		t.Fatalf("info body: %s", got)
	}
}

func TestSmokeWebsiteCarouselNewsCalendarFooter(t *testing.T) {
	engine := newEnv(t)
	cases := []struct {
		path string
		want string
	}{
		{"/api/website/carousel", `{"code":0,"data":["/uploads/image/01.jpg"],"msg":"success"}`},
		{"/api/website/news?source=baidu", `{"code":0,"data":{"source":"baidu","update_time":"t","hot_list":[{"index":1,"title":"热搜","description":"","image":"","popularity":"","url":""}]},"msg":"success"}`},
		{"/api/website/calendar", `{"code":0,"data":{"date":"2026/1002","lunar_date":"八月廿二","ganzhi":"","zodiac":"","day_of_year":"","solar_term":"","auspicious":"","inauspicious":""},"msg":"success"}`},
		{"/api/website/footerLink", `{"code":0,"data":[{"id":1,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","logo":"","link":"https://x","name":"n","description":""}],"msg":"success"}`},
	}
	for _, c := range cases {
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, c.path, nil))
		if got := w.Body.String(); got != c.want {
			t.Errorf("%s body mismatch:\n got: %s\nwant: %s", c.path, got, c.want)
		}
	}
}
