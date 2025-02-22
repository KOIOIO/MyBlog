package service

type ServiceGroup struct {
	EsService
	BaseService
	JwtService
	GaodeService
	UserService
	QQService
	ImageService
}

var ServiceGroupApp = new(ServiceGroup)
