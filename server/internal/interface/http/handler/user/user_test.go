package user_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"server/config"
	authapp "server/internal/application/auth"
	userapp "server/internal/application/user"
	jwtutil "server/internal/common/jwt"
	"server/internal/domain/auth"
	"server/internal/domain/shared"
	userdomain "server/internal/domain/user"
	userhandler "server/internal/interface/http/handler/user"
	"server/internal/interface/http/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

// ---- 手写 stub（不引入 mock 框架） ----

var errNotFound = errors.New("record not found")

type stubUserRepo struct {
	mu      sync.Mutex
	byID    map[uint]*userdomain.User
	byEmail map[string]*userdomain.User
}

func newStubUserRepo(users ...*userdomain.User) *stubUserRepo {
	r := &stubUserRepo{byID: make(map[uint]*userdomain.User), byEmail: make(map[string]*userdomain.User)}
	for _, u := range users {
		r.byID[u.ID] = u
		r.byEmail[u.Email] = u
	}
	return r
}

func (r *stubUserRepo) FindByEmail(ctx context.Context, email string) (*userdomain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if u, ok := r.byEmail[email]; ok {
		return u, nil
	}
	return nil, errNotFound
}

func (r *stubUserRepo) FindByUUID(ctx context.Context, uid uuid.UUID) (*userdomain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, u := range r.byID {
		if u.UUID == uid {
			return u, nil
		}
	}
	return nil, errNotFound
}

func (r *stubUserRepo) FindByID(ctx context.Context, id uint) (*userdomain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if u, ok := r.byID[id]; ok {
		return u, nil
	}
	return nil, errNotFound
}

func (r *stubUserRepo) Create(ctx context.Context, u *userdomain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
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
		return nil, errNotFound
	}
	prev := *u
	u.SetFrozen(frozen)
	return &prev, nil
}

func (r *stubUserRepo) UpdateInfo(ctx context.Context, id uint, username, address, signature string) error {
	return nil
}

func (r *stubUserRepo) UpdateAvatar(ctx context.Context, id uint, avatarURL string) (*userdomain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, ok := r.byID[id]
	if !ok {
		return nil, errNotFound
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
		out = append(out, u)
	}
	return out, int64(len(out)), nil
}

func (r *stubUserRepo) FindIDByUUID(ctx context.Context, uuidVal uuid.UUID) (uint, error) {
	return 0, errNotFound
}

type stubLoginRepo struct {
	mu    sync.Mutex
	items []*userdomain.LoginRecord
}

func (r *stubLoginRepo) Create(ctx context.Context, rec *userdomain.LoginRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items = append(r.items, rec)
	return nil
}

func (r *stubLoginRepo) CountByDate(ctx context.Context, days int) (map[string]int, error) {
	return map[string]int{}, nil
}

func (r *stubLoginRepo) Page(ctx context.Context, cond userdomain.LoginListCond) ([]*userdomain.LoginRecord, int64, error) {
	return nil, 0, nil
}

type stubCache struct {
	mu   sync.Mutex
	data map[string]string
}

func (c *stubCache) Get(ctx context.Context, key string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if v, ok := c.data[key]; ok {
		return v, nil
	}
	return "", errNotFound
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

type stubGeo struct{}

func (g *stubGeo) LocationByIP(ctx context.Context, ip string) (userdomain.GeoInfo, error) {
	return userdomain.GeoInfo{Province: "河南省", City: "郑州市", Adcode: "410100"}, nil
}

func (g *stubGeo) WeatherByAdcode(ctx context.Context, adcode string) (userdomain.WeatherInfo, error) {
	return userdomain.WeatherInfo{Province: "河南省", City: "郑州市", Weather: "晴", Temperature: "25"}, nil
}

type stubSessionStore struct {
	mu   sync.Mutex
	data map[string]string
}

func (s *stubSessionStore) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
	return nil
}

func (s *stubSessionStore) Get(ctx context.Context, key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.data[key]; ok {
		return v, nil
	}
	return "", auth.ErrSessionNotFound
}

func (s *stubSessionStore) Del(ctx context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
	return nil
}

type stubBlacklist struct {
	mu  sync.Mutex
	set map[string]struct{}
}

func (b *stubBlacklist) Add(ctx context.Context, jwt string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.set[jwt] = struct{}{}
	return nil
}

func (b *stubBlacklist) Contains(ctx context.Context, jwt string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	_, ok := b.set[jwt]
	return ok
}

func (b *stubBlacklist) LoadAll(ctx context.Context) error { return nil }

// ---- 测试装配 ----

var (
	testAdminUUID = uuid.Must(uuid.FromString("fbd5364d-bb15-11f1-b230-16b7a2303b52"))
	testUserUUID  = uuid.Must(uuid.FromString("aaaaaaaa-bb15-11f1-b230-16b7a2303b52"))
)

