package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/AndreyAD1/spaceship/internal/config"
	"github.com/charmbracelet/log"
)

func NewLogger(configuration config.StartupConfig) (*log.Logger, error) {
	var logFile io.Writer
	var err error
	if configuration.LogFile == "" {
		logFile, err = os.CreateTemp("", "spaceship-")
	} else {
		logFile, err = os.Create(configuration.LogFile)
	}
	if err != nil {
		return nil, fmt.Errorf("can not create a new logger: %w", err)
	}
	logger := log.New(logFile)
	logger.SetLevel(log.InfoLevel)
	if configuration.Debug {
		logger.SetLevel(log.DebugLevel)
	}
	return logger, nil
}
