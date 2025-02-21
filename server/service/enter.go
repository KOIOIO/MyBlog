package service

type ServiceGroup struct {
	EsService
	BaseService
	JWTService
	GaodeService
	UserService
}

var ServiceGroupApp = new(ServiceGroup)
