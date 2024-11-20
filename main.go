package main

import (
	"fmt"
	"os"
	"src/post_relay/cmd"
	"src/post_relay/config"
	"src/post_relay/internal/logger"
)

func main() {
	log := logger.GetLogger()
	log.Infof("Version: %s Commit: %s \n", config.Version, config.Commit)
	log.Info("Application started")
	if err := cmd.Execute(); err != nil {
		fmt.Println("Error executing command:", err)
		os.Exit(1)
	}
}
