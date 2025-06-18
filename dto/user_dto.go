package dto

type CreateUserRequest struct {
	Email       string `json:"email" validate:"required,email"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	AvatarURL   string `json:"avatar_url"`
}

type ListUsersRequest struct {
	Search string `form:"search"`
	Offset int32  `form:"offset" validate:"gte=0"`
	Limit  int32  `form:"limit" validate:"gte=1,lte=100"`
}