func newSmokeEnv(t *testing.T) (*gin.Engine, *userapp.UserService, *authapp.AuthService) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	log := zap.NewNop()

	cfg := &config.Config{}
	cfg.System.UseMultipoint = false
	cfg.System.OssType = "local"
	cfg.Upload.Size = 20
	cfg.Upload.Path = "/tmp/blog-test-uploads"

	admin := &userdomain.User{ID: 1, UUID: testAdminUUID, Username: "wwy", Email: "admin@blog.dev", RoleID: shared.Admin, Register: shared.Email}
	normal := &userdomain.User{ID: 2, UUID: testUserUUID, Username: "guest", Email: "guest@blog.dev", RoleID: shared.User, Register: shared.Email}
	userRepo := newStubUserRepo(admin, normal)

	jwtUtil, err := jwtutil.New("smoke-access", "smoke-refresh", "2h", "2h", "smoke")
	if err != nil {
		t.Fatal(err)
	}
	sessionStore := &stubSessionStore{data: make(map[string]string)}
	blacklist := &stubBlacklist{set: make(map[string]struct{})}
	authApp := authapp.NewAuthService(jwtUtil, sessionStore, blacklist, userRepo, 2*time.Hour)

	userApp := userapp.NewUserService(userRepo, &stubLoginRepo{}, &stubCache{data: make(map[string]string)}, &stubGeo{}, authApp, cfg)

	h := userhandler.NewHandler(userApp, authApp, base64Captcha.DefaultMemStore, cfg, log)

	engine := gin.New()
	engine.Use(middleware.GinLogger(log), middleware.GinRecovery(log, true))
	cookieStore := cookie.NewStore([]byte("smoke-secret"))
	engine.Use(sessions.Sessions("session", cookieStore))

	const prefix = "api"
	publicGroup := engine.Group(prefix)
	privateGroup := engine.Group(prefix)
	privateGroup.Use(middleware.JWTAuth(authApp))
	loginGroup := engine.Group(prefix)
	loginGroup.Use(middleware.LoginRecord(&stubGeo{}, &stubLoginRepo{}, log))
	adminGroup := engine.Group(prefix)
	adminGroup.Use(middleware.JWTAuth(authApp)).Use(middleware.AdminAuth())

	userRouter := privateGroup.Group("user")
	userPublicRouter := publicGroup.Group("user")
	userLoginRouter := loginGroup.Group("user")
	userAdminRouter := adminGroup.Group("user")

	userRouter.POST("logout", h.Logout)
	userRouter.PUT("resetPassword", h.UserResetPassword)
	userRouter.GET("info", h.UserInfo)
	userRouter.PUT("changeInfo", h.UserChangeInfo)
	userRouter.POST("avatar", h.UploadAvatar)
	userRouter.GET("weather", h.UserWeather)
	userRouter.GET("chart", h.UserChart)

	userPublicRouter.POST("forgotPassword", h.ForgotPassword)
	userPublicRouter.GET("card", h.UserCard)

	userLoginRouter.POST("register", h.Register)
	userLoginRouter.POST("login", h.Login)

	userAdminRouter.GET("list", h.UserList)
	userAdminRouter.PUT("freeze", h.UserFreeze)
	userAdminRouter.PUT("unfreeze", h.UserUnfreeze)
	userAdminRouter.GET("loginList", h.UserLoginList)

	return engine, userApp, authApp
}

func doRequest(engine *gin.Engine, method, path string, body string, headers map[string]string) *httptest.ResponseRecorder {
	var reader *strings.Reader
	if body == "" {
		reader = strings.NewReader("")
	} else {
		reader = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, path, reader)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	return w
}

func decodeResp(t *testing.T, w *httptest.ResponseRecorder) (code float64, data interface{}, msg string) {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		t.Fatalf("bad json body %q: %v", w.Body.String(), err)
	}
	if v, ok := m["code"].(float64); ok {
		code = v
	}
	msg, _ = m["msg"].(string)
	data = m["data"]
	return
}

// ---- 冒烟测试 ----

