# MyBlog AI Agent 功能实施计划（Plan）

> **For agentic workers:**
>
>  REQUIRED SUB-SKILL: Use superpowers:executing-plans 按任务逐项执行。步骤用 
>
> `- [ ]`
>
>  勾选跟踪。

**Goal:** 为 MyBlog 新增与论坛同级的 AI Agent 栏目，支持带记忆的多轮对话、可携带站内文章提问、按用户隔离的历史会话列表（左侧点击恢复），MVP 不依赖向量数据库。

**Architecture:** 后端新增独立 BC `agent`（DDD 四层：domain/application/interface/infrastructure），LLM 走阿里云百炼 OpenAI 兼容接口 + SSE 流式（不引第三方 SDK）；记忆分三层（会话级 / 历史恢复 / 画像记忆）全部落 MySQL。前端新增 `/agent` 路由、导航菜单、Pinia store 与页面组件。向量数据库仅作为 Phase 2 可选增强，架构上预留检索端口。

**Tech Stack:** Go 1.22 + Gin + GORM（MySQL）+ SSE；Vue3 + Vite + Pinia + Element Plus + TypeScript + i18n；DashScope OpenAI 兼容 API（模型 `qwen-max`）。

**Spec:** `docs/ai-agent-spec.md`（本计划依据，执行时两文档同读）

## Global Constraints



* 后端新增代码必须遵循现有 DDD 分层（domain 不 import 基础设施库），BC 目录结构与 forum 一致。

* 新表必须注册进 `server/internal/flag/flag_sql.go` 的 AutoMigrate 列表，并同步 SQL 迁移文件。

* 所有 `/agent/*` 接口必须挂在 `privateGroup`（JWT），用户归属一律用 `middleware.GetUserID(c)`，禁止信任前端传的 user\_id。

* API Key 只存在于 `server/config.yaml` / 环境变量，禁止硬编码、禁止写入数据库与前端代码。

* 文章正文注入前截断：每篇 ≤ `MaxArticleChars`（默认 8000 字符）；历史上下文保留最近 `MaxHistory`（默认 20）条消息。

* 前端流式接口用原生 `fetch` + `ReadableStream` 实现，不经过 axios 拦截器。

* 新增文案必须接入现有 i18n（zh-CN /en/zh-TW），菜单项 key 为 `nav.agent`。



***

## Phase 1（MVP，无向量库）

### Task 1: 后端数据模型与建表注册

**Files:**



* Create: `server/internal/model/database/agent_conversation.go`

* Create: `server/internal/model/database/agent_message.go`

* Create: `server/internal/model/database/agent_memory.go`

* Modify: `server/internal/flag/flag_sql.go`（AutoMigrate 列表追加 3 个模型）

**Interfaces:**



* Produces: `database.AgentConversation{ MODEL; UserID uint; UUID string; Title string; Summary string; Status int }`

* Produces: `database.AgentMessage{ MODEL; ConversationID uint; Role string; Content string; ArticleIDs datatypes.JSON; Tokens int }`

* Produces: `database.AgentMemory{ MODEL; UserID uint; Content string; Source string }`

- [ ] **Step 1: 编写 3 个 GORM 模型**



```
// server/internal/model/database/agent_conversation.go
package database

// AgentConversation 会话表
type AgentConversation struct {
	MODEL
	UserID  uint   `json:"user_id" gorm:"index"`
	UUID    string `json:"uuid" gorm:"type:char(36);unique"`
	Title   string `json:"title" gorm:"size:100"`
	Summary string `json:"summary" gorm:"type:text"`
	Status  int    `json:"status" gorm:"default:0"`
}
```

`agent_message.go`：`ConversationID uint `gorm:"index"``、`Role string `gorm:"size:10"``、`Content string `gorm:"type:text"``、`ArticleIDs datatypes.JSON`（用 `gorm.io/datatypes`，已在项目依赖中核对，若缺失则 `go get gorm.io/datatypes`）、`Tokens int`。``

`` `agent_memory.go`：`UserID uint `gorm:"index" ``、`Content string `gorm:"size:500"``、`Source string `gorm:"size:20"``。



* [ ] **Step 2: 注册 AutoMigrate**

在 `server/internal/flag/flag_sql.go` 的 AutoMigrate 调用参数列表末尾追加：



```
&database.AgentConversation{},
&database.AgentMessage{},
&database.AgentMemory{},
```



* [ ] **Step 3: 同步 SQL 迁移文件**

在 `server/mysql_20250218.sql`（或项目现行迁移文件）追加 3 张表的 `CREATE TABLE`（字段对应模型；`article_ids` 为 `JSON`）。



