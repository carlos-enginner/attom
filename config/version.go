package config

import "fmt"

var (
	Version string
)

func LogVersion() {
	fmt.Println("version=", Version)
}
