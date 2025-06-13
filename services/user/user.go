package services

import (
	"context"
	"time"

	errorWrap "user-service/common/error"
	"user-service/config"
	errConstant "user-service/constants/error"
	"user-service/domain/dto"
	userRepositories "user-service/repositories/user"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	Update(ctx context.Context, req *dto.UpdateUserRequest, uuid string) (*dto.UserResponse, error)
	// GetUserLogin(context.Context) (*dto.UserResponse, error)
	FindByUUID(ctx context.Context, uuid string) (*dto.UserResponse, error)
}

type userService struct {
	userRepository userRepositories.UserRepository
}

func NewUserService(userRepository userRepositories.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

type Claims struct {
	User *dto.UserResponse
	jwt.RegisteredClaims
}

func (us *userService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	if us.isEmailExist(ctx, req.Email) {
		logrus.Warnf("email already exist %s\n", req.Email)
		return nil, errorWrap.WrapError(errConstant.ErrEmailAlreadyExists)
	}

	if us.isUsernameExist(ctx, req.Username) {
		logrus.Warnf("username already exist %s\n", req.Username)
		return nil, errorWrap.WrapError(errConstant.ErrUsernameAlreadyExists)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Errorf("failed to hashing password: %v\n", err)
		return nil, errorWrap.WrapError(errConstant.ErrInternalServerErr)
	}

	req.Password = string(hashed)
	user, err := us.userRepository.Register(ctx, req)
	if err != nil {
		return nil, err
	}

	userResponse := &dto.RegisterResponse{
		User: dto.UserResponse{
			UUID:        user.UUID,
			Name:        user.Name,
			PhoneNumber: user.PhoneNumber,
			Email:       user.Email,
			Role:        user.Role,
		},
	}

	return userResponse, nil
}

func (us *userService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := us.userRepository.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errorWrap.WrapError(errConstant.ErrUserInvalidCredentials)
	}

	expirationTime := time.Now().Add(time.Duration(config.Config.JwtExpirationTime) * time.Minute)
	data := &dto.UserResponse{
		UUID:        user.UUID,
		Name:        user.Name,
		Username:    user.Username,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		Role:        user.Role,
	}
	claims := &Claims{
		User: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.Config.JwtSecretKey)
	if err != nil {
		logrus.Errorf("failed signedString token : %v\n", err)
		return nil, errorWrap.WrapError(errConstant.ErrInternalServerErr)
	}

	return &dto.LoginResponse{
		User:  *data,
		Token: tokenString,
	}, nil
}

func (us *userService) Update(ctx context.Context, req *dto.UpdateUserRequest, uuid string) (*dto.UserResponse, error) {
	user, err := us.userRepository.FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	// cek jika username baru(req.username) exist? dan kalo beda sama user.username berarti dia mau ganti kan
	// nah tolak pergantian tersebut misal isUsernameExist bernilai true
	if us.isUsernameExist(ctx, req.Username) && user.Username != req.Username {
		return nil, errorWrap.WrapError(errConstant.ErrUsernameAlreadyExists)
	}
	if us.isEmailExist(ctx, req.Email) && user.Email != req.Email {
		return nil, errorWrap.WrapError(errConstant.ErrEmailAlreadyExists)
	}

	if req.Password != nil {
		bytes, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			logrus.Errorf("failed to hashing password: %v\n", err)
			return nil, errorWrap.WrapError(errConstant.ErrInternalServerErr)
		}
		hashed := string(bytes)
		req.Password = &hashed
	}

	result, err := us.userRepository.Update(ctx, req, uuid)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		UUID:        result.UUID,
		Name:        result.Name,
		Username:    result.Username,
		PhoneNumber: result.PhoneNumber,
		Email:       result.Email,
		Role:        result.Role,
	}, nil
}

// func (us *userService) GetUserLogin(ctx context.Context) (*dto.UserResponse, error) {
// 	ctx.
// }

func (us *userService) FindByUUID(ctx context.Context, uuid string) (*dto.UserResponse, error) {
	user, err := us.userRepository.FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		UUID:        user.UUID,
		Name:        user.Name,
		Username:    user.Username,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
	}, nil
}

// helper
func (us *userService) isUsernameExist(ctx context.Context, username string) bool {
	_, err := us.userRepository.FindByUsername(ctx, username)
	return err == nil
}

func (us *userService) isEmailExist(ctx context.Context, email string) bool {
	_, err := us.userRepository.FindByEmail(ctx, email)
	return err == nil
}
