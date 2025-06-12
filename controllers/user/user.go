package controllers

import (
	"errors"
	"net/http"

	"user-service/common/response"
	"user-service/domain/dto"
	"user-service/services"

	errWrap "user-service/common/error"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserController interface {
	Login(*gin.Context)
	Register(*gin.Context)
	Update(*gin.Context)
	GetUserLogin(*gin.Context)
	GetUserByUUID(*gin.Context)
}

type userController struct {
	service   services.ServiceRegistry
	validator *validator.Validate
}

func NewUserController(service services.ServiceRegistry, validator *validator.Validate) UserController {
	return &userController{service, validator}
}

func (uc *userController) Login(ctx *gin.Context) {
	var req dto.LoginRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	if err = uc.validator.Struct(&req); err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errWrap.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHTTPResponse{
			Code:    http.StatusBadRequest,
			Message: &errMessage,
			Err:     err,
			Gin:     ctx,
			Data:    errResponse,
		})
		return
	}

	res, err := uc.service.UserServices().Login(ctx, &req)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}
	response.HttpResponse(response.ParamHTTPResponse{
		Code:  http.StatusOK,
		Gin:   ctx,
		Data:  res.User,
		Token: &res.Token,
	})
}

func (uc *userController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	// confirm password check
	if req.ConfirmPassword != req.Password {
		errMessage := "password not match"
		response.HttpResponse(response.ParamHTTPResponse{
			Code:    http.StatusBadRequest,
			Err:     errors.New("password not match"),
			Gin:     ctx,
			Message: &errMessage,
		})
		return
	}

	if err = uc.validator.Struct(&req); err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errWrap.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHTTPResponse{
			Code:    http.StatusBadRequest,
			Message: &errMessage,
			Err:     err,
			Gin:     ctx,
			Data:    errResponse,
		})
		return
	}

	res, err := uc.service.UserServices().Register(ctx, &req)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}
	response.HttpResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Gin:  ctx,
		Data: res.User,
	})
}

func (uc *userController) Update(ctx *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (uc *userController) GetUserLogin(ctx *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (uc *userController) GetUserByUUID(ctx *gin.Context) {
	panic("not implemented") // TODO: Implement
}
