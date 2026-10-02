package config_test

import (
	"context"
	"testing"

	"server/config"
	cconfig "server/internal/application/config"
	"server/internal/domain/website"
)

type stubStore struct {
	saved int
}

func (s *stubStore) Save(context.Context) error {
	s.saved++
	return nil
}

type stubImages struct {
	initCalled   []string
	changeCalled []string
	changeCat    string
}

func (s *stubImages) CarouselURLs(context.Context) ([]string, error) { return nil, nil }
func (s *stubImages) ChangeCategory(_ context.Context, urls []string, category string) error {
	s.changeCalled = append(s.changeCalled, urls...)
	s.changeCat = category
	return nil
}
func (s *stubImages) InitCategory(_ context.Context, urls []string) error {
	s.initCalled = append(s.initCalled, urls...)
	return nil
}

var _ website.WebsiteImagePort = (*stubImages)(nil)

func TestUpdateSystemSaves(t *testing.T) {
	cfg := &config.Config{}
	store := &stubStore{}
	svc := cconfig.NewService(cfg, store, &stubImages{})

	if err := svc.UpdateSystem(context.Background(), config.System{UseMultipoint: true, SessionsSecret: "s", OssType: "qiniu"}); err != nil {
		t.Fatal(err)
	}
	if cfg.System.OssType != "qiniu" || cfg.System.SessionsSecret != "s" || !cfg.System.UseMultipoint {
		t.Fatalf("config not mutated: %+v", cfg.System)
	}
	if store.saved != 1 {
		t.Fatalf("store not saved, saved=%d", store.saved)
	}
}

func TestUpdateWebsiteImageCategorySync(t *testing.T) {
	cfg := &config.Config{}
	cfg.Website.Logo = "/old/logo.png"
	cfg.Website.FullLogo = "/keep/full.png"
	store := &stubStore{}
	imgs := &stubImages{}
	svc := cconfig.NewService(cfg, store, imgs)

	newW := cfg.Website
	newW.Logo = "/new/logo.png"
	if err := svc.UpdateWebsite(context.Background(), newW); err != nil {
		t.Fatal(err)
	}
	// removed: /old/logo.png → InitCategory；added: /new/logo.png → ChangeCategory("系统")
	if len(imgs.initCalled) != 1 || imgs.initCalled[0] != "/old/logo.png" {
		t.Fatalf("init want [/old/logo.png], got %v", imgs.initCalled)
	}
	if len(imgs.changeCalled) != 1 || imgs.changeCalled[0] != "/new/logo.png" || imgs.changeCat != "系统" {
		t.Fatalf("change want [/new/logo.png 系统], got %v %q", imgs.changeCalled, imgs.changeCat)
	}
	if cfg.Website.Logo != "/new/logo.png" {
		t.Fatalf("website not mutated: %+v", cfg.Website)
	}
	if store.saved != 1 {
		t.Fatalf("store not saved, saved=%d", store.saved)
	}
}

func TestUpdateJwtAndGaode(t *testing.T) {
	cfg := &config.Config{}
	store := &stubStore{}
	svc := cconfig.NewService(cfg, store, &stubImages{})

	if err := svc.UpdateJwt(context.Background(), config.Jwt{Issuer: "x"}); err != nil {
		t.Fatal(err)
	}
	if err := svc.UpdateGaode(context.Background(), config.Gaode{Key: "k"}); err != nil {
		t.Fatal(err)
	}
	if store.saved != 2 {
		t.Fatalf("want 2 saves, got %d", store.saved)
	}
}
