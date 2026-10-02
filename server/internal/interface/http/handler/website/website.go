// Package website 提供网站 HTTP 处理器。
package website

import (
	"context"
	"net/http"
	"time"

	"server/config"
	websiteapp "server/internal/application/website"
	websitedomain "server/internal/domain/website"
	"server/model/request"
	"server/model/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler 网站处理器。
type Handler struct {
	svc *websiteapp.Service
	cfg *config.Config
	log *zap.Logger
}

// NewHandler 构造网站处理器。
func NewHandler(svc *websiteapp.Service, cfg *config.Config, log *zap.Logger) *Handler {
	return &Handler{svc: svc, cfg: cfg, log: log}
}

// Logo 网站 Logo 链接。
func (h *Handler) Logo(c *gin.Context) {
	if h.cfg.Website.Logo != "" {
		c.Redirect(http.StatusMovedPermanently, h.cfg.Website.Logo)
	} else {
		c.Redirect(http.StatusMovedPermanently, "/image/logo.png")
	}
}

// Title 网站标题栏。
func (h *Handler) Title(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"title": h.cfg.Website.Title})
}

// Info 获取网站信息。
func (h *Handler) Info(c *gin.Context) {
	response.OkWithData(h.cfg.Website, c)
}

// Carousel 获取首页背景。
func (h *Handler) Carousel(c *gin.Context) {
	urls, err := h.svc.Carousel(context.Background())
	if err != nil {
		h.log.Error("Failed to get carousel:", zap.Error(err))
		response.FailWithMessage("Failed to get carousel", c)
		return
	}
	response.OkWithData(urls, c)
}

// News 获取新闻。
func (h *Handler) News(c *gin.Context) {
	sourceStr := c.Query("source")
	hotSearchData, err := h.svc.News(context.Background(), sourceStr)
	if err != nil {
		h.log.Error("Failed to get news:", zap.Error(err))
		response.FailWithMessage("Failed to get news", c)
		return
	}
	response.OkWithData(hotSearchData, c)
}

// Calendar 获取日历。
func (h *Handler) Calendar(c *gin.Context) {
	dateStr := time.Now().Format("2006/0102")
	calendar, err := h.svc.Calendar(context.Background(), dateStr)
	if err != nil {
		h.log.Error("Failed to get calendar:", zap.Error(err))
		response.FailWithMessage("Failed to get calendar", c)
		return
	}
	response.OkWithData(calendar, c)
}

// FooterLink 获取页脚链接。
func (h *Handler) FooterLink(c *gin.Context) {
	footerLinks, err := h.svc.FooterLink(context.Background())
	if err != nil {
		h.log.Error("Failed to get footer link:", zap.Error(err))
		response.FailWithMessage("Failed to get footer link", c)
		return
	}
	response.OkWithData(footerLinks, c)
}

// AddCarousel 添加首页背景。
func (h *Handler) AddCarousel(c *gin.Context) {
	var req request.WebsiteCarouselOperation
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.svc.AddCarousel(context.Background(), req.Url); err != nil {
		h.log.Error("Failed to add carousel:", zap.Error(err))
		response.FailWithMessage("Failed to add carousel", c)
		return
	}
	response.OkWithMessage("Successfully added carousel", c)
}

// CancelCarousel 移除首页背景。
func (h *Handler) CancelCarousel(c *gin.Context) {
	var req request.WebsiteCarouselOperation
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.svc.CancelCarousel(context.Background(), req.Url); err != nil {
		h.log.Error("Failed to cancel carousel:", zap.Error(err))
		response.FailWithMessage("Failed to cancel carousel", c)
		return
	}
	response.OkWithMessage("Successfully canceled carousel", c)
}

// CreateFooterLink 创建页脚链接。
func (h *Handler) CreateFooterLink(c *gin.Context) {
	var req websitedomain.FooterLink
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.svc.CreateFooterLink(context.Background(), &req); err != nil {
		h.log.Error("Failed to create footer link:", zap.Error(err))
		response.FailWithMessage("Failed to create footer link", c)
		return
	}
	response.OkWithMessage("Successfully created footer link", c)
}

// DeleteFooterLink 删除页脚链接。
func (h *Handler) DeleteFooterLink(c *gin.Context) {
	var req websitedomain.FooterLink
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.svc.DeleteFooterLink(context.Background(), &req); err != nil {
		h.log.Error("Failed to delete footer link:", zap.Error(err))
		response.FailWithMessage("Failed to delete footer link", c)
		return
	}
	response.OkWithMessage("Successfully deleted footer link", c)
}
