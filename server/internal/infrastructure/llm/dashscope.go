// Package llm 提供大模型调用基础设施实现（阿里云百炼 DashScope OpenAI 兼容接口，SSE 流式）。
package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"server/config"
	"server/internal/domain/agent"

	"go.uber.org/zap"
)

// dashScopeChunk 流式响应中的单个 data 块
type dashScopeChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
	Usage *struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	} `json:"error"`
}

// DashScopeProvider 基于 DashScope OpenAI 兼容接口的大模型调用器。
type DashScopeProvider struct {
	baseURL string
	apiKey  string
	model   string
	client  *http.Client
	log     *zap.Logger
}

// NewDashScopeProvider 构造 DashScope 流式调用器。
func NewDashScopeProvider(cfg *config.LLM, log *zap.Logger) *DashScopeProvider {
	return &DashScopeProvider{
		baseURL: strings.TrimRight(cfg.BaseURL, "/"),
		apiKey:  cfg.APIKey,
		model:   cfg.Model,
		client:  &http.Client{Timeout: 5 * time.Minute},
		log:     log,
	}
}

// ChatStream 发起流式对话：逐块回调 onDelta，返回完整文本与 token 用量。
// ctx 取消时返回 ctx.Err()。
func (p *DashScopeProvider) ChatStream(ctx context.Context, req agent.ChatRequest, onDelta func(string) error) (*agent.ChatResult, error) {
	messages := make([]map[string]string, 0, len(req.Messages)+1)
	if req.System != "" {
		messages = append(messages, map[string]string{"role": "system", "content": req.System})
	}
	for _, m := range req.Messages {
		messages = append(messages, map[string]string{"role": m.Role, "content": m.Content})
	}

	body := map[string]interface{}{
		"model":    p.model,
		"messages": messages,
		"stream":   true,
		// 请求 usage 统计，流式结束时返回 usage 块
		"stream_options": map[string]bool{"include_usage": true},
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("dashscope: marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("dashscope: new request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("dashscope: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return nil, fmt.Errorf("dashscope: status %d: %s", resp.StatusCode, strings.TrimSpace(string(b)))
	}

	// ctx 取消时关闭响应体，中断阻塞的读取
	go func() {
		<-ctx.Done()
		resp.Body.Close()
	}()

	var full strings.Builder
	tokens := 0
	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "[DONE]" {
			break
		}
		var chunk dashScopeChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			// 忽略无法解析的块（部分代理可能夹带注释行）
			continue
		}
		if chunk.Error != nil {
			return nil, fmt.Errorf("dashscope: %s (code: %s)", chunk.Error.Message, chunk.Error.Code)
		}
		if len(chunk.Choices) > 0 {
			delta := chunk.Choices[0].Delta.Content
			if delta != "" {
				full.WriteString(delta)
				if onDelta != nil {
					if err := onDelta(delta); err != nil {
						return nil, err
					}
				}
			}
		}
		if chunk.Usage != nil {
			tokens = chunk.Usage.TotalTokens
		}
	}
	if err := scanner.Err(); err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, fmt.Errorf("dashscope: read stream: %w", err)
	}
	return &agent.ChatResult{FullText: full.String(), Tokens: tokens}, nil
}
