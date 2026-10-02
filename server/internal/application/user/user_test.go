package user

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"server/config"
	"server/internal/common/crypto"
	"server/internal/common/errs"
	"server/internal/domain/auth"
	"server/internal/domain/shared"
	userdomain "server/internal/domain/user"

	"github.com/gofrs/uuid"
)

// ---- 手写 stub（不引入 mock 框架） ----

type stubUserRepo struct {
	mu      sync.Mutex
	byID    map[uint]*userdomain.User
	byEmail map[string]*userdomain.User
	nextID  uint
}

func newStubUserRepo() *stubUserRepo {
	return &stubUserRepo{byID: make(map[uint]*userdomain.User), byEmail: make(map[string]*userdomain.User), nextID: 1}
}

func (r *stubUserRepo) FindByEmail(ctx context.Context, email string) (*userdomain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if u, ok := r.byEmail[email]; ok {
		return u, nil
	}
	return nil, errs.ErrNotFound
}

func (r *stubUserRepo) FindByUUID(ctx context.Context, uid uuid.UUID) (*userdomain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, u := range r.byID {
		if u.UUID == uid {
			return u, nil
		}
	}
	return nil, errs.ErrNotFound
}

func (r *stubUserRepo) FindByID(ctx context.Context, id uint) (*userdomain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if u, ok := r.byID[id]; ok {
		return u, nil
	}
	return nil, errs.ErrNotFound
}

func (r *stubUserRepo) Create(ctx context.Context, u *userdomain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	u.ID = r.nextID
	r.nextID++
	r.byID[u.ID] = u
	r.byEmail[u.Email] = u
	return nil
}

func (r *stubUserRepo) Update(ctx context.Context, u *userdomain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byID[u.ID] = u
	r.byEmail[u.Email] = u
	return nil
}

func (r *stubUserRepo) UpdateFreeze(ctx context.Context, id uint, frozen bool) (*userdomain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, ok := r.byID[id]
	if !ok {
		return nil, errs.ErrNotFound
	}
	prev := *u
	u.SetFrozen(frozen)
	return &prev, nil
}

func (r *stubUserRepo) UpdateInfo(ctx context.Context, id uint, username, address, signature string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, ok := r.byID[id]
	if !ok {
		return errs.ErrNotFound
	}
	u.Username = username
	u.Address = address
	u.Signature = signature
	return nil
}

func (r *stubUserRepo) UpdateAvatar(ctx context.Context, id uint, avatarURL string) (*userdomain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, ok := r.byID[id]
	if !ok {
		return nil, errs.ErrNotFound
	}
	prev := *u
	u.Avatar = avatarURL
	return &prev, nil
}

func (r *stubUserRepo) CountByDate(ctx context.Context, days int) (map[string]int, error) {
	return map[string]int{}, nil
}

func (r *stubUserRepo) Page(ctx context.Context, cond userdomain.ListCond) ([]*userdomain.User, int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []*userdomain.User
	for _, u := range r.byID {
		if cond.UUID != nil && !strings.Contains(u.UUID.String(), *cond.UUID) {
			continue
		}
		if cond.Register != nil && !strings.EqualFold(strconv.Itoa(int(u.Register)), *cond.Register) {
			continue
		}
		out = append(out, u)
	}
	return out, int64(len(out)), nil
}

func (r *stubUserRepo) FindIDByUUID(ctx context.Context, uuidVal uuid.UUID) (uint, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for id, u := range r.byID {
		if u.UUID == uuidVal {
			return id, nil
		}
	}
	return 0, errs.ErrNotFound
}

type stubLoginRepo struct {
	mu    sync.Mutex
	items []*userdomain.LoginRecord
}

func (r *stubLoginRepo) Create(ctx context.Context, rec *userdomain.LoginRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	rec.ID = uint(len(r.items)) + 1
	r.items = append(r.items, rec)
	return nil
}

func (r *stubLoginRepo) CountByDate(ctx context.Context, days int) (map[string]int, error) {
	return map[string]int{}, nil
}

func (r *stubLoginRepo) Page(ctx context.Context, cond userdomain.LoginListCond) ([]*userdomain.LoginRecord, int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.items, int64(len(r.items)), nil
}

type stubCache struct {
	mu   sync.Mutex
	data map[string]string
}

func newStubCache() *stubCache { return &stubCache{data: make(map[string]string)} }

func (c *stubCache) Get(ctx context.Context, key string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if v, ok := c.data[key]; ok {
		return v, nil
	}
	return "", errs.ErrNotFound
}

func (c *stubCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
	return nil
}

func (c *stubCache) Del(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
	return nil
}

type stubGeo struct {
	province string
	city     string
	adcode   string
}