* [ ] **Step 4: 验证**

Run: `cd server && go build ./...`

Expected: 编译通过。再执行 `go run . --sql`（项目既有建表命令）确认 3 张表可创建；不重复建已存在表。



* [ ] **Step 5: Commit**



```
git add server/internal/model/database/agent_*.go server/internal/flag/flag_sql.go server/mysql_20250218.sql
git commit -m "feat(agent): add agent conversation/message/memory models and migrate"
```

### Task 2: 后端 domain/agent 实体与端口

**Files:**



* Create: `server/internal/domain/agent/agent.go`

* Create: `server/internal/domain/agent/agent_test.go`

**Interfaces:**



* Consumes: `database.AgentConversation / AgentMessage / AgentMemory`（Task 1 形状，仅用于领域实体字段对齐）

* Produces（本任务核心，后续任务全部依赖）:


  * `agent.Conversation`, `agent.Message`, `agent.Memory` 实体

  * `agent.ArticleChunk{ ID uint; Title string; Content string }`

  * `agent.ChatMessage{ Role, Content string }`

  * `agent.ChatRequest{ System string; Messages []agent.ChatMessage; Stream bool }`

  * `agent.ChatResult{ FullText string; Tokens int }`

  * `agent.ConversationRepository` / `agent.MessageRepository` / `agent.MemoryRepository` / `agent.ArticleContentReader` / `agent.ModelProvider`（签名见下）

  * `agent.ErrConversationNotFound`, `agent.ErrForbidden` 错误值

- [ ] **Step 1: 写失败测试（title 截取与归属规则）**



```
// server/internal/domain/agent/agent_test.go
package agent

import "testing"

func TestMakeTitle(t *testing.T) {
	got := MakeTitle("这是一段超过三十个字的消息内容用来验证标题截取逻辑是否正确执行")
	if len([]rune(got)) != 30 {
		t.Fatalf("want 30 runes, got %d: %q", len([]rune(got)), got)
	}
	if got := MakeTitle("短消息"); got != "短消息" {
		t.Fatalf("want original, got %q", got)
	}
}

func TestValidateAccess(t *testing.T) {
	if err := ValidateAccess(1, 1); err != nil { t.Fatalf("same user should pass: %v", err) }
	if err := ValidateAccess(1, 2); err != ErrForbidden { t.Fatalf("want ErrForbidden, got %v", err) }
}
```



* [ ] **Step 2: 运行确认失败**

Run: `cd server && go test ./internal/domain/agent/ -run 'TestMakeTitle|TestValidateAccess'`

Expected: FAIL（包 / 函数不存在）。



* [ ] **Step 3: 实现实体与端口**



```
// server/internal/domain/agent/agent.go
package agent

import (
	"context"
	"errors"
	"time"
)

var (
	ErrConversationNotFound = errors.New("conversation not found")
	ErrForbidden            = errors.New("forbidden")
)

// MakeTitle 取消息前 30 个 rune 作为会话标题。
func MakeTitle(content string) string {
	r := []rune(content)
	if len(r) > 30 {
		return string(r[:30])
	}
	return content
}

// ValidateAccess 校验资源归属，防止越权。
func ValidateAccess(resourceUserID, callerUserID uint) error {
	if resourceUserID != callerUserID {
		return ErrForbidden
	}
	return nil
}

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

type Message struct {
	ID             uint      `json:"id"`
	ConversationID uint      `json:"conversation_id"`
	Role           string    `json:"role"`
	Content        string    `json:"content"`
	ArticleIDs     []uint    `json:"article_ids,omitempty"`
	Tokens         int       `json:"-"`
	CreatedAt      time.Time `json:"created_at"`
}

type Memory struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"-"`
	Content   string    `json:"content"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"created_at"`
}

// ---- 端口 ----

type ConversationRepository interface {
	Create(ctx context.Context, c *Conversation) error
	ListByUser(ctx context.Context, userID uint) ([]*Conversation, error)
	GetByID(ctx context.Context, id, userID uint) (*Conversation, error)
	Delete(ctx context.Context, id, userID uint) error
}

type MessageRepository interface {
	Append(ctx context.Context, m *Message) error
	ListByConversation(ctx context.Context, conversationID uint) ([]*Message, error)
	DeleteByConversation(ctx context.Context, conversationID uint) error
}

type MemoryRepository interface {
	ListByUser(ctx context.Context, userID uint) ([]*Memory, error)
	Append(ctx context.Context, m *Memory) error
}

