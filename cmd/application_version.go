package cmd

import (
	"fmt"
	"src/post_relay/config"

	"github.com/spf13/cobra"
)

func ApplicationGetVersion() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version application",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("\nVersion: %s\nCommit: %s\n\n", config.Version, config.Commit)
		},
	}
}
