// ForumRepo 论坛仓储（GORM）。
package mysql

import (
	"context"
	"errors"

	"server/internal/common/page"
	"server/internal/domain/forum"
	"server/internal/domain/shared"
	"server/model/database"
	"server/model/other"

	"gorm.io/gorm"
)

// ForumRepo 论坛仓储（GORM）。
type ForumRepo struct {
	db *gorm.DB
}

// NewForumRepo 构造论坛仓储。
func NewForumRepo(db *gorm.DB) *ForumRepo {
	return &ForumRepo{db: db}
}

var _ forum.ForumRepository = (*ForumRepo)(nil)

// detailUserSelect 详情评论用户预加载字段。
const detailUserSelect = "id, uuid, username, avatar, signature"

func toDomainForumPost(p *database.ForumPost) *forum.ForumPost {
	d := &forum.ForumPost{
		ID:           p.ID,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
		UserID:       p.UserID,
		Title:        p.Title,
		Content:      p.Content,
		Category:     p.Category,
		Tags:         []string(p.Tags),
		Images:       []string(p.Images),
		LikeCount:    p.LikeCount,
		CommentCount: p.CommentCount,
		ViewCount:    p.ViewCount,
	}
	u := p.User
	d.User = forum.UserBrief{
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

func toDomainForumComment(c *database.ForumComment) *forum.ForumComment {
	d := &forum.ForumComment{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		PostID:    c.PostID,
		ParentID:  c.ParentID,
		UserID:    c.UserID,
		Content:   c.Content,
	}
	if c.Children != nil {
		d.Children = make([]*forum.ForumComment, 0, len(c.Children))
		for i := range c.Children {
			d.Children = append(d.Children, toDomainForumComment(&c.Children[i]))
		}
	}
	u := c.User
	d.User = forum.UserBrief{
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

// Tags 固定标签库全部标签。
func (r *ForumRepo) Tags(ctx context.Context) ([]forum.Tag, error) {
	var tags []database.BlogTag
	if err := r.db.WithContext(ctx).Order("`group` ASC, number DESC, tag ASC").Find(&tags).Error; err != nil {
		return nil, err
	}
	out := make([]forum.Tag, 0, len(tags))
	for _, t := range tags {
		out = append(out, forum.Tag{Tag: t.Tag, Group: t.Group, Number: t.Number})
	}
	return out, nil
}

// checkTagsExist 校验标签都存在于固定标签库。
func (r *ForumRepo) checkTagsExist(ctx context.Context, tags []string) error {
	if len(tags) == 0 {
		return nil
	}
	var count int64
	if err := r.db.WithContext(ctx).Model(&database.BlogTag{}).Where("tag IN ?", tags).Count(&count).Error; err != nil {
		return err
	}
	if int(count) != len(tags) {
		return forum.ErrTagNotExist
	}
	return nil
}

// incrTagCount 标签引用数 +1。
func incrTagCount(tx *gorm.DB, tags []string) error {
	for _, tag := range tags {
		if err := tx.Model(&database.BlogTag{}).Where("tag = ?", tag).
			Update("number", gorm.Expr("number + ?", 1)).Error; err != nil {
			return err
		}
	}
	return nil
}

// decrTagCount 标签引用数 -1。
func decrTagCount(tx *gorm.DB, tags []string) error {
	for _, tag := range tags {
		if err := tx.Model(&database.BlogTag{}).Where("tag = ?", tag).
			Update("number", gorm.Expr("number - ?", 1)).Error; err != nil {
			return err
		}
	}
	return nil
}

// Publish 发布帖子（事务：创建 + 标签引用数 +1）。
func (r *ForumRepo) Publish(ctx context.Context, p *forum.ForumPost) (uint, error) {
	if err := r.checkTagsExist(ctx, p.Tags); err != nil {
		return 0, err
	}
	var postID uint
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		post := database.ForumPost{
			UserID:   p.UserID,
			Title:    p.Title,
			Content:  p.Content,
			Category: p.Category,
			Tags:     database.JSONStringArray(p.Tags),
			Images:   database.JSONStringArray(p.Images),
		}
		if err := tx.Create(&post).Error; err != nil {
			return err
		}
		postID = post.ID
		return incrTagCount(tx, p.Tags)
	})
	return postID, err
}

// List 帖子分页列表。
func (r *ForumRepo) List(ctx context.Context, cond forum.ListCond) ([]*forum.ForumPost, int64, error) {
	db := r.db.WithContext(ctx).Model(&database.ForumPost{})
	if cond.Category != nil && *cond.Category != "" {
		db = db.Where("category = ?", *cond.Category)
	}
	if cond.Tag != nil && *cond.Tag != "" {
		db = db.Where("tags LIKE ?", "%\""+*cond.Tag+"\"%")
	}
	option := other.MySQLOption{
		PageInfo: pageInfoOf(cond.Page, cond.PageSize),
		Order:    "created_at DESC",
		Where:    db,
		Preload:  []string{"User"},
	}
	list, total, err := page.MySQLPagination(r.db.WithContext(ctx), &database.ForumPost{}, option)
	if err != nil {
		return nil, 0, err
	}
	out := make([]*forum.ForumPost, 0, len(list))
	for i := range list {
		out = append(out, toDomainForumPost(&list[i]))
	}
	return out, total, nil
}

// Detail 帖子详情：浏览数 +1（忽略错误），一级评论树含二级回复用户。
func (r *ForumRepo) Detail(ctx context.Context, id uint) (*forum.ForumPost, []*forum.ForumComment, error) {
	var post database.ForumPost
	if err := r.db.WithContext(ctx).Preload("User").First(&post, id).Error; err != nil {
		return nil, nil, err
	}
	r.db.WithContext(ctx).Model(&database.ForumPost{}).Where("id = ?", id).
		Update("view_count", gorm.Expr("view_count + ?", 1))

	var comments []database.ForumComment
	if err := r.db.WithContext(ctx).Where("post_id = ? AND parent_id = 0", id).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select(detailUserSelect)
		}).
		Preload("Children.User", func(db *gorm.DB) *gorm.DB {
			return db.Select(detailUserSelect)
		}).
		Order("created_at ASC").
		Find(&comments).Error; err != nil {
		return nil, nil, err
	}
	out := make([]*forum.ForumComment, 0, len(comments))
	for i := range comments {
		out = append(out, toDomainForumComment(&comments[i]))
	}
	return toDomainForumPost(&post), out, nil
}

