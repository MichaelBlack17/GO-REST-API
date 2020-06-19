package apiserver

type Config struct {
	BindAddr     string `toml:"bind_addr"`
	LogLevel     string `toml:"log_level"`
	DatabaseURL  string `toml:"database_url"`
	QueueLength  int    `toml:"queue_length"`
	ValidTimeOut int    `toml:"valid_timeout"`
}

// New Config ...
func NewConfig() *Config {
	return &Config{
		BindAddr:     `:8080`,
		LogLevel:     "debug",
		QueueLength:  3,
		ValidTimeOut: 1,
	}
}
