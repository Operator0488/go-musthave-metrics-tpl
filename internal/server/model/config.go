package models

// ServerConfig - ReportInterval и PollInterval секунды
type ServerConfig struct {
	Port            string `env:"ADDRESS"`
	StoreInterval   *int   `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         *bool  `env:"RESTORE"`
	DbDsn           string `env:"DATABASE_DSN"`
}
