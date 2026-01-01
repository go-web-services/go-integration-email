package dto

// Auth related email DTOs

type AuthForgotPasswordInputDTO struct {
	BaseSendEmailDTO
	Params AuthForgotPasswordParams `json:"params" binding:"required"`
}

type AuthForgotPasswordParams struct {
	ForgotPasswordLink      string `json:"forgotPasswordLink" binding:"required"`
	ExpirationTimeInMinutes int    `json:"expirationTimeInMinutes" binding:"required"`
}

type AuthEmailConfirmInputDTO struct {
	BaseSendEmailDTO
	Params AuthEmailConfirmParams `json:"params" binding:"required"`
}

type AuthEmailConfirmParams struct {
	EmailConfirmLink        string `json:"emailConfirmLink" binding:"required"`
	ExpirationTimeInMinutes int    `json:"expirationTimeInMinutes" binding:"required"`
}

type AuthOTPSigninInputDTO struct {
	BaseSendEmailDTO
	Params AuthOTPSigninParams `json:"params" binding:"required"`
}

type AuthOTPSigninParams struct {
	OTPCode                 string `json:"otpCode" binding:"required"`
	ExpirationTimeInMinutes int    `json:"expirationTimeInMinutes" binding:"required"`
}
