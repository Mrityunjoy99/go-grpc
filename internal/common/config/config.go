// Package config provides a configuration loading and validation functionality.
package config

// Config represents the application configuration. This will contain all secrets and configs for the application.
type Config struct {
	Server ServerConfig `validate:"required"`
	App    AppConfig    `validate:"required"`
}

// ServerConfig represents the server configuration.
type ServerConfig struct {
	Port string `default:"50051" validate:"required,numeric"`
}

// AppConfig represents the application configuration.
type AppConfig struct {
	LogToFile bool `default:"false"`
}
