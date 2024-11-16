package config

import "fmt"

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func GetVersion() {
	fmt.Printf("Attom %s, commit %s, built at %s", Version, Commit, Date)

	fmt.Printf("Version: %s Commit: %s, Date: ", Version, Commit, Date)
}
