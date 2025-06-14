package cmd

import (
	"fmt"
	"net/http"
	"time"

	"user-service/common/response"
	"user-service/config"
	"user-service/constants"
	"user-service/controllers"
	"user-service/database/seeders"
	"user-service/domain/models"
	"user-service/middlewares"
	"user-service/repositories"
	"user-service/routes"
	"user-service/services"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func Run() {
	err := command.Execute()
	if err != nil {
		panic(err)
	}
}

var command = &cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		// load env with godotenv
		_ = godotenv.Load()

		// init config
		config.Init()

		// init database
		db, err := config.InitDatabase()
		if err != nil {
			panic(err)
		}

		// ubah timezone ke asia/jakarta
		loc, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			panic(err)
		}
		time.Local = loc

		// run migration
		err = db.AutoMigrate(
			&models.Role{},
			&models.User{},
		)
		if err != nil {
			panic(err)
		}

		// run seeders
		seed := seeders.NewSeederRegistry(db)
		seed.Run()

		// golang validator
		validator := validator.New()

		// repository
		repositories := repositories.NewRepositoryRegistry(db)

		// service
		services := services.NewServiceRegistry(repositories)

		// controller
		controllers := controllers.NewControllerRegistry(services, validator)

		// setuprouter
		router := gin.New()
		router.Use(gin.Logger())
		router.Use(middlewares.CustomRecovery())
		router.NoRoute(func(ctx *gin.Context) {
			ctx.JSON(http.StatusNotFound, response.Response{
				Status:  constants.Error,
				Message: http.StatusText(http.StatusNotFound),
			})
		})
		router.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, response.Response{
				Status:  constants.Success,
				Message: http.StatusText(http.StatusOK),
			})
		})
		router.Use(func(ctx *gin.Context) {
			ctx.Writer.Header().Set("Access-Control-Allow-Origin", "localhost:3000")
			ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			ctx.Writer.Header().Set("Access-Control-Allow-Headers",
				"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, x-service-name, x-api-key, x-request-at, Authorization, accept, origin, Cache-Control, X-Requested-With")
			ctx.Next()
		})

		// middleware limiter
		lmt := tollbooth.NewLimiter(
			float64(config.Config.RateLimiterMaxRequest),
			&limiter.ExpirableOptions{
				DefaultExpirationTTL: time.Duration(config.Config.RateLimiterTimeSecond) * time.Second,
			})
		router.Use(middlewares.RateLimiter(lmt))

		group := router.Group("/api/v1")

		route := routes.NewRouteRegistry(controllers, group)
		route.Serve()

		router.Run(fmt.Sprintf(":%d", config.Config.Port))
	},
}
