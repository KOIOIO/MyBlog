package service

type ServiceGroup struct {
	EsService
	BaseService
	JWTService
	GaodeService
}

var ServiceGroupApp = new(ServiceGroup)
