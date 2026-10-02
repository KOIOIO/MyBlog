package agent

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"

	"server/config"
	"server/internal/domain/agent"

	"go.uber.org/zap"
)

// ---- fake 仓储 ----

type fakeConvs struct {
	mu   sync.Mutex
	seq  uint
	rows map[uint]*agent.Conversation
}

func newFakeConvs() *fakeConvs {
	return &fakeConvs{rows: map[uint]*agent.Conversation{}}
}

func (f *fakeConvs) Create(_ context.Context, c *agent.Conversation) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.seq++
	c.ID = f.seq
	f.rows[c.ID] = c
	return nil
}

func (f *fakeConvs) ListByUser(_ context.Context, userID uint) ([]*agent.Conversation, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []*agent.Conversation
	for _, c := range f.rows {
		if c.UserID == userID {
			out = append(out, c)
		}
	}
	return out, nil
}

func (f *fakeConvs) GetByID(_ context.Context, id, userID uint) (*agent.Conversation, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	c, ok := f.rows[id]
	if !ok {
		return nil, agent.ErrConversationNotFound
	}
	if c.UserID != userID {
		return nil, agent.ErrForbidden
	}
	return c, nil
}

func (f *fakeConvs) UpdateTitle(_ context.Context, id, userID uint, title string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	c, ok := f.rows[id]
	if !ok {
		return agent.ErrConversationNotFound
	}
	if c.UserID != userID {
		return agent.ErrForbidden
	}
	c.Title = title
	return nil
}

func (f *fakeConvs) Delete(_ context.Context, id, userID uint) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	c, ok := f.rows[id]
	if !ok {
		return agent.ErrConversationNotFound
	}
	if c.UserID != userID {
		return agent.ErrForbidden
	}
	delete(f.rows, id)
	return nil
}

type fakeMsgs struct {
	mu   sync.Mutex
	seq  uint
	rows map[uint][]*agent.Message
}

func newFakeMsgs() *fakeMsgs {
	return &fakeMsgs{rows: map[uint][]*agent.Message{}}
}

func (f *fakeMsgs) Append(_ context.Context, m *agent.Message) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.seq++
	m.ID = f.seq
	f.rows[m.ConversationID] = append(f.rows[m.ConversationID], m)
	return nil
}

func (f *fakeMsgs) ListByConversation(_ context.Context, conversationID uint) ([]*agent.Message, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]*agent.Message{}, f.rows[conversationID]...), nil
}

func (f *fakeMsgs) DeleteByConversation(_ context.Context, conversationID uint) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.rows, conversationID)
	return nil
}

type fakeMems struct {
	rows map[uint][]*agent.Memory
}

func newFakeMems() *fakeMems {
	return &fakeMems{rows: map[uint][]*agent.Memory{}}
}

func (f *fakeMems) ListByUser(_ context.Context, userID uint) ([]*agent.Memory, error) {
	return append([]*agent.Memory{}, f.rows[userID]...), nil
}

func (f *fakeMems) Append(_ context.Context, m *agent.Memory) error {
	f.rows[m.UserID] = append(f.rows[m.UserID], m)
	return nil
}

type fakeArticles struct {
	chunks map[uint]*agent.ArticleChunk
}

func newFakeArticles() *fakeArticles {
	return &fakeArticles{chunks: map[uint]*agent.ArticleChunk{
		1: {ID: 1, Title: "Go 并发模型", Content: "goroutine 与 channel 是核心。"},
		2: {ID: 2, Title: "MySQL 索引", Content: "B+ 树结构详解。"},
	}}
}

func (f *fakeArticles) ReadByIDs(_ context.Context, ids []uint) ([]*agent.ArticleChunk, error) {
	var out []*agent.ArticleChunk
	for _, id := range ids {
		if c, ok := f.chunks[id]; ok {
			out = append(out, c)
		}
	}
	return out, nil
}

type fakeModel struct {
	mu      sync.Mutex
	calls   []agent.ChatRequest
	reply   string
	tokens  int
	err     error
	onDelta bool
}

func newFakeModel(reply string, tokens int) *fakeModel {
	return &fakeModel{reply: reply, tokens: tokens}
}

