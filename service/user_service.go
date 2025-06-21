package service

import (
	"context"

	db "user-service/db/sqlc"
	"user-service/dto"
	logger "user-service/pkg"
	"user-service/pkg/helper"

	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(ctx context.Context, req dto.CreateUserRequest) (dto.UserResponse, error)
	ListUsers(ctx context.Context, arg db.ListUsersParams) ([]dto.UserResponse, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (dto.UserResponse, error)
	GetUserWithMetadata(ctx context.Context, id uuid.UUID) (db.GetUserWithMetadataRow, error)
}

type userService struct {
	store db.Store
}

func NewUserService(store db.Store) UserService {
	return &userService{
		store: store,
	}
}

func (us *userService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (dto.UserResponse, error) {
	arg := db.CreateuserWithMetadataParams{
		CreateUserParams: db.CreateUserParams{
			Email:       req.Email,
			FullName:    helper.StringToPGTextValid(req.FullName),
			PhoneNumber: helper.StringToPGText(req.PhoneNumber),
			Role:        "user",
			AvatarUrl:   helper.StringToPGText(req.AvatarURL),
		},
		UserMetadata: db.UserMetadata{
			Device: "web",
		},
	}

	result, err := us.store.CreateUserWithMetadata(ctx, arg)
	if err != nil {
		logger.Log.Errorf("failed to create user with metadata: %v", err)
		return dto.UserResponse{}, err
	}

	response := dto.UserResponse{
		ID:          result.User.ID,
		Email:       result.User.Email,
		FullName:    helper.PGTextToStringOrNil(result.User.FullName),
		PhoneNumber: helper.PGTextToStringOrNil(result.User.PhoneNumber),
		Role:        result.User.Role,
		AvatarUrl:   helper.PGTextToStringOrNil(result.User.AvatarUrl),
		CreatedAt:   helper.PGTimestamptzToTime(result.User.CreatedAt),
		UpdatedAt:   helper.PGTimestamptzToTime(result.User.UpdatedAt),
		DeletedAt:   helper.PGTimestamptzToTimePtr(result.DeletedAt),
	}

	return response, nil
}

func (us *userService) GetUserByID(ctx context.Context, id uuid.UUID) (dto.UserResponse, error) {
	result, err := us.store.GetUserByID(ctx, id)
	if err != nil {
		logger.Log.Errorf("failed to get user by id: %v", err)
		return dto.UserResponse{}, err
	}

	response := dto.UserResponse{
		ID:          result.ID,
		Email:       result.Email,
		FullName:    helper.PGTextToStringOrNil(result.FullName),
		PhoneNumber: helper.PGTextToStringOrNil(result.PhoneNumber),
		Role:        result.Role,
		AvatarUrl:   helper.PGTextToStringOrNil(result.AvatarUrl),
		CreatedAt:   helper.PGTimestamptzToTime(result.CreatedAt),
		UpdatedAt:   helper.PGTimestamptzToTime(result.UpdatedAt),
		DeletedAt:   helper.PGTimestamptzToTimePtr(result.DeletedAt),
	}

	return response, nil
}

func (us *userService) GetUserByEmail(ctx context.Context, email string) (dto.UserResponse, error) {
	result, err := us.store.GetUserByEmail(ctx, email)
	if err != nil {
		logger.Log.Errorf("failed to get user by email: %v", err)
		return dto.UserResponse{}, err
	}

	response := dto.UserResponse{
		ID:          result.ID,
		Email:       result.Email,
		FullName:    helper.PGTextToStringOrNil(result.FullName),
		PhoneNumber: helper.PGTextToStringOrNil(result.PhoneNumber),
		Role:        result.Role,
		AvatarUrl:   helper.PGTextToStringOrNil(result.AvatarUrl),
		CreatedAt:   helper.PGTimestamptzToTime(result.CreatedAt),
		UpdatedAt:   helper.PGTimestamptzToTime(result.UpdatedAt),
		DeletedAt:   helper.PGTimestamptzToTimePtr(result.DeletedAt),
	}

	return response, nil
}

func (us *userService) GetUserWithMetadata(ctx context.Context, id uuid.UUID) (db.GetUserWithMetadataRow, error) {
	return us.store.GetUserWithMetadata(ctx, id)
}

func (us *userService) ListUsers(ctx context.Context, arg db.ListUsersParams) ([]dto.UserResponse, error) {
	result, err := us.store.ListUsers(ctx, arg)
	if err != nil {
		logger.Log.Errorf("failed to get list users: %v", err)
		return nil, err
	}

	var response []dto.UserResponse
	for _, item := range result {
		response = append(response, dto.UserResponse{
			ID:          item.ID,
			Email:       item.Email,
			FullName:    helper.PGTextToStringOrNil(item.FullName),
			PhoneNumber: helper.PGTextToStringOrNil(item.PhoneNumber),
			Role:        item.Role,
			AvatarUrl:   helper.PGTextToStringOrNil(item.AvatarUrl),
			CreatedAt:   helper.PGTimestamptzToTime(item.CreatedAt),
			UpdatedAt:   helper.PGTimestamptzToTime(item.UpdatedAt),
			DeletedAt:   helper.PGTimestamptzToTimePtr(item.DeletedAt),
		})
	}

	return response, nil
}
