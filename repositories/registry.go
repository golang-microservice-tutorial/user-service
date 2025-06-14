package repositories

import (
	userRepository "user-service/repositories/user"

	"gorm.io/gorm"
)

type repositoryRegistry struct {
	db *gorm.DB
}

type RepositoryRegistry interface {
	UserRepositories() userRepository.UserRepository
}

func NewRepositoryRegistry(db *gorm.DB) RepositoryRegistry {
	return &repositoryRegistry{db: db}
}

func (r *repositoryRegistry) UserRepositories() userRepository.UserRepository {
	return userRepository.NewUserRepository(r.db)
}
