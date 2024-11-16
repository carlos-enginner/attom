package config

import "fmt"

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func GetVersion() {
	fmt.Printf("Version: %s Commit: %s, Date: ", Version, Commit, Date)
}
