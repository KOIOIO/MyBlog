package llm

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"server/config"
	"server/internal/domain/agent"

	"go.uber.org/zap"
)

func TestDashScopeChatStream(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			t.Error("missing Authorization header")
		}
		if !strings.Contains(r.URL.Path, "chat/completions") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Write([]byte("data: {\"choices\":[{\"delta\":{\"content\":\"你\"}}]}\n\n"))
		w.Write([]byte("data: {\"choices\":[{\"delta\":{\"content\":\"好\"}}]}\n\n"))
		w.Write([]byte("data: {\"choices\":[],\"usage\":{\"prompt_tokens\":16,\"completion_tokens\":7,\"total_tokens\":23}}\n\n"))
		w.Write([]byte("data: [DONE]\n\n"))
	}))
	defer srv.Close()

	p := NewDashScopeProvider(&config.LLM{BaseURL: srv.URL, APIKey: "test-key", Model: "qwen-max"}, zap.NewNop())
	var got strings.Builder
	res, err := p.ChatStream(context.Background(), agent.ChatRequest{
		System:   "你是助手",
		Messages: []agent.ChatMessage{{Role: "user", Content: "hi"}},
		Stream:   true,
	}, func(d string) error { got.WriteString(d); return nil })
	if err != nil {
		t.Fatal(err)
	}
	if got.String() != "你好" {
		t.Fatalf("want 你好, got %q", got.String())
	}
	if res.FullText != "你好" {
		t.Fatalf("full text mismatch: %q", res.FullText)
	}
	if res.Tokens != 23 {
		t.Fatalf("want tokens 23, got %d", res.Tokens)
	}
}

func TestDashScopeChatStreamError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":{"message":"The model does not exist","code":"model_not_found"}}`))
	}))
	defer srv.Close()

	p := NewDashScopeProvider(&config.LLM{BaseURL: srv.URL, APIKey: "test-key", Model: "qwen-max"}, zap.NewNop())
	_, err := p.ChatStream(context.Background(), agent.ChatRequest{
		Messages: []agent.ChatMessage{{Role: "user", Content: "hi"}},
		Stream:   true,
	}, nil)
	if err == nil {
		t.Fatal("want error, got nil")
	}
	if !strings.Contains(err.Error(), "model_not_found") {
		t.Fatalf("want error mentioning model_not_found, got: %v", err)
	}
}

func TestDashScopeChatStreamContextCancel(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		// 模拟持续输出，等待 context 取消
		w.(http.Flusher).Flush()
		<-r.Context().Done()
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	p := NewDashScopeProvider(&config.LLM{BaseURL: srv.URL, APIKey: "test-key", Model: "qwen-max"}, zap.NewNop())
	done := make(chan error, 1)
	go func() {
		_, err := p.ChatStream(ctx, agent.ChatRequest{
			Messages: []agent.ChatMessage{{Role: "user", Content: "hi"}},
			Stream:   true,
		}, nil)
		done <- err
	}()
	cancel()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("want context error, got nil")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("chat stream did not return after context cancel")
	}
}
