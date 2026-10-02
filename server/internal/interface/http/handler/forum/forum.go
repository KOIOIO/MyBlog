// Package forum 提供论坛 HTTP 处理器。
package forum

import (
	"context"

	"server/config"
	forumapp "server/internal/application/forum"
	forumdomain "server/internal/domain/forum"
	"server/internal/domain/shared"
	"server/internal/interface/http/middleware"
	"server/model/request"
	"server/model/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler 论坛处理器。
type Handler struct {
	svc *forumapp.Service
	cfg *config.Config
	log *zap.Logger
}

// NewHandler 构造论坛处理器。
func NewHandler(svc *forumapp.Service, cfg *config.Config, log *zap.Logger) *Handler {
	return &Handler{svc: svc, cfg: cfg, log: log}
}

// Publish 发布帖子。
func (h *Handler) Publish(c *gin.Context) {
	var req request.ForumPublish
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.UserID = middleware.GetUserID(c)
	id, err := h.svc.Publish(context.Background(), &forumdomain.ForumPost{
		UserID:   req.UserID,
		Title:    req.Title,
		Content:  req.Content,
		Category: req.Category,
		Tags:     req.Tags,
		Images:   req.Images,
	})
	if err != nil {
		h.log.Error("Failed to publish forum post:", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"id": id}, "Successfully published post", c)
}

// List 帖子列表。
func (h *Handler) List(c *gin.Context) {
	var req request.ForumList
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := h.svc.List(context.Background(), forumdomain.ListCond{
		Category: req.Category,
		Tag:      req.Tag,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		h.log.Error("Failed to get forum list:", zap.Error(err))
		response.FailWithMessage("Failed to get forum list", c)
		return
	}
	response.OkWithData(response.PageResult{List: list, Total: total}, c)
}

// Detail 帖子详情。
func (h *Handler) Detail(c *gin.Context) {
	var req request.ForumDetail
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := h.svc.Detail(context.Background(), req.ID)
	if err != nil {
		h.log.Error("Failed to get forum detail:", zap.Error(err))
		response.FailWithMessage("Failed to get forum detail", c)
		return
	}
	response.OkWithData(data, c)
}

// Like 点赞 / 取消点赞。
func (h *Handler) Like(c *gin.Context) {
	var req request.ForumLike
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.UserID = middleware.GetUserID(c)
	liked, likeCount, err := h.svc.Like(context.Background(), req.PostID, req.UserID)
	if err != nil {
		h.log.Error("Failed to like forum post:", zap.Error(err))
		response.FailWithMessage("Failed to like forum post", c)
		return
	}
	response.OkWithData(gin.H{"liked": liked, "like_count": likeCount}, c)
}

// Comment 发表评论。
func (h *Handler) Comment(c *gin.Context) {
	var req request.ForumComment
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.UserID = middleware.GetUserID(c)
	if err := h.svc.Comment(context.Background(), &forumdomain.ForumComment{
		PostID:   req.PostID,
		ParentID: req.ParentID,
		UserID:   req.UserID,
		Content:  req.Content,
	}); err != nil {
		h.log.Error("Failed to comment forum post:", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("Successfully commented", c)
}

// Tags 固定标签库。
func (h *Handler) Tags(c *gin.Context) {
	tags, err := h.svc.Tags(context.Background())
	if err != nil {
		h.log.Error("Failed to get forum tags:", zap.Error(err))
		response.FailWithMessage("Failed to get forum tags", c)
		return
	}
	response.OkWithData(tags, c)
}

// Upload 论坛图片上传。
func (h *Handler) Upload(c *gin.Context) {
	_, header, err := c.Request.FormFile("image")
	if err != nil {
		h.log.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	url, err := h.svc.Upload(header)
	if err != nil {
		h.log.Error("Failed to upload forum image:", zap.Error(err))
		response.FailWithMessage("Failed to upload forum image", c)
		return
	}
	response.OkWithDetailed(response.ImageUpload{
		Url:     url,
		OssType: h.cfg.System.OssType,
	}, "Successfully uploaded image", c)
}

// ManageList 帖子管理列表。
func (h *Handler) ManageList(c *gin.Context) {
	var req request.ForumManageList
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.UserID = middleware.GetUserID(c)
	req.RoleID = middleware.GetRoleID(c)
	list, total, err := h.svc.ManageList(context.Background(), forumdomain.ManageListCond{
		Title:    req.Title,
		UserID:   req.UserID,
		RoleID:   shared.RoleID(req.RoleID),
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		h.log.Error("Failed to get forum manage list:", zap.Error(err))
		response.FailWithMessage("Failed to get forum manage list", c)
		return
	}
	response.OkWithData(response.PageResult{List: list, Total: total}, c)
}

// Delete 删除帖子。
func (h *Handler) Delete(c *gin.Context) {
	var req request.ForumDelete
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.UserID = middleware.GetUserID(c)
	req.RoleID = middleware.GetRoleID(c)
	if err := h.svc.Delete(context.Background(), req.IDs, req.UserID, shared.RoleID(req.RoleID)); err != nil {
		h.log.Error("Failed to delete forum posts:", zap.Error(err))
		response.FailWithMessage("Failed to delete forum posts", c)
		return
	}
	response.OkWithMessage("Successfully deleted forum posts", c)
}

// ManageComments 评论管理列表。
func (h *Handler) ManageComments(c *gin.Context) {
	var req request.ForumManageComments
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.UserID = middleware.GetUserID(c)
	req.RoleID = middleware.GetRoleID(c)
	list, total, err := h.svc.ManageComments(context.Background(), forumdomain.ManageCommentCond{
		PostID:   req.PostID,
		UserID:   req.UserID,
		RoleID:   shared.RoleID(req.RoleID),
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		h.log.Error("Failed to get forum comments:", zap.Error(err))
		response.FailWithMessage("Failed to get forum comments", c)
		return
	}
	response.OkWithData(response.PageResult{List: list, Total: total}, c)
}

// CommentDelete 删除评论。
func (h *Handler) CommentDelete(c *gin.Context) {
	var req request.ForumCommentDelete
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.UserID = middleware.GetUserID(c)
	req.RoleID = middleware.GetRoleID(c)
	if err := h.svc.DeleteComments(context.Background(), req.IDs, req.UserID, shared.RoleID(req.RoleID)); err != nil {
		h.log.Error("Failed to delete forum comments:", zap.Error(err))
		response.FailWithMessage("Failed to delete forum comments", c)
		return
	}
	response.OkWithMessage("Successfully deleted forum comments", c)
}
