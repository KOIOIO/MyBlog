package database

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"server/global"
)

// JSONStringArray 存储 JSON 数组文本，读写自动编解码（如 ["Go","Vue"]）
type JSONStringArray []string

// Value 实现 driver.Valuer，写入时序列化为 JSON 文本
func (a JSONStringArray) Value() (driver.Value, error) {
	if a == nil {
		return "[]", nil
	}
	return json.Marshal(a)
}

// Scan 实现 sql.Scanner，读取时解析 JSON 文本为字符串数组
func (a *JSONStringArray) Scan(value interface{}) error {
	if value == nil {
		*a = JSONStringArray{}
		return nil
	}
	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return json.Unmarshal(data, a)
}

// ForumPost 论坛帖子表
type ForumPost struct {
	global.MODEL
	UserID       uint            `json:"user_id"`                       // 发帖用户 ID
	User         User            `json:"user" gorm:"foreignKey:UserID"` // 关联的用户
	Title        string          `json:"title"`                         // 标题
	Content      string          `json:"content"`                       // 内容
	Category     string          `json:"category"`                      // 分类：技术 / 生活
	Tags         JSONStringArray `json:"tags"`                          // 标签数组
	Images       JSONStringArray `json:"images"`                        // 图片 URL 数组
	LikeCount    int             `json:"like_count"`                    // 点赞数
	CommentCount int             `json:"comment_count"`                 // 评论数
	ViewCount    int             `json:"view_count"`                    // 浏览数
}
