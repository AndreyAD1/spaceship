package logger

import (
	"io"
	"os"

	"github.com/charmbracelet/log"
	"github.com/AndreyAD1/spaceship/internal/config"
)

func GetNewLogger(configuration config.StartupConfig) (*log.Logger, error) {
	var logFile io.Writer
	var err error
	if configuration.LogFile == "" {
		logFile, err = os.CreateTemp("", "spaceship-")
	} else {
		logFile, err = os.Create(configuration.LogFile)
	}
	if err != nil {
		return nil, err
	}
	logger := log.New(logFile)
	logger.SetLevel(log.InfoLevel)
	if configuration.Debug {
		logger.SetLevel(log.DebugLevel)
	}
	return logger, nil
}