type ArticleChunk struct {
	ID      uint
	Title   string
	Content string
}

type ArticleContentReader interface {
	ReadByIDs(ctx context.Context, ids []uint) ([]*ArticleChunk, error)
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	System   string        `json:"system"`
	Messages []ChatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

type ChatResult struct {
	FullText string
	Tokens   int
}

type ModelProvider interface {
	ChatStream(ctx context.Context, req ChatRequest, onDelta func(string) error) (*ChatResult, error)
}
```



* [ ] **Step 4: 运行通过**

Run: `cd server && go test ./internal/domain/agent/`

Expected: PASS。



* [ ] **Step 5: Commit**



```
git add server/internal/domain/agent/
git commit -m "feat(agent): domain entities, ports and access rules"
```

### Task 3: 后端 infrastructure LLM Provider（DashScope + SSE）

**Files:**



* Create: `server/config/conf_llm.go`

* Modify: `server/config.yaml`（追加 `llm:` 段）

* Create: `server/internal/infrastructure/llm/dashscope.go`

* Create: `server/internal/infrastructure/llm/dashscope_test.go`

**Interfaces:**



* Consumes: `agent.ChatRequest / ChatResult / ChatMessage`、`agent.ModelProvider`（Task 2）

* Produces: `llm.NewDashScopeProvider(cfg *config.LLM, log *zap.Logger) *DashScopeProvider`（实现 `agent.ModelProvider`）

* Produces: `config.Config` 新增字段 `LLM config.LLM`

- [ ] **Step 1: 配置结构**



```
// server/config/conf_llm.go
package config

// LLM 大模型配置（DashScope OpenAI 兼容）。
type LLM struct {
	BaseURL        string `mapstructure:"base-url" json:"base-url"`
	APIKey         string `mapstructure:"api-key" json:"api-key"`
	Model          string `mapstructure:"model" json:"model"`
	MaxHistory     int    `mapstructure:"max-history" json:"max-history"`
	MaxArticleChars int   `mapstructure:"max-article-chars" json:"max-article-chars"`
}
```

`config.yaml` 追加：



```
llm:
  base-url: "https://dashscope.aliyuncs.com/compatible-mode/v1"
  api-key: "${DASHSCOPE_API_KEY}"   # 或用已存在 secrets 方式注入，禁止硬编码真实 Key
  model: "qwen-max"
  max-history: 20
  max-article-chars: 8000
```

在 `server/config` 配置加载处将 `llm` 段绑定进 `Config`（沿用现有 Viper 绑定模式）。



* [ ] **Step 2: 写失败测试（httptest 模拟 SSE 响应）**



```
// server/internal/infrastructure/llm/dashscope_test.go
package llm

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"server/internal/domain/agent"
	"go.uber.org/zap"
)

func TestDashScopeChatStream(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		// 校验 Authorization 与 model
		if r.Header.Get("Authorization") == "" { t.Error("missing auth") }
		w.Write([]byte("data: {\"choices\":[{\"delta\":{\"content\":\"你\"}}]}\n\n"))
		w.Write([]byte("data: {\"choices\":[{\"delta\":{\"content\":\"好\"}}]}\n\n"))
		w.Write([]byte("data: [DONE]\n\n"))
	}))
	defer srv.Close()

	p := NewDashScopeProvider(&Config{BaseURL: srv.URL, APIKey: "test-key", Model: "qwen-max"}, zap.NewNop())
	var got strings.Builder
	res, err := p.ChatStream(context.Background(), agent.ChatRequest{
		System: "你是助手",
		Messages: []agent.ChatMessage{{Role: "user", Content: "hi"}},
		Stream: true,
	}, func(d string) error { got.WriteString(d); return nil })
	if err != nil { t.Fatal(err) }
	if got.String() != "你好" { t.Fatalf("want 你好, got %q", got.String()) }
	if res.FullText != "你好" { t.Fatalf("full text mismatch: %q", res.FullText) }
}
```



* [ ] **Step 3: 运行确认失败**

Run: `cd server && go test ./internal/infrastructure/llm/`

Expected: FAIL（包不存在）。



* [ ] **Step 4: 实现 DashScope Provider**

关键实现（`dashscope.go`）：



* `ChatStream` 构造 `POST {BaseURL}/chat/completions`，Header `Authorization: Bearer {APIKey}`、`Content-Type: application/json`。

* 请求体：`{"model":..., "messages":[{"role":"system","content":System}]+Messages, "stream":true}`。

* 用 `http.Client` 发起，`bufio.Scanner` 按行读 `data: ` 前缀 JSON：取 `choices[0].delta.content` 逐块回调 `onDelta` 并累计 `FullText`；遇 `[DONE]` 结束；错误码（如 `model_not_found`/ 鉴权失败）包装为带前缀错误返回。

* 每次 `onDelta` 后 `flush`；context 取消时返回 `ctx.Err()`（中断处理由上层兜底）。

- [ ] **Step 5: 运行通过**

Run: `cd server && go test ./internal/infrastructure/llm/`

Expected: PASS。



* [ ] **Step 6: Commit**



```
git add server/config/conf_llm.go server/config.yaml server/internal/infrastructure/llm/
git commit -m "feat(agent): DashScope streaming LLM provider"
```

### Task 4: 后端 infrastructure MySQL 仓储 + 文章内容读取

**Files:**



* Create: `server/internal/infrastructure/mysql/agent.go`

* Create: `server/internal/infrastructure/mysql/agent_test.go`（用 sqlmock 或项目既有测试模式；若项目无 sqlmock 依赖，则用真实 MySQL 集成测试标记 `//go:build integration` 并提供 docker 运行说明）

