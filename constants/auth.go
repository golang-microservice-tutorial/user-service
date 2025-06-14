package constants

type ContextKey string

const (
	UserLogin        = ContextKey("user_login")
	Token            = "token"
	PasswordNotmatch = "password not match"
)
