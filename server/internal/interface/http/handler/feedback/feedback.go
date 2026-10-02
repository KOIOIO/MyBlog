// Package feedback 提供反馈 HTTP 处理器。
package feedback

import (
	"context"

	fbapp "server/internal/application/feedback"
	fbdomain "server/internal/domain/feedback"
	"server/internal/interface/http/middleware"
	"server/internal/model/request"
	"server/internal/model/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler 反馈处理器。
type Handler struct {
	svc *fbapp.Service
	log *zap.Logger
}

// NewHandler 构造反馈处理器。
func NewHandler(svc *fbapp.Service, log *zap.Logger) *Handler {
	return &Handler{svc: svc, log: log}
}

// New 获取最新反馈。
func (h *Handler) New(c *gin.Context) {
	list, err := h.svc.Newest(context.Background())
	if err != nil {
		h.log.Error("Failed to get new feedback:", zap.Error(err))
		response.FailWithMessage("Failed to get new feedback", c)
		return
	}
	response.OkWithData(list, c)
}

// Create 创建反馈。
func (h *Handler) Create(c *gin.Context) {
	var req request.FeedbackCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.svc.Create(context.Background(), &fbdomain.Feedback{
		UserUUID: middleware.GetUUID(c),
		Content:  req.Content,
	}); err != nil {
		h.log.Error("Failed to create feedback:", zap.Error(err))
		response.FailWithMessage("Failed to create feedback", c)
		return
	}
	response.OkWithMessage("Successfully created feedback", c)
}

// Info 获取用户反馈信息。
func (h *Handler) Info(c *gin.Context) {
	list, err := h.svc.Info(context.Background(), middleware.GetUUID(c))
	if err != nil {
		h.log.Error("Failed to get feedback information:", zap.Error(err))
		response.FailWithMessage("Failed to get feedback information", c)
		return
	}
	response.OkWithData(list, c)
}

// Delete 删除反馈。
func (h *Handler) Delete(c *gin.Context) {
	var req request.FeedbackDelete
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.svc.Delete(context.Background(), req.IDs); err != nil {
		h.log.Error("Failed to delete feedback:", zap.Error(err))
		response.FailWithMessage("Failed to delete feedback", c)
		return
	}
	response.OkWithMessage("Successfully deleted feedback", c)
}

// Reply 回复反馈。
func (h *Handler) Reply(c *gin.Context) {
	var req request.FeedbackReply
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.svc.Reply(context.Background(), req.ID, req.Reply); err != nil {
		h.log.Error("Failed to update feedback:", zap.Error(err))
		response.FailWithMessage("Failed to update feedback", c)
		return
	}
	response.OkWithMessage("Successfully updated feedback", c)
}

// List 获取反馈列表。
func (h *Handler) List(c *gin.Context) {
	var pageInfo request.PageInfo
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := h.svc.List(context.Background(), pageInfo.Page, pageInfo.PageSize)
	if err != nil {
		h.log.Error("Failed to get feedback list:", zap.Error(err))
		response.FailWithMessage("Failed to get feedback list", c)
		return
	}
	response.OkWithData(response.PageResult{List: list, Total: total}, c)
}
