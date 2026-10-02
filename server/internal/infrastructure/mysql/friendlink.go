// FriendLinkRepo 友链仓储（GORM）。
package mysql

import (
	"context"

	"server/internal/common/page"
	"server/internal/domain/friendlink"
	"server/internal/model/appTypes"
	"server/internal/model/database"
	"server/internal/model/other"

	"gorm.io/gorm"
)

// FriendLinkRepo 友链仓储。
type FriendLinkRepo struct {
	db *gorm.DB
}

// NewFriendLinkRepo 构造友链仓储。
func NewFriendLinkRepo(db *gorm.DB) *FriendLinkRepo {
	return &FriendLinkRepo{db: db}
}

var _ friendlink.FriendLinkRepository = (*FriendLinkRepo)(nil)

func toDomainFriendLink(m *database.FriendLink) *friendlink.FriendLink {
	return &friendlink.FriendLink{
		ID:          m.ID,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		Logo:        m.Logo,
		Link:        m.Link,
		Name:        m.Name,
		Description: m.Description,
	}
}

// Info 友链信息（全量 + 总数）。
func (r *FriendLinkRepo) Info(ctx context.Context) ([]*friendlink.FriendLink, int64, error) {
	var links []database.FriendLink
	var total int64
	err := r.db.WithContext(ctx).Model(&database.FriendLink{}).Count(&total).Find(&links).Error
	if err != nil {
		return nil, 0, err
	}
	out := make([]*friendlink.FriendLink, 0, len(links))
	for i := range links {
		out = append(out, toDomainFriendLink(&links[i]))
	}
	return out, total, nil
}

// Create 创建友链（事务：图片类别改为"友链" + 创建）。
func (r *FriendLinkRepo) Create(ctx context.Context, link *friendlink.FriendLink) error {
	toCreate := database.FriendLink{
		Logo:        link.Logo,
		Link:        link.Link,
		Name:        link.Name,
		Description: link.Description,
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := ChangeImagesCategory(tx, []string{toCreate.Logo}, appTypes.Logo); err != nil {
			return err
		}
		return tx.Create(&toCreate).Error
	})
}

// Delete 删除友链（事务：图片类别重置 + 删除）。
func (r *FriendLinkRepo) Delete(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, id := range ids {
			var link database.FriendLink
			if err := tx.Take(&link, id).Error; err != nil {
				return err
			}
			if err := InitImagesCategory(tx, []string{link.Logo}); err != nil {
				return err
			}
			if err := tx.Delete(&link).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// Update 更新友链（仅 link/name/description）。
func (r *FriendLinkRepo) Update(ctx context.Context, link *friendlink.FriendLink) error {
	updates := struct {
		Link        string `json:"link"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}{
		Link:        link.Link,
		Name:        link.Name,
		Description: link.Description,
	}
	return r.db.WithContext(ctx).Take(&database.FriendLink{}, link.ID).Updates(updates).Error
}

// List 友链分页列表。
func (r *FriendLinkRepo) List(ctx context.Context, cond friendlink.ListCond) ([]*friendlink.FriendLink, int64, error) {
	db := r.db.WithContext(ctx)
	if cond.Name != nil {
		db = db.Where("name LIKE ?", "%"+*cond.Name+"%")
	}
	if cond.Description != nil {
		db = db.Where("description LIKE ?", "%"+*cond.Description+"%")
	}
	option := other.MySQLOption{
		PageInfo: pageInfoOf(cond.Page, cond.PageSize),
		Where:    db,
	}
	list, total, err := page.MySQLPagination(r.db.WithContext(ctx), &database.FriendLink{}, option)
	if err != nil {
		return nil, 0, err
	}
	out := make([]*friendlink.FriendLink, 0, len(list))
	for i := range list {
		out = append(out, toDomainFriendLink(&list[i]))
	}
	return out, total, nil
}
