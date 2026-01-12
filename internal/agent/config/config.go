package config

// AgentConfig - ReportInterval и PollInterval секунды
type AgentConfig struct {
	Port           string
	ReportInterval int
	PollInterval   int
}

func NewAgentConfig() AgentConfig {
	return AgentConfig{}
}
