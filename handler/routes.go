package handler

import (
	"user-service/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func NewRegisterRoutes(service service.ServiceRegistry, r chi.Router, validator *validator.Validate) {
	userHandler := NewUserHandler(service.UserService(), validator)

	r.Route("/users", func(r chi.Router) {
		r.Get("/", userHandler.ListUsers)
		r.Post("/", userHandler.CreateUser)
		r.Get("/{id}", userHandler.GetUserByID)
	})
}
