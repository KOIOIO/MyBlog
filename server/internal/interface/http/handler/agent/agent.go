// Package agent 提供 AI Agent HTTP 处理器。
package agent

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	agentapp "server/internal/application/agent"
	"server/internal/domain/agent"
	"server/internal/interface/http/middleware"
	"server/internal/model/request"
	"server/internal/model/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler AI Agent 处理器。
type Handler struct {
	svc *agentapp.Service
	log *zap.Logger
}

// NewHandler 构造 AI Agent 处理器。
func NewHandler(svc *agentapp.Service, log *zap.Logger) *Handler {
	return &Handler{svc: svc, log: log}
}

// Create 新建会话。
func (h *Handler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	conv, err := h.svc.CreateConversation(context.Background(), userID)
	if err != nil {
		h.log.Error("create conversation failed", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"id": conv.ID, "uuid": conv.UUID, "title": conv.Title}, "created", c)
}

// List 会话列表。
func (h *Handler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	convs, err := h.svc.ListConversations(context.Background(), userID)
	if err != nil {
		h.log.Error("list conversations failed", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(convs, c)
}

// Messages 会话消息列表。
func (h *Handler) Messages(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMessage("invalid conversation id", c)
		return
	}
	msgs, err := h.svc.ListMessages(context.Background(), uint(id), userID)
	if err != nil {
		h.writeDomainError(err, c)
		return
	}
	response.OkWithData(msgs, c)
}

// Delete 删除会话。
func (h *Handler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMessage("invalid conversation id", c)
		return
	}
	if err := h.svc.DeleteConversation(context.Background(), uint(id), userID); err != nil {
		h.writeDomainError(err, c)
		return
	}
	response.Ok(c)
}

// Chat 发送消息（SSE 流式返回）。
func (h *Handler) Chat(c *gin.Context) {
	var req request.AgentChat
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if req.Message == "" {
		response.FailWithMessage("message is required", c)
		return
	}
	userID := middleware.GetUserID(c)

	// SSE 响应头先设置，后续所有输出（含预校验错误）统一走 SSE 协议
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	writeEvent := func(payload gin.H) {
		data, _ := json.Marshal(payload)
		c.Writer.Write([]byte("data: " + string(data) + "\n\n"))
		c.Writer.Flush()
	}

	// 归属预校验（会话存在性），失败以 SSE 错误事件返回
	if req.ConversationID != 0 {
		if _, err := h.svc.PeekConversation(context.Background(), req.ConversationID, userID); err != nil {
			writeEvent(gin.H{"error": err.Error()})
			return
		}
	}

	_, err := h.svc.Chat(c.Request.Context(), userID, req.ConversationID, req.Message, req.ArticleIDs,
		func(delta string) error {
			writeEvent(gin.H{"delta": delta})
			return nil
		})
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			h.log.Error("chat stream failed", zap.Error(err))
		}
		writeEvent(gin.H{"error": err.Error()})
		return
	}
	writeEvent(gin.H{"done": true})
}

// writeDomainError 领域错误 → HTTP 状态码。
func (h *Handler) writeDomainError(err error, c *gin.Context) {
	switch {
	case errors.Is(err, agent.ErrConversationNotFound):
		response.FailWithDetailed(gin.H{"code": 404}, err.Error(), c)
	case errors.Is(err, agent.ErrForbidden):
		response.Forbidden(err.Error(), c)
	default:
		h.log.Error("agent handler error", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	}
}
