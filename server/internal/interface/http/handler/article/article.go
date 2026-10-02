// Package article 提供 article 路由 Handler（迁移自 api/article.go，依赖构造注入）。
package article

import (
	articleapp "server/internal/application/article"
	articledomain "server/internal/domain/article"
	"server/internal/interface/http/middleware"
	"server/internal/model/request"
	"server/internal/model/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler article 接口处理器。
type Handler struct {
	articles *articleapp.Service
	log      *zap.Logger
}

// NewHandler 构造 article 处理器。
func NewHandler(articles *articleapp.Service, log *zap.Logger) *Handler {
	return &Handler{articles: articles, log: log}
}

// InfoByID 根据文章 ID 获取文章详细信息。
func (h *Handler) InfoByID(c *gin.Context) {
	var req request.ArticleInfoByID
	if err := c.ShouldBindUri(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	articleDoc, err := h.articles.InfoByID(c.Request.Context(), req.ID)
	if err != nil {
		h.log.Error("Failed to get article information:", zap.Error(err))
		response.FailWithMessage("Failed to get article information", c)
		return
	}
	response.OkWithData(articleDoc, c)
}

// Search 文章搜索接口。
func (h *Handler) Search(c *gin.Context) {
	var info request.ArticleSearch
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := h.articles.Search(c.Request.Context(), articledomain.SearchSpec{
		Query:    info.Query,
		Category: info.Category,
		Tag:      info.Tag,
		Sort:     info.Sort,
		Order:    info.Order,
		Page:     info.Page,
		PageSize: info.PageSize,
	})
	if err != nil {
		h.log.Error("Failed to get article search results:", zap.Error(err))
		response.FailWithMessage("Failed to get article search results", c)
		return
	}
	response.OkWithData(response.PageResult{List: list, Total: total}, c)
}

// Category 获取所有文章类别及数量。
func (h *Handler) Category(c *gin.Context) {
	category, err := h.articles.Categories(c.Request.Context())
	if err != nil {
		h.log.Error("Failed to get article category:", zap.Error(err))
		response.FailWithMessage("Failed to get article category", c)
		return
	}
	response.OkWithData(category, c)
}

// Tags 获取所有文章标签及数量。
func (h *Handler) Tags(c *gin.Context) {
	tags, err := h.articles.Tags(c.Request.Context())
	if err != nil {
		h.log.Error("Failed to get article tags:", zap.Error(err))
		response.FailWithMessage("Failed to get article tags", c)
		return
	}
	response.OkWithData(tags, c)
}

// Like 文章收藏/取消收藏。
func (h *Handler) Like(c *gin.Context) {
	var req request.ArticleLike
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.UserID = middleware.GetUserID(c)
	if err := h.articles.Like(c.Request.Context(), req.UserID, req.ArticleID); err != nil {
		h.log.Error("Failed to complete the operation:", zap.Error(err))
		response.FailWithMessage("Failed to complete the operation", c)
		return
	}
	response.OkWithMessage("Successfully completed the operation", c)
}

// IsLike 返回文章收藏状态。
func (h *Handler) IsLike(c *gin.Context) {
	var req request.ArticleLike
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.UserID = middleware.GetUserID(c)
	isLike, err := h.articles.IsLike(c.Request.Context(), req.UserID, req.ArticleID)
	if err != nil {
		h.log.Error("Failed to get like status:", zap.Error(err))
		response.FailWithMessage("Failed to get like status", c)
		return
	}
	response.OkWithData(isLike, c)
}

// LikesList 获取文章收藏列表。
func (h *Handler) LikesList(c *gin.Context) {
	var pageInfo request.ArticleLikesList
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	pageInfo.UserID = middleware.GetUserID(c)
	list, total, err := h.articles.LikesList(c.Request.Context(), pageInfo.UserID, pageInfo.Page, pageInfo.PageSize)
	if err != nil {
		h.log.Error("Failed to get likes list:", zap.Error(err))
		response.FailWithMessage("Failed to get likes list", c)
		return
	}
	response.OkWithData(response.PageResult{List: list, Total: total}, c)
}

// Create 发布文章。
func (h *Handler) Create(c *gin.Context) {
	var req request.ArticleCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err := h.articles.Create(c.Request.Context(), &articledomain.Article{
		Cover:    req.Cover,
		Title:    req.Title,
		Category: req.Category,
		Tags:     req.Tags,
		Abstract: req.Abstract,
		Content:  req.Content,
	})
	if err != nil {
		h.log.Error("Failed to create article:", zap.Error(err))
		response.FailWithMessage("Failed to create article", c)
		return
	}
	response.OkWithMessage("Successfully created article", c)
}

// Delete 删除文章。
func (h *Handler) Delete(c *gin.Context) {
	var req request.ArticleDelete
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.articles.Delete(c.Request.Context(), req.IDs); err != nil {
		h.log.Error("Failed to delete article:", zap.Error(err))
		response.FailWithMessage("Failed to delete article", c)
		return
	}
	response.OkWithMessage("Successfully deleted article", c)
}

// Update 更新文章。
func (h *Handler) Update(c *gin.Context) {
	var req request.ArticleUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err := h.articles.Update(c.Request.Context(), req.ID, &articledomain.Article{
		Cover:    req.Cover,
		Title:    req.Title,
		Category: req.Category,
		Tags:     req.Tags,
		Abstract: req.Abstract,
		Content:  req.Content,
	})
	if err != nil {
		h.log.Error("Failed to update article:", zap.Error(err))
		response.FailWithMessage("Failed to update article", c)
		return
	}
	response.OkWithMessage("Successfully updated article", c)
}

// SetTop 设置文章置顶。
func (h *Handler) SetTop(c *gin.Context) {
	var req request.ArticleSetTop
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.articles.SetTop(c.Request.Context(), req.ID, req.IsTop); err != nil {
		h.log.Error("Failed to set article top:", zap.Error(err))
		response.FailWithMessage("Failed to set article top", c)
		return
	}
	response.OkWithMessage("Successfully set article top", c)
}

// List 获取文章列表（后台）。
func (h *Handler) List(c *gin.Context) {
	var pageInfo request.ArticleList
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := h.articles.List(c.Request.Context(), articledomain.SearchSpec{
		Title:    pageInfo.Title,
		Abstract: pageInfo.Abstract,
		Category: valueOrEmpty(pageInfo.Category),
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	})
	if err != nil {
		h.log.Error("Failed to get article list:", zap.Error(err))
		response.FailWithMessage("Failed to get article list", c)
		return
	}
	response.OkWithData(response.PageResult{List: list, Total: total}, c)
}

func valueOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
