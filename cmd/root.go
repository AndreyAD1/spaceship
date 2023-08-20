package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)


var debug *bool
var logFile *string
var cpuProfile *string

var rootCmd = &cobra.Command{
	Use:   "spaceship.exe",
	Short: "A small console game",
	Long: `This is an old-fashioned console game. 
A user controls a spaceship flying through a meteor shower. 
A user's goal is to destroy meteorites and not collide with them.`,
	Run: run,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	debug = rootCmd.Flags().BoolP("debug", "d", false, "Run in a debug mode")
	logFile = rootCmd.Flags().StringP(
		"log_file", 
		"l", 
		"", 
		"write logs to this file",
	)
	cpuProfile = rootCmd.Flags().StringP(
		"cpuprofile", 
		"c", 
		"", 
		"write a cpu profile to this file",
	)
}


