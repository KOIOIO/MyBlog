// Package agent 提供 AI Agent 应用服务（用例编排，依赖领域端口与基础设施抽象）。
package agent

import (
	"context"
	"strings"

	"server/config"
	"server/internal/domain/agent"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

// memoryLimit 注入 system 的用户画像记忆条数上限。
const memoryLimit = 10

// Service AI Agent 应用服务。
type Service struct {
	cfg      *config.Config
	log      *zap.Logger
	convs    agent.ConversationRepository
	msgs     agent.MessageRepository
	mems     agent.MemoryRepository
	articles agent.ArticleContentReader
	model    agent.ModelProvider
}

// NewService 构造 Agent 应用服务。
func NewService(cfg *config.Config, log *zap.Logger, convs agent.ConversationRepository, msgs agent.MessageRepository, mems agent.MemoryRepository, articles agent.ArticleContentReader, model agent.ModelProvider) *Service {
	return &Service{cfg: cfg, log: log, convs: convs, msgs: msgs, mems: mems, articles: articles, model: model}
}

// CreateConversation 新建会话（标题留空，首条消息时填充）。
func (s *Service) CreateConversation(ctx context.Context, userID uint) (*agent.Conversation, error) {
	c := &agent.Conversation{
		UUID:   uuid.Must(uuid.NewV4()).String(),
		UserID: userID,
		Status: 0,
	}
	if err := s.convs.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

// ListConversations 查询用户会话列表。
func (s *Service) ListConversations(ctx context.Context, userID uint) ([]*agent.Conversation, error) {
	return s.convs.ListByUser(ctx, userID)
}

// ListMessages 查询会话全部消息（校验归属）。
func (s *Service) ListMessages(ctx context.Context, conversationID, userID uint) ([]*agent.Message, error) {
	if _, err := s.convs.GetByID(ctx, conversationID, userID); err != nil {
		return nil, err
	}
	return s.msgs.ListByConversation(ctx, conversationID)
}

// DeleteConversation 删除会话及其消息（校验归属）。
func (s *Service) DeleteConversation(ctx context.Context, conversationID, userID uint) error {
	if err := s.convs.Delete(ctx, conversationID, userID); err != nil {
		return err
	}
	return s.msgs.DeleteByConversation(ctx, conversationID)
}

// Chat 发送消息并流式返回模型回复。
// conversationID 为 0 时自动新建会话；文章 ID 列表注入上下文；onDelta 逐块回调。
func (s *Service) Chat(ctx context.Context, userID, conversationID uint, content string, articleIDs []uint, onDelta func(string) error) (string, error) {
	var cid = conversationID
	if cid == 0 {
		c, err := s.CreateConversation(ctx, userID)
		if err != nil {
			return "", err
		}
		cid = c.ID
		// 首条消息：以内容前 30 字作为会话标题
		if err := s.convs.UpdateTitle(ctx, c.ID, userID, agent.MakeTitle(content)); err != nil {
			return "", err
		}
	} else {
		c, err := s.convs.GetByID(ctx, cid, userID)
		if err != nil {
			return "", err
		}
		// 历史会话但标题为空（异常兜底），补填标题
		if c.Title == "" {
			if err := s.convs.UpdateTitle(ctx, c.ID, userID, agent.MakeTitle(content)); err != nil {
				return "", err
			}
		}
	}

	// 先取历史（不含当前消息），再落库当前 user 消息，保证失败可恢复
	history, err := s.msgs.ListByConversation(ctx, cid)
	if err != nil {
		return "", err
	}
	if err := s.msgs.Append(ctx, &agent.Message{
		ConversationID: cid,
		Role:           "user",
		Content:        content,
		ArticleIDs:     articleIDs,
	}); err != nil {
		return "", err
	}

	system, err := s.buildSystem(ctx, userID, articleIDs)
	if err != nil {
		return "", err
	}
	messages := s.buildMessages(history, content)

	result, err := s.model.ChatStream(ctx, agent.ChatRequest{System: system, Messages: messages, Stream: true}, onDelta)
	if err != nil {
		s.log.Error("agent chat stream failed", zap.Uint("conversation", cid), zap.Error(err))
		return "", err
	}

	if err := s.msgs.Append(ctx, &agent.Message{
		ConversationID: cid,
		Role:           "assistant",
		Content:        result.FullText,
		Tokens:         result.Tokens,
	}); err != nil {
		return "", err
	}
	return result.FullText, nil
}

// buildSystem 组装 system 提示：助手设定 + 用户画像记忆 + 站内文章上下文。
func (s *Service) buildSystem(ctx context.Context, userID uint, articleIDs []uint) (string, error) {
	var sb strings.Builder
	sb.WriteString("你是 MyBlog 站点的 AI 助手，回答使用简体中文，简洁、准确。")

	// L3 用户画像记忆
	mems, err := s.mems.ListByUser(ctx, userID)
	if err != nil {
		return "", err
	}
	if len(mems) > 0 {
		sb.WriteString("\n\n[用户画像记忆]")
		for i, m := range mems {
			if i >= memoryLimit {
				break
			}
			sb.WriteString("\n- " + m.Content)
		}
	}

	// 站内文章上下文
	if len(articleIDs) > 0 {
		chunks, err := s.articles.ReadByIDs(ctx, articleIDs)
		if err != nil {
			return "", err
		}
		if len(chunks) > 0 {
			sb.WriteString("\n\n[站内文章上下文]（回答请优先基于以下站内文章内容）")
			for _, ch := range chunks {
				sb.WriteString("\n\n文章《" + ch.Title + "》：\n" + truncateRunes(ch.Content, s.cfg.LLM.MaxArticleChars))
			}
		}
	}
	return sb.String(), nil
}

// buildMessages 组装对话消息：历史（最近 MaxHistory 条）+ 当前用户消息。
func (s *Service) buildMessages(history []*agent.Message, current string) []agent.ChatMessage {
	max := s.cfg.LLM.MaxHistory
	if max <= 0 {
		max = 20
	}
	msgs := make([]agent.ChatMessage, 0, len(history)+1)
	start := 0
	if len(history) > max {
		start = len(history) - max
	}
	for _, m := range history[start:] {
		msgs = append(msgs, agent.ChatMessage{Role: m.Role, Content: m.Content})
	}
	msgs = append(msgs, agent.ChatMessage{Role: "user", Content: current})
	return msgs
}

// truncateRunes 按 rune 数截断字符串。
func truncateRunes(s string, max int) string {
	if max <= 0 {
		return s
	}
	r := []rune(s)
	if len(r) > max {
		return string(r[:max])
	}
	return s
}
