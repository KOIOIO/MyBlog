// Package config 提供系统配置 HTTP 处理器。
package config

import (
	"context"

	"server/config"
	cfgapp "server/internal/application/config"
	"server/model/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler 配置处理器。
type Handler struct {
	svc *cfgapp.Service
	cfg *config.Config
	log *zap.Logger
}

// NewHandler 构造配置处理器。
func NewHandler(svc *cfgapp.Service, cfg *config.Config, log *zap.Logger) *Handler {
	return &Handler{svc: svc, cfg: cfg, log: log}
}

// GetWebsite 获取网站配置。
func (h *Handler) GetWebsite(c *gin.Context) {
	response.OkWithData(h.cfg.Website, c)
}

// UpdateWebsite 更新网站配置。
func (h *Handler) UpdateWebsite(c *gin.Context) {
	var req config.Website
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.svc.UpdateWebsite(context.Background(), req); err != nil {
		h.log.Error("Failed to update website:", zap.Error(err))
		response.FailWithMessage("Failed to update website", c)
		return
	}
	response.OkWithMessage("Successfully updated website", c)
}

// GetSystem 获取系统配置。
func (h *Handler) GetSystem(c *gin.Context) {
	response.OkWithData(h.cfg.System, c)
}

// UpdateSystem 更新系统配置。
func (h *Handler) UpdateSystem(c *gin.Context) {
	var req config.System
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.svc.UpdateSystem(context.Background(), req); err != nil {
		h.log.Error("Failed to update system:", zap.Error(err))
		response.FailWithMessage("Failed to update system", c)
		return
	}
	response.OkWithMessage("Successfully updated system", c)
}

// GetEmail 获取邮箱配置。
func (h *Handler) GetEmail(c *gin.Context) {
	response.OkWithData(h.cfg.Email, c)
}

// UpdateEmail 更新邮箱配置。
func (h *Handler) UpdateEmail(c *gin.Context) {
	var req config.Email
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.svc.UpdateEmail(context.Background(), req); err != nil {
		h.log.Error("Failed to update email:", zap.Error(err))
		response.FailWithMessage("Failed to update email", c)
		return
	}
	response.OkWithMessage("Successfully updated email", c)
}

// GetQiniu 获取七牛云配置。
func (h *Handler) GetQiniu(c *gin.Context) {
	response.OkWithData(h.cfg.Qiniu, c)
}

// UpdateQiniu 更新七牛云配置。
func (h *Handler) UpdateQiniu(c *gin.Context) {
	var req config.Qiniu
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.svc.UpdateQiniu(context.Background(), req); err != nil {
		h.log.Error("Failed to update qiniu:", zap.Error(err))
		response.FailWithMessage("Failed to update qiniu", c)
		return
	}
	response.OkWithMessage("Successfully updated qiniu", c)
}

// GetJwt 获取 JWT 配置。
func (h *Handler) GetJwt(c *gin.Context) {
	response.OkWithData(h.cfg.Jwt, c)
}

// UpdateJwt 更新 JWT 配置。
func (h *Handler) UpdateJwt(c *gin.Context) {
	var req config.Jwt
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.svc.UpdateJwt(context.Background(), req); err != nil {
		h.log.Error("Failed to update jwt:", zap.Error(err))
		response.FailWithMessage("Failed to update jwt", c)
		return
	}
	response.OkWithMessage("Successfully updated jwt", c)
}

// GetGaode 获取高德配置。
func (h *Handler) GetGaode(c *gin.Context) {
	response.OkWithData(h.cfg.Gaode, c)
}

// UpdateGaode 更新高德配置。
func (h *Handler) UpdateGaode(c *gin.Context) {
	var req config.Gaode
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.svc.UpdateGaode(context.Background(), req); err != nil {
		h.log.Error("Failed to update gaode:", zap.Error(err))
		response.FailWithMessage("Failed to update gaode", c)
		return
	}
	response.OkWithMessage("Successfully updated gaode", c)
}
