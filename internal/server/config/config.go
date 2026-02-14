package config

// ServerConfig - ReportInterval и PollInterval секунды
type ServerConfig struct {
	Port            string `env:"ADDRESS"`
	StoreInterval   *int   `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         *bool  `env:"RESTORE"`
}

func NewServerConfig() ServerConfig {
	return ServerConfig{}
}
