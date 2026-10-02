package feedback

import (
	"context"
	"testing"

	"server/internal/domain/feedback"

	"github.com/gofrs/uuid"
)

type stubRepo struct {
	newest  []*feedback.Feedback
	created []*feedback.Feedback
	info    []*feedback.Feedback
	deleted []uint
	replied []string
	list    []*feedback.Feedback
	total   int64
}

func (s *stubRepo) Newest(_ context.Context, _ int) ([]*feedback.Feedback, error) {
	return s.newest, nil
}
func (s *stubRepo) Create(_ context.Context, f *feedback.Feedback) error {
	s.created = append(s.created, f)
	return nil
}
func (s *stubRepo) Info(_ context.Context, _ uuid.UUID) ([]*feedback.Feedback, error) {
	return s.info, nil
}
func (s *stubRepo) Delete(_ context.Context, ids []uint) error {
	s.deleted = append(s.deleted, ids...)
	return nil
}
func (s *stubRepo) Reply(_ context.Context, _ uint, reply string) error {
	s.replied = append(s.replied, reply)
	return nil
}
func (s *stubRepo) List(_ context.Context, _, _ int) ([]*feedback.Feedback, int64, error) {
	return s.list, s.total, nil
}

func TestServiceLifecycle(t *testing.T) {
	repo := &stubRepo{}
	svc := NewService(repo)
	u := uuid.FromStringOrNil("fbd5364d-bb15-11f1-b230-16b7a2303b52")

	if err := svc.Create(context.Background(), &feedback.Feedback{UserUUID: u, Content: "c"}); err != nil {
		t.Fatal(err)
	}
	if len(repo.created) != 1 || repo.created[0].Content != "c" {
		t.Fatalf("create: %v", repo.created)
	}
	if err := svc.Delete(context.Background(), []uint{1}); err != nil {
		t.Fatal(err)
	}
	if len(repo.deleted) != 1 {
		t.Fatalf("delete: %v", repo.deleted)
	}
	if err := svc.Reply(context.Background(), 1, "r"); err != nil {
		t.Fatal(err)
	}
	if len(repo.replied) != 1 || repo.replied[0] != "r" {
		t.Fatalf("reply: %v", repo.replied)
	}
	if _, err := svc.Newest(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Info(context.Background(), u); err != nil {
		t.Fatal(err)
	}
	list, total, err := svc.List(context.Background(), 1, 10)
	if err != nil || total != 0 || len(list) != 0 {
		t.Fatalf("list: %v %d %v", list, total, err)
	}
}
