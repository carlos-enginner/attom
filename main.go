package main

import (
	"fmt"
	"os"
	"src/post_relay/cmd"
	"src/post_relay/config"
	"src/post_relay/internal/logger"

	figure "github.com/common-nighthawk/go-figure"
)

func printAsciiHeader() {
	art := figure.NewFigure("ATTOM", "", true)
	art.Print()

	fmt.Printf("This application listens for notifications from the database and forwards them to API.\n\n")
}

func main() {
	printAsciiHeader()
	log := logger.GetLogger()
	log.Info("Application started")
	log.Infof("Version: %s Commit: %s \n", config.Version, config.Commit)
	if err := cmd.Execute(); err != nil {
		fmt.Println("Error executing command:", err)
		os.Exit(1)
	}
}