func (g *stubGeo) LocationByIP(ctx context.Context, ip string) (userdomain.GeoInfo, error) {
	return userdomain.GeoInfo{Province: g.province, City: g.city, Adcode: g.adcode}, nil
}

func (g *stubGeo) WeatherByAdcode(ctx context.Context, adcode string) (userdomain.WeatherInfo, error) {
	return userdomain.WeatherInfo{Province: g.province, City: g.city, Weather: "晴", Temperature: "25", WindDirection: "东南风", WindPower: "3级", Humidity: "40"}, nil
}

type stubAuth struct {
	mu       sync.Mutex
	sessions map[string]string
	black    map[string]struct{}
}

func newStubAuth() *stubAuth {
	return &stubAuth{sessions: make(map[string]string), black: make(map[string]struct{})}
}

func (a *stubAuth) GenerateToken(ctx context.Context, sub auth.TokenSubject, useMultipoint bool) (auth.TokenPair, error) {
	return auth.TokenPair{Access: "a", Refresh: "r", AccessExpiresAt: 1, RefreshMaxAge: 1}, nil
}
func (a *stubAuth) ParseAccess(ctx context.Context, token string) (auth.AccessClaims, error) {
	return auth.AccessClaims{}, nil
}
func (a *stubAuth) ParseRefresh(ctx context.Context, token string) (auth.RefreshClaims, error) {
	return auth.RefreshClaims{}, nil
}
func (a *stubAuth) Blacklist(ctx context.Context, jwt string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.black[jwt] = struct{}{}
	return nil
}
func (a *stubAuth) IsBlacklisted(ctx context.Context, jwt string) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	_, ok := a.black[jwt]
	return ok
}
func (a *stubAuth) SetSession(ctx context.Context, uuidVal, refresh string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.sessions[uuidVal] = refresh
	return nil
}
func (a *stubAuth) GetSession(ctx context.Context, uuidVal string) (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if v, ok := a.sessions[uuidVal]; ok {
		return v, nil
	}
	return "", auth.ErrSessionNotFound
}
func (a *stubAuth) DelSession(ctx context.Context, uuidVal string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.sessions, uuidVal)
	return nil
}
func (a *stubAuth) Authenticate(ctx context.Context, accessToken, refreshToken string) (*auth.AuthResult, error) {
	return nil, nil
}
func (a *stubAuth) LoadAll(ctx context.Context) error { return nil }

func newTestService(t *testing.T) (*UserService, *stubUserRepo, *stubAuth, *stubCache) {
	t.Helper()
	cfg := &config.Config{}
	cfg.Upload.Size = 20
	cfg.Upload.Path = "/tmp/blog-test-uploads"
	cfg.System.UseMultipoint = true
	repo := newStubUserRepo()
	authStub := newStubAuth()
	cache := newStubCache()
	svc := NewUserService(repo, &stubLoginRepo{}, cache, &stubGeo{province: "河南省", city: "郑州市", adcode: "410100"}, authStub, cfg)
	return svc, repo, authStub, cache
}

func TestRegisterHappyPath(t *testing.T) {
	svc, repo, _, _ := newTestService(t)
	u, err := svc.Register(context.Background(), "wwy", "pwd123", "a@b.c")
	if err != nil {
		t.Fatal(err)
	}
	if u.ID == 0 || u.UUID == (uuid.UUID{}) {
		t.Fatal("register must assign id and uuid")
	}
	if !crypto.BcryptCheck("pwd123", u.Password) {
		t.Fatal("password must be bcrypt-hashed")
	}
	if u.RoleID != shared.User || u.Register != shared.Email || u.Avatar != "/image/avatar.jpg" {
		t.Fatalf("default info mismatch: %+v", u)
	}
	if _, err := repo.FindByEmail(context.Background(), "a@b.c"); err != nil {
		t.Fatal("user must be persisted")
	}
}

func TestRegisterDuplicateEmail(t *testing.T) {
	svc, _, _, _ := newTestService(t)
	_, err := svc.Register(context.Background(), "wwy", "pwd123", "a@b.c")
	if err != nil {
		t.Fatal(err)
	}
	_, err = svc.Register(context.Background(), "other", "pwd456", "a@b.c")
	if err == nil || err.Error() != "this email address is already registered, please check the information you filled in, or retrieve your password" {
		t.Fatalf("want duplicate-email message, got %v", err)
	}
}

