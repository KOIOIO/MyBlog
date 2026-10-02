// Package mysql 提供评论与论坛的 GORM 仓储实现。
package mysql

import (
	"context"

	"server/internal/common/page"
	"server/internal/domain/comment"
	"server/internal/domain/shared"
	"server/internal/model/database"
	esmodel "server/internal/model/elasticsearch"
	"server/internal/model/other"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/scriptlanguage"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// CommentRepo 评论仓储（GORM 持久化，ES 评论数同步显式执行，替代原模型钩子）。
type CommentRepo struct {
	db *gorm.DB
	es *elasticsearch.TypedClient
}

// NewCommentRepo 构造评论仓储。
func NewCommentRepo(db *gorm.DB, es *elasticsearch.TypedClient) *CommentRepo {
	return &CommentRepo{db: db, es: es}
}

// updateESCommentCount 同步 ES 文章评论数（原 database.Comment 钩子逻辑）。
func (r *CommentRepo) updateESCommentCount(ctx context.Context, articleID string, delta int) error {
	source := "ctx._source.comments += 1"
	if delta < 0 {
		source = "ctx._source.comments -= 1"
	}
	script := types.Script{Source: &source, Lang: &scriptlanguage.Painless}
	_, err := r.es.Update(esmodel.ArticleIndex(), articleID).Script(&script).Do(ctx)
	return err
}

var _ comment.CommentRepository = (*CommentRepo)(nil)

// userSelect 评论关联用户预加载字段（与原 service 一致）。
const userSelect = "uuid, username, avatar, address, signature"

func toDomainComment(c *database.Comment) *comment.Comment {
	d := &comment.Comment{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		ArticleID: c.ArticleID,
		PID:       c.PID,
		UserUUID:  c.UserUUID,
		Content:   c.Content,
	}
	if c.Children != nil {
		d.Children = make([]*comment.Comment, 0, len(c.Children))
		for i := range c.Children {
			d.Children = append(d.Children, toDomainComment(&c.Children[i]))
		}
	}
	u := c.User
	d.User = comment.UserBrief{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		UUID:      u.UUID,
		Username:  u.Username,
		Email:     u.Email,
		Openid:    u.Openid,
		Avatar:    u.Avatar,
		Address:   u.Address,
		Signature: u.Signature,
		RoleID:    shared.RoleID(u.RoleID),
		Register:  shared.Register(u.Register),
		Freeze:    u.Freeze,
	}
	return d
}

// ByArticle 文章一级评论并递归加载子评论。
func (r *CommentRepo) ByArticle(ctx context.Context, articleID string) ([]*comment.Comment, error) {
	var roots []database.Comment
	if err := r.db.WithContext(ctx).Where("article_id = ? AND p_id IS NULL", articleID).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select(userSelect)
		}).Find(&roots).Error; err != nil {
		return nil, err
	}
	out := make([]*comment.Comment, 0, len(roots))
	for i := range roots {
		if err := r.loadChildren(ctx, &roots[i]); err != nil {
			return nil, err
		}
		out = append(out, toDomainComment(&roots[i]))
	}
	return out, nil
}

// loadChildren 递归加载子评论。
func (r *CommentRepo) loadChildren(ctx context.Context, c *database.Comment) error {
	var children []database.Comment
	if err := r.db.WithContext(ctx).Where("p_id = ?", c.ID).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select(userSelect)
		}).Find(&children).Error; err != nil {
		return err
	}
	for i := range children {
		if err := r.loadChildren(ctx, &children[i]); err != nil {
			return err
		}
	}
	c.Children = children
	return nil
}

// Newest 最新 N 条评论（不加载子评论）。
func (r *CommentRepo) Newest(ctx context.Context, limit int) ([]*comment.Comment, error) {
	var comments []database.Comment
	if err := r.db.WithContext(ctx).Order("id desc").Limit(limit).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select(userSelect)
		}).Find(&comments).Error; err != nil {
		return nil, err
	}
	out := make([]*comment.Comment, 0, len(comments))
	for i := range comments {
		out = append(out, toDomainComment(&comments[i]))
	}
	return out, nil
}

// Create 创建评论并同步 ES 评论数。
func (r *CommentRepo) Create(ctx context.Context, c *comment.Comment) error {
	if err := r.db.WithContext(ctx).Create(&database.Comment{
		ArticleID: c.ArticleID,
		PID:       c.PID,
		UserUUID:  c.UserUUID,
		Content:   c.Content,
	}).Error; err != nil {
		return err
	}
	return r.updateESCommentCount(ctx, c.ArticleID, 1)
}

// DeleteTree 校验权限并事务级联删除。
func (r *CommentRepo) DeleteTree(ctx context.Context, ids []uint, userUUID uuid.UUID, roleID shared.RoleID) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, id := range ids {
			var c database.Comment
			if err := tx.Take(&c, id).Error; err != nil {
				return err
			}
			if userUUID != c.UserUUID && roleID != shared.Admin {
				return comment.ErrNoPermission
			}
			if err := r.deleteTree(tx, id); err != nil {
				return err
			}
		}
		return nil
	})
}

// deleteTree 递归删除评论及其子评论（显式同步 ES 评论数）。
func (r *CommentRepo) deleteTree(tx *gorm.DB, id uint) error {
	var children []database.Comment
	if err := tx.Where("p_id = ?", id).Find(&children).Error; err != nil {
		return err
	}
	for _, child := range children {
		if err := r.deleteTree(tx, child.ID); err != nil {
			return err
		}
	}
	var c database.Comment
	if err := tx.Take(&c, id).Error; err != nil {
		return err
	}
	if err := tx.Delete(&database.Comment{}, id).Error; err != nil {
		return err
	}
	return r.updateESCommentCount(context.Background(), c.ArticleID, -1)
}

// ByUser 用户全部评论（含子评论树）。
func (r *CommentRepo) ByUser(ctx context.Context, userUUID uuid.UUID) ([]*comment.Comment, error) {
	var raw []database.Comment
	if err := r.db.WithContext(ctx).Order("id desc").Where("user_uuid = ?", userUUID).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select(userSelect)
		}).Find(&raw).Error; err != nil {
		return nil, err
	}
	out := make([]*comment.Comment, 0, len(raw))
	for i := range raw {
		if err := r.loadChildren(ctx, &raw[i]); err != nil {
			return nil, err
		}
		out = append(out, toDomainComment(&raw[i]))
	}
	return out, nil
}

// List 管理列表分页（不预加载用户，保持旧行为）。
func (r *CommentRepo) List(ctx context.Context, cond comment.ListCond) ([]*comment.Comment, int64, error) {
	db := r.db.WithContext(ctx)
	if cond.ArticleID != nil {
		db = db.Where("article_id = ?", *cond.ArticleID)
	}
	if cond.UserUUID != nil {
		db = db.Where("user_uuid = ?", *cond.UserUUID)
	}
	if cond.Content != nil {
		db = db.Where("content LIKE ?", "%"+*cond.Content+"%")
	}
	option := other.MySQLOption{
		PageInfo: pageInfoOf(cond.Page, cond.PageSize),
		Where:    db,
	}
	list, total, err := page.MySQLPagination(db, &database.Comment{}, option)
	if err != nil {
		return nil, 0, err
	}
	out := make([]*comment.Comment, 0, len(list))
	for i := range list {
		out = append(out, toDomainComment(&list[i]))
	}
	return out, total, nil
}