// Like 切换点赞（事务）。
func (r *ForumRepo) Like(ctx context.Context, postID, userID uint) (bool, int, error) {
	var liked bool
	var likeCount int
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var like database.ForumLike
		dlErr := tx.Where("post_id = ? AND user_id = ?", postID, userID).First(&like).Error
		if errors.Is(dlErr, gorm.ErrRecordNotFound) {
			if cErr := tx.Create(&database.ForumLike{PostID: postID, UserID: userID}).Error; cErr != nil {
				return cErr
			}
			liked = true
			if uErr := tx.Model(&database.ForumPost{}).Where("id = ?", postID).
				Update("like_count", gorm.Expr("like_count + ?", 1)).Error; uErr != nil {
				return uErr
			}
		} else if dlErr != nil {
			return dlErr
		} else {
			if dErr := tx.Unscoped().Delete(&like).Error; dErr != nil {
				return dErr
			}
			liked = false
			if uErr := tx.Model(&database.ForumPost{}).Where("id = ?", postID).
				Update("like_count", gorm.Expr("like_count - ?", 1)).Error; uErr != nil {
				return uErr
			}
		}
		var post database.ForumPost
		if pErr := tx.Select("like_count").First(&post, postID).Error; pErr != nil {
			return pErr
		}
		likeCount = post.LikeCount
		return nil
	})
	return liked, likeCount, err
}

// Comment 发表评论（事务：校验帖子/父评论存在 + 创建 + comment_count +1）。
func (r *ForumRepo) Comment(ctx context.Context, c *forum.ForumComment) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var post database.ForumPost
		if err := tx.Select("id").First(&post, c.PostID).Error; err != nil {
			return forum.ErrPostNotExist
		}
		if c.ParentID != 0 {
			var parent database.ForumComment
			if err := tx.Select("id").First(&parent, c.ParentID).Error; err != nil {
				return forum.ErrParentNotExist
			}
		}
		if err := tx.Create(&database.ForumComment{
			PostID:   c.PostID,
			ParentID: c.ParentID,
			UserID:   c.UserID,
			Content:  c.Content,
		}).Error; err != nil {
			return err
		}
		return tx.Model(&database.ForumPost{}).Where("id = ?", c.PostID).
			Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error
	})
}

