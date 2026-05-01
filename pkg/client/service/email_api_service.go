package service

import (
	"github.com/gin-gonic/gin"

	"github.com/go-web-services/go-integration-email/pkg/client/dto"

	platformUtils "github.com/go-web-services/go-web-platform/utils"

	"fmt"
)

type EmailAPIService interface {
	SendEmailV1(context *gin.Context, payload dto.EmailPayload) (dto.SendEmailOutputDTO, error)
}

type emailAPIService struct {
	apiURL string
}

func NewEmailAPIService(host string) EmailAPIService {
	return &emailAPIService{apiURL: fmt.Sprintf("%s/api/v1", host)}
}

func (s *emailAPIService) SendEmailV1(context *gin.Context, payload dto.EmailPayload) (dto.SendEmailOutputDTO, error) {
	url := fmt.Sprintf("%s/send", s.apiURL)
	var responseBody dto.SendEmailOutputDTO

	err := platformUtils.SendRequest("POST", url, payload, &responseBody, context)
	return responseBody, err
}
