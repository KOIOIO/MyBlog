// Package comment 提供评论 HTTP 处理器。
package comment

import (
	"context"

	commentapp "server/internal/application/comment"
	commentdomain "server/internal/domain/comment"
	"server/internal/domain/shared"
	"server/internal/interface/http/middleware"
	"server/model/request"
	"server/model/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler 评论处理器。
type Handler struct {
	svc *commentapp.Service
	log *zap.Logger
}

// NewHandler 构造评论处理器。
func NewHandler(svc *commentapp.Service, log *zap.Logger) *Handler {
	return &Handler{svc: svc, log: log}
}

// InfoByArticleID 根据文章 ID 获取评论信息。
func (h *Handler) InfoByArticleID(c *gin.Context) {
	var req request.CommentInfoByArticleID
	if err := c.ShouldBindUri(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, err := h.svc.InfoByArticleID(context.Background(), req.ArticleID)
	if err != nil {
		h.log.Error("Failed to get comment information:", zap.Error(err))
		response.FailWithMessage("Failed to get comment information", c)
		return
	}
	response.OkWithData(list, c)
}

// New 获取最新评论。
func (h *Handler) New(c *gin.Context) {
	list, err := h.svc.Newest(context.Background())
	if err != nil {
		h.log.Error("Failed to get new comment:", zap.Error(err))
		response.FailWithMessage("Failed to get new comment", c)
		return
	}
	response.OkWithData(list, c)
}

// Create 创建评论。
func (h *Handler) Create(c *gin.Context) {
	var req request.CommentCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.UserUUID = middleware.GetUUID(c)
	if err := h.svc.Create(context.Background(), &commentdomain.Comment{
		ArticleID: req.ArticleID,
		PID:       req.PID,
		UserUUID:  req.UserUUID,
		Content:   req.Content,
	}); err != nil {
		h.log.Error("Failed to create comment:", zap.Error(err))
		response.FailWithMessage("Failed to create comment", c)
		return
	}
	response.OkWithMessage("Successfully created comment", c)
}

// Delete 删除评论。
func (h *Handler) Delete(c *gin.Context) {
	var req request.CommentDelete
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.svc.Delete(context.Background(), req.IDs, middleware.GetUUID(c), shared.RoleID(middleware.GetRoleID(c))); err != nil {
		h.log.Error("Failed to delete comment:", zap.Error(err))
		response.FailWithMessage("Failed to delete comment", c)
		return
	}
	response.OkWithMessage("Successfully deleted comment", c)
}

// Info 获取用户评论。
func (h *Handler) Info(c *gin.Context) {
	list, err := h.svc.Info(context.Background(), middleware.GetUUID(c))
	if err != nil {
		h.log.Error("Failed to get comment information:", zap.Error(err))
		response.FailWithMessage("Failed to get comment information", c)
		return
	}
	response.OkWithData(list, c)
}

// List 评论管理列表。
func (h *Handler) List(c *gin.Context) {
	var pageInfo request.CommentList
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := h.svc.List(context.Background(), commentdomain.ListCond{
		ArticleID: pageInfo.ArticleID,
		UserUUID:  pageInfo.UserUUID,
		Content:   pageInfo.Content,
		Page:      pageInfo.Page,
		PageSize:  pageInfo.PageSize,
	})
	if err != nil {
		h.log.Error("Failed to get comment list:", zap.Error(err))
		response.FailWithMessage("Failed to get comment list", c)
		return
	}
	response.OkWithData(response.PageResult{List: list, Total: total}, c)
}
