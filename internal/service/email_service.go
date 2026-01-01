package service

import (
	"bytes"
	"fmt"
	"html/template"
	"mime/multipart"
	"net/smtp"
	"time"

	"github.com/Lomank123/go-integration-email/config"
	"github.com/Lomank123/go-integration-email/internal/mappings"
	"github.com/Lomank123/go-integration-email/internal/utils"

	"github.com/Lomank123/go-web-platform/logger"

	clientTypes "github.com/Lomank123/go-integration-email/pkg/client/types"
)

type EmailService interface {
	SendEmail(
		emailType clientTypes.EmailType,
		recipients []string,
		params map[string]any,
	) error
}

type emailService struct {
	log logger.Logger
}

func NewEmailService(log logger.Logger) EmailService {
	return &emailService{log: log}
}

// SendEmail sends an email to the specified recipients with the parameters.
// 2 types of templates are used: HTML and plain text.
func (s *emailService) SendEmail(
	emailType clientTypes.EmailType,
	recipients []string,
	params map[string]any,
) error {
	var subject bytes.Buffer
	var htmlBody bytes.Buffer
	var textBody bytes.Buffer
	fileNames := mappings.EmailTypeMapping[emailType]

	// Inject current year for copyright notices
	params["Year"] = time.Now().Year()

	err := s.processTemplate(&subject, utils.BuildTemplatePath(fileNames.SubjectFileName), params)
	if err != nil {
		return err
	}
	err = s.processTemplate(&htmlBody, utils.BuildTemplatePath(fileNames.HTMLTemplateFileName), params)
	if err != nil {
		return err
	}
	err = s.processTemplate(&textBody, utils.BuildTemplatePath(fileNames.TextTemplateFileName), params)
	if err != nil {
		return err
	}

	msg, err := s.buildEmailBody(&subject, &textBody, &htmlBody)
	if err != nil {
		return err
	}

	host := fmt.Sprintf("%s:%d", config.Cfg.App.EmailServer, config.Cfg.App.EmailPort)
	auth := smtp.PlainAuth(
		"",
		config.Cfg.App.EmailUsername,
		config.Cfg.App.EmailPassword,
		config.Cfg.App.EmailServer,
	)
	err = smtp.SendMail(
		host,
		auth,
		config.Cfg.App.EmailFrom,
		recipients,
		msg.Bytes(),
	)
	if err != nil {
		return err
	}

	return nil
}

// Insert params inside the template
// Watch out for extra newlines in template files.
// Even a single extra one can ruin the whole email.
func (s *emailService) processTemplate(
	buffer *bytes.Buffer,
	path string,
	params map[string]any,
) error {
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return err
	}
	err = tmpl.Execute(buffer, params)
	if err != nil {
		return err
	}

	return nil
}

// Build email body with multiple representations (plain text and HTML)
func (s *emailService) buildEmailBody(
	subject *bytes.Buffer,
	text *bytes.Buffer,
	html *bytes.Buffer,
) (bytes.Buffer, error) {
	// Create a MIME multipart message
	var msg bytes.Buffer
	writer := multipart.NewWriter(&msg)

	// Write the MIME headers including subject
	msg.WriteString(fmt.Sprintf("Subject: %s\n", subject.String()))
	msg.WriteString("MIME-Version: 1.0\n")
	msg.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=%s\n\n", writer.Boundary()))

	// Write the plain text part
	textPart, err := writer.CreatePart(map[string][]string{
		"Content-Type": {"text/plain; charset=UTF-8"},
	})
	if err != nil {
		return bytes.Buffer{}, err
	}

	_, err = textPart.Write(text.Bytes())
	if err != nil {
		return bytes.Buffer{}, err
	}

	// Write the HTML part
	htmlPart, err := writer.CreatePart(map[string][]string{
		"Content-Type": {"text/html; charset=UTF-8"},
	})
	if err != nil {
		return bytes.Buffer{}, err
	}

	_, err = htmlPart.Write(html.Bytes())
	if err != nil {
		return bytes.Buffer{}, err
	}

	err = writer.Close()
	if err != nil {
		return bytes.Buffer{}, err
	}

	return msg, nil
}
