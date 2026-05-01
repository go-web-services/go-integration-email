package constants

import (
	"github.com/go-web-services/go-integration-email/pkg/client/types"
)

var (
	AuthForgotPassword types.EmailType = "AuthForgotPassword"
	AuthEmailConfirm   types.EmailType = "AuthEmailConfirm"
	AuthOTPSignin      types.EmailType = "AuthOTPSignin"
)
