package http

import (
	"go-integration-email/internal/service"
	http "go-integration-email/internal/transport/http/handler"

	"github.com/Lomank123/go-web-platform/logger"

	"github.com/gin-gonic/gin"
)

// SetupRouter setups gin router with all routes
func SetupRouter(router *gin.Engine, emailService service.EmailService, log logger.Logger) *gin.Engine {
	emailHandler := http.NewEmailHandler(emailService, log)
	v1 := router.Group("/api/v1")
	{
		v1.POST("/send", emailHandler.SendEmailV1)
	}

	return router
}
