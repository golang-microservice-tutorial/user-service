package controllers

import (
	userController "user-service/controllers/user"
	"user-service/services"

	"github.com/go-playground/validator/v10"
)

type ControllerRegistry interface {
	UserController() userController.UserController
}

type controllerRegistry struct {
	services  services.ServiceRegistry
	validator *validator.Validate
}

func NewControllerRegistry(services services.ServiceRegistry, validator *validator.Validate) ControllerRegistry {
	return &controllerRegistry{services, validator}
}

func (cr *controllerRegistry) UserController() userController.UserController {
	return userController.NewUserController(cr.services, cr.validator)
}
