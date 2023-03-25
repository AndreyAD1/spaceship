package logger

import (
	"log"
	"os"

	"github.com/AndreyAD1/spaceship/internal/config"
)

func GetNewLogger(configuration config.StartupConfig) (*log.Logger, error) {
	logger := log.Default()
	if configuration.LogFile == "" {
		return logger, nil
	}
	logFile, err := os.Create(configuration.LogFile)
	if err != nil {
		return nil, err
	}
	logger.SetOutput(logFile)
	return logger, nil
}