**Interfaces:**



* Consumes: `agent.ConversationRepository / MessageRepository / MemoryRepository / ArticleContentReader`、`database.*`（Task 1/2）

* Produces:


  * `mysql.NewAgentConversationRepo(db *gorm.DB) agent.ConversationRepository`

  * `mysql.NewAgentMessageRepo(db *gorm.DB) agent.MessageRepository`

  * `mysql.NewAgentMemoryRepo(db *gorm.DB) agent.MemoryRepository`

  * `mysql.NewAgentArticleReader(db *gorm.DB) agent.ArticleContentReader`（读 `database.Article` 的 ID/Title/Content，截断 `MaxArticleChars` 由 application 层处理，本层只读全量正文）

- [ ] **Step 1: 写失败测试（仓储行为）**

以 `ConversationRepository.Create + ListByUser + GetByID(归属校验) + Delete` 为核心断言：Create 后 ID 非零；ListByUser 只返回该用户；`GetByID` 的 userID 不匹配时返回 `agent.ErrForbidden`；Delete 后 List 为空。



* [ ] **Step 2: 运行确认失败 → Step 3 实现 → Step 4 通过**

实现要点：全部查询带 `WHERE user_id = ?`；软删除用 GORM `Delete`；`ArticleIDs` 存取用 `datatypes.JSON` 与 `[]uint` 互转；`MessageRepository.ListByConversation` 按 `created_at ASC` 排序；`DeleteByConversation` 批量软删。



* [ ] **Step 5: Commit**



```
git add server/internal/infrastructure/mysql/agent*.go
git commit -m "feat(agent): MySQL repositories and article reader"
```

### Task 5: 后端 application/agent 服务（用例编排）

**Files:**



* Create: `server/internal/application/agent/agent.go`

* Create: `server/internal/application/agent/agent_test.go`

**Interfaces:**



* Consumes: Task 2 全部端口 + Task 3 `llm` provider 实例 + `config.Config`

* Produces:


  * `agent.NewService(cfg *config.Config, log *zap.Logger, convs agent.ConversationRepository, msgs agent.MessageRepository, mems agent.MemoryRepository, articles agent.ArticleContentReader, model agent.ModelProvider) *Service`

  * `(*Service).CreateConversation(ctx, userID uint) (*agent.Conversation, error)`

  * `(*Service).ListConversations(ctx, userID uint) ([]*agent.Conversation, error)`

  * `(*Service).ListMessages(ctx, conversationID, userID uint) ([]*agent.Message, error)`

  * `(*Service).DeleteConversation(ctx, conversationID, userID uint) error`

  * `(*Service).Chat(ctx, userID, conversationID uint, content string, articleIDs []uint, onDelta func(string) error) (string, error)`

- [ ] **Step 1: 写失败测试（fake 仓储 + fake provider，核心为 Chat 用例）**



```
// server/internal/application/agent/agent_test.go
// 使用内存 fake 实现 Task 2 的 4 个仓储/Reader 端口，断言：
// 1) conversationID==0 时自动建会话，title = MakeTitle(content)，返回会话 UUID；
// 2) user 消息先落库（含 article_ids），assistant 消息流式结束后落库；
// 3) 文章被读入 System 提示（fake ArticleContentReader 返回固定 chunk，断言 ChatRequest.System 包含文章标题）；
// 4) conversationID!=0 且归属不符 → ErrForbidden，不调用 provider；
// 5) onDelta 收到的拼接 == provider 返回 FullText == 返回 reply。
```



