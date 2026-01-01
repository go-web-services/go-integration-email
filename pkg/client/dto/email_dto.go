package dto

import (
	"github.com/Lomank123/go-integration-email/pkg/client/types"
)

// All email DTOs must implement this interface
type EmailPayload any

// Used to map email send response payload
type SendEmailOutputDTO struct {
	Message string `json:"message"`
}

// Contains the base fields for sending an email
type BaseSendEmailDTO struct {
	EmailType  types.EmailType `json:"emailType" binding:"required"`
	Recipients []string        `json:"recipients" binding:"required"`
}

// Used to ensure the incoming params are correct
type GeneralSendEmailDTO struct {
	BaseSendEmailDTO
	Params map[string]any `json:"params,omitempty"`
}
