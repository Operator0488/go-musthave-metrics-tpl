package config

// AgentConfig - ReportInterval и PollInterval секунды
type AgentConfig struct {
	Port           string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

func NewAgentConfig() AgentConfig {
	return AgentConfig{}
}
