package utils

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/Lomank123/go-integration-email/config"
	internalConstants "github.com/Lomank123/go-integration-email/internal/constants"
)

// PingEmailServer checks if the SMTP server is reachable.
func PingEmailServer() error {
	host := fmt.Sprintf("%s:%d", config.Cfg.App.EmailServer, config.Cfg.App.EmailPort)

	// Create a custom dialer with TLS config
	tlsConfig := &tls.Config{
		ServerName: config.Cfg.App.EmailServer,
		MinVersion: tls.VersionTLS12,
	}

	c, err := smtp.Dial(host)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %v", host, err)
	}
	defer c.Close()

	// Start TLS
	err = c.StartTLS(tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to start TLS: %v", err)
	}

	err = c.Noop()
	if err != nil {
		return fmt.Errorf("NOOP check failed: %v", err)
	}

	return nil
}

func BuildTemplatePath(templateName string) string {
	return fmt.Sprintf("%s/%s", internalConstants.EmailTemplateDir, templateName)
}
