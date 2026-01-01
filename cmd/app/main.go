package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	platform "github.com/Lomank123/go-web-platform/entrypoint"
	"github.com/Lomank123/go-web-platform/logger"
	platformMiddleware "github.com/Lomank123/go-web-platform/middleware"
	"github.com/gin-gonic/gin"

	"github.com/Lomank123/go-integration-email/config"
	"github.com/Lomank123/go-integration-email/docs"
	"github.com/Lomank123/go-integration-email/internal/service"
	emailHTTP "github.com/Lomank123/go-integration-email/internal/transport/http"
	"github.com/Lomank123/go-integration-email/internal/utils"
)

// @title           Email Integration API
// @version         1.0
// @basePath        /api

func main() {
	// Load configuration from environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize custom logger
	logg := logger.NewLogger(cfg.App.Env)

	// Prepare services
	emailService := service.NewEmailService(logg)

	// Prepare HTTP server (router)
	router := gin.New()
	// Platform integration
	platform.SetupPlatform(
		router,
		logg,
		utils.PingEmailServer,
		platformMiddleware.DefaultLoggingConfig(),
		nil,
		cfg.App.Env,
	)
	emailHTTP.SetupRouter(router, emailService, logg)

	// Swagger docs
	swaggerBasePath := "/api"
	if cfg.App.SwaggerBasePath != "" {
		swaggerBasePath = "/" + cfg.App.SwaggerBasePath + swaggerBasePath
	}
	docs.SwaggerInfo.BasePath = swaggerBasePath

	// Start HTTP server
	serverAddr := fmt.Sprintf(":%d", cfg.App.Port)
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}
	logg.Info("Starting server on port ", cfg.App.Port)
	go func() {
		if e := router.Run(serverAddr); e != nil {
			logg.Fatal("Failed to start HTTP server: ", e)
		}
	}()

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	logg.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if e := srv.Shutdown(ctx); e != nil {
		logg.Fatal("Server forced to shutdown: ", e)
	}

	logg.Info("Server stopped.")
}
