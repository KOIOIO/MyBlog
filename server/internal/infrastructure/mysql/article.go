// Package mysql 提供文章 MySQL 仓储（分类/标签计数、收藏、图片类别联动均在此实现）。
package mysql

import (
	"context"
	"errors"

	"server/internal/domain/article"
	"server/internal/model/appTypes"
	"server/internal/model/database"

	"gorm.io/gorm"
)

// ArticleRepository 文章仓储。
type ArticleRepository struct {
	db *gorm.DB
}

// NewArticleRepository 构造文章仓储。
func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

// Categories 获取全部文章类别。
func (r *ArticleRepository) Categories(ctx context.Context) ([]article.ArticleCategory, error) {
	var rows []database.ArticleCategory
	if err := r.db.WithContext(ctx).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]article.ArticleCategory, 0, len(rows))
	for _, row := range rows {
		out = append(out, article.ArticleCategory{Category: row.Category, Number: row.Number})
	}
	return out, nil
}

// Tags 获取全部固定标签（按 `group` ASC, number DESC, tag ASC 排序）。
func (r *ArticleRepository) Tags(ctx context.Context) ([]article.BlogTag, error) {
	var rows []database.BlogTag
	if err := r.db.WithContext(ctx).Order("`group` ASC, number DESC, tag ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]article.BlogTag, 0, len(rows))
	for _, row := range rows {
		out = append(out, article.BlogTag{Tag: row.Tag, Group: row.Group, Number: row.Number})
	}
	return out, nil
}

// CheckTagsExist 校验标签全部存在于固定标签库。
func (r *ArticleRepository) CheckTagsExist(ctx context.Context, tags []string) error {
	if len(tags) == 0 {
		return nil
	}
	var count int64
	if err := r.db.WithContext(ctx).Model(&database.BlogTag{}).Where("tag IN ?", tags).Count(&count).Error; err != nil {
		return err
	}
	if int(count) != len(tags) {
		return article.ErrTagNotExist
	}
	return nil
}

// LikeTx 收藏/取消收藏（DB 事务），返回 delta。
func (r *ArticleRepository) LikeTx(ctx context.Context, userID uint, articleID string) (int, error) {
	var delta int
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var al database.ArticleLike
		if errors.Is(tx.Where("user_id = ? AND article_id = ?", userID, articleID).First(&al).Error, gorm.ErrRecordNotFound) {
			if err := tx.Create(&database.ArticleLike{UserID: userID, ArticleID: articleID}).Error; err != nil {
				return err
			}
			delta = 1
		} else {
			if err := tx.Delete(&al).Error; err != nil {
				return err
			}
			delta = -1
		}
		return nil
	})
	return delta, err
}

// IsLike 用户是否已收藏。
func (r *ArticleRepository) IsLike(ctx context.Context, userID uint, articleID string) (bool, error) {
	err := r.db.WithContext(ctx).Where("user_id = ? AND article_id = ?", userID, articleID).First(&database.ArticleLike{}).Error
	return !errors.Is(err, gorm.ErrRecordNotFound), nil
}

