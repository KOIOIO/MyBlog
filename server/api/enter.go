package api

import "server/service"

type ApiGroup struct {
	BaseApi
	UserApi
}

var ApiGroupApp = new(ApiGroup)

var baseService = service.ServiceGroupApp.BaseService
var userService = service.ServiceGroupApp.UserService
var qqService = service.ServiceGroupApp.QQService
var jwtService = service.ServiceGroupApp.JwtService
