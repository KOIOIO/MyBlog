package article

import (
	"context"
	"sync"
	"testing"
	"time"

	"server/internal/domain/article"
)

// ---- 手写 stub（不引入 mock 框架） ----

type fakeEsStore struct {
	mu      sync.Mutex
	docs    map[string]*article.Article
	updates []struct {
		id  string
		doc any
	}
	likesAdded []int
	viewsAdded []int
	indexed    []string
	deleted    []string
}

func newFakeEsStore() *fakeEsStore {
	return &fakeEsStore{docs: make(map[string]*article.Article)}
}

func (s *fakeEsStore) Index(ctx context.Context, a *article.Article) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.indexed = append(s.indexed, a.Title)
	s.docs["doc-"+a.Title] = a
	return nil
}

func (s *fakeEsStore) Update(ctx context.Context, id string, doc any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.updates = append(s.updates, struct {
		id  string
		doc any
	}{id, doc})
	return nil
}

func (s *fakeEsStore) Get(ctx context.Context, id string) (*article.Article, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if a, ok := s.docs[id]; ok {
		cp := *a
		return &cp, nil
	}
	return nil, article.ErrDocumentNotFound
}

func (s *fakeEsStore) DeleteByIDs(ctx context.Context, ids []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.deleted = append(s.deleted, ids...)
	for _, id := range ids {
		delete(s.docs, id)
	}
	return nil
}

func (s *fakeEsStore) Exists(ctx context.Context, title string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.docs["doc-"+title]
	return ok, nil
}

func (s *fakeEsStore) AddLikes(ctx context.Context, id string, delta int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.likesAdded = append(s.likesAdded, delta)
	return nil
}

func (s *fakeEsStore) AddViews(ctx context.Context, id string, num int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.viewsAdded = append(s.viewsAdded, num)
	return nil
}

func (s *fakeEsStore) Search(ctx context.Context, spec article.SearchSpec) (article.SearchResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// 校验查询规格后返回空命中
	if spec.ListMode && spec.Title != nil && *spec.Title != "" {
		if _, ok := s.docs["doc-"+*spec.Title]; ok {
			return article.SearchResult{Hits: []any{}, Total: 1}, nil
		}
	}
	return article.SearchResult{Hits: []any{}, Total: 0}, nil
}

type stubArticleRepo struct {
	mu              sync.Mutex
	categories      []article.ArticleCategory
	tags            []article.BlogTag
	tagExists       map[string]bool
	likes           map[string]bool // key: userID:articleID
	createdCalls    []string
	updatedCalls    []string
	deletedCalls    []string
	likeTxCalls     int
	categoriesCalls int
	tagsCalls       int
	checkTagsCalls  int
	likesListCalls  int
}

func newStubArticleRepo() *stubArticleRepo {
	return &stubArticleRepo{
		tagExists: map[string]bool{"Gin": true, "Go": true},
		likes:     make(map[string]bool),
	}
}

func (r *stubArticleRepo) Categories(ctx context.Context) ([]article.ArticleCategory, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.categoriesCalls++
	return r.categories, nil
}

func (r *stubArticleRepo) Tags(ctx context.Context) ([]article.BlogTag, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tagsCalls++
	return r.tags, nil
}

func (r *stubArticleRepo) CheckTagsExist(ctx context.Context, tags []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.checkTagsCalls++
	for _, tag := range tags {
		if !r.tagExists[tag] {
			return article.ErrTagNotExist
		}
	}
	return nil
}

func (r *stubArticleRepo) LikeTx(ctx context.Context, userID uint, articleID string) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.likeTxCalls++
	key := userIDKey(userID, articleID)
	if !r.likes[key] {
		r.likes[key] = true
		return 1, nil
	}
	delete(r.likes, key)
	return -1, nil
}

func (r *stubArticleRepo) IsLike(ctx context.Context, userID uint, articleID string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.likes[userIDKey(userID, articleID)], nil
}

func (r *stubArticleRepo) LikesList(ctx context.Context, userID uint, page, pageSize int) ([]article.ArticleLike, int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.likesListCalls++
	return nil, 0, nil
}

func (r *stubArticleRepo) CreateWithCounts(ctx context.Context, a *article.Article, cover string, illustrations []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.createdCalls = append(r.createdCalls, a.Title)
	return nil
}

func (r *stubArticleRepo) UpdateWithCounts(ctx context.Context, a *article.Article, old *article.Article, newCover string, newIllustrations []string, addedIllustrations, removedIllustrations []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.updatedCalls = append(r.updatedCalls, a.Title)
	return nil
}

func (r *stubArticleRepo) DeleteWithCounts(ctx context.Context, a *article.Article, cover string, illustrations []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.deletedCalls = append(r.deletedCalls, a.Title)
	return nil
}

