package config

// ServerConfig - ReportInterval и PollInterval секунды
type ServerConfig struct {
	Port string
}

func NewServerConfig() ServerConfig {
	return ServerConfig{}
}
