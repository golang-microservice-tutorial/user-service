package dto

import (
	"user-service/domain/models"

	"github.com/google/uuid"
)

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	UUID        uuid.UUID   `json:"uuid"`
	Name        string      `json:"name"`
	Username    string      `json:"username"`
	PhoneNumber string      `json:"phoneNumber"`
	Email       string      `json:"email"`
	Role        models.Role `json:"role"`
}

type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

type RegisterRequest struct {
	Name            string `json:"name" validate:"required"`
	Username        string `json:"username" validate:"required"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirmPassword" validate:"required"`
	PhoneNumber     string `json:"phoneNumber" validate:"required,min=10,max=15"`
	Email           string `json:"email" validate:"required,email"`
	RoleID          uint
}

type RegisterResponse struct {
	User UserResponse `json:"user"`
}

type UpdateUserRequest struct {
	Name            string  `json:"name" validate:"required"`
	Username        string  `json:"username" validate:"required"`
	Password        *string `json:"password,omitempty"`
	ConfirmPassword *string `json:"confirmPassword,omitempty"`
	PhoneNumber     string  `json:"phoneNumber" validate:"required,min=10,max=15"`
	Email           string  `json:"email" validate:"required,email"`
	RoleID          uint
}
