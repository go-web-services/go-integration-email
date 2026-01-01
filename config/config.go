package config

import (
	"strconv"

	platformUtils "github.com/Lomank123/go-web-platform/utils"

	platformTypes "github.com/Lomank123/go-web-platform/types"
)

type AppConfig struct {
	Port            int
	Env             platformTypes.Environment
	EmailServer     string
	EmailPort       int
	EmailUsername   string
	EmailPassword   string
	EmailFrom       string
	SwaggerBasePath string
}

type Config struct {
	App AppConfig
}

var Cfg Config

func LoadConfig() (*Config, error) {
	// Load from env variables or fallback to defaults
	portStr := platformUtils.GetEnv("APP_PORT", "8080")
	port, _ := strconv.Atoi(portStr)

	emailPortStr := platformUtils.GetEnv("EMAIL_PORT", "587")
	emailPort, _ := strconv.Atoi(emailPortStr)

	Cfg = Config{
		App: AppConfig{
			Port:          port,
			Env:           platformTypes.Environment(platformUtils.GetEnv("APP_ENV", "development")),
			EmailPort:     emailPort,
			EmailServer:   platformUtils.GetEnv("EMAIL_SERVER", ""),
			EmailUsername: platformUtils.GetEnv("EMAIL_USERNAME", ""),
			EmailPassword: platformUtils.GetEnv("EMAIL_PASSWORD", ""),
			EmailFrom:     platformUtils.GetEnv("EMAIL_FROM", ""),
		},
	}

	return &Cfg, nil
}