* [ ] **Step 2: 运行确认失败 → Step 3: 实现 Service → Step 4: 通过**

`Chat` 实现流程：



1. 若 `conversationID == 0`：`CreateConversation`（UUID 用 `gofrs/uuid.Must(uuid.NewV4()).String()`，title=`MakeTitle(content)`）；否则 `GetByID` 校验归属（`ErrConversationNotFound` 透传）。

2. `msgs.Append(user消息)`（含 `articleIDs`）。

3. 组装上下文：

* `memories` → 拼接画像条目（`ListByUser`，基础版最多取 10 条）进 system。

* `articles`（若 `len(articleIDs)>0`）→ `ReadByIDs`，每篇截断 `cfg.LLM.MaxArticleChars` rune，拼「\[站内文章上下文]」块进 system。

* 历史：`msgs.ListByConversation` 取最近 `cfg.LLM.MaxHistory` 条转 `agent.ChatMessage`。

1. `model.ChatStream(system+历史+当前user消息, onDelta)`；流式失败时返回错误（用户消息已落库，可重试）。

2. 落库 assistant 消息（含 `Tokens`），返回 `FullText`。

* [ ] **Step 5: Commit**



```
git add server/internal/application/agent/
git commit -m "feat(agent): application service for conversations and chat"
```

### Task 6: 后端 interface handler + 路由注册 + bootstrap 装配

**Files:**



* Create: `server/internal/interface/http/handler/agent/agent.go`

* Modify: `server/internal/interface/http/router.go`（Deps + 路由段）

* Modify: `server/internal/bootstrap/app.go`（装配 agent service/handler 并注入 Deps）

**Interfaces:**



* Consumes: `agentapp.Service`（Task 5）、`middleware.GetUserID(c)`、`middleware.JWTAuth`

* Produces: `agent.NewHandler(svc *agentapp.Service) *Handler`，方法：


  * `Create / List / Messages / Delete / Chat`

* Routes（privateGroup）：


  * `POST /agent/conversations` → Create

  * `GET /agent/conversations` → List

  * `GET /agent/conversations/:id/messages` → Messages

  * `DELETE /agent/conversations/:id` → Delete

  * `POST /agent/chat` → Chat（SSE）

- [ ] **Step 1: 实现 Handler**

* 所有 handler 第一行取 `userID := middleware.GetUserID(c)`。

* `Chat`：解析 body `{conversation_id, message, article_ids}`；设置 `Content-Type: text/event-stream`、`Cache-Control: no-cache`；`c.Stream` 逐块写 `data: {"delta": "..."}\n\n`；结束写 `data: [DONE]\n\n`；错误写 `data: {"error":"..."}\n\n` 并中止。

* `Messages`/`Delete` 解析 `:id` 为 uint，透传 `ErrConversationNotFound`（404）与 `ErrForbidden`（403）。

* 统一走现有响应包装（`server/internal/model/response`，沿用项目 Result 风格；SSE 接口除外）。

- [ ] **Step 2: 注册路由（router.go）**



```
// Deps 追加
AgentHandler *agenthandler.Handler

// NewRouter privateGroup 段追加
{
	agentRouter := privateGroup.Group("agent")
	agentRouter.POST("conversations", deps.AgentHandler.Create)
	agentRouter.GET("conversations", deps.AgentHandler.List)
	agentRouter.GET("conversations/:id/messages", deps.AgentHandler.Messages)
	agentRouter.DELETE("conversations/:id", deps.AgentHandler.Delete)
	agentRouter.POST("chat", deps.AgentHandler.Chat)
}
```



* [ ] **Step 3: bootstrap 装配**

在 `server/internal/bootstrap/app.go` 构造：`llm.NewDashScopeProvider(&cfg.LLM, log)` → 4 个仓储（`mysql.NewAgentXxx(db)`）→ `agentapp.NewService(...)` → `agent.NewHandler(svc)`，写入 router `Deps`。



* [ ] **Step 4: 验证（编译 + 冒烟）**

Run: `cd server && go build ./... && go vet ./...`

然后启动（项目 `start.sh`），用 curl 冒烟：



```
# 登录拿 token（沿用现有登录接口），再：
curl -X POST -H "x-access-token: $TOKEN" http://localhost:8080/agent/conversations
curl -N -X POST -H "x-access-token: $TOKEN" -H "Content-Type: application/json" \
  -d '{"message":"你好"}' http://localhost:8080/agent/chat   # 期望 SSE 流
```