// ManageList 帖子管理列表。
func (r *ForumRepo) ManageList(ctx context.Context, cond forum.ManageListCond) ([]*forum.ForumPost, int64, error) {
	db := r.db.WithContext(ctx).Model(&database.ForumPost{})
	if cond.RoleID != shared.Admin {
		db = db.Where("user_id = ?", cond.UserID)
	}
	if cond.Title != nil && *cond.Title != "" {
		db = db.Where("title LIKE ?", "%"+*cond.Title+"%")
	}
	option := other.MySQLOption{
		PageInfo: pageInfoOf(cond.Page, cond.PageSize),
		Order:    "created_at DESC",
		Where:    db,
		Preload:  []string{"User"},
	}
	list, total, err := page.MySQLPagination(r.db.WithContext(ctx), &database.ForumPost{}, option)
	if err != nil {
		return nil, 0, err
	}
	out := make([]*forum.ForumPost, 0, len(list))
	for i := range list {
		out = append(out, toDomainForumPost(&list[i]))
	}
	return out, total, nil
}

// DeletePosts 删除帖子（事务）。
func (r *ForumRepo) DeletePosts(ctx context.Context, ids []uint, userID uint, roleID shared.RoleID) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var posts []database.ForumPost
		postQuery := tx.Where("id IN ?", ids)
		if roleID != shared.Admin {
			postQuery = postQuery.Where("user_id = ?", userID)
		}
		if err := postQuery.Find(&posts).Error; err != nil {
			return err
		}
		if len(posts) != len(ids) {
			return forum.ErrDeleteNotOwned
		}
		if err := tx.Unscoped().Where("post_id IN ?", ids).Delete(&database.ForumLike{}).Error; err != nil {
			return err
		}
		if err := tx.Where("post_id IN ?", ids).Delete(&database.ForumComment{}).Error; err != nil {
			return err
		}
		if err := tx.Where("id IN ?", ids).Delete(&database.ForumPost{}).Error; err != nil {
			return err
		}
		for _, post := range posts {
			if err := decrTagCount(tx, []string(post.Tags)); err != nil {
				return err
			}
		}
		return nil
	})
}

// ManageComments 评论管理列表（含帖子标题）。
func (r *ForumRepo) ManageComments(ctx context.Context, cond forum.ManageCommentCond) ([]*forum.ManageComment, int64, error) {
	db := r.db.WithContext(ctx).Model(&database.ForumComment{})
	if cond.PostID != nil {
		db = db.Where("post_id = ?", *cond.PostID)
	}
	if cond.RoleID != shared.Admin {
		db = db.Where("post_id IN (SELECT id FROM forum_posts WHERE user_id = ?)", cond.UserID)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	pageNum := cond.Page
	if pageNum < 1 {
		pageNum = 1
	}
	pageSize := cond.PageSize
	if pageSize < 1 {
		pageSize = 10
	}

	var comments []database.ForumComment
	if err := db.Preload("User", func(d *gorm.DB) *gorm.DB {
		return d.Select("id, uuid, username, avatar")
	}).Order("created_at DESC").
		Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	// 补充帖子标题
	postIDSet := map[uint]struct{}{}
	for i := range comments {
		postIDSet[comments[i].PostID] = struct{}{}
	}
	postIDs := make([]uint, 0, len(postIDSet))
	for id := range postIDSet {
		postIDs = append(postIDs, id)
	}
	var posts []database.ForumPost
	r.db.WithContext(ctx).Select("id, title").Where("id IN ?", postIDs).Find(&posts)
	postTitleMap := map[uint]string{}
	for _, p := range posts {
		postTitleMap[p.ID] = p.Title
	}

	list := make([]*forum.ManageComment, 0, len(comments))
	for i := range comments {
		item := &forum.ManageComment{ForumComment: *toDomainForumComment(&comments[i])}
		item.PostTitle = postTitleMap[comments[i].PostID]
		list = append(list, item)
	}
	return list, total, nil
}

// DeleteComments 删除评论（事务：软删 + comment_count 回退）。
func (r *ForumRepo) DeleteComments(ctx context.Context, ids []uint, userID uint, roleID shared.RoleID) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var comments []database.ForumComment
		commentQuery := tx.Where("id IN ?", ids)
		if roleID != shared.Admin {
			commentQuery = commentQuery.Where("post_id IN (SELECT id FROM forum_posts WHERE user_id = ?)", userID)
		}
		if err := commentQuery.Find(&comments).Error; err != nil {
			return err
		}
		if len(comments) != len(ids) {
			return forum.ErrCommentNotOwned
		}
		postCount := map[uint]int{}
		for _, c := range comments {
			postCount[c.PostID]++
		}
		if err := tx.Where("id IN ?", ids).Delete(&database.ForumComment{}).Error; err != nil {
			return err
		}
		for postID, n := range postCount {
			if err := tx.Model(&database.ForumPost{}).Where("id = ?", postID).
				Update("comment_count", gorm.Expr("comment_count - ?", n)).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
