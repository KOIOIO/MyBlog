// Package storage 提供文件存储端口与 local/qiniu 实现（替代 utils/upload）。
// 全部依赖构造注入，不读取 global。
package storage

import (
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

// FileStorage 文件存储端口。
type FileStorage interface {
	// UploadImage 上传图片，返回（访问 URL, 文件名, 错误）。
	// local 返回格式 /uploads/image/fileName；qiniu 返回格式 http(s)://image.xxx.xx/fileName。
	UploadImage(file *multipart.FileHeader) (url, filename string, err error)
	// DeleteImage 按文件名删除图片。
	DeleteImage(key string) error
}

// New 根据 ossType 构造存储实现。
func New(ossType string, cfg *config.Config) FileStorage {
	switch ossType {
	case "qiniu":
		return &Qiniu{cfg: cfg}
	default:
		return &Local{cfg: cfg}
	}
}

// validateAndFilename 公共校验与文件名生成（大小/类型/MD5+时间戳）。
func validateAndFilename(cfg *config.Config, file *multipart.FileHeader) (string, error) {
	size := float64(file.Size) / float64(1024*1024)
	if size >= float64(cfg.Upload.Size) {
		return "", fmt.Errorf("the image size exceeds the set size, the current size is: %.2f MB, the set size is: %d MB", size, cfg.Upload.Size)
	}

	ext := filepath.Ext(file.Filename)
	name := strings.TrimSuffix(file.Filename, ext)
	if _, exists := whiteImageList[ext]; !exists {
		return "", errors.New("don't upload files that aren't image types")
	}

	return crypto.MD5V([]byte(name)) + "-" + time.Now().Format("20060102150405") + ext, nil
}

// Local 本地磁盘存储。
type Local struct {
	cfg *config.Config
}

// UploadImage 上传图片到本地 uploads/image/ 目录。
func (s *Local) UploadImage(file *multipart.FileHeader) (string, string, error) {
	filename, err := validateAndFilename(s.cfg, file)
	if err != nil {
		return "", "", err
	}
	path := s.cfg.Upload.Path + "/image/"
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return "", "", err
	}

	filepath := path + filename
	out, err := os.Create(filepath)
	if err != nil {
		return "", "", err
	}
	defer out.Close()

	f, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer f.Close()

	if _, err = io.Copy(out, f); err != nil {
		return "", "", err
	}

	return "/" + filepath, filename, nil
}

// DeleteImage 删除本地图片。
func (s *Local) DeleteImage(key string) error {
	return os.Remove(s.cfg.Upload.Path + "/image/" + key)
}
