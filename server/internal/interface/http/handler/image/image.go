// Package image 提供图片 HTTP 处理器。
package image

import (
	"context"

	"server/config"
	imageapp "server/internal/application/image"
	imagedomain "server/internal/domain/image"
	"server/model/request"
	"server/model/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler 图片处理器。
type Handler struct {
	svc *imageapp.Service
	cfg *config.Config
	log *zap.Logger
}

// NewHandler 构造图片处理器。
func NewHandler(svc *imageapp.Service, cfg *config.Config, log *zap.Logger) *Handler {
	return &Handler{svc: svc, cfg: cfg, log: log}
}

// Upload 上传图片。
func (h *Handler) Upload(c *gin.Context) {
	_, header, err := c.Request.FormFile("image")
	if err != nil {
		h.log.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	url, err := h.svc.Upload(context.Background(), header)
	if err != nil {
		h.log.Error("Failed to upload image:", zap.Error(err))
		response.FailWithMessage("Failed to upload image", c)
		return
	}
	response.OkWithDetailed(response.ImageUpload{
		Url:     url,
		OssType: h.cfg.System.OssType,
	}, "Successfully uploaded image", c)
}

// Delete 删除图片。
func (h *Handler) Delete(c *gin.Context) {
	var req request.ImageDelete
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.svc.Delete(context.Background(), req.IDs); err != nil {
		h.log.Error("Failed to delete image:", zap.Error(err))
		response.FailWithMessage("Failed to delete image", c)
		return
	}
	response.OkWithMessage("Successfully deleted image", c)
}

// List 图片列表。
func (h *Handler) List(c *gin.Context) {
	var pageInfo request.ImageList
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	imageList, total, err := h.svc.List(context.Background(), imagedomain.ListCond{
		Name:     pageInfo.Name,
		Category: pageInfo.Category,
		Storage:  pageInfo.Storage,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	})
	if err != nil {
		h.log.Error("Failed to get image list:", zap.Error(err))
		response.FailWithMessage("Failed to get image list", c)
		return
	}
	response.OkWithData(response.PageResult{List: imageList, Total: total}, c)
}
