package service

type ServiceGroup struct {
	EsService
	BaseService
	JwtService
	GaodeService
	UserService
	QQService
}

var ServiceGroupApp = new(ServiceGroup)
