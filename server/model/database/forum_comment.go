package database

import "server/global"

// ForumComment 论坛评论表（二级嵌套）
type ForumComment struct {
	global.MODEL
	PostID   uint           `json:"post_id"`                       // 所属帖子 ID
	ParentID uint           `json:"parent_id"`                     // 父评论 ID，0 表示一级评论
	UserID   uint           `json:"user_id"`                       // 评论用户 ID
	User     User           `json:"user" gorm:"foreignKey:UserID"` // 关联的用户
	Content  string         `json:"content"`                       // 评论内容
	Children []ForumComment `json:"children" gorm:"foreignKey:ParentID"` // 二级回复
}
