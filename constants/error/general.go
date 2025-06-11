package error

import "errors"

var (
	ErrInternalServerErr  = errors.New("internal server error")
	ErrBadRequest         = errors.New("bad request")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrForbidden          = errors.New("forbidden")
	ErrTooManyRequests    = errors.New("too many requests")

	ErrSqlQueryFailed       = errors.New("sql query failed")
	ErrSqlExecFailed        = errors.New("sql exec failed")
	ErrSqlScanFailed        = errors.New("sql scan failed")
	ErrSqlTransactionFailed = errors.New("sql transaction failed")
)

var GeneralErrors = []error{
	ErrInternalServerErr,
	ErrBadRequest,
	ErrUnauthorized,
	ErrInvalidCredentials,
	ErrForbidden,
	ErrTooManyRequests,
	ErrSqlQueryFailed,
	ErrSqlExecFailed,
	ErrSqlScanFailed,
	ErrSqlTransactionFailed,
}
