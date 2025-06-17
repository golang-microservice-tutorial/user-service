package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"auth-service/config"
	logger "auth-service/pkg"

	"github.com/go-chi/chi/v5"
)

type ServerOptions struct {
	Config *config.AppConfig
}

func Run(opts ServerOptions) {
	router := chi.NewRouter()
	if router == nil {
		logger.Log.Fatal("router cannot be nil")
	}

	port := opts.Config.AppPort
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

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Log.Errorf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	logger.Log.Infof("🚀 Server running on http://localhost%s", addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.Log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}
