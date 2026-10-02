package database

// AgentConversation AI Agent 会话表
type AgentConversation struct {
	MODEL
	UserID  uint   `json:"user_id" gorm:"index"`                 // 所属用户 ID
	UUID    string `json:"uuid" gorm:"type:char(36);unique"`     // 会话 UUID（对外标识）
	Title   string `json:"title" gorm:"size:100"`                // 会话标题（首条消息截取）
	Summary string `json:"summary" gorm:"type:text"`             // 会话摘要（画像提取来源，可空）
	Status  int    `json:"status" gorm:"default:0"`              // 0=活跃 1=归档
}