Expected: 会话创建成功；chat 返回 `data: {"delta":...}` 流并 `[DONE]`。



* [ ] **Step 5: Commit**



```
git add server/internal/interface/http/handler/agent/ server/internal/interface/http/router.go server/internal/bootstrap/app.go
git commit -m "feat(agent): HTTP handlers, routes and wiring"
```

### Task 7: 前端 API 层与 Pinia store

**Files:**



* Create: `web/src/api/agent.ts`

* Create: `web/src/stores/agent.ts`

**Interfaces:**



* Consumes: `web/src/utils/request.ts`（现有 axios 实例）、`web/src/stores/user.ts`

* Produces:


  * `web/src/api/agent.ts`:


    * `createConversation(): Promise<AgentConversation>`

    * `listConversations(): Promise<AgentConversation[]>`

    * `listMessages(conversationId: number): Promise<AgentMessage[]>`

    * `deleteConversation(conversationId: number): Promise<void>`

    * `chatStream(payload: {conversationId?: number; message: string; articleIds: number[]}, onDelta: (d: string) => void, signal?: AbortSignal): Promise<void>`（原生 fetch + ReadableStream 解析 SSE）

  * 类型：`AgentConversation{ id, uuid, title, updated_at }`、`AgentMessage{ id, conversation_id, role, content, article_ids, created_at }`

  * `web/src/stores/agent.ts`（Pinia）: `conversations / currentConversationId / messages / streaming` 状态 + `loadConversations / createConversation / openConversation / removeConversation / sendMessage` actions

- [ ] **Step 1: 实现 api/agent.ts（流式解析是关键）**



```
export async function chatStream(
  payload: { conversationId?: number; message: string; articleIds: number[] },
  onDelta: (d: string) => void,
  signal?: AbortSignal,
): Promise<void> {
  const token = useUserStore().accessToken // 与现有 request.ts 的取 token 方式保持一致
  const res = await fetch('/agent/chat', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', 'x-access-token': token },
    body: JSON.stringify({ conversation_id: payload.conversationId, message: payload.message, article_ids: payload.articleIds }),
    signal,
  })
  if (!res.ok || !res.body) throw new Error(`chat failed: ${res.status}`)
  const reader = res.body.getReader()
  const decoder = new TextDecoder()
  let buf = ''
  while (true) {
    const { done, value } = await reader.read()
    if (done) break
    buf += decoder.decode(value, { stream: true })
    const lines = buf.split('\n')
    buf = lines.pop() ?? ''
    for (const line of lines) {
      const s = line.trim()
      if (!s.startsWith('data:')) continue
      const payload = s.slice(5).trim()
      if (payload === '[DONE]') return
      try {
        const json = JSON.parse(payload)
        if (json.delta) onDelta(json.delta)
        if (json.error) throw new Error(json.error)
      } catch (e) { if (e instanceof Error && e.message !== '') throw e }
    }
  }
}
```



* [ ] **Step 2: store 实现**

`sendMessage`：`streaming=true` → 调 `chatStream` 追加流式文本到当前 assistant 消息 → 完成 / 中断后 `streaming=false`；`conversationId` 为 0 时先 `createConversation` 并刷新列表；结束后追加完整消息到 `messages`（服务端已落库，前端不必二次拉取）。



* [ ] **Step 3: 验证**

Run: `cd web && npx vue-tsc --noEmit && npm run build`

Expected: 类型检查与构建通过。



* [ ] **Step 4: Commit**



```
git add web/src/api/agent.ts web/src/stores/agent.ts
git commit -m "feat(agent): frontend API client and store"
```

### Task 8: 前端 Agent 页面组件

**Files:**



* Create: `web/src/views/web/agent/index.vue`

* Create: `web/src/components/agent/ArticlePicker.vue`

* Create: `web/src/components/agent/MessageList.vue`

* Modify: `web/src/views/web/index.vue`（挂载子路由容器，若现有结构为 `<router-view/>` 自动生效则无需改动）

**Interfaces:**



* Consumes: Task 7 store/api、现有 `web/src/api/article.ts`（`GET /article/search`）、`AuthPopover.vue`、`useUserStore`

* Produces: `/agent` 页面（左：会话列表 + 新建 / 删除；右：消息区 + 输入区 + 文章选择器）

- [ ] **Step 1: 页面布局（index.vue）**



