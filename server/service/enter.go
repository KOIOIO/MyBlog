package service

type ServiceGroup struct {
	EsService
	BaseService
	JwtService
	GaodeService
	UserService
	ImageService
	ArticleService
	CommentService
	AdvertisementService
	FriendLinkService
	FeedbackService
	WebsiteService
	HotSearchService
	CalendarService
	ConfigService
	ForumService
}

var ServiceGroupApp = new(ServiceGroup)
