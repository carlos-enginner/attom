package main

import (
	"fmt"
	"os"
	"src/post_relay/cmd"
)

var Version, Commit, Date string

func main() {
	fmt.Printf("Version: %s Commit: %s, Date: ", Version, Commit, Date)
	if err := cmd.Execute(); err != nil {
		fmt.Println("Error executing command:", err)
		os.Exit(1)
	}
}