```
┌──────────────────────────────────────────────┐
│ 左侧栏(260px, 可收起)                          │
│  [+ 新建会话]                                  │
│  ┌──────────────────┐  ┌─────────────────┐   │
│  │ 会话1 (标题/时间)  │  │ 消息区(MessageList)│   │
│  │ 会话2            │  │                 │   │
│  │ ...             │  │                 │   │
│  └──────────────────┘  ├─────────────────┤   │
│  (移动端时隐藏为抽屉)    │ 输入区 + [选择文章] │   │
└──────────────────────────────────────────────┘
```



* 未登录：显示 AuthPopover 引导，不渲染对话区。

* 交互：点击左侧会话 → `openConversation(id)` 拉消息；「+」新建；hover 删除按钮（el-popconfirm 确认）。

* 输入区：`el-input type="textarea"` + 发送按钮；发送中禁用并显示停止按钮（AbortController 中断流式）。

- [ ] **Step 2: ArticlePicker.vue（复用文章搜索接口）**

* 打开 `el-dialog`，内嵌 `el-select multiple filterable remote`，`remote-method` 调 `GET /article/search?keyword=`（复用现有 `web/src/api/article.ts` 中搜索方法，若无则新增 `searchArticles(keyword)`），选项显示标题，值取文章 id。

* 已选文章在输入区上方以 `el-tag` 展示（可移除），随消息一起提交 `article_ids`。

- [ ] **Step 3: MessageList.vue**

* 渲染 `messages`：user 消息右侧气泡（若含 `article_ids`，气泡上方展示「📄 引用 N 篇文章」小标签）；assistant 左侧气泡。

* 流式时当前 assistant 消息原地追加（store 的 `streaming` 状态），支持 Markdown 轻量渲染（复用项目已有的 md 渲染依赖；若无则先纯文本，Phase 2 换 md-editor-v3 的预览模式）。

* 自动滚动到底部；空态引导文案。

- [ ] **Step 4: 验证**

Run: `cd web && npm run build`

手动：登录后进入 `/agent`，新建会话 → 自由对话（流式）→ 新建会话选文章提问 → 刷新页面恢复历史。



* [ ] **Step 5: Commit**



```
git add web/src/views/web/agent/ web/src/components/agent/
git commit -m "feat(agent): agent page, article picker and message list"
```

### Task 9: 前端路由 + 导航菜单 + i18n

**Files:**



* Modify: `web/src/router/index.ts`（web children 追加 agent）

* Modify: `web/src/components/layout/WebNavbar.vue`（menuList 追加）

* Modify: `web/src/i18n/locales/zh-CN/base.ts`、`web/src/i18n/locales/en/base.ts`、`web/src/i18n/locales/zh-TW/base.ts`（nav.agent）

- [ ] **Step 1: 路由**



```
{
  path: "agent",
  name: "agent",
  component: () => import('@/views/web/agent/index.vue'),
  meta: { title: "AI Agent", requiresAuth: true }
},
```

（children 中，与 forum 同级；`requiresAuth` 复用现有路由守卫逻辑。）



* [ ] **Step 2: 导航菜单**



```
{titleKey: "nav.agent", name: "/agent"},
```



* [ ] **Step 3: i18n**

`zh-CN`：`agent: "AI Agent"`；`en`：`agent: "AI Agent"`；`zh-TW`：`agent: "AI Agent"`（按三语言惯例；页面内文案 key 一并补齐：新建会话 / 选择文章 / 发送 / 停止 / 暂无会话等）。



* [ ] **Step 4: 验证**

Run: `cd web && npm run build`

手动：导航栏出现 AI Agent；点击进入页面；未登录显示登录引导；登录后正常使用。



* [ ] **Step 5: Commit**



```
git add web/src/router/index.ts web/src/components/layout/WebNavbar.vue web/src/i18n/locales/
git commit -m "feat(agent): nav menu, route and i18n"
```

### Task 10: 端到端联调与验证

**Files:** 无新增（验证任务）



* [ ] **Step 1: 后端全量测试**

Run: `cd server && go test ./...`

Expected: 全绿（含既有测试）。



* [ ] **Step 2: 前端构建 + 类型检查**

Run: `cd web && npx vue-tsc --noEmit && npm run build`



* [ ] **Step 3: 端到端手工验证清单（对照 spec 验收标准逐条过）**

1. 导航栏出现 AI Agent，与论坛同级；未登录访问 `/agent` 显示登录引导。

2. 登录后新建会话 → 自由对话，流式逐字返回（首 token < 2s）。

3. 选择 2 篇文章提问 → 回答内容基于所选文章（用「这篇文章讲了什么」验证）。

4. 多轮追问 → 模型能引用前文（会话级记忆）。