func (f *fakeModel) ChatStream(_ context.Context, req agent.ChatRequest, onDelta func(string) error) (*agent.ChatResult, error) {
	f.mu.Lock()
	f.calls = append(f.calls, req)
	f.mu.Unlock()
	if f.err != nil {
		return nil, f.err
	}
	// 模拟真实流式：逐字回调
	if onDelta != nil {
		for _, r := range f.reply {
			if err := onDelta(string(r)); err != nil {
				return nil, err
			}
		}
	}
	return &agent.ChatResult{FullText: f.reply, Tokens: f.tokens}, nil
}

func newTestService(convs *fakeConvs, msgs *fakeMsgs, mems *fakeMems, articles *fakeArticles, model *fakeModel) *Service {
	cfg := &config.Config{}
	cfg.LLM = config.LLM{MaxHistory: 20, MaxArticleChars: 8000}
	return NewService(cfg, zap.NewNop(), convs, msgs, mems, articles, model)
}

// ---- 测试用例 ----

func TestChatAutoCreateConversation(t *testing.T) {
	convs, msgs, mems, articles := newFakeConvs(), newFakeMsgs(), newFakeMems(), newFakeArticles()
	model := newFakeModel("你好", 5)
	svc := newTestService(convs, msgs, mems, articles, model)

	content := "你好，介绍一下你自己"
	reply, err := svc.Chat(context.Background(), 100, 0, content, nil, func(string) error { return nil })
	if err != nil {
		t.Fatal(err)
	}
	if reply != "你好" {
		t.Fatalf("want reply 你好, got %q", reply)
	}
	// 自动建会话，标题为前 30 字
	list, _ := convs.ListByUser(context.Background(), 100)
	if len(list) != 1 {
		t.Fatalf("want 1 conversation, got %d", len(list))
	}
	if list[0].Title != content {
		t.Fatalf("title mismatch: %q", list[0].Title)
	}
	if list[0].UUID == "" {
		t.Fatal("uuid should be set")
	}
	// user 与 assistant 消息均已落库
	msgs2, _ := msgs.ListByConversation(context.Background(), list[0].ID)
	if len(msgs2) != 2 || msgs2[0].Role != "user" || msgs2[1].Role != "assistant" {
		t.Fatalf("want [user, assistant], got %d msgs", len(msgs2))
	}
}

func TestChatWithArticlesInjectsSystem(t *testing.T) {
	convs, msgs, mems, articles := newFakeConvs(), newFakeMsgs(), newFakeMems(), newFakeArticles()
	model := newFakeModel("基于文章回答", 10)
	svc := newTestService(convs, msgs, mems, articles, model)

	_, err := svc.Chat(context.Background(), 100, 0, "这篇文章讲了什么", []uint{1, 2}, func(string) error { return nil })
	if err != nil {
		t.Fatal(err)
	}
	if len(model.calls) != 1 {
		t.Fatalf("want 1 model call, got %d", len(model.calls))
	}
	sys := model.calls[0].System
	if !strings.Contains(sys, "Go 并发模型") || !strings.Contains(sys, "MySQL 索引") {
		t.Fatalf("system should contain article titles, got: %s", sys)
	}
	if !strings.Contains(sys, "站内文章上下文") {
		t.Fatalf("system should contain article context block, got: %s", sys)
	}
}