func TestEmailLogin(t *testing.T) {
	svc, _, _, _ := newTestService(t)
	_, err := svc.Register(context.Background(), "wwy", "pwd123", "a@b.c")
	if err != nil {
		t.Fatal(err)
	}

	// 正确密码
	u, err := svc.EmailLogin(context.Background(), "a@b.c", "pwd123")
	if err != nil || u.Email != "a@b.c" {
		t.Fatalf("login failed: %v", err)
	}
	// 错误密码
	if _, err := svc.EmailLogin(context.Background(), "a@b.c", "wrong"); err == nil || err.Error() != "incorrect email or password" {
		t.Fatalf("want wrong-password message, got %v", err)
	}
	// 未知邮箱
	if _, err := svc.EmailLogin(context.Background(), "nobody@x.y", "pwd123"); !errors.Is(err, errs.ErrNotFound) {
		t.Fatalf("want not-found, got %v", err)
	}
}

func TestForgotPassword(t *testing.T) {
	svc, repo, _, _ := newTestService(t)
	_, err := svc.Register(context.Background(), "wwy", "pwd123", "a@b.c")
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.ForgotPassword(context.Background(), "a@b.c", "newpwd"); err != nil {
		t.Fatal(err)
	}
	u, _ := repo.FindByEmail(context.Background(), "a@b.c")
	if !crypto.BcryptCheck("newpwd", u.Password) {
		t.Fatal("password must be updated")
	}
}

func TestResetPassword(t *testing.T) {
	svc, _, _, _ := newTestService(t)
	_, err := svc.Register(context.Background(), "wwy", "pwd123", "a@b.c")
	if err != nil {
		t.Fatal(err)
	}
	// 原密码错误
	if err := svc.ResetPassword(context.Background(), 1, "wrong", "newpwd"); err == nil || err.Error() != "original password does not match the current account" {
		t.Fatalf("want original-password mismatch, got %v", err)
	}
	// 正确修改
	if err := svc.ResetPassword(context.Background(), 1, "pwd123", "newpwd"); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.EmailLogin(context.Background(), "a@b.c", "newpwd"); err != nil {
		t.Fatal("new password must work")
	}
}

func TestFreezeUnfreeze(t *testing.T) {
	svc, repo, authStub, _ := newTestService(t)
	_, err := svc.Register(context.Background(), "wwy", "pwd123", "a@b.c")
	if err != nil {
		t.Fatal(err)
	}
	u, _ := repo.FindByID(context.Background(), 1)

	// 预置会话（模拟已登录）
	_ = authStub.SetSession(context.Background(), u.UUID.String(), "old-refresh")

	if err := svc.Freeze(context.Background(), 1); err != nil {
		t.Fatal(err)
	}
	if !authStub.IsBlacklisted(context.Background(), "old-refresh") {
		t.Fatal("freeze must blacklist the session refresh token")
	}
	frozen, _ := repo.FindByID(context.Background(), 1)
	if !frozen.IsFrozen() {
		t.Fatal("user must be frozen after Freeze")
	}
	if err := svc.Unfreeze(context.Background(), 1); err != nil {
		t.Fatal(err)
	}
	unfrozen, _ := repo.FindByID(context.Background(), 1)
	if unfrozen.IsFrozen() {
		t.Fatal("user must be unfrozen after Unfreeze")
	}
}

func TestWeatherCache(t *testing.T) {
	svc, _, _, cache := newTestService(t)
	w, err := svc.Weather(context.Background(), "1.2.3.4")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(w, "河南省-郑州市") || !strings.Contains(w, "天气：晴") {
		t.Fatalf("weather format mismatch: %q", w)
	}
	// 命中缓存
	w2, err := svc.Weather(context.Background(), "1.2.3.4")
	if err != nil || w2 != w {
		t.Fatalf("cache hit mismatch: %q vs %q (%v)", w2, w, err)
	}
	// 缓存 Key 带 IP
	if _, err := cache.Get(context.Background(), "weather-1.2.3.4"); err != nil {
		t.Fatalf("cache key must be weather-<ip>: %v", err)
	}
}

func TestCard(t *testing.T) {
	svc, repo, _, _ := newTestService(t)
	_, err := svc.Register(context.Background(), "wwy", "pwd123", "a@b.c")
	if err != nil {
		t.Fatal(err)
	}
	u, _ := repo.FindByID(context.Background(), 1)
	card, err := svc.Card(context.Background(), u.UUID.String())
	if err != nil {
		t.Fatal(err)
	}
	if card.Username != "wwy" {
		t.Fatalf("card username mismatch: %q", card.Username)
	}
}

func TestChart(t *testing.T) {
	svc, _, _, _ := newTestService(t)
	res, err := svc.Chart(context.Background(), 7)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.DateList) != 7 || len(res.LoginData) != 7 || len(res.RegisterData) != 7 {
		t.Fatalf("chart length mismatch: %d", len(res.DateList))
	}
}
