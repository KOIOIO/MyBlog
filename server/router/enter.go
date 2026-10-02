package router

type RouterGroup struct {
	BaseRouter
	UserRouter
	ImageRouter
	ArticleRouter
	CommentRouter
	AdvertisementRouter
	FriendLinkRouter
	FeedbackRouter
	WebsiteRouter
	ConfigRouter
	ForumRouter
}

var RouterGroupApp = new(RouterGroup)
