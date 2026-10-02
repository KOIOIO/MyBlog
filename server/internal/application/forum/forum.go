// Package forum 提供论坛应用服务（用例编排，依赖领域端口）。
package forum

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"server/config"
	"server/internal/common/crypto"
	"server/internal/domain/forum"
	"server/internal/domain/shared"

	"go.uber.org/zap"
)

// whiteImageList 允许上传的图片扩展名白名单（与 utils/upload.WhiteImageList 一致）。
var whiteImageList = map[string]struct{}{
	".jpg":  {},
	".png":  {},
	".jpeg": {},
	".ico":  {},
	".tiff": {},
	".gif":  {},
	".svg":  {},
	".webp": {},
}

// Service 论坛应用服务。
type Service struct {
	posts forum.ForumRepository
	cfg   *config.Config
	log   *zap.Logger
}

// NewService 构造论坛应用服务。
func NewService(posts forum.ForumRepository, cfg *config.Config, log *zap.Logger) *Service {
	return &Service{posts: posts, cfg: cfg, log: log}
}

// Tags 固定标签库。
func (s *Service) Tags(ctx context.Context) ([]forum.Tag, error) {
	return s.posts.Tags(ctx)
}

// Publish 发布帖子，返回新帖子 ID。
func (s *Service) Publish(ctx context.Context, p *forum.ForumPost) (uint, error) {
	if _, ok := forum.AllowedCategory[p.Category]; !ok {
		return 0, forum.ErrCategoryNotAllowed
	}
	return s.posts.Publish(ctx, p)
}

// List 帖子分页列表。
func (s *Service) List(ctx context.Context, cond forum.ListCond) ([]*forum.ForumPost, int64, error) {
	return s.posts.List(ctx, cond)
}

// Detail 帖子详情（浏览数 +1）。
func (s *Service) Detail(ctx context.Context, id uint) (*forum.ForumDetail, error) {
	post, comments, err := s.posts.Detail(ctx, id)
	if err != nil {
		return nil, err
	}
	return &forum.ForumDetail{ForumPost: *post, Comments: comments}, nil
}

// Like 切换点赞。
func (s *Service) Like(ctx context.Context, postID, userID uint) (bool, int, error) {
	return s.posts.Like(ctx, postID, userID)
}

// Comment 发表评论。
func (s *Service) Comment(ctx context.Context, c *forum.ForumComment) error {
	return s.posts.Comment(ctx, c)
}

// ManageList 帖子管理列表。
func (s *Service) ManageList(ctx context.Context, cond forum.ManageListCond) ([]*forum.ForumPost, int64, error) {
	return s.posts.ManageList(ctx, cond)
}

// Delete 删除帖子。
func (s *Service) Delete(ctx context.Context, ids []uint, userID uint, roleID shared.RoleID) error {
	return s.posts.DeletePosts(ctx, ids, userID, roleID)
}

// ManageComments 评论管理列表。
func (s *Service) ManageComments(ctx context.Context, cond forum.ManageCommentCond) ([]*forum.ManageComment, int64, error) {
	return s.posts.ManageComments(ctx, cond)
}

// DeleteComments 删除评论。
func (s *Service) DeleteComments(ctx context.Context, ids []uint, userID uint, roleID shared.RoleID) error {
	return s.posts.DeleteComments(ctx, ids, userID, roleID)
}

// Upload 论坛图片上传，存到 uploads/forum/ 目录。
func (s *Service) Upload(file *multipart.FileHeader) (string, error) {
	size := float64(file.Size) / float64(1024*1024)
	if size >= float64(s.cfg.Upload.Size) {
		return "", fmt.Errorf("the image size exceeds the set size, the current size is: %.2f MB, the set size is: %d MB", size, s.cfg.Upload.Size)
	}

	ext := filepath.Ext(file.Filename)
	name := strings.TrimSuffix(file.Filename, ext)
	if _, exists := whiteImageList[ext]; !exists {
		return "", errors.New("don't upload files that aren't image types")
	}

	filename := crypto.MD5V([]byte(name)) + "-" + time.Now().Format("20060102150405") + ext
	dir := s.cfg.Upload.Path + "/forum/"
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", err
	}

	dst := dir + filename
	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	if _, err = io.Copy(out, src); err != nil {
		return "", err
	}

	return "/" + dst, nil
}
