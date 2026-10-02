package auth

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	jwtutil "server/internal/common/jwt"
	"server/internal/domain/auth"
	"server/internal/domain/shared"
	userdomain "server/internal/domain/user"

	"github.com/gofrs/uuid"
)

const (
	accessSecret  = "test-access-secret"
	refreshSecret = "test-refresh-secret"
	issuer        = "test-issuer"
)

// ---- 手写 stub（不引入 mock 框架） ----

type stubSessionStore struct {
	mu   sync.Mutex
	data map[string]sessionEntry
}

type sessionEntry struct {
	value string
	ttl   time.Time
}

func newStubSessionStore() *stubSessionStore {
	return &stubSessionStore{data: make(map[string]sessionEntry)}
}

func (s *stubSessionStore) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = sessionEntry{value: value, ttl: time.Now().Add(ttl)}
	return nil
}

func (s *stubSessionStore) Get(ctx context.Context, key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.data[key]
	if !ok || time.Now().After(e.ttl) {
		return "", auth.ErrSessionNotFound
	}
	return e.value, nil
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

func newStubBlacklist() *stubBlacklist {
	return &stubBlacklist{set: make(map[string]struct{})}
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
	return nil, 0, nil
}
func (r *stubUserRepo) FindIDByUUID(ctx context.Context, uuidVal uuid.UUID) (uint, error) {
	return 0, errNotFound
}

var errNotFound = errors.New("record not found")

func newTestJWT(t *testing.T, accessExpiry, refreshExpiry string) *jwtutil.JWT {
	t.Helper()
	j, err := jwtutil.New(accessSecret, refreshSecret, accessExpiry, refreshExpiry, issuer)
	if err != nil {
		t.Fatal(err)
	}
	return j
}

func newTestService(t *testing.T, repo userdomain.UserRepository, accessExpiry string) (*AuthService, *stubSessionStore, *stubBlacklist) {
	t.Helper()
	j := newTestJWT(t, accessExpiry, "1h")
	sessions := newStubSessionStore()
	blacklist := newStubBlacklist()
	svc := NewAuthService(j, sessions, blacklist, repo, time.Hour)
	return svc, sessions, blacklist
}

func TestGenerateTokenNoMultipoint(t *testing.T) {
	uid := uuid.Must(uuid.NewV4())
	repo := newStubUserRepo(&userdomain.User{ID: 1, UUID: uid, Email: "a@b.c"})
	svc, sessions, _ := newTestService(t, repo, "2h")

	pair, err := svc.GenerateToken(context.Background(), auth.TokenSubject{UserID: 1, UUID: uid, RoleID: shared.User}, false)
	if err != nil {
		t.Fatal(err)
	}
	if pair.Access == "" || pair.Refresh == "" {
		t.Fatal("tokens must be non-empty")
	}
	if pair.AccessExpiresAt <= 0 || pair.RefreshMaxAge <= 0 {
		t.Fatalf("expiry values must be positive: %d %d", pair.AccessExpiresAt, pair.RefreshMaxAge)
	}
	// 单点登录关闭时不应写入会话
	if _, err := sessions.Get(context.Background(), uid.String()); !errors.Is(err, auth.ErrSessionNotFound) {
		t.Fatal("single-point login must not create session")
	}
}

func TestGenerateTokenMultipointBlacklistsOldRefresh(t *testing.T) {
	uid := uuid.Must(uuid.NewV4())
	repo := newStubUserRepo(&userdomain.User{ID: 1, UUID: uid, Email: "a@b.c"})
	svc, _, blacklist := newTestService(t, repo, "2h")
	sub := auth.TokenSubject{UserID: 1, UUID: uid, RoleID: shared.User}

	first, err := svc.GenerateToken(context.Background(), sub, true)
	if err != nil {
		t.Fatal(err)
	}
	// 第一次登录：仅写会话，不产生黑名单
	if blacklist.Contains(context.Background(), first.Refresh) {
		t.Fatal("first login must not blacklist its own refresh token")
	}

	// 注意：JWT 时间戳精度为秒，同一秒内两次签发会产生相同令牌（重构前亦如此）。
	// 等待 2 秒确保两次签发的 refresh 令牌可区分。
	time.Sleep(2 * time.Second)
	second, err := svc.GenerateToken(context.Background(), sub, true)
	if err != nil {
		t.Fatal(err)
	}
	if first.Refresh == second.Refresh {
		t.Fatal("tokens must differ after waiting 2s")
	}
	// 第二次登录：旧 refresh 应被加入黑名单（多地点登录拦截）
	if !blacklist.Contains(context.Background(), first.Refresh) {
		t.Fatal("second login must blacklist the old refresh token")
	}
	if blacklist.Contains(context.Background(), second.Refresh) {
		t.Fatal("new refresh token must not be blacklisted")
	}
}

func TestAuthenticateValidAccess(t *testing.T) {
	uid := uuid.Must(uuid.NewV4())
	repo := newStubUserRepo(&userdomain.User{ID: 1, UUID: uid, Email: "a@b.c", RoleID: shared.Admin})
	svc, _, _ := newTestService(t, repo, "2h")

	pair, _ := svc.GenerateToken(context.Background(), auth.TokenSubject{UserID: 1, UUID: uid, RoleID: shared.Admin}, false)
	res, err := svc.Authenticate(context.Background(), pair.Access, pair.Refresh)
	if err != nil {
		t.Fatal(err)
	}
	if res.NewAccessToken != "" {
		t.Fatal("valid access token must not trigger refresh")
	}
	if res.Claims.UserID != 1 || res.Claims.UUID != uid || res.Claims.RoleID != shared.Admin {
		t.Fatalf("claims mismatch: %+v", res.Claims)
	}
}

func TestAuthenticateExpiredAccessRefreshes(t *testing.T) {
	uid := uuid.Must(uuid.NewV4())
	repo := newStubUserRepo(&userdomain.User{ID: 1, UUID: uid, Email: "a@b.c"})
	// 1 秒过期的 access，用于触发续期
	svc, _, _ := newTestService(t, repo, "1s")

	pair, _ := svc.GenerateToken(context.Background(), auth.TokenSubject{UserID: 1, UUID: uid, RoleID: shared.User}, false)
	time.Sleep(1100 * time.Millisecond)

	res, err := svc.Authenticate(context.Background(), pair.Access, pair.Refresh)
	if err != nil {
		t.Fatal(err)
	}
	if res.NewAccessToken == "" {
		t.Fatal("expired access with valid refresh must issue a new access token")
	}
	if res.Claims.UserID != 1 {
		t.Fatalf("claims user mismatch: %+v", res.Claims)
	}
}

func TestAuthenticateRefreshBlacklisted(t *testing.T) {
	uid := uuid.Must(uuid.NewV4())
	repo := newStubUserRepo(&userdomain.User{ID: 1, UUID: uid, Email: "a@b.c"})
	svc, _, blacklist := newTestService(t, repo, "2h")

	pair, _ := svc.GenerateToken(context.Background(), auth.TokenSubject{UserID: 1, UUID: uid, RoleID: shared.User}, false)
	_ = blacklist.Add(context.Background(), pair.Refresh)

	if _, err := svc.Authenticate(context.Background(), pair.Access, pair.Refresh); !errors.Is(err, auth.ErrRefreshBlacklisted) {
		t.Fatalf("want ErrRefreshBlacklisted, got %v", err)
	}
}

func TestAuthenticateRefreshInvalid(t *testing.T) {
	uid := uuid.Must(uuid.NewV4())
	repo := newStubUserRepo(&userdomain.User{ID: 1, UUID: uid, Email: "a@b.c"})
	svc, _, _ := newTestService(t, repo, "2h")

	pair, _ := svc.GenerateToken(context.Background(), auth.TokenSubject{UserID: 1, UUID: uid, RoleID: shared.User}, false)
	// 空 access + 非法 refresh → 续期路径解析失败
	if _, err := svc.Authenticate(context.Background(), "", "garbage-refresh"); !errors.Is(err, auth.ErrRefreshInvalid) {
		t.Fatalf("want ErrRefreshInvalid, got %v", err)
	}
	// 有效 access + 非法 refresh → access 直接通过，不触发续期
	if _, err := svc.Authenticate(context.Background(), pair.Access, "garbage-refresh"); err != nil {
		t.Fatalf("valid access must pass regardless of refresh, got %v", err)
	}
}

func TestAuthenticateUserNotFound(t *testing.T) {
	uid := uuid.Must(uuid.NewV4())
	repo := newStubUserRepo(&userdomain.User{ID: 1, UUID: uid, Email: "a@b.c"})
	svc, _, _ := newTestService(t, repo, "1s")

	pair, _ := svc.GenerateToken(context.Background(), auth.TokenSubject{UserID: 1, UUID: uid, RoleID: shared.User}, false)
	time.Sleep(1100 * time.Millisecond)

	// 空用户表场景：access 过期走续期时查不到用户
	emptySvc, _, _ := newTestService(t, newStubUserRepo(), "2h")
	if _, err := emptySvc.Authenticate(context.Background(), "", pair.Refresh); !errors.Is(err, auth.ErrUserNotFound) {
		t.Fatalf("want ErrUserNotFound, got %v", err)
	}
}

func TestAuthenticateInvalidAccess(t *testing.T) {
	uid := uuid.Must(uuid.NewV4())
	repo := newStubUserRepo(&userdomain.User{ID: 1, UUID: uid, Email: "a@b.c"})
	svc, _, _ := newTestService(t, repo, "2h")

	// 格式非法且非空的 access token
	if _, err := svc.Authenticate(context.Background(), "not-a-jwt", "garbage-refresh"); !errors.Is(err, auth.ErrAccessInvalid) {
		t.Fatalf("want ErrAccessInvalid, got %v", err)
	}
}
