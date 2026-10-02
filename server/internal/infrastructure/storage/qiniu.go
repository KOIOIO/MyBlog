// Qiniu 七牛云对象存储。
package storage

import (
	"context"
	"mime/multipart"

	"server/config"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	qstorage "github.com/qiniu/go-sdk/v7/storage"
)

// Qiniu 七牛云存储实现。
type Qiniu struct {
	cfg *config.Config
}

// UploadImage 上传图片到七牛云。
func (s *Qiniu) UploadImage(file *multipart.FileHeader) (string, string, error) {
	filename, err := validateAndFilename(s.cfg, file)
	if err != nil {
		return "", "", err
	}

	q := s.cfg.Qiniu
	putPolicy := qstorage.PutPolicy{Scope: q.Bucket}
	mac := qbox.NewMac(q.AccessKey, q.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := qiniuConfig(q)
	formUploader := qstorage.NewFormUploader(cfg)
	putRet := qstorage.PutRet{}
	putExtra := qstorage.PutExtra{Params: map[string]string{}}

	data, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer data.Close()

	if err = formUploader.Put(context.Background(), &putRet, upToken, filename, data, file.Size, &putExtra); err != nil {
		return "", "", err
	}
	return q.ImgPath + putRet.Key, putRet.Key, nil
}

// DeleteImage 删除七牛云图片。
func (s *Qiniu) DeleteImage(key string) error {
	q := s.cfg.Qiniu
	mac := qbox.NewMac(q.AccessKey, q.SecretKey)
	cfg := qiniuConfig(q)
	bucketManager := qstorage.NewBucketManager(mac, cfg)
	return bucketManager.Delete(q.Bucket, key)
}

// qiniuConfig 从配置构造七牛存储配置。
func qiniuConfig(q config.Qiniu) *qstorage.Config {
	cfg := qstorage.Config{
		UseHTTPS:      q.UseHTTPS,
		UseCdnDomains: q.UseCdnDomains,
	}
	switch q.Zone {
	case "z0", "ZoneHuadong":
		cfg.Zone = &qstorage.ZoneHuadong
	case "z1", "ZoneHuabei":
		cfg.Zone = &qstorage.ZoneHuabei
	case "z2", "ZoneHuanan":
		cfg.Zone = &qstorage.ZoneHuanan
	case "na0", "ZoneBeimei":
		cfg.Zone = &qstorage.ZoneBeimei
	case "as0", "ZoneXinjiapo":
		cfg.Zone = &qstorage.ZoneXinjiapo
	case "ZoneHuadongZheJiang2":
		cfg.Zone = &qstorage.ZoneHuadongZheJiang2
	}
	return &cfg
}
