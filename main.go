package main

import (
	"fmt"
	"os"
	"src/post_relay/cmd"
)

var Version string

func main() {
	fmt.Printf("Version: %s", Version)
	if err := cmd.Execute(); err != nil {
		fmt.Println("Error executing command:", err)
		os.Exit(1)
	}
}

// o self update precisa pegar a versão que foi buildada no proprio aplicativo
// parei aqui.
