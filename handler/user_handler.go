package handler

import (
	"net/http"

	"user-service/constants"
	db "user-service/db/sqlc"
	"user-service/dto"
	"user-service/pkg/helper"
	"user-service/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type userHandler struct {
	userService service.UserService
	validate    *validator.Validate
}

func NewUserHandler(us service.UserService, validator *validator.Validate) *userHandler {
	return &userHandler{userService: us, validate: validator}
}

func (h *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := helper.BindRequest(r, &req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.validate.Struct(&req); err != nil {
		err := helper.GenerateMessage(err, constants.FromRequestBody)
		helper.WriteError(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.userService.CreateUser(r.Context(), req)
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.WriteCreated(w, user)
}

func (h *userHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	uuid, err := helper.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, "UUID is not valid")
		return
	}
	user, err := h.userService.GetUserByID(r.Context(), uuid)
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.WriteSuccess(w, user)
}

func (h *userHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	req := dto.ListUsersRequest{
		Search: r.URL.Query().Get("search"),
		Offset: helper.ParseInt32(r.URL.Query().Get("offset"), 0),
		Limit:  helper.ParseInt32(r.URL.Query().Get("limit"), 10),
	}
	if err := h.validate.Struct(&req); err != nil {
		err := helper.GenerateMessage(err, constants.FromQueryParams)
		helper.WriteError(w, http.StatusBadRequest, err)
		return
	}

	arg := db.ListUsersParams{
		Search: helper.StringToPGTextValid(req.Search),
		Offset: req.Offset,
		Limit:  req.Limit,
	}
	users, err := h.userService.ListUsers(r.Context(), arg)
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.WriteSuccess(w, users)
}
