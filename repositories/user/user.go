package repositories

import (
	"context"

	"user-service/constants"
	"user-service/domain/dto"
	"user-service/domain/models"

	errorWrap "user-service/common/error"
	errConstants "user-service/constants/error"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository interface {
	Register(context.Context, *dto.RegisterRequest) (*models.User, error)
	Update(context.Context, *dto.UpdateUserRequest, string) (*models.User, error)
	FindByUsername(context.Context, string) (*models.User, error)
	FindByEmail(context.Context, string) (*models.User, error)
	FindByUUID(context.Context, string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (ur *userRepository) Register(ctx context.Context, req *dto.RegisterRequest) (*models.User, error) {
	user := &models.User{
		UUID:        uuid.New(),
		Name:        req.Name,
		Username:    req.Username,
		Password:    req.Password,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		RoleID:      constants.Customer,
	}

	if err := ur.db.WithContext(ctx).Create(user).Error; err != nil {
		logrus.Errorf("Failed to register user, data: %+v,\n error: %v", user, err)
		return nil, errorWrap.WrapError(errConstants.ErrSqlExecFailed)
	}
	return user, nil
}

func (ur *userRepository) Update(ctx context.Context, req *dto.UpdateUserRequest, uuid string) (*models.User, error) {
	user := &models.User{
		Name:        req.Name,
		Username:    req.Username,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Password:    *req.Password,
	}

	err := ur.db.WithContext(ctx).
		Where("uuid = ?", uuid).
		Updates(user).Error
	if err != nil {
		logrus.Errorf("Failed to update user with uuid %s, data: %+v, error: %v", uuid, user, err)
		return nil, errorWrap.WrapError(errConstants.ErrSqlExecFailed)
	}

	return user, nil
}

func (ur *userRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user *models.User
	if err := ur.db.WithContext(ctx).Preload("Role").Where("username = ?", username).First(user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logrus.Warnf("User with username %s not found", username)
			return nil, errorWrap.WrapError(errConstants.ErrUserNotFound)
		}
		logrus.Errorf("Failed to find user by username %s, error: %v", username, err)
		return nil, errorWrap.WrapError(errConstants.ErrSqlQueryFailed)
	}
	return user, nil
}

func (ur *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user *models.User
	if err := ur.db.WithContext(ctx).Preload("Role").Where("email = ?", email).First(user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logrus.Warnf("User with email %s not found", email)
			return nil, errorWrap.WrapError(errConstants.ErrUserNotFound)
		}
		logrus.Errorf("Failed to find user by email %s, error: %v", email, err)
		return nil, errorWrap.WrapError(errConstants.ErrSqlQueryFailed)
	}
	return user, nil
}

func (ur *userRepository) FindByUUID(ctx context.Context, uuid string) (*models.User, error) {
	var user *models.User
	if err := ur.db.WithContext(ctx).Preload("Role").Where("email = ?", uuid).First(user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logrus.Warnf("User with uuid %s not found", uuid)
			return nil, errorWrap.WrapError(errConstants.ErrUserNotFound)
		}
		logrus.Errorf("Failed to find user by uuid %s, error: %v", uuid, err)
		return nil, errorWrap.WrapError(errConstants.ErrSqlQueryFailed)
	}
	return user, nil
}
