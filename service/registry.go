package service

import db "user-service/db/sqlc"

type ServiceRegistry interface {
	UserService() UserService
}

type serviceRegistry struct {
	store db.Store
}

func NewServiceRegistry(store db.Store) ServiceRegistry {
	return &serviceRegistry{
		store: store,
	}
}

func (sr *serviceRegistry) UserService() UserService {
	return NewUserService(sr.store)
}