func TestChatForbiddenAccess(t *testing.T) {
	convs, msgs, mems, articles := newFakeConvs(), newFakeMsgs(), newFakeMems(), newFakeArticles()
	model := newFakeModel("x", 1)
	svc := newTestService(convs, msgs, mems, articles, model)

	// 用户 100 建会话
	_, err := svc.Chat(context.Background(), 100, 0, "hi", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	list, _ := convs.ListByUser(context.Background(), 100)
	cid := list[0].ID

	// 用户 200 尝试访问该会话 → ErrForbidden，且不调用模型
	model.calls = nil
	_, err = svc.Chat(context.Background(), 200, cid, "偷看", nil, nil)
	if !errors.Is(err, agent.ErrForbidden) {
		t.Fatalf("want ErrForbidden, got %v", err)
	}
	if len(model.calls) != 0 {
		t.Fatal("model should not be called on forbidden access")
	}
	// 消息不应落库到该会话
	msgs2, _ := msgs.ListByConversation(context.Background(), cid)
	if len(msgs2) != 2 {
		t.Fatalf("messages should stay 2, got %d", len(msgs2))
	}
}

func TestChatHistoryContext(t *testing.T) {
	convs, msgs, mems, articles := newFakeConvs(), newFakeMsgs(), newFakeMems(), newFakeArticles()
	model := newFakeModel("ok", 1)
	svc := newTestService(convs, msgs, mems, articles, model)

	// 第一轮
	if _, err := svc.Chat(context.Background(), 100, 0, "第一问", nil, nil); err != nil {
		t.Fatal(err)
	}
	list, _ := convs.ListByUser(context.Background(), 100)
	cid := list[0].ID
	model.calls = nil

	// 第二轮：历史应包含第一轮的 user+assistant
	if _, err := svc.Chat(context.Background(), 100, cid, "第二问", nil, nil); err != nil {
		t.Fatal(err)
	}
	if len(model.calls) != 1 {
		t.Fatalf("want 1 call, got %d", len(model.calls))
	}
	got := model.calls[0].Messages
	if len(got) != 3 {
		t.Fatalf("want 3 messages (hist user, hist assistant, current), got %d", len(got))
	}
	if got[0].Content != "第一问" || got[1].Role != "assistant" || got[2].Content != "第二问" {
		t.Fatalf("history order wrong: %+v", got)
	}
}

func TestMemoryInjected(t *testing.T) {
	convs, msgs, mems, articles := newFakeConvs(), newFakeMsgs(), newFakeMems(), newFakeArticles()
	model := newFakeModel("ok", 1)
	svc := newTestService(convs, msgs, mems, articles, model)

	mems.Append(context.Background(), &agent.Memory{UserID: 100, Content: "用户偏好 Go 语言", Source: "auto"})

	if _, err := svc.Chat(context.Background(), 100, 0, "推荐语言", nil, nil); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(model.calls[0].System, "用户偏好 Go 语言") {
		t.Fatalf("system should contain memory, got: %s", model.calls[0].System)
	}
}

func TestStreamDeltaMatchesReply(t *testing.T) {
	convs, msgs, mems, articles := newFakeConvs(), newFakeMsgs(), newFakeMems(), newFakeArticles()
	model := newFakeModel("流式回复内容", 9)
	svc := newTestService(convs, msgs, mems, articles, model)

	var got strings.Builder
	reply, err := svc.Chat(context.Background(), 100, 0, "hi", nil, func(d string) error { got.WriteString(d); return nil })
	if err != nil {
		t.Fatal(err)
	}
	if reply != "流式回复内容" || got.String() != reply {
		t.Fatalf("delta=%q reply=%q mismatch", got.String(), reply)
	}
	// assistant 消息 tokens 已落库
	list, _ := convs.ListByUser(context.Background(), 100)
	msgs2, _ := msgs.ListByConversation(context.Background(), list[0].ID)
	if msgs2[1].Tokens != 9 {
		t.Fatalf("want tokens 9, got %d", msgs2[1].Tokens)
	}
}

func TestDeleteConversation(t *testing.T) {
	convs, msgs, mems, articles := newFakeConvs(), newFakeMsgs(), newFakeMems(), newFakeArticles()
	model := newFakeModel("ok", 1)
	svc := newTestService(convs, msgs, mems, articles, model)

	if _, err := svc.Chat(context.Background(), 100, 0, "hi", nil, nil); err != nil {
		t.Fatal(err)
	}
	list, _ := convs.ListByUser(context.Background(), 100)
	cid := list[0].ID

	// 其他用户删除 → ErrForbidden
	if err := svc.DeleteConversation(context.Background(), cid, 200); !errors.Is(err, agent.ErrForbidden) {
		t.Fatalf("want ErrForbidden for other user, got %v", err)
	}
	// 所有者删除成功
	if err := svc.DeleteConversation(context.Background(), cid, 100); err != nil {
		t.Fatal(err)
	}
	// 不存在的会话 → ErrConversationNotFound
	if err := svc.DeleteConversation(context.Background(), 99999, 100); !errors.Is(err, agent.ErrConversationNotFound) {
		t.Fatalf("want ErrConversationNotFound, got %v", err)
	}
}