// LikesList 用户收藏分页列表。
func (r *ArticleRepository) LikesList(ctx context.Context, userID uint, pageNo, pageSize int) ([]article.ArticleLike, int64, error) {
	db := r.db.WithContext(ctx).Where("user_id = ?", userID)
	var rows []database.ArticleLike
	var total int64
	if err := db.Model(&database.ArticleLike{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if pageNo < 1 {
		pageNo = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if err := db.Order("id desc").Limit(pageSize).Offset((pageNo - 1) * pageSize).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	out := make([]article.ArticleLike, 0, len(rows))
	for _, row := range rows {
		out = append(out, article.ArticleLike{ID: row.ID, ArticleID: row.ArticleID, UserID: row.UserID})
	}
	return out, total, nil
}

// CreateWithCounts 创建文章事务联动（分类计数 + 标签计数 + 图片类别）。
func (r *ArticleRepository) CreateWithCounts(ctx context.Context, a *article.Article, cover string, illustrations []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := updateCategoryCount(tx, "", a.Category); err != nil {
			return err
		}
		if err := updateTagsCount(tx, []string{}, a.Tags); err != nil {
			return err
		}
		if err := changeImagesCategory(tx, []string{cover}, appTypes.Cover); err != nil {
			return err
		}
		return changeImagesCategory(tx, illustrations, appTypes.Illustration)
	})
}

// UpdateWithCounts 更新文章事务联动（分类/标签/封面/插图）。
func (r *ArticleRepository) UpdateWithCounts(ctx context.Context, a *article.Article, old *article.Article, newCover string, newIllustrations []string, addedIllustrations, removedIllustrations []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := updateCategoryCount(tx, old.Category, a.Category); err != nil {
			return err
		}
		if err := updateTagsCount(tx, old.Tags, a.Tags); err != nil {
			return err
		}
		if a.Cover != old.Cover {
			if err := initImagesCategory(tx, []string{old.Cover}); err != nil {
				return err
			}
			if err := changeImagesCategory(tx, []string{newCover}, appTypes.Cover); err != nil {
				return err
			}
		}
		if err := initImagesCategory(tx, removedIllustrations); err != nil {
			return err
		}
		return changeImagesCategory(tx, addedIllustrations, appTypes.Illustration)
	})
}

// DeleteWithCounts 删除文章事务联动（分类/标签/图片初始化）。
func (r *ArticleRepository) DeleteWithCounts(ctx context.Context, a *article.Article, cover string, illustrations []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := updateCategoryCount(tx, a.Category, ""); err != nil {
			return err
		}
		if err := updateTagsCount(tx, a.Tags, []string{}); err != nil {
			return err
		}
		if err := initImagesCategory(tx, []string{cover}); err != nil {
			return err
		}
		return initImagesCategory(tx, illustrations)
	})
}

// ---- 事务内辅助（与原 service/article_helpers.go 行为一致） ----

// updateCategoryCount 更新类别计数。
func updateCategoryCount(tx *gorm.DB, oldCategory, newCategory string) error {
	if newCategory == oldCategory {
		return nil
	}
	if newCategory != "" {
		var newArticleCategory database.ArticleCategory
		if errors.Is(tx.Where("category = ?", newCategory).First(&newArticleCategory).Error, gorm.ErrRecordNotFound) {
			if err := tx.Create(&database.ArticleCategory{Category: newCategory, Number: 1}).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Model(&newArticleCategory).Update("number", gorm.Expr("number + ?", 1)).Error; err != nil {
				return err
			}
		}
	}
	if oldCategory != "" {
		var oldArticleCategory database.ArticleCategory
		if err := tx.Where("category = ?", oldCategory).First(&oldArticleCategory).Update("number", gorm.Expr("number - ?", 1)).Error; err != nil {
			return err
		}
		if oldArticleCategory.Number == 1 {
			if err := tx.Delete(&oldArticleCategory).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

// updateTagsCount 更新标签引用计数。
func updateTagsCount(tx *gorm.DB, oldTags, newTags []string) error {
	addedTags, removedTags := article.DiffTags(oldTags, newTags)
	for _, addedTag := range addedTags {
		if err := tx.Model(&database.BlogTag{}).Where("tag = ?", addedTag).
			Update("number", gorm.Expr("number + ?", 1)).Error; err != nil {
			return err
		}
	}
	for _, removedTag := range removedTags {
		if err := tx.Model(&database.BlogTag{}).Where("tag = ?", removedTag).
			Update("number", gorm.Expr("number - ?", 1)).Error; err != nil {
			return err
		}
	}
	return nil
}

// initImagesCategory 初始化图片类别。
func initImagesCategory(tx *gorm.DB, urls []string) error {
	return tx.Model(&database.Image{}).Where("url IN ?", urls).Update("category", appTypes.Null).Error
}

// changeImagesCategory 修改图片类别。
func changeImagesCategory(tx *gorm.DB, urls []string, category appTypes.Category) error {
	return tx.Model(&database.Image{}).Where("url IN ?", urls).Update("category", category).Error
}

// 编译期断言。
var _ article.ArticleRepository = (*ArticleRepository)(nil)
