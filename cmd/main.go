package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"user-service/config"
	db "user-service/db/sqlc"
	"user-service/handler"
	logger "user-service/pkg"
	"user-service/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServerOptions struct {
	Config *config.AppConfig
	DB     *pgxpool.Pool
}

func Run(opts ServerOptions) {
	// store
	store := db.NewStore(opts.DB)

	// service
	service := service.NewServiceRegistry(store)

	// server
	router := chi.NewRouter()
	if router == nil {
		logger.Log.Fatal("router cannot be nil")
	}
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	validator := validator.New()
	// routes
	handler.NewRegisterRoutes(service, router, validator)

	port := opts.Config.AppPort
	logger.Log.Infof("port: %s", port)
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)

	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// Graceful shutdown
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		logger.Log.Info("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Log.Errorf("HTTP server Shutdown: %v", err)
		}

		opts.DB.Close()
		close(idleConnsClosed)
	}()

	logger.Log.Infof("🚀 Server running on http://localhost%s", addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.Log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}
