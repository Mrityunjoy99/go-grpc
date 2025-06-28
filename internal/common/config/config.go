package config

type Config struct {
	Server ServerConfig `validate:"required"`
	App    AppConfig    `validate:"required"`
}

type ServerConfig struct {
	Port string `default:"50051" validate:"required,numeric"`
}

type AppConfig struct {
	LogToFile bool `default:"false" validate:"required"`
}
