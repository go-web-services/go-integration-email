package mappings

import (
	"github.com/go-web-services/go-integration-email/internal/types"

	clientConsts "github.com/go-web-services/go-integration-email/pkg/client/constants"
	clientTypes "github.com/go-web-services/go-integration-email/pkg/client/types"
)

var EmailTypeMapping = map[clientTypes.EmailType]types.EmailTemplateMapping{
	clientConsts.AuthForgotPassword: {
		HTMLTemplateFileName: "auth/forgot-password/body.html",
		TextTemplateFileName: "auth/forgot-password/body.txt",
		SubjectFileName:      "auth/forgot-password/subject.txt",
	},
	clientConsts.AuthEmailConfirm: {
		HTMLTemplateFileName: "auth/email-confirm/body.html",
		TextTemplateFileName: "auth/email-confirm/body.txt",
		SubjectFileName:      "auth/email-confirm/subject.txt",
	},
	clientConsts.AuthOTPSignin: {
		HTMLTemplateFileName: "auth/otp-signin/body.html",
		TextTemplateFileName: "auth/otp-signin/body.txt",
		SubjectFileName:      "auth/otp-signin/subject.txt",
	},
}
