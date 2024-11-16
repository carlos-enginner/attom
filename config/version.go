package config

import "fmt"

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func LogVersion() {
	fmt.Printf("Attom %s, commit %s, built at %s", version, commit, date)
}
