package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"server/global"
	"server/model/request"
	"server/model/response"
	"server/utils"
)

type ForumApi struct {
}

// ForumPublish 发布帖子
func (forumApi *ForumApi) ForumPublish(c *gin.Context) {
	var req request.ForumPublish
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.UserID = utils.GetUserID(c)
	id, err := forumService.ForumPublish(req)
	if err != nil {
		global.Log.Error("Failed to publish forum post:", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"id": id}, "Successfully published post", c)
}

// ForumList 帖子列表
func (forumApi *ForumApi) ForumList(c *gin.Context) {
	var req request.ForumList
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := forumService.ForumList(req)
	if err != nil {
		global.Log.Error("Failed to get forum list:", zap.Error(err))
		response.FailWithMessage("Failed to get forum list", c)
		return
	}
	response.OkWithData(response.PageResult{List: list, Total: total}, c)
}

// ForumDetail 帖子详情
func (forumApi *ForumApi) ForumDetail(c *gin.Context) {
	var req request.ForumDetail
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := forumService.ForumDetail(req.ID)
	if err != nil {
		global.Log.Error("Failed to get forum detail:", zap.Error(err))
		response.FailWithMessage("Failed to get forum detail", c)
		return
	}
	response.OkWithData(data, c)
}

// ForumLike 点赞 / 取消点赞
func (forumApi *ForumApi) ForumLike(c *gin.Context) {
	var req request.ForumLike
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.UserID = utils.GetUserID(c)
	liked, likeCount, err := forumService.ForumLike(req)
	if err != nil {
		global.Log.Error("Failed to like forum post:", zap.Error(err))
		response.FailWithMessage("Failed to like forum post", c)
		return
	}
	response.OkWithData(gin.H{"liked": liked, "like_count": likeCount}, c)
}

// ForumComment 发表评论
func (forumApi *ForumApi) ForumComment(c *gin.Context) {
	var req request.ForumComment
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.UserID = utils.GetUserID(c)
	if err := forumService.ForumComment(req); err != nil {
		global.Log.Error("Failed to comment forum post:", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("Successfully commented", c)
}

// ForumTags 固定标签库
func (forumApi *ForumApi) ForumTags(c *gin.Context) {
	tags, err := forumService.ForumTags()
	if err != nil {
		global.Log.Error("Failed to get forum tags:", zap.Error(err))
		response.FailWithMessage("Failed to get forum tags", c)
		return
	}
	response.OkWithData(tags, c)
}

// ForumUpload 论坛图片上传
func (forumApi *ForumApi) ForumUpload(c *gin.Context) {
	_, header, err := c.Request.FormFile("image")
	if err != nil {
		global.Log.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	url, err := forumService.ForumUpload(header)
	if err != nil {
		global.Log.Error("Failed to upload forum image:", zap.Error(err))
		response.FailWithMessage("Failed to upload forum image", c)
		return
	}
	response.OkWithDetailed(response.ImageUpload{
		Url:     url,
		OssType: global.Config.System.OssType,
	}, "Successfully uploaded image", c)
}

// ForumManageList 管理员帖子列表
func (forumApi *ForumApi) ForumManageList(c *gin.Context) {
	var req request.ForumManageList
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.UserID = utils.GetUserID(c)
	req.RoleID = utils.GetRoleID(c)
	list, total, err := forumService.ForumManageList(req)
	if err != nil {
		global.Log.Error("Failed to get forum manage list:", zap.Error(err))
		response.FailWithMessage("Failed to get forum manage list", c)
		return
	}
	response.OkWithData(response.PageResult{List: list, Total: total}, c)
}

// ForumDelete 删除帖子
func (forumApi *ForumApi) ForumDelete(c *gin.Context) {
	var req request.ForumDelete
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.UserID = utils.GetUserID(c)
	req.RoleID = utils.GetRoleID(c)
	if err := forumService.ForumDelete(req.IDs, req.UserID, req.RoleID); err != nil {
		global.Log.Error("Failed to delete forum posts:", zap.Error(err))
		response.FailWithMessage("Failed to delete forum posts", c)
		return
	}
	response.OkWithMessage("Successfully deleted forum posts", c)
}

// ForumManageComments 管理员评论列表
func (forumApi *ForumApi) ForumManageComments(c *gin.Context) {
	var req request.ForumManageComments
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.UserID = utils.GetUserID(c)
	req.RoleID = utils.GetRoleID(c)
	list, total, err := forumService.ForumManageComments(req)
	if err != nil {
		global.Log.Error("Failed to get forum comments:", zap.Error(err))
		response.FailWithMessage("Failed to get forum comments", c)
		return
	}
	response.OkWithData(response.PageResult{List: list, Total: total}, c)
}

// ForumCommentDelete 删除评论
func (forumApi *ForumApi) ForumCommentDelete(c *gin.Context) {
	var req request.ForumCommentDelete
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.UserID = utils.GetUserID(c)
	req.RoleID = utils.GetRoleID(c)
	if err := forumService.ForumCommentDelete(req.IDs, req.UserID, req.RoleID); err != nil {
		global.Log.Error("Failed to delete forum comments:", zap.Error(err))
		response.FailWithMessage("Failed to delete forum comments", c)
		return
	}
	response.OkWithMessage("Successfully deleted forum comments", c)
}
