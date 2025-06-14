package routes

import (
	"user-service/controllers"
	userRoute "user-service/routes/user"

	"github.com/gin-gonic/gin"
)

type RouteRegistry interface {
	Serve()
}

type routeRegistry struct {
	controller controllers.ControllerRegistry
	group      *gin.RouterGroup
}

func NewRouteRegistry(controller controllers.ControllerRegistry, group *gin.RouterGroup) RouteRegistry {
	return &routeRegistry{controller, group}
}

func (rr *routeRegistry) Serve() {
	rr.userRoute().Mount()
}

func (rr *routeRegistry) userRoute() userRoute.UserRoute {
	return userRoute.NewUserRoute(rr.controller, rr.group)
}
