package request

import "server/model/appTypes"

// ForumPublish 发布帖子
type ForumPublish struct {
	Title    string   `json:"title" binding:"required"`
	Content  string   `json:"content" binding:"required"`
	Category string   `json:"category" binding:"required"`
	Tags     []string `json:"tags"`
	Images   []string `json:"images"`
	UserID   uint     `json:"-"`
}

// ForumList 帖子列表
type ForumList struct {
	Category *string `json:"category" form:"category"`
	Tag      *string `json:"tag" form:"tag"`
	PageInfo
}

// ForumDetail 帖子详情
type ForumDetail struct {
	ID uint `json:"id" form:"id" uri:"id" binding:"required"`
}

// ForumLike 点赞/取消点赞
type ForumLike struct {
	PostID uint `json:"post_id" form:"post_id" binding:"required"`
	UserID uint `json:"-"`
}

// ForumComment 发表评论
type ForumComment struct {
	PostID   uint   `json:"post_id" binding:"required"`
	ParentID uint   `json:"parent_id"`
	Content  string `json:"content" binding:"required"`
	UserID   uint   `json:"-"`
}

// ForumDelete 批量删除帖子
type ForumDelete struct {
	IDs    []uint           `json:"ids" binding:"required"`
	UserID uint             `json:"-"`
	RoleID appTypes.RoleID  `json:"-"`
}

// ForumManageList 帖子管理列表（管理员全量，普通用户仅自己的）
type ForumManageList struct {
	Title  *string          `json:"title" form:"title"`
	UserID uint             `json:"-"`
	RoleID appTypes.RoleID  `json:"-"`
	PageInfo
}

// ForumManageComments 评论管理列表（管理员全量，普通用户仅自己帖子下的）
type ForumManageComments struct {
	PostID *uint           `json:"post_id" form:"post_id"`
	UserID uint            `json:"-"`
	RoleID appTypes.RoleID `json:"-"`
	PageInfo
}

// ForumCommentDelete 批量删除评论
type ForumCommentDelete struct {
	IDs    []uint           `json:"ids" binding:"required"`
	UserID uint             `json:"-"`
	RoleID appTypes.RoleID  `json:"-"`
}
