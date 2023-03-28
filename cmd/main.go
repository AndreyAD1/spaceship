package main

import (
	"flag"
	"log"

	"github.com/AndreyAD1/spaceship/internal/application"
	"github.com/AndreyAD1/spaceship/internal/config"
	"github.com/AndreyAD1/spaceship/internal/logger"
	"github.com/caarlos0/env/v7"
)

func main() {
	debug := flag.String("debug", "", "run in a debug mode")
	logFile := flag.String("log_file", "", "write logs to this file")
	flag.Parse()
	configuration := config.StartupConfig{}
	err := env.Parse(&configuration)
	if err != nil {
		log.Fatal(err)
	}
	if *debug != "" {
		configuration.Debug = *debug == "true"
	}
	if *logFile != "" {
		configuration.LogFile = *logFile
	}
	newLogger, err := logger.GetNewLogger(configuration)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if p := recover(); p != nil {
			newLogger.Errorf("Internal error: %v", p)
		}
	}()
	app := application.GetApplication(newLogger)
	newLogger.Debug("run application")
	err = app.Run()
	newLogger.Debug("finish application: %v", err)
	newLogger.Print(err)
}
