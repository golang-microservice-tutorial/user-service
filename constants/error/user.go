package error

import "errors"

var (
	ErrUserNotFound           = errors.New("user not found")
	ErrUsernameAlreadyExists  = errors.New("username already exists")
	ErrUserInvalidCredentials = errors.New("username or password is invalid")
	ErrUserUnauthorized       = errors.New("user unauthorized")
)

var UserErrors = []error{
	ErrUserNotFound,
	ErrUsernameAlreadyExists,
	ErrUserInvalidCredentials,
	ErrUserUnauthorized,
}