func userIDKey(userID uint, articleID string) string {
	return string(rune(userID)) + ":" + articleID
}

type stubViewCounter struct {
	mu      sync.Mutex
	count   map[string]int
	cleared bool
}

func (v *stubViewCounter) Set(ctx context.Context, id string) error {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.count[id]++
	return nil
}

func (v *stubViewCounter) GetInfo(ctx context.Context) map[string]int {
	v.mu.Lock()
	defer v.mu.Unlock()
	out := make(map[string]int, len(v.count))
	for k, n := range v.count {
		out[k] = n
	}
	return out
}

func (v *stubViewCounter) Clear(ctx context.Context) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.cleared = true
	v.count = map[string]int{}
}

func newTestService() (*Service, *stubArticleRepo, *fakeEsStore, *stubViewCounter) {
	repo := newStubArticleRepo()
	esStore := newFakeEsStore()
	views := &stubViewCounter{count: map[string]int{}}
	svc := NewService(repo, esStore, views)
	return svc, repo, esStore, views
}

func TestCreateValidationOrder(t *testing.T) {
	svc, repo, esStore, _ := newTestService()

	// 标题已存在 → "the article already exists"
	esStore.docs["doc-dup"] = &article.Article{Title: "dup"}
	if err := svc.Create(context.Background(), &article.Article{Title: "dup"}); err == nil || err.Error() != "the article already exists" {
		t.Fatalf("want duplicate title error, got %v", err)
	}
	// 分类不合法 → "分类只能是 技术 或 生活"
	if err := svc.Create(context.Background(), &article.Article{Title: "t1", Category: "科技"}); err == nil || err.Error() != "分类只能是 技术 或 生活" {
		t.Fatalf("want category error, got %v", err)
	}
	// 标签不存在 → "标签不存在"
	if err := svc.Create(context.Background(), &article.Article{Title: "t1", Category: "技术", Tags: []string{"Nope"}}); err == nil || err.Error() != "标签不存在" {
		t.Fatalf("want tag error, got %v", err)
	}
	// 校验失败时不触发任何计数联动
	if len(repo.createdCalls) != 0 || len(esStore.indexed) != 0 {
		t.Fatalf("no linkage expected on validation failure: %v %v", repo.createdCalls, esStore.indexed)
	}
}

func TestCreateHappyPath(t *testing.T) {
	svc, repo, esStore, _ := newTestService()
	a := &article.Article{Title: "新文章", Category: "技术", Tags: []string{"Gin"}, Cover: "/uploads/c.png",
		Abstract: "ab", Content: "正文 ![i](/uploads/1.png)"}
	if err := svc.Create(context.Background(), a); err != nil {
		t.Fatal(err)
	}
	if len(repo.createdCalls) != 1 || repo.createdCalls[0] != "新文章" {
		t.Fatalf("repo linkage mismatch: %v", repo.createdCalls)
	}
	if len(esStore.indexed) != 1 {
		t.Fatal("article must be indexed")
	}
	if a.CreatedAt == "" || a.UpdatedAt == "" || a.Keyword != "新文章" {
		t.Fatalf("timestamps/keyword not filled: %+v", a)
	}
}

func TestDeleteLinkagePerArticle(t *testing.T) {
	svc, repo, esStore, _ := newTestService()
	esStore.docs["4"] = &article.Article{Title: "A", Category: "技术", Content: "![x](/uploads/1.png)"}
	if err := svc.Delete(context.Background(), []string{"4"}); err != nil {
		t.Fatal(err)
	}
	if len(repo.deletedCalls) != 1 || repo.deletedCalls[0] != "A" {
		t.Fatalf("delete linkage mismatch: %v", repo.deletedCalls)
	}
	if len(esStore.deleted) != 1 || esStore.deleted[0] != "4" {
		t.Fatalf("es delete mismatch: %v", esStore.deleted)
	}
	// 空 ID 列表不执行任何操作
	if err := svc.Delete(context.Background(), []string{}); err != nil {
		t.Fatal(err)
	}
	if len(repo.deletedCalls) != 1 {
		t.Fatalf("empty ids must not trigger linkage: %v", repo.deletedCalls)
	}
}

