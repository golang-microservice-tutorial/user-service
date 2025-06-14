package middlewares

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"user-service/common/response"
	"user-service/config"
	"user-service/constants"
	errConstants "user-service/constants/error"
	userService "user-service/services/user"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

func CustomRecovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("recovered from panic: %v\n", r)
				msg := errConstants.ErrInternalServerErr.Error()
				response.HttpResponse(response.ParamHTTPResponse{
					Code:    http.StatusInternalServerError,
					Message: &msg,
					Gin:     ctx,
				})
				ctx.Abort()
				return
			}
		}()

		ctx.Next()
	}
}

func RateLimiter(lmt *limiter.Limiter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := tollbooth.LimitByRequest(lmt, ctx.Writer, ctx.Request)
		if err != nil {
			msg := err.Error()
			response.HttpResponse(response.ParamHTTPResponse{
				Code:    http.StatusTooManyRequests,
				Message: &msg,
				Gin:     ctx,
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func WithAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error

		// validate api key
		if err = validateApiKey(ctx); err != nil {
			sendUnauthorized(ctx)
			return
		}

		// get token from header
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			logrus.Warn("Authorization header not found")
			sendUnauthorized(ctx)
			return
		}

		token := extractBearerToken(authHeader)
		if token == "" {
			logrus.Warn("Authorization token not found")
			sendUnauthorized(ctx)
			return
		}

		// validate token
		if err := validateBearerToken(ctx, token); err != nil {
			sendUnauthorized(ctx)
			return
		}

		ctx.Next()
	}
}

// ? #################### helper ####################
func extractBearerToken(token string) string {
	parts := strings.SplitN(token, " ", 2)
	if len(parts) == 2 || strings.ToLower(parts[0]) == "bearer" || parts[1] != "" {
		return parts[1]
	}
	return ""
}

func sendUnauthorized(ctx *gin.Context) {
	msg := errConstants.ErrUnauthorized.Error()
	response.HttpResponse(response.ParamHTTPResponse{
		Code:    http.StatusUnauthorized,
		Message: &msg,
		Gin:     ctx,
	})
	ctx.Abort()
}

func validateApiKey(ctx *gin.Context) error {
	apiKey := ctx.GetHeader(constants.XApiKey)
	requestAt := ctx.GetHeader(constants.XRequestAt)
	serviceName := ctx.GetHeader(constants.XServiceName)
	signatureKey := config.Config.SignatureKey

	vaidateKey := fmt.Sprintf("%s:%s:%s", serviceName, signatureKey, requestAt)
	hash := sha256.New()
	hash.Write([]byte(vaidateKey))
	resultHash := hex.EncodeToString(hash.Sum(nil))

	if apiKey != resultHash {
		return errConstants.ErrUnauthorized
	}
	return nil
}

func validateBearerToken(ctx *gin.Context, token string) error {
	claims := &userService.Claims{}
	tokenJwt, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			logrus.Warn("unexpected signing method")
			return nil, errConstants.ErrInvalidToken
		}
		secret := []byte(config.Config.JwtSecretKey)
		return secret, nil
	})

	if err != nil || !tokenJwt.Valid {
		logrus.Warnf("invalid token or error: %v", err)
		return errConstants.ErrUnauthorized
	}

	// set context yang ada di net/http (standard library)
	ctx.Request = ctx.Request.WithContext(
		context.WithValue(ctx.Request.Context(),
			constants.UserLogin,
			claims.User,
		))
	// set context yang ada di gin
	ctx.Set(constants.Token, token)

	return nil
}
