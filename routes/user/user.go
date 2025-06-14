package routes

import (
	"user-service/controllers"
	"user-service/middlewares"

	"github.com/gin-gonic/gin"
)

type UserRoute interface {
	Mount()
}

type userRoute struct {
	controller controllers.ControllerRegistry
	group      *gin.RouterGroup
}

func NewUserRoute(controller controllers.ControllerRegistry, group *gin.RouterGroup) UserRoute {
	return &userRoute{controller, group}
}

func (ur *userRoute) Mount() {
	group := ur.group.Group("/auth")

	group.GET("/user", middlewares.WithAuth(), ur.controller.UserController().GetUserLogin)
	group.GET("/:uuid", middlewares.WithAuth(), ur.controller.UserController().GetUserByUUID)
	group.PUT("/update", middlewares.WithAuth(), ur.controller.UserController().Update)

	group.POST("/login", ur.controller.UserController().Login)
	group.POST("/register", ur.controller.UserController().Register)
}
