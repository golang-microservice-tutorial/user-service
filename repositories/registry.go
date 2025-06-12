package repositories

import (
	userRepositories "user-service/repositories/user"

	"gorm.io/gorm"
)

type repositoryRegistry struct {
	db *gorm.DB
}

type RepositoryRegistry interface {
	UserRepositories() userRepositories.UserRepository
}

func NewRepositoryRegistry(db *gorm.DB) RepositoryRegistry {
	return &repositoryRegistry{db: db}
}

func (r *repositoryRegistry) UserRepositories() userRepositories.UserRepository {
	return userRepositories.NewUserRepository(r.db)
}
