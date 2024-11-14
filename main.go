package main

import (
	"fmt"
	"os"
	"src/post_relay/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println("Error executing command:", err)
		os.Exit(1)
	}
}
