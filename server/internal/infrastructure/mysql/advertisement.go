// AdvertisementRepo 广告仓储（GORM）。
package mysql

import (
	"context"

	"server/internal/common/page"
	"server/internal/domain/advertisement"
	"server/internal/model/appTypes"
	"server/internal/model/database"
	"server/internal/model/other"

	"gorm.io/gorm"
)

// AdvertisementRepo 广告仓储。
type AdvertisementRepo struct {
	db *gorm.DB
}

// NewAdvertisementRepo 构造广告仓储。
func NewAdvertisementRepo(db *gorm.DB) *AdvertisementRepo {
	return &AdvertisementRepo{db: db}
}

var _ advertisement.AdvertisementRepository = (*AdvertisementRepo)(nil)

func toDomainAdvertisement(m *database.Advertisement) *advertisement.Advertisement {
	return &advertisement.Advertisement{
		ID:        m.ID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		AdImage:   m.AdImage,
		Link:      m.Link,
		Title:     m.Title,
		Content:   m.Content,
	}
}

// Info 广告信息（全量 + 总数）。
func (r *AdvertisementRepo) Info(ctx context.Context) ([]*advertisement.Advertisement, int64, error) {
	var ads []database.Advertisement
	var total int64
	err := r.db.WithContext(ctx).Model(&database.Advertisement{}).Count(&total).Find(&ads).Error
	if err != nil {
		return nil, 0, err
	}
	out := make([]*advertisement.Advertisement, 0, len(ads))
	for i := range ads {
		out = append(out, toDomainAdvertisement(&ads[i]))
	}
	return out, total, nil
}

// Create 创建广告（事务：图片类别改为"广告" + 创建）。
func (r *AdvertisementRepo) Create(ctx context.Context, ad *advertisement.Advertisement) error {
	toCreate := database.Advertisement{
		AdImage: ad.AdImage,
		Link:    ad.Link,
		Title:   ad.Title,
		Content: ad.Content,
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := ChangeImagesCategory(tx, []string{toCreate.AdImage}, appTypes.AdImage); err != nil {
			return err
		}
		return tx.Create(&toCreate).Error
	})
}

// Delete 删除广告（事务：图片类别重置 + 删除）。
func (r *AdvertisementRepo) Delete(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, id := range ids {
			var ad database.Advertisement
			if err := tx.Take(&ad, id).Error; err != nil {
				return err
			}
			if err := InitImagesCategory(tx, []string{ad.AdImage}); err != nil {
				return err
			}
			if err := tx.Delete(&ad).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// Update 更新广告（仅 link/title/content）。
func (r *AdvertisementRepo) Update(ctx context.Context, ad *advertisement.Advertisement) error {
	updates := struct {
		Link    string `json:"link"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}{
		Link:    ad.Link,
		Title:   ad.Title,
		Content: ad.Content,
	}
	return r.db.WithContext(ctx).Take(&database.Advertisement{}, ad.ID).Updates(updates).Error
}

// List 广告分页列表。
func (r *AdvertisementRepo) List(ctx context.Context, cond advertisement.ListCond) ([]*advertisement.Advertisement, int64, error) {
	db := r.db.WithContext(ctx)
	if cond.Title != nil {
		db = db.Where("title LIKE ?", "%"+*cond.Title+"%")
	}
	if cond.Content != nil {
		db = db.Where("content LIKE ?", "%"+*cond.Content+"%")
	}
	option := other.MySQLOption{
		PageInfo: pageInfoOf(cond.Page, cond.PageSize),
		Where:    db,
	}
	list, total, err := page.MySQLPagination(r.db.WithContext(ctx), &database.Advertisement{}, option)
	if err != nil {
		return nil, 0, err
	}
	out := make([]*advertisement.Advertisement, 0, len(list))
	for i := range list {
		out = append(out, toDomainAdvertisement(&list[i]))
	}
	return out, total, nil
}
