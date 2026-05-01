package http

import (
	"net/http"

	"github.com/go-web-services/go-integration-email/internal/service"

	"github.com/go-web-services/go-web-platform/logger"

	"github.com/gin-gonic/gin"

	clientDTO "github.com/go-web-services/go-integration-email/pkg/client/dto"

	platformError "github.com/go-web-services/go-web-platform/error"
)

type EmailHandler interface {
	SendEmailV1(c *gin.Context)
}

type emailHandler struct {
	service service.EmailService
	log     logger.Logger
}

func NewEmailHandler(service service.EmailService, log logger.Logger) EmailHandler {
	return &emailHandler{service: service, log: log}
}

// SendEmailV1
// @Summary Send an email
// @Description Send an email to the specified recipients with the parameters that will be inserted
// @Accept json
// @Produce json
// @Param SendEmailRequest body clientDTO.GeneralSendEmailDTO true "Send Email Request"
// @Success 200 {object} clientDTO.SendEmailOutputDTO
// @Router /v1/send [post]
func (eh *emailHandler) SendEmailV1(c *gin.Context) {
	var payload clientDTO.GeneralSendEmailDTO

	// Parsing request body
	if err := c.ShouldBindJSON(&payload); err != nil {
		_ = c.Error(platformError.ErrInvalidRequestPayload)
		return
	}

	err := eh.service.SendEmail(payload.EmailType, payload.Recipients, payload.Params)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(
		http.StatusOK,
		clientDTO.SendEmailOutputDTO{Message: "Email sent successfully"},
	)
}
