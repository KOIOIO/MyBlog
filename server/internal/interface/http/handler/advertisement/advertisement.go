// Package advertisement 提供广告 HTTP 处理器。
package advertisement

import (
	"context"

	advapp "server/internal/application/advertisement"
	advdomain "server/internal/domain/advertisement"
	"server/model/request"
	"server/model/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler 广告处理器。
type Handler struct {
	svc *advapp.Service
	log *zap.Logger
}

// NewHandler 构造广告处理器。
func NewHandler(svc *advapp.Service, log *zap.Logger) *Handler {
	return &Handler{svc: svc, log: log}
}

// Info 获取广告信息。
func (h *Handler) Info(c *gin.Context) {
	list, total, err := h.svc.Info(context.Background())
	if err != nil {
		h.log.Error("Failed to get advertisement information:", zap.Error(err))
		response.FailWithMessage("Failed to get advertisement information", c)
		return
	}
	response.OkWithData(gin.H{"list": list, "total": total}, c)
}

// Create 创建广告。
func (h *Handler) Create(c *gin.Context) {
	var req request.AdvertisementCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.svc.Create(context.Background(), &advdomain.Advertisement{
		AdImage: req.AdImage,
		Link:    req.Link,
		Title:   req.Title,
		Content: req.Content,
	}); err != nil {
		h.log.Error("Failed to create advertisement:", zap.Error(err))
		response.FailWithMessage("Failed to create advertisement", c)
		return
	}
	response.OkWithMessage("Successfully created advertisement", c)
}

// Delete 删除广告。
func (h *Handler) Delete(c *gin.Context) {
	var req request.AdvertisementDelete
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.svc.Delete(context.Background(), req.IDs); err != nil {
		h.log.Error("Failed to delete advertisement:", zap.Error(err))
		response.FailWithMessage("Failed to delete advertisement", c)
		return
	}
	response.OkWithMessage("Successfully deleted advertisement", c)
}

// Update 更新广告。
func (h *Handler) Update(c *gin.Context) {
	var req request.AdvertisementUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.svc.Update(context.Background(), &advdomain.Advertisement{
		ID:      req.ID,
		Link:    req.Link,
		Title:   req.Title,
		Content: req.Content,
	}); err != nil {
		h.log.Error("Failed to update advertisement:", zap.Error(err))
		response.FailWithMessage("Failed to update advertisement", c)
		return
	}
	response.OkWithMessage("Successfully updated advertisement", c)
}

// List 获取广告列表。
func (h *Handler) List(c *gin.Context) {
	var pageInfo request.AdvertisementList
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := h.svc.List(context.Background(), advdomain.ListCond{
		Title:    pageInfo.Title,
		Content:  pageInfo.Content,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	})
	if err != nil {
		h.log.Error("Failed to get advertisement list:", zap.Error(err))
		response.FailWithMessage("Failed to get advertisement list", c)
		return
	}
	response.OkWithData(response.PageResult{List: list, Total: total}, c)
}
