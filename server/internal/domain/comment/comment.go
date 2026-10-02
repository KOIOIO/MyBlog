// Package comment 提供评论领域模型（实体、规则、端口）。
// 本包不依赖任何基础设施库。
package comment

import (
	"context"
	"errors"
	"time"

	"server/internal/domain/shared"

	"github.com/gofrs/uuid"
)

// UserBrief 评论关联的用户（JSON 与 database.User 全字段形状一致，未加载字段为零值）。
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

// Comment 评论实体（JSON 字段顺序与 database.Comment 一致）。
// Children 未加载时为 nil（序列化 null），加载后为 []（序列化 []）。
type Comment struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	ArticleID string     `json:"article_id"`
	PID       *uint      `json:"p_id"`
	Children  []*Comment `json:"children"`
	UserUUID  uuid.UUID  `json:"user_uuid"`
	User      UserBrief  `json:"user"`
	Content   string     `json:"content"`
}

// ListCond 评论管理列表查询条件。
type ListCond struct {
	ArticleID *string
	UserUUID  *string
	Content   *string
	Page      int
	PageSize  int
}

// ErrNoPermission 删除权限不足。
var ErrNoPermission = errors.New("you do not have permission to delete this comment")

// CommentRepository 评论仓储端口（树加载与级联删除行为收敛于此）。
type CommentRepository interface {
	// ByArticle 文章一级评论并递归加载子评论（children 序列化 []）。
	ByArticle(ctx context.Context, articleID string) ([]*Comment, error)
	// Newest 最新 N 条评论（不加载子评论，children 为 nil）。
	Newest(ctx context.Context, limit int) ([]*Comment, error)
	// Create 创建评论（database.Comment 钩子会同步 ES 文章评论数）。
	Create(ctx context.Context, c *Comment) error
	// DeleteTree 校验权限并事务级联删除（含 ES 评论数同步钩子）。
	DeleteTree(ctx context.Context, ids []uint, userUUID uuid.UUID, roleID shared.RoleID) error
	// ByUser 用户全部评论（含子评论树）。
	ByUser(ctx context.Context, userUUID uuid.UUID) ([]*Comment, error)
	// List 管理列表分页（不预加载用户，保持旧行为）。
	List(ctx context.Context, cond ListCond) ([]*Comment, int64, error)
}

// DedupByRootUser 评论去重规则：若某评论作为另一条评论（同用户）的子孙评论出现，
// 则该评论从根列表中剔除（行为与原 FindChildCommentsIDByRootCommentUserUUID 一致，
// 仅过滤根评论自身 ID，不修改树的 children 结构）。
func DedupByRootUser(rootComments []*Comment) []*Comment {
	excluded := map[uint]struct{}{}
	for _, root := range rootComments {
		var walk func(children []*Comment)
		walk = func(children []*Comment) {
			for _, child := range children {
				if child.UserUUID == root.UserUUID {
					excluded[child.ID] = struct{}{}
				}
				if len(child.Children) > 0 {
					walk(child.Children)
				}
			}
		}
		walk(root.Children)
	}

	out := make([]*Comment, 0, len(rootComments))
	for _, c := range rootComments {
		if _, exists := excluded[c.ID]; !exists {
			out = append(out, c)
		}
	}
	return out
}