func TestUpdateFillsKeywordAndDiff(t *testing.T) {
	svc, repo, esStore, _ := newTestService()
	esStore.docs["4"] = &article.Article{Title: "旧标题", Category: "生活", Content: "![a](/uploads/old.png)"}
	a := &article.Article{Title: "新标题", Category: "技术", Tags: []string{"Gin"}, Cover: "/uploads/new.png",
		Abstract: "ab", Content: "![b](/uploads/new.png)"}
	if err := svc.Update(context.Background(), "4", a); err != nil {
		t.Fatal(err)
	}
	if a.UpdatedAt == "" || a.Keyword != "新标题" {
		t.Fatalf("updated_at/keyword not filled: %+v", a)
	}
	if len(repo.updatedCalls) != 1 {
		t.Fatalf("update linkage mismatch: %v", repo.updatedCalls)
	}
	// ES 更新以整篇文档提交（Doc=article）
	if len(esStore.updates) != 1 || esStore.updates[0].id != "4" {
		t.Fatalf("es update mismatch: %v", esStore.updates)
	}
}

func TestLikeToggles(t *testing.T) {
	svc, repo, esStore, _ := newTestService()
	if err := svc.Like(context.Background(), 1, "4"); err != nil {
		t.Fatal(err)
	}
	if !repo.likes[userIDKey(1, "4")] {
		t.Fatal("like must be set")
	}
	if len(esStore.likesAdded) != 1 || esStore.likesAdded[0] != 1 {
		t.Fatalf("es likes delta mismatch: %v", esStore.likesAdded)
	}
	if err := svc.Like(context.Background(), 1, "4"); err != nil {
		t.Fatal(err)
	}
	if repo.likes[userIDKey(1, "4")] {
		t.Fatal("second like must cancel")
	}
	if len(esStore.likesAdded) != 2 || esStore.likesAdded[1] != -1 {
		t.Fatalf("es likes delta mismatch: %v", esStore.likesAdded)
	}
	isLike, err := svc.IsLike(context.Background(), 1, "4")
	if err != nil || isLike {
		t.Fatalf("isLike mismatch: %v %v", isLike, err)
	}
}

func TestSetTop(t *testing.T) {
	svc, _, esStore, _ := newTestService()
	if err := svc.SetTop(context.Background(), "4", true); err != nil {
		t.Fatal(err)
	}
	if len(esStore.updates) != 1 {
		t.Fatal("setTop must call es update")
	}
	doc := esStore.updates[0].doc.(map[string]any)
	if doc["is_top"] != 1 {
		t.Fatalf("is_top must be 1: %v", doc)
	}
	if err := svc.SetTop(context.Background(), "4", false); err != nil {
		t.Fatal(err)
	}
	doc = esStore.updates[1].doc.(map[string]any)
	if doc["is_top"] != 0 {
		t.Fatalf("is_top must be 0: %v", doc)
	}
}

func TestSearchSpecPassThrough(t *testing.T) {
	svc, _, esStore, _ := newTestService()
	esStore.docs["doc-Gin"] = &article.Article{Title: "Gin"}
	// ListMode 由 List 设置
	_, total, err := svc.List(context.Background(), article.SearchSpec{Title: ptr("Gin")})
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 {
		t.Fatalf("list must route to list-mode search, total=%d", total)
	}
	// Search 保持前台模式
	_, total, err = svc.Search(context.Background(), article.SearchSpec{Query: "Gin", Page: 1, PageSize: 10})
	if err != nil {
		t.Fatal(err)
	}
	if total != 0 {
		t.Fatalf("search total mismatch: %d", total)
	}
}

func ptr(s string) *string { return &s }

func TestSyncViewsSyncsAllAndClears(t *testing.T) {
	svc, _, esStore, views := newTestService()
	views.count = map[string]int{"4": 3, "5": 1}
	if err := svc.SyncViews(context.Background()); err != nil {
		t.Fatal(err)
	}
	seen := map[int]bool{}
	for _, v := range esStore.viewsAdded {
		seen[v] = true
	}
	if len(esStore.viewsAdded) != 2 || !seen[3] || !seen[1] {
		t.Fatalf("views must be synced for both articles: %v", esStore.viewsAdded)
	}
	if !views.cleared {
		t.Fatal("views must be cleared after sync")
	}
	if len(views.count) != 0 {
		t.Fatalf("counter must be cleared: %v", views.count)
	}
}

func TestInfoByIDAsyncView(t *testing.T) {
	svc, _, esStore, views := newTestService()
	esStore.docs["4"] = &article.Article{Title: "A"}
	a, err := svc.InfoByID(context.Background(), "4")
	if err != nil {
		t.Fatal(err)
	}
	if a.Title != "A" {
		t.Fatalf("info mismatch: %v", a)
	}
	// 异步浏览量累加（goroutine 调度后生效）
	waitViewCounter(t, views, 1)
}

func waitViewCounter(t *testing.T, views *stubViewCounter, want int) {
	t.Helper()
	for i := 0; i < 100; i++ {
		views.mu.Lock()
		n := views.count["4"]
		views.mu.Unlock()
		if n == want {
			return
		}
		time.Sleep(2 * time.Millisecond)
	}
	t.Fatal("view counter not incremented")
}
