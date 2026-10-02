package friendlink

import (
	"context"
	"testing"

	"server/internal/domain/friendlink"
)

type stubRepo struct {
	info      []*friendlink.FriendLink
	total     int64
	created   []*friendlink.FriendLink
	deleted   []uint
	updated   []*friendlink.FriendLink
	list      []*friendlink.FriendLink
	listTotal int64
}

func (s *stubRepo) Info(context.Context) ([]*friendlink.FriendLink, int64, error) {
	return s.info, s.total, nil
}
func (s *stubRepo) Create(_ context.Context, l *friendlink.FriendLink) error {
	s.created = append(s.created, l)
	return nil
}
func (s *stubRepo) Delete(_ context.Context, ids []uint) error {
	s.deleted = append(s.deleted, ids...)
	return nil
}
func (s *stubRepo) Update(_ context.Context, l *friendlink.FriendLink) error {
	s.updated = append(s.updated, l)
	return nil
}
func (s *stubRepo) List(_ context.Context, _ friendlink.ListCond) ([]*friendlink.FriendLink, int64, error) {
	return s.list, s.listTotal, nil
}

func TestServiceLifecycle(t *testing.T) {
	repo := &stubRepo{}
	svc := NewService(repo)

	if err := svc.Create(context.Background(), &friendlink.FriendLink{Name: "a", Logo: "/l.png"}); err != nil {
		t.Fatal(err)
	}
	if len(repo.created) != 1 || repo.created[0].Name != "a" {
		t.Fatalf("create: %v", repo.created)
	}
	if err := svc.Delete(context.Background(), []uint{1}); err != nil {
		t.Fatal(err)
	}
	if len(repo.deleted) != 1 {
		t.Fatalf("delete: %v", repo.deleted)
	}
	if err := svc.Update(context.Background(), &friendlink.FriendLink{ID: 1, Name: "b"}); err != nil {
		t.Fatal(err)
	}
	if len(repo.updated) != 1 || repo.updated[0].Name != "b" {
		t.Fatalf("update: %v", repo.updated)
	}
	info, total, err := svc.Info(context.Background())
	if err != nil || total != 0 {
		t.Fatalf("info: %v %d %v", info, total, err)
	}
	list, total, err := svc.List(context.Background(), friendlink.ListCond{})
	if err != nil || total != 0 || len(list) != 0 {
		t.Fatalf("list: %v %d %v", list, total, err)
	}
}
