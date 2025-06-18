package service

import (
	"context"

	db "user-service/db/sqlc"
	"user-service/dto"
	"user-service/pkg/helper"

	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(ctx context.Context, req dto.CreateUserRequest) (db.User, error)
	ListUsers(ctx context.Context, arg db.ListUsersParams) ([]db.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (db.User, error)
}

type userService struct {
	store db.Store
}

func NewUserService(store db.Store) UserService {
	return &userService{
		store: store,
	}
}

func (us *userService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (db.User, error) {
	arg := db.CreateuserWithMetadataParams{
		CreateUserParams: db.CreateUserParams{
			Email:       req.Email,
			FullName:    helper.ToPGText(req.FullName),
			PhoneNumber: helper.ToPGText(req.PhoneNumber),
			Role:        "user",
			AvatarUrl:   helper.ToPGText(req.AvatarURL),
		},
		UserMetadata: db.UserMetadata{
			Device: "web",
		},
	}

	result, err := us.store.CreateUserWithMetadata(ctx, arg)
	if err != nil {
		return db.User{}, err
	}

	return result.User, nil
}

func (us *userService) GetUserByID(ctx context.Context, id uuid.UUID) (db.User, error) {
	return us.store.GetUserByID(ctx, id)
}

func (us *userService) ListUsers(ctx context.Context, arg db.ListUsersParams) ([]db.User, error) {
	return us.store.ListUsers(ctx, arg)
}
