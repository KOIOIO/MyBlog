package database

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// JSONUintArray 存储 JSON 数字数组文本，读写自动编解码（如 [1,5,9]）
type JSONUintArray []uint

// Value 实现 driver.Valuer，写入时序列化为 JSON 文本
func (a JSONUintArray) Value() (driver.Value, error) {
	if a == nil {
		return "[]", nil
	}
	return json.Marshal(a)
}

// Scan 实现 sql.Scanner，读取时解析 JSON 文本为数字数组
func (a *JSONUintArray) Scan(value interface{}) error {
	if value == nil {
		*a = JSONUintArray{}
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

// AgentMessage AI Agent 会话消息表
type AgentMessage struct {
	MODEL
	ConversationID uint           `json:"conversation_id" gorm:"index"` // 所属会话 ID
	Role           string         `json:"role" gorm:"size:10"`          // user / assistant
	Content        string         `json:"content" gorm:"type:text"`     // 消息正文
	ArticleIDs     JSONUintArray  `json:"article_ids"`                  // 本条消息携带的文章 ID 列表
	Tokens         int            `json:"tokens"`                       // 消耗 token（统计用）
}
