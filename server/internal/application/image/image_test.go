package image

import (
	"context"
	"errors"
	"mime/multipart"
	"testing"

	"server/internal/domain/image"
)

type stubRepo struct {
	created       []string
	deleted       []uint
	deletedImages []image.Image
	list          []image.Image
	total         int64
	err           error
}

func (s *stubRepo) Create(_ context.Context, name, url, storage string) error {
	s.created = append(s.created, name)
	return s.err
}
func (s *stubRepo) Delete(_ context.Context, ids []uint) ([]image.Image, error) {
	s.deleted = append(s.deleted, ids...)
	return s.deletedImages, s.err
}
func (s *stubRepo) List(_ context.Context, _ image.ListCond) ([]image.Image, int64, error) {
	return s.list, s.total, s.err
}

type stubStorage struct {
	uploaded bool
	deleted  []string
}

func (s *stubStorage) UploadImage(_ *multipart.FileHeader) (string, string, error) {
	s.uploaded = true
	return "/uploads/image/x.png", "x.png", nil
}
func (s *stubStorage) DeleteImage(key string) error {
	s.deleted = append(s.deleted, key)
	return nil
}

func TestDeleteCallsStorage(t *testing.T) {
	repo := &stubRepo{deletedImages: []image.Image{{ID: 1, Name: "a.png"}, {ID: 2, Name: "b.png"}}}
	st := &stubStorage{}
	svc := NewService(repo, st, "本地")

	if err := svc.Delete(context.Background(), []uint{1, 2}); err != nil {
		t.Fatal(err)
	}
	if len(st.deleted) != 2 || st.deleted[0] != "a.png" || st.deleted[1] != "b.png" {
		t.Fatalf("storage deletes: %v", st.deleted)
	}
	// 空 ids 不触发
	if err := svc.Delete(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	if len(st.deleted) != 2 {
		t.Fatalf("empty ids must not delete, got %v", st.deleted)
	}
}

func TestListPassesThrough(t *testing.T) {
	repo := &stubRepo{list: []image.Image{{ID: 5}}, total: 1}
	svc := NewService(repo, &stubStorage{}, "本地")
	list, total, err := svc.List(context.Background(), image.ListCond{})
	if err != nil || total != 1 || len(list) != 1 {
		t.Fatalf("list: %v %d %v", list, total, err)
	}
}

func TestDeleteErrorPropagates(t *testing.T) {
	repo := &stubRepo{err: errors.New("db down")}
	svc := NewService(repo, &stubStorage{}, "本地")
	if err := svc.Delete(context.Background(), []uint{1}); err == nil {
		t.Fatal("want error")
	}
}
