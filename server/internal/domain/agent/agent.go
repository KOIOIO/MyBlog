// Package agent 提供 AI Agent 领域模型（实体、规则、端口）。
// 本包不依赖任何基础设施库。
package agent

import (
	"context"
	"errors"
	"time"
)

var (
	// ErrConversationNotFound 会话不存在
	ErrConversationNotFound = errors.New("conversation not found")
	// ErrForbidden 资源归属校验失败（越权访问）
	ErrForbidden = errors.New("forbidden")
)

// MakeTitle 取消息前 30 个 rune 作为会话标题。
func MakeTitle(content string) string {
	r := []rune(content)
	if len(r) > 30 {
		return string(r[:30])
	}
	return content
}

// ValidateAccess 校验资源归属，防止越权访问他人数据。
func ValidateAccess(resourceUserID, callerUserID uint) error {
	if resourceUserID != callerUserID {
		return ErrForbidden
	}
	return nil
}

// Conversation 会话实体
type Conversation struct {
	ID        uint      `json:"id"`
	UUID      string    `json:"uuid"`
	UserID    uint      `json:"-"`
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Message 会话消息实体
type Message struct {
	ID             uint      `json:"id"`
	ConversationID uint      `json:"conversation_id"`
	Role           string    `json:"role"`
	Content        string    `json:"content"`
	ArticleIDs     []uint    `json:"article_ids,omitempty"`
	Tokens         int       `json:"-"`
	CreatedAt      time.Time `json:"created_at"`
}

// Memory 用户画像记忆条目
type Memory struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"-"`
	Content   string    `json:"content"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"created_at"`
}

// ---- 端口（Ports） ----

// ConversationRepository 会话仓储端口
type ConversationRepository interface {
	Create(ctx context.Context, c *Conversation) error
	ListByUser(ctx context.Context, userID uint) ([]*Conversation, error)
	GetByID(ctx context.Context, id, userID uint) (*Conversation, error)
	Delete(ctx context.Context, id, userID uint) error
}

// MessageRepository 消息仓储端口
type MessageRepository interface {
	Append(ctx context.Context, m *Message) error
	ListByConversation(ctx context.Context, conversationID uint) ([]*Message, error)
	DeleteByConversation(ctx context.Context, conversationID uint) error
}

// MemoryRepository 用户画像记忆仓储端口
type MemoryRepository interface {
	ListByUser(ctx context.Context, userID uint) ([]*Memory, error)
	Append(ctx context.Context, m *Memory) error
}

// ArticleChunk 注入上下文的文章片段
type ArticleChunk struct {
	ID      uint
	Title   string
	Content string
}

// ArticleContentReader 站内文章内容读取端口（按 ID 批量）
type ArticleContentReader interface {
	ReadByIDs(ctx context.Context, ids []uint) ([]*ArticleChunk, error)
}

// ChatMessage 发给模型的单条消息
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest 模型调用请求
type ChatRequest struct {
	System   string        `json:"system"`
	Messages []ChatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

// ChatResult 模型调用结果
type ChatResult struct {
	FullText string
	Tokens   int
}

// ModelProvider 大模型调用端口（流式）
type ModelProvider interface {
	ChatStream(ctx context.Context, req ChatRequest, onDelta func(string) error) (*ChatResult, error)
}
