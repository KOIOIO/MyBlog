package database

// AgentMemory 用户画像记忆表（长期记忆，L3 基础版）
type AgentMemory struct {
	MODEL
	UserID  uint   `json:"user_id" gorm:"index"` // 所属用户 ID
	Content string `json:"content" gorm:"size:500"` // 记忆条目（如「用户偏好 Go 语言」）
	Source  string `json:"source" gorm:"size:20"`   // 来源：manual / auto
}
