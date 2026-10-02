// Package friendlink 提供友链 HTTP 处理器。
package friendlink

import (
	"context"

	flapp "server/internal/application/friendlink"
	fldomain "server/internal/domain/friendlink"
	"server/model/request"
	"server/model/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler 友链处理器。
type Handler struct {
	svc *flapp.Service
	log *zap.Logger
}

// NewHandler 构造友链处理器。
func NewHandler(svc *flapp.Service, log *zap.Logger) *Handler {
	return &Handler{svc: svc, log: log}
}

// Info 获取友链信息。
func (h *Handler) Info(c *gin.Context) {
	list, total, err := h.svc.Info(context.Background())
	if err != nil {
		h.log.Error("Failed to get friend link information:", zap.Error(err))
		response.FailWithMessage("Failed to get friend link information", c)
		return
	}
	response.OkWithData(gin.H{"list": list, "total": total}, c)
}

// Create 创建友链。
func (h *Handler) Create(c *gin.Context) {
	var req request.FriendLinkCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.svc.Create(context.Background(), &fldomain.FriendLink{
		Logo:        req.Logo,
		Link:        req.Link,
		Name:        req.Name,
		Description: req.Description,
	}); err != nil {
		h.log.Error("Failed to create friend link:", zap.Error(err))
		response.FailWithMessage("Failed to create friend link", c)
		return
	}
	response.OkWithMessage("Successfully created friend link", c)
}

// Delete 删除友链。
func (h *Handler) Delete(c *gin.Context) {
	var req request.FriendLinkDelete
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.svc.Delete(context.Background(), req.IDs); err != nil {
		h.log.Error("Failed to delete friend link:", zap.Error(err))
		response.FailWithMessage("Failed to delete friend link", c)
		return
	}
	response.OkWithMessage("Successfully deleted friend link", c)
}

// Update 更新友链。
func (h *Handler) Update(c *gin.Context) {
	var req request.FriendLinkUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.svc.Update(context.Background(), &fldomain.FriendLink{
		ID:          req.ID,
		Link:        req.Link,
		Name:        req.Name,
		Description: req.Description,
	}); err != nil {
		h.log.Error("Failed to update friend link:", zap.Error(err))
		response.FailWithMessage("Failed to update friend link", c)
		return
	}
	response.OkWithMessage("Successfully updated friend link", c)
}

// List 获取友链列表。
func (h *Handler) List(c *gin.Context) {
	var pageInfo request.FriendLinkList
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := h.svc.List(context.Background(), fldomain.ListCond{
		Name:        pageInfo.Name,
		Description: pageInfo.Description,
		Page:        pageInfo.Page,
		PageSize:    pageInfo.PageSize,
	})
	if err != nil {
		h.log.Error("Failed to get friend link list:", zap.Error(err))
		response.FailWithMessage("Failed to get friend link list", c)
		return
	}
	response.OkWithData(response.PageResult{List: list, Total: total}, c)
}
