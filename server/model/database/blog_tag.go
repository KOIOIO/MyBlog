package database

// BlogTag 固定标签库
type BlogTag struct {
	Tag    string `json:"tag" gorm:"primaryKey"` // 标签名
	Group  string `json:"group"`                 // 分组：tech / life
	Number int    `json:"number"`                // 被引用次数
}