5. 刷新页面 → 左侧列表保留会话；点击恢复全部历史并继续对话（L2 记忆）。

6. 删除会话 → 列表消失，重新拉取不再出现。

7. 用户 A 的会话在用户 B 的列表中不可见（换账号验证越权隔离）。

8. 对话中携带 `article_ids` 的 user 消息，回复后刷新仍显示「引用文章」标识。

9. 流式中断（停止按钮）→ 不崩溃，已落库的 user 消息保留，assistant 部分可重发。

10. 三语言切换下菜单与页面文案正常。

* [ ] **Step 4: 更新项目文档**

在 `Readme.md` 功能列表追加 Agent 栏目说明；`record.md` 记录本次改造（沿用项目记录惯例）。



* [ ] **Step 5: Commit**



```
git add Readme.md record.md
git commit -m "docs(agent): update readme and record"
```



***

## Phase 2（可选增强：向量数据库）

### Task 11: 向量检索增强（docker 部署）

**Files（示意，进入 P2 时按 Task 1-10 同粒度细化）:**



* Modify: `server/internal/domain/agent/agent.go`（新增 `ArticleRetriever` / `MemoryRetriever` 端口）

* Create: `server/internal/infrastructure/vector/qdrant.go`（或 pgvector 仓储）

* Create: `server/internal/application/agent/retrieve.go`（无 `article_ids` 时自动召回 top-k 文章；记忆按语义召回）

* Modify: `server/config/conf_llm.go`（embedding 模型配置）

* Create: `deploy/docker-compose.vector.yml`

- [ ] **Step 1: 部署向量库**



```
docker run -d --name myblog-qdrant -p 6333:6333 -v qdrant_data:/qdrant/storage qdrant/qdrant
```

（或 `docker run -d --name myblog-pgvector -e POSTGRES_PASSWORD=... -p 5432:5432 pgvector/pgvector:pg16`；个人博客量级两者皆可，Qdrant 免建表更省事。）



* [ ] **Step 2: 文章向量化与索引**

- 用 DashScope `text-embedding-v3`（或 `bge-m3`）对全站文章（标题 + 正文）批量 embedding，入库 Qdrant collection `articles`（payload: article\_id, title）。

- 定时 / 增量同步：文章发布时同步向量（可挂现有 cron 或发布 hook）。

* [ ] **Step 3: 语义召回接入 Chat**

- `Chat` 中 `len(articleIDs)==0` 时：`ArticleRetriever.Search(ctx, userContent, topK=3)` → 命中文章进 system 上下文（标注「自动检索」）。

- 记忆增强：`MemoryRetriever` 按 query 语义召回历史画像条目替换「取最近 10 条」策略。

* [ ] **Step 4: 验证**

- 提问「讲一下站内关于 MySQL 索引的文章」且不手动选文章 → 回复引用相关文章内容。

- 两个不同主题会话 → 新会话中旧主题相关提问能唤起对应记忆条目。

- `docker compose -f deploy/docker-compose.vector.yml up -d` 一键拉起（供部署脚本复用）。

* [ ] **Step 5: Commit**



```
git add server/internal/domain/agent/agent.go server/internal/infrastructure/vector/ server/internal/application/agent/retrieve.go server/config/conf_llm.go deploy/
git commit -m "feat(agent): vector retrieval enhancement (qdrant)"
```



***

## 任务依赖关系



```
Task1(表) → Task2(domain) → Task3(LLM) ─┐
                              Task4(仓储) ─┼→ Task5(application) → Task6(接口+装配) → Task10(联调)
Task3、Task4 可并行                               ↓
Task7(前端 API/store) → Task8(页面) → Task9(导航/i18n) ──┘
Task11(Phase 2) 独立可选，依赖 Task2 端口预留与 Task5 Chat 逻辑
```

## Self-Review 结果（写计划时自查）



* **Spec 覆盖**：F1→Task9；F2/F3/F6/F7→Task5+8；F4/F8/F9/F10→Task4+5+8；F5→Task5（基础画像注入）；F11/F12→Task11；NFR（SSE / 越权 / 截断 /i18n）→Global Constraints + Task6/8/9。无遗漏。

* **占位符检查**：无 TBD；关键代码均已给出签名或实现片段。

* **类型一致性**：`agent.ChatRequest/ChatResult`、`MakeTitle`、`ValidateAccess`、`middleware.GetUserID`、`chatStream(payload, onDelta, signal)` 在 Task 2/3/5/6/7 间签名一致；`article_ids` 字段在 Go 侧 `[]uint` 与前端 `number[]` 对齐。