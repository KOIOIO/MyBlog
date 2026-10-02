package router

import (
	"github.com/gin-gonic/gin"
	"server/api"
)

type ForumRouter struct {
}

func (f *ForumRouter) InitForumRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup, AdminRouter *gin.RouterGroup) {
	forumRouter := Router.Group("forum")            // 登录用户
	forumPublicRouter := PublicRouter.Group("forum") // 公开

	forumApi := api.ApiGroupApp.ForumApi
	{
		// 公开接口
		forumPublicRouter.GET("list", forumApi.ForumList)
		forumPublicRouter.GET("detail", forumApi.ForumDetail)
		forumPublicRouter.GET("tags", forumApi.ForumTags)
	}
	{
		// 登录用户接口（管理员全量，普通用户仅本人内容，权限在 service 层校验）
		forumRouter.POST("publish", forumApi.ForumPublish)
		forumRouter.POST("upload", forumApi.ForumUpload)
		forumRouter.POST("like", forumApi.ForumLike)
		forumRouter.POST("comment", forumApi.ForumComment)
		forumRouter.GET("manageList", forumApi.ForumManageList)
		forumRouter.DELETE("delete", forumApi.ForumDelete)
		forumRouter.GET("manageComments", forumApi.ForumManageComments)
		forumRouter.DELETE("comment", forumApi.ForumCommentDelete)
	}
}
