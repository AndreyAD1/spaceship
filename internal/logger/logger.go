package logger

import (
	"io"
	"log"
	"os"

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
	logger := log.New(logFile, "", log.LstdFlags)
	return logger, nil
}
