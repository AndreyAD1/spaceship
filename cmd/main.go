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
	// TODO fix the launch without a log_file
	newLogger, err := logger.GetNewLogger(configuration)
	if err != nil {
		log.Fatal(err)
	}
	app := application.GetApplication(newLogger)
	err = app.Run()
	log.Println(err)
}
