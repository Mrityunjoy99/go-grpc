package config

type Config struct {
	Server ServerConfig
}

type ServerConfig struct {
	Port string
}
