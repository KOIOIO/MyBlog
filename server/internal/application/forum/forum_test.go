package forum

import (
	"context"
	"errors"
	"testing"

	"server/config"
	"server/internal/domain/forum"
	"server/internal/domain/shared"

	"go.uber.org/zap"
)

// stubRepo 内存论坛仓储。
type stubRepo struct {
	tags       []forum.Tag
	published  []*forum.ForumPost
	postID     uint
	list       []*forum.ForumPost
	total      int64
	detail     *forum.ForumDetail
	liked      bool
	likeCount  int
	comments   []*forum.ForumComment
	manageList []*forum.ForumPost
	manageCmt  []*forum.ManageComment
	err        error
}

func (s *stubRepo) Tags(_ context.Context) ([]forum.Tag, error) { return s.tags, s.err }
func (s *stubRepo) Publish(_ context.Context, p *forum.ForumPost) (uint, error) {
	s.published = append(s.published, p)
	return s.postID, s.err
}
func (s *stubRepo) List(_ context.Context, _ forum.ListCond) ([]*forum.ForumPost, int64, error) {
	return s.list, s.total, s.err
}
func (s *stubRepo) Detail(_ context.Context, id uint) (*forum.ForumPost, []*forum.ForumComment, error) {
	if s.detail == nil {
		return nil, nil, errors.New("not found")
	}
	return &s.detail.ForumPost, s.detail.Comments, s.err
}
func (s *stubRepo) Like(_ context.Context, _, _ uint) (bool, int, error) {
	return s.liked, s.likeCount, s.err
}
func (s *stubRepo) Comment(_ context.Context, c *forum.ForumComment) error {
	s.comments = append(s.comments, c)
	return s.err
}
func (s *stubRepo) ManageList(_ context.Context, _ forum.ManageListCond) ([]*forum.ForumPost, int64, error) {
	return s.manageList, s.total, s.err
}
func (s *stubRepo) DeletePosts(_ context.Context, _ []uint, _ uint, _ shared.RoleID) error {
	return s.err
}
func (s *stubRepo) ManageComments(_ context.Context, _ forum.ManageCommentCond) ([]*forum.ManageComment, int64, error) {
	return s.manageCmt, s.total, s.err
}
func (s *stubRepo) DeleteComments(_ context.Context, _ []uint, _ uint, _ shared.RoleID) error {
	return s.err
}

func newService(stub *stubRepo) *Service {
	return NewService(stub, &config.Config{}, zap.NewNop())
}

func TestPublishCategoryValidation(t *testing.T) {
	stub := &stubRepo{}
	svc := newService(stub)

	_, err := svc.Publish(context.Background(), &forum.ForumPost{Category: "娱乐"})
	if !errors.Is(err, forum.ErrCategoryNotAllowed) {
		t.Fatalf("want ErrCategoryNotAllowed, got %v", err)
	}
	if len(stub.published) != 0 {
		t.Fatal("repo must not be called for invalid category")
	}

	id, err := svc.Publish(context.Background(), &forum.ForumPost{Category: "技术", Tags: []string{"Go"}})
	if err != nil || id != 0 {
		t.Fatalf("want success, got %v %v", id, err)
	}
	if len(stub.published) != 1 || stub.published[0].Category != "技术" {
		t.Fatalf("publish not recorded: %+v", stub.published)
	}
}

func TestDetailBuildsForumDetail(t *testing.T) {
	stub := &stubRepo{}
	stub.detail = &forum.ForumDetail{
		ForumPost: forum.ForumPost{ID: 3, Title: "t"},
		Comments:  []*forum.ForumComment{{ID: 1}},
	}
	svc := newService(stub)
	d, err := svc.Detail(context.Background(), 3)
	if err != nil || d.ID != 3 || len(d.Comments) != 1 {
		t.Fatalf("unexpected detail: %+v %v", d, err)
	}
}

func TestCommentAndLike(t *testing.T) {
	stub := &stubRepo{liked: true, likeCount: 5}
	svc := newService(stub)

	if err := svc.Comment(context.Background(), &forum.ForumComment{PostID: 3, Content: "x"}); err != nil {
		t.Fatal(err)
	}
	if len(stub.comments) != 1 {
		t.Fatal("comment not recorded")
	}

	liked, count, err := svc.Like(context.Background(), 3, 1)
	if err != nil || !liked || count != 5 {
		t.Fatalf("unexpected like result: %v %v %v", liked, count, err)
	}
}

func TestTagsAndDeletes(t *testing.T) {
	stub := &stubRepo{}
	stub.tags = []forum.Tag{{Tag: "Go"}}
	stub.manageCmt = []*forum.ManageComment{{PostTitle: "t"}}
	svc := newService(stub)

	tags, err := svc.Tags(context.Background())
	if err != nil || len(tags) != 1 || tags[0].Tag != "Go" {
		t.Fatalf("unexpected tags: %+v %v", tags, err)
	}
	if err := svc.Delete(context.Background(), []uint{1}, 1, shared.User); err != nil {
		t.Fatal(err)
	}
	cmts, total, err := svc.ManageComments(context.Background(), forum.ManageCommentCond{})
	if err != nil || total != 0 || len(cmts) != 1 {
		t.Fatalf("unexpected manage comments: %+v %d %v", cmts, total, err)
	}
	if err := svc.DeleteComments(context.Background(), []uint{1}, 1, shared.User); err != nil {
		t.Fatal(err)
	}
}
