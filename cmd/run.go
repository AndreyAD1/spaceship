package cmd

import (
	"fmt"
	"log"
	"os"
	"runtime"
	runtimeDebug "runtime/debug"
	"runtime/pprof"

	"github.com/AndreyAD1/spaceship/internal/application"
	"github.com/AndreyAD1/spaceship/internal/config"
	"github.com/AndreyAD1/spaceship/internal/logger"
	"github.com/caarlos0/env/v7"
	"github.com/spf13/cobra"
)

func run(cmd *cobra.Command, args []string) {
	configuration := config.StartupConfig{}
	err := env.Parse(&configuration)
	if err != nil {
		log.Fatal(err)
	}
	if *debug {
		configuration.Debug = true
	}
	if *logFile != "" {
		configuration.LogFile = *logFile
	}
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	newLogger, err := logger.NewLogger(configuration)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if p := recover(); p != nil {
			newLogger.Errorf("Internal error: %v", p)
			stackTrace := runtimeDebug.Stack()
			newLogger.Errorf("Error stack trace: %s", stackTrace)
			fmt.Printf("Critical error: %s", stackTrace)
		}
		newLogger.Debugf("final goroutine number: %v", runtime.NumGoroutine())
	}()
	app := application.NewApplication(newLogger)
	newLogger.Debug("run application")
	err = app.Run()
	if err != nil {
		newLogger.Error(err)
	}
}
