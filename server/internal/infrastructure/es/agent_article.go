package es

import (
	"context"
	"errors"
	"strconv"

	"server/internal/domain/agent"
	"server/internal/domain/article"
)

// AgentArticleReader 站内文章内容读取器（ES 实现，供 AI Agent 注入上下文）。
type AgentArticleReader struct {
	store article.EsArticleStore
}

// NewAgentArticleReader 构造文章读取器。
func NewAgentArticleReader(store article.EsArticleStore) *AgentArticleReader {
	return &AgentArticleReader{store: store}
}

// ReadByIDs 按 ID 批量读取文章标题与正文；已删除/不存在的文章自动跳过。
func (r *AgentArticleReader) ReadByIDs(ctx context.Context, ids []uint) ([]*agent.ArticleChunk, error) {
	out := make([]*agent.ArticleChunk, 0, len(ids))
	for _, id := range ids {
		a, err := r.store.Get(ctx, strconv.FormatUint(uint64(id), 10))
		if err != nil {
			if errors.Is(err, article.ErrDocumentNotFound) {
				continue
			}
			return nil, err
		}
		out = append(out, &agent.ArticleChunk{ID: id, Title: a.Title, Content: a.Content})
	}
	return out, nil
}
