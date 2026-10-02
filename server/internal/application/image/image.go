// Package image 提供图片应用服务。
package image

import (
	"context"
	"mime/multipart"

	"server/internal/domain/image"
	"server/internal/infrastructure/storage"
)

// Service 图片应用服务。
type Service struct {
	images      image.ImageRepository
	storage     storage.FileStorage
	storageName string // 存储类型中文名（本地/七牛云）
}

// NewService 构造图片应用服务。
func NewService(images image.ImageRepository, storage storage.FileStorage, storageName string) *Service {
	return &Service{images: images, storage: storage, storageName: storageName}
}

// Upload 上传图片：存储到对象存储 + 落库。
func (s *Service) Upload(ctx context.Context, file *multipart.FileHeader) (string, error) {
	url, filename, err := s.storage.UploadImage(file)
	if err != nil {
		return "", err
	}
	if err := s.images.Create(ctx, filename, url, s.storageName); err != nil {
		return "", err
	}
	return url, nil
}

// Delete 删除图片：落库删除 + 存储侧删除文件。
func (s *Service) Delete(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	images, err := s.images.Delete(ctx, ids)
	if err != nil {
		return err
	}
	for _, img := range images {
		if err := s.storage.DeleteImage(img.Name); err != nil {
			return err
		}
	}
	return nil
}

// List 图片分页列表。
func (s *Service) List(ctx context.Context, cond image.ListCond) ([]image.Image, int64, error) {
	return s.images.List(ctx, cond)
}
