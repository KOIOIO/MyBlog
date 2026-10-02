package comment

import (
	"context"
	"testing"

	"server/internal/domain/comment"
	"server/internal/domain/shared"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

// stubRepo 内存评论仓储。
type stubRepo struct {
	byArticle []*comment.Comment
	byUser    []*comment.Comment
	created   []*comment.Comment
	deleted   bool
	list      []*comment.Comment
	total     int64
	err       error
}

func (s *stubRepo) ByArticle(_ context.Context, _ string) ([]*comment.Comment, error) {
	return s.byArticle, s.err
}
func (s *stubRepo) Newest(_ context.Context, _ int) ([]*comment.Comment, error) {
	return s.byArticle, s.err
}
func (s *stubRepo) Create(_ context.Context, c *comment.Comment) error {
	s.created = append(s.created, c)
	return s.err
}
func (s *stubRepo) DeleteTree(_ context.Context, ids []uint, _ uuid.UUID, _ shared.RoleID) error {
	s.deleted = len(ids) > 0
	return s.err
}
func (s *stubRepo) ByUser(_ context.Context, _ uuid.UUID) ([]*comment.Comment, error) {
	return s.byUser, s.err
}
func (s *stubRepo) List(_ context.Context, _ comment.ListCond) ([]*comment.Comment, int64, error) {
	return s.list, s.total, s.err
}

func newStub() *stubRepo { return &stubRepo{} }

func TestServiceInfoByArticleID(t *testing.T) {
	stub := newStub()
	stub.byArticle = []*comment.Comment{{ID: 1, ArticleID: "15"}}
	svc := NewService(stub, zap.NewNop())

	out, err := svc.InfoByArticleID(context.Background(), "15")
	if err != nil || len(out) != 1 || out[0].ArticleID != "15" {
		t.Fatalf("unexpected result: %v %v", out, err)
	}
}

func TestServiceInfoDedup(t *testing.T) {
	stub := newStub()
	// 用户评论树：root(id 1) 的子评论 child(id 2) 与 root 同用户 → 根列表剔除 child
	u := uuid.FromStringOrNil("fbd5364d-bb15-11f1-b230-16b7a2303b52")
	stub.byUser = []*comment.Comment{
		{ID: 1, UserUUID: u, Children: []*comment.Comment{{ID: 2, UserUUID: u}}},
		{ID: 2, UserUUID: u},
	}
	svc := NewService(stub, zap.NewNop())
	out, err := svc.Info(context.Background(), u)
	if err != nil || len(out) != 1 || out[0].ID != 1 {
		t.Fatalf("want only root 1, got %+v err=%v", out, err)
	}
}

func TestServiceCreateAndDelete(t *testing.T) {
	stub := newStub()
	svc := NewService(stub, zap.NewNop())

	c := &comment.Comment{ArticleID: "15", Content: "x"}
	if err := svc.Create(context.Background(), c); err != nil {
		t.Fatal(err)
	}
	if len(stub.created) != 1 || stub.created[0].Content != "x" {
		t.Fatalf("create not recorded: %+v", stub.created)
	}

	if err := svc.Delete(context.Background(), []uint{1}, uuid.Nil, shared.Admin); err != nil {
		t.Fatal(err)
	}
	if !stub.deleted {
		t.Fatal("delete not called")
	}

	// 空 ids 不触发删除
	if err := svc.Delete(context.Background(), nil, uuid.Nil, shared.Admin); err != nil {
		t.Fatal(err)
	}
}

func TestServiceList(t *testing.T) {
	stub := newStub()
	stub.list = []*comment.Comment{{ID: 7}}
	stub.total = 1
	svc := NewService(stub, zap.NewNop())
	list, total, err := svc.List(context.Background(), comment.ListCond{})
	if err != nil || total != 1 || len(list) != 1 {
		t.Fatalf("unexpected list: %v %d %v", list, total, err)
	}
}
