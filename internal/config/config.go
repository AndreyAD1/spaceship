package config

type StartupConfig struct {
	Debug   bool   `env:"DEBUG"`
	LogFile string `env:"LOGFILE"`
}
