package advertisement

import (
	"context"
	"testing"

	"server/internal/domain/advertisement"
)

type stubRepo struct {
	info      []*advertisement.Advertisement
	total     int64
	created   []*advertisement.Advertisement
	deleted   []uint
	updated   []*advertisement.Advertisement
	list      []*advertisement.Advertisement
	listTotal int64
}

func (s *stubRepo) Info(context.Context) ([]*advertisement.Advertisement, int64, error) {
	return s.info, s.total, nil
}
func (s *stubRepo) Create(_ context.Context, ad *advertisement.Advertisement) error {
	s.created = append(s.created, ad)
	return nil
}
func (s *stubRepo) Delete(_ context.Context, ids []uint) error {
	s.deleted = append(s.deleted, ids...)
	return nil
}
func (s *stubRepo) Update(_ context.Context, ad *advertisement.Advertisement) error {
	s.updated = append(s.updated, ad)
	return nil
}
func (s *stubRepo) List(_ context.Context, _ advertisement.ListCond) ([]*advertisement.Advertisement, int64, error) {
	return s.list, s.listTotal, nil
}

func TestServiceLifecycle(t *testing.T) {
	repo := &stubRepo{}
	svc := NewService(repo)

	if err := svc.Create(context.Background(), &advertisement.Advertisement{Title: "a", AdImage: "/a.png"}); err != nil {
		t.Fatal(err)
	}
	if len(repo.created) != 1 || repo.created[0].Title != "a" {
		t.Fatalf("create: %v", repo.created)
	}
	if err := svc.Delete(context.Background(), []uint{1}); err != nil {
		t.Fatal(err)
	}
	if len(repo.deleted) != 1 {
		t.Fatalf("delete: %v", repo.deleted)
	}
	if err := svc.Update(context.Background(), &advertisement.Advertisement{ID: 1, Title: "b"}); err != nil {
		t.Fatal(err)
	}
	if len(repo.updated) != 1 || repo.updated[0].Title != "b" {
		t.Fatalf("update: %v", repo.updated)
	}
	list, total, err := svc.List(context.Background(), advertisement.ListCond{})
	if err != nil || len(list) != 0 || total != 0 {
		t.Fatalf("list: %v %d %v", list, total, err)
	}
	info, total, err := svc.Info(context.Background())
	if err != nil || total != 0 {
		t.Fatalf("info: %v %d %v", info, total, err)
	}
}
