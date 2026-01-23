package config

// ServerConfig - ReportInterval и PollInterval секунды
type ServerConfig struct {
	Port string `env:"ADDRESS"`
}

func NewServerConfig() ServerConfig {
	return ServerConfig{}
}
