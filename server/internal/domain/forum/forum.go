// Package forum 提供论坛领域模型（实体、规则、端口）。
// 本包不依赖任何基础设施库。
package forum

import (
	"context"
	"errors"
	"time"

	"server/internal/domain/article"
	"server/internal/domain/shared"

	"github.com/gofrs/uuid"
)

// UserBrief 帖子/评论关联的用户（JSON 与 database.User 全字段形状一致）。
type UserBrief struct {
	ID        uint            `json:"id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	UUID      uuid.UUID       `json:"uuid"`
	Username  string          `json:"username"`
	Email     string          `json:"email"`
	Openid    string          `json:"openid"`
	Avatar    string          `json:"avatar"`
	Address   string          `json:"address"`
	Signature string          `json:"signature"`
	RoleID    shared.RoleID   `json:"role_id"`
	Register  shared.Register `json:"register"`
	Freeze    bool            `json:"freeze"`
}

// ForumPost 论坛帖子实体（JSON 字段顺序与 database.ForumPost 一致）。
type ForumPost struct {
	ID           uint      `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	UserID       uint      `json:"user_id"`
	User         UserBrief `json:"user"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	Category     string    `json:"category"`
	Tags         []string  `json:"tags"`
	Images       []string  `json:"images"`
	LikeCount    int       `json:"like_count"`
	CommentCount int       `json:"comment_count"`
	ViewCount    int       `json:"view_count"`
}

// ForumComment 论坛评论实体（二级嵌套）。
type ForumComment struct {
	ID        uint            `json:"id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	PostID    uint            `json:"post_id"`
	ParentID  uint            `json:"parent_id"`
	UserID    uint            `json:"user_id"`
	User      UserBrief       `json:"user"`
	Content   string          `json:"content"`
	Children  []*ForumComment `json:"children"`
}

// ForumLike 论坛点赞记录。
type ForumLike struct {
	PostID uint `json:"post_id"`
	UserID uint `json:"user_id"`
}

// ForumDetail 帖子详情返回结构（帖子 + 评论树）。
type ForumDetail struct {
	ForumPost
	Comments []*ForumComment `json:"comments"`
}

// ManageComment 评论管理列表项（含帖子标题）。
type ManageComment struct {
	ForumComment
	PostTitle string `json:"post_title"`
}

// ListCond 帖子列表条件。
type ListCond struct {
	Category *string
	Tag      *string
	Page     int
	PageSize int
}

// ManageListCond 帖子管理列表条件。
type ManageListCond struct {
	Title    *string
	UserID   uint
	RoleID   shared.RoleID
	Page     int
	PageSize int
}

// ManageCommentCond 评论管理列表条件。
type ManageCommentCond struct {
	PostID   *uint
	UserID   uint
	RoleID   shared.RoleID
	Page     int
	PageSize int
}

// 领域错误。
var (
	ErrCategoryNotAllowed = errors.New("分类只能是 技术 或 生活")
	ErrTagNotExist        = errors.New("标签不存在")
	ErrPostNotExist       = errors.New("帖子不存在")
	ErrParentNotExist     = errors.New("父评论不存在")
	ErrDeleteNotOwned     = errors.New("无权删除非本人发布的帖子")
	ErrCommentNotOwned    = errors.New("无权删除非本人帖子下的评论")
)

// AllowedCategory 允许的帖子分类集合。
var AllowedCategory = map[string]struct{}{"技术": {}, "生活": {}}

// Tag 固定标签库标签（复用文章域标签模型）。
type Tag = article.BlogTag

// ForumRepository 论坛仓储端口（事务行为收敛于此）。
type ForumRepository interface {
	// Tags 固定标签库全部标签（group ASC, number DESC, tag ASC）。
	Tags(ctx context.Context) ([]Tag, error)
	// Publish 发布帖子（事务：创建帖子 + 标签引用数 +1），返回新帖子 ID。
	Publish(ctx context.Context, post *ForumPost) (uint, error)
	// List 帖子分页列表（category/tag 过滤 + created_at DESC + 预加载用户）。
	List(ctx context.Context, cond ListCond) ([]*ForumPost, int64, error)
	// Detail 帖子详情：预加载用户、浏览数 +1（忽略错误）、一级评论树（含二级回复用户预加载）。
	Detail(ctx context.Context, id uint) (*ForumPost, []*ForumComment, error)
	// Like 切换点赞（事务），返回是否已点赞与最新点赞数。
	Like(ctx context.Context, postID, userID uint) (bool, int, error)
	// Comment 发表评论（事务：校验帖子/父评论存在 + 创建 + comment_count +1）。
	Comment(ctx context.Context, c *ForumComment) error
	// ManageList 帖子管理列表（非管理员仅本人帖子 + title 过滤）。
	ManageList(ctx context.Context, cond ManageListCond) ([]*ForumPost, int64, error)
	// DeletePosts 删除帖子（事务：所有权校验 + 物理删点赞 + 软删评论/帖子 + 标签计数回退）。
	DeletePosts(ctx context.Context, ids []uint, userID uint, roleID shared.RoleID) error
	// ManageComments 评论管理列表（非管理员仅本人帖子下的评论 + 补充帖子标题）。
	ManageComments(ctx context.Context, cond ManageCommentCond) ([]*ManageComment, int64, error)
	// DeleteComments 删除评论（事务：所有权校验 + 软删 + 帖子 comment_count 回退）。
	DeleteComments(ctx context.Context, ids []uint, userID uint, roleID shared.RoleID) error
}