func TestSmokeUserInfoNoToken(t *testing.T) {
	engine, _, _ := newSmokeEnv(t)
	w := doRequest(engine, http.MethodGet, "/api/user/info", "", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	code, data, msg := decodeResp(t, w)
	if code != 7 {
		t.Fatalf("want code 7, got %v", code)
	}
	// 无 token 时（access 与 refresh 均为空）：与重构前一致返回 refresh 失效文案（快照 private_noauth_1 锁定）。
	if msg != "Refresh token expired or invalid" {
		t.Fatalf("want NoAuth msg, got %q", msg)
	}
	if reload, ok := data.(map[string]interface{})["reload"].(bool); !ok || !reload {
		t.Fatalf("want data.reload=true, got %v", data)
	}
}

func TestSmokeLoginWrongCaptcha(t *testing.T) {
	engine, _, _ := newSmokeEnv(t)
	// 字段满足校验（password>=8、captcha len=6），但验证码不匹配存储
	body := `{"email":"admin@blog.dev","password":"password123","captcha":"000000","captcha_id":"none"}`
	w := doRequest(engine, http.MethodPost, "/api/user/login", body, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	code, _, msg := decodeResp(t, w)
	if code == 0 || msg != "Incorrect verification code" {
		t.Fatalf("want wrong-captcha msg, got code=%v msg=%q", code, msg)
	}
}

func TestSmokeRegisterWithoutSession(t *testing.T) {
	engine, _, _ := newSmokeEnv(t)
	body := `{"username":"newbie","password":"password123","email":"new@blog.dev","verification_code":"123456"}`
	w := doRequest(engine, http.MethodPost, "/api/user/register", body, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	_, _, msg := decodeResp(t, w)
	if msg != "This email doesn't match the email to be verified" {
		t.Fatalf("want no-session msg, got %q", msg)
	}
}

func TestSmokeUserCard(t *testing.T) {
	engine, _, _ := newSmokeEnv(t)
	w := doRequest(engine, http.MethodGet, "/api/user/card?uuid="+testAdminUUID.String(), "", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	code, data, _ := decodeResp(t, w)
	if code != 0 {
		t.Fatalf("want code 0, got %v", code)
	}
	d, ok := data.(map[string]interface{})
	if !ok || d["username"] != "wwy" {
		t.Fatalf("want card username wwy, got %v", data)
	}
}

func TestSmokeInfoWithToken(t *testing.T) {
	engine, _, authApp := newSmokeEnv(t)
	pair, err := authApp.GenerateToken(context.Background(), auth.TokenSubject{UserID: 2, UUID: testUserUUID, RoleID: shared.User}, false)
	if err != nil {
		t.Fatal(err)
	}
	w := doRequest(engine, http.MethodGet, "/api/user/info", "", map[string]string{"x-access-token": pair.Access})
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	code, data, _ := decodeResp(t, w)
	if code != 0 {
		t.Fatalf("want code 0, got %v", code)
	}
	d, ok := data.(map[string]interface{})
	if !ok || d["username"] != "guest" || d["id"] != float64(2) {
		t.Fatalf("want guest info, got %v", data)
	}
}

func TestSmokeLogoutClearsCookieAndBlacklists(t *testing.T) {
	engine, _, authApp := newSmokeEnv(t)
	pair, err := authApp.GenerateToken(context.Background(), auth.TokenSubject{UserID: 2, UUID: testUserUUID, RoleID: shared.User}, false)
	if err != nil {
		t.Fatal(err)
	}
	w := doRequest(engine, http.MethodPost, "/api/user/logout", "", map[string]string{
		"x-access-token": pair.Access,
		"Cookie":         "x-refresh-token=" + pair.Refresh,
	})
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	_, _, msg := decodeResp(t, w)
	if msg != "Successful logout" {
		t.Fatalf("want logout msg, got %q", msg)
	}
	if !strings.Contains(w.Header().Get("Set-Cookie"), "x-refresh-token=;") {
		t.Fatalf("refresh cookie must be cleared, got %q", w.Header().Get("Set-Cookie"))
	}
	if !authApp.IsBlacklisted(context.Background(), pair.Refresh) {
		t.Fatal("refresh token must be blacklisted after logout")
	}
}

func TestSmokeAdminListForbiddenForNormalUser(t *testing.T) {
	engine, _, authApp := newSmokeEnv(t)
	pair, err := authApp.GenerateToken(context.Background(), auth.TokenSubject{UserID: 2, UUID: testUserUUID, RoleID: shared.User}, false)
	if err != nil {
		t.Fatal(err)
	}
	w := doRequest(engine, http.MethodGet, "/api/user/list", "", map[string]string{"x-access-token": pair.Access})
	if w.Code != http.StatusForbidden {
		t.Fatalf("want 403, got %d", w.Code)
	}
	code, _, msg := decodeResp(t, w)
	if code != 7 || msg != "Access denied. Adimin privileges are required" {
		t.Fatalf("want forbidden msg, got code=%v msg=%q", code, msg)
	}
}

func TestSmokeAdminListForAdmin(t *testing.T) {
	engine, _, authApp := newSmokeEnv(t)
	pair, err := authApp.GenerateToken(context.Background(), auth.TokenSubject{UserID: 1, UUID: testAdminUUID, RoleID: shared.Admin}, false)
	if err != nil {
		t.Fatal(err)
	}
	w := doRequest(engine, http.MethodGet, "/api/user/list", "", map[string]string{"x-access-token": pair.Access})
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	code, _, _ := decodeResp(t, w)
	if code != 0 {
		t.Fatalf("want code 0, got %v", code)
	}
}
