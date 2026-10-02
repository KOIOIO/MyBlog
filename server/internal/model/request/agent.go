package request

// AgentChat AI Agent 对话请求
type AgentChat struct {
	ConversationID uint   `json:"conversation_id"` // 会话 ID，0 表示新建会话
	Message        string `json:"message"`         // 用户消息
	ArticleIDs     []uint `json:"article_ids"`     // 携带的站内文章 ID 列表（可空）
}
