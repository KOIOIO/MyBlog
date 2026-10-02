// ImageRepo 图片仓储（GORM）。
package mysql

import (
	"context"

	"server/internal/common/page"
	"server/internal/domain/image"
	"server/model/appTypes"
	"server/model/database"
	"server/model/other"

	"gorm.io/gorm"
)

// ImageRepo 图片仓储。
type ImageRepo struct {
	db *gorm.DB
}

// NewImageRepo 构造图片仓储。
func NewImageRepo(db *gorm.DB) *ImageRepo {
	return &ImageRepo{db: db}
}

var _ image.ImageRepository = (*ImageRepo)(nil)

func toDomainImage(m *database.Image) *image.Image {
	return &image.Image{
		ID:        m.ID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Name:      m.Name,
		URL:       m.URL,
		Category:  m.Category.String(),
		Storage:   m.Storage.String(),
	}
}

// ChangeImagesCategory 修改图片类别（可接受事务内调用，保持旧 utils 行为）。
func ChangeImagesCategory(tx *gorm.DB, urls []string, category appTypes.Category) error {
	return tx.Model(&database.Image{}).Where("url IN ?", urls).Update("category", category).Error
}

// InitImagesCategory 重置图片类别为"未使用"（可接受事务内调用）。
func InitImagesCategory(tx *gorm.DB, urls []string) error {
	return tx.Model(&database.Image{}).Where("url IN ?", urls).Update("category", appTypes.Null).Error
}

// Create 记录上传图片。
func (r *ImageRepo) Create(ctx context.Context, name, url, storage string) error {
	return r.db.WithContext(ctx).Create(&database.Image{
		Name:     name,
		URL:      url,
		Category: appTypes.Null,
		Storage:  appTypes.ToStorage(storage),
	}).Error
}

// Delete 查找并删除图片记录。
func (r *ImageRepo) Delete(ctx context.Context, ids []uint) ([]image.Image, error) {
	var images []database.Image
	if err := r.db.WithContext(ctx).Find(&images, ids).Error; err != nil {
		return nil, err
	}
	out := make([]image.Image, 0, len(images))
	for i := range images {
		out = append(out, *toDomainImage(&images[i]))
	}
	for i := range images {
		if err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			return tx.Delete(&images[i]).Error
		}); err != nil {
			return nil, err
		}
	}
	return out, nil
}

// List 图片分页列表。
func (r *ImageRepo) List(ctx context.Context, cond image.ListCond) ([]image.Image, int64, error) {
	db := r.db.WithContext(ctx)
	if cond.Name != nil {
		db = db.Where("name LIKE ?", "%"+*cond.Name+"%")
	}
	if cond.Category != nil {
		db = db.Where("category = ?", appTypes.ToCategory(*cond.Category))
	}
	if cond.Storage != nil {
		db = db.Where("storage = ?", appTypes.ToStorage(*cond.Storage))
	}

	option := other.MySQLOption{
		PageInfo: pageInfoOf(cond.Page, cond.PageSize),
		Where:    db,
	}
	list, total, err := page.MySQLPagination(r.db.WithContext(ctx), &database.Image{}, option)
	if err != nil {
		return nil, 0, err
	}
	out := make([]image.Image, 0, len(list))
	for i := range list {
		out = append(out, *toDomainImage(&list[i]))
	}
	return out, total, nil
}
