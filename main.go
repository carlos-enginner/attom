package main

import (
	"fmt"
	"os"
	"src/post_relay/cmd"
	"src/post_relay/config"
)

// var Version, Commit, Date string

func main() {
	fmt.Printf("Version: %s Commit: %s, Date: ", config.Version, config.Commit)
	if err := cmd.Execute(); err != nil {
		fmt.Println("Error executing command:", err)
		os.Exit(1)
	}
}
