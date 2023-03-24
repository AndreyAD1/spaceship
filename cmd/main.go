package main

import (
	"flag"
	"log"

	"github.com/AndreyAD1/spaceship/internal/application"
	"github.com/AndreyAD1/spaceship/internal/config"
	"github.com/caarlos0/env/v7"
)


func main() {
	debug := flag.String("debug", "", "run in a debug mode")
	flag.Parse()
	configuration := config.StartupConfig{}
	err := env.Parse(&configuration)
	if err != nil {
		log.Fatal(err)
	}
	if *debug != "" {
		configuration.Debug = *debug == "true"
	}
	app := application.GetApplication(configuration)
	err = app.Run()
	log.Println(err)
}