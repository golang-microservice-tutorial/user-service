package services

import (
	"user-service/repositories"
	userService "user-service/services/user"
)

type ServiceRegistry interface {
	UserServices() userService.UserService
}

type serviceRegistry struct {
	repositories repositories.RepositoryRegistry
}

func NewServiceRegistry(repositories repositories.RepositoryRegistry) ServiceRegistry {
	return &serviceRegistry{repositories}
}

func (sr *serviceRegistry) UserServices() userService.UserService {
	return userService.NewUserService(sr.repositories.UserRepositories())
}
