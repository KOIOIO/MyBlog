package mysql

import (
	"context"
	"errors"

	"server/internal/domain/agent"
	"server/internal/model/database"

	"gorm.io/gorm"
)

// AgentConversationRepository 会话仓储（GORM 实现）。
type AgentConversationRepository struct {
	db *gorm.DB
}

// NewAgentConversationRepo 构造会话仓储。
func NewAgentConversationRepo(db *gorm.DB) *AgentConversationRepository {
	return &AgentConversationRepository{db: db}
}

// Create 新建会话。
func (r *AgentConversationRepository) Create(ctx context.Context, c *agent.Conversation) error {
	row := database.AgentConversation{
		UserID:  c.UserID,
		UUID:    c.UUID,
		Title:   c.Title,
		Summary: c.Summary,
		Status:  c.Status,
	}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return err
	}
	c.ID = row.ID
	c.CreatedAt = row.CreatedAt
	c.UpdatedAt = row.UpdatedAt
	return nil
}

// ListByUser 查询指定用户的全部会话（按更新时间倒序）。
func (r *AgentConversationRepository) ListByUser(ctx context.Context, userID uint) ([]*agent.Conversation, error) {
	var rows []database.AgentConversation
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("updated_at DESC").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]*agent.Conversation, 0, len(rows))
	for i := range rows {
		out = append(out, &agent.Conversation{
			ID:        rows[i].ID,
			UUID:      rows[i].UUID,
			UserID:    rows[i].UserID,
			Title:     rows[i].Title,
			Summary:   rows[i].Summary,
			Status:    rows[i].Status,
			CreatedAt: rows[i].CreatedAt,
			UpdatedAt: rows[i].UpdatedAt,
		})
	}
	return out, nil
}

// GetByID 按 ID 查询会话；不存在返回 ErrConversationNotFound，归属不符返回 ErrForbidden。
func (r *AgentConversationRepository) GetByID(ctx context.Context, id, userID uint) (*agent.Conversation, error) {
	var row database.AgentConversation
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, agent.ErrConversationNotFound
		}
		return nil, err
	}
	if row.UserID != userID {
		return nil, agent.ErrForbidden
	}
	return &agent.Conversation{
		ID:        row.ID,
		UUID:      row.UUID,
		UserID:    row.UserID,
		Title:     row.Title,
		Summary:   row.Summary,
		Status:    row.Status,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}, nil
}

// Delete 删除会话（软删）。
func (r *AgentConversationRepository) Delete(ctx context.Context, id, userID uint) error {
	res := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&database.AgentConversation{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return agent.ErrConversationNotFound
	}
	return nil
}

// AgentMessageRepository 消息仓储（GORM 实现）。
type AgentMessageRepository struct {
	db *gorm.DB
}

// NewAgentMessageRepo 构造消息仓储。
func NewAgentMessageRepo(db *gorm.DB) *AgentMessageRepository {
	return &AgentMessageRepository{db: db}
}

// Append 追加一条消息。
func (r *AgentMessageRepository) Append(ctx context.Context, m *agent.Message) error {
	row := database.AgentMessage{
		ConversationID: m.ConversationID,
		Role:           m.Role,
		Content:        m.Content,
		ArticleIDs:     database.JSONUintArray(m.ArticleIDs),
		Tokens:         m.Tokens,
	}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return err
	}
	m.ID = row.ID
	m.CreatedAt = row.CreatedAt
	return nil
}

// ListByConversation 查询会话全部消息（按创建时间正序）。
func (r *AgentMessageRepository) ListByConversation(ctx context.Context, conversationID uint) ([]*agent.Message, error) {
	var rows []database.AgentMessage
	if err := r.db.WithContext(ctx).
		Where("conversation_id = ?", conversationID).
		Order("created_at ASC, id ASC").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]*agent.Message, 0, len(rows))
	for i := range rows {
		out = append(out, &agent.Message{
			ID:             rows[i].ID,
			ConversationID: rows[i].ConversationID,
			Role:           rows[i].Role,
			Content:        rows[i].Content,
			ArticleIDs:     []uint(rows[i].ArticleIDs),
			Tokens:         rows[i].Tokens,
			CreatedAt:      rows[i].CreatedAt,
		})
	}
	return out, nil
}

// DeleteByConversation 删除会话下的全部消息（软删）。
func (r *AgentMessageRepository) DeleteByConversation(ctx context.Context, conversationID uint) error {
	return r.db.WithContext(ctx).
		Where("conversation_id = ?", conversationID).
		Delete(&database.AgentMessage{}).Error
}

// AgentMemoryRepository 用户画像记忆仓储（GORM 实现）。
type AgentMemoryRepository struct {
	db *gorm.DB
}

// NewAgentMemoryRepo 构造画像记忆仓储。
func NewAgentMemoryRepo(db *gorm.DB) *AgentMemoryRepository {
	return &AgentMemoryRepository{db: db}
}

// ListByUser 查询用户画像记忆（按创建时间倒序）。
func (r *AgentMemoryRepository) ListByUser(ctx context.Context, userID uint) ([]*agent.Memory, error) {
	var rows []database.AgentMemory
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]*agent.Memory, 0, len(rows))
	for i := range rows {
		out = append(out, &agent.Memory{
			ID:        rows[i].ID,
			UserID:    rows[i].UserID,
			Content:   rows[i].Content,
			Source:    rows[i].Source,
			CreatedAt: rows[i].CreatedAt,
		})
	}
	return out, nil
}

// Append 追加一条画像记忆。
func (r *AgentMemoryRepository) Append(ctx context.Context, m *agent.Memory) error {
	row := database.AgentMemory{
		UserID:  m.UserID,
		Content: m.Content,
		Source:  m.Source,
	}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return err
	}
	m.ID = row.ID
	m.CreatedAt = row.CreatedAt
	return nil
}
