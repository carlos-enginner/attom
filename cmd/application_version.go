package cmd

import (
	"src/post_relay/config"

	"github.com/spf13/cobra"
)

func ApplicationGetVersion() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version application",
		Run: func(cmd *cobra.Command, args []string) {
			config.LogVersion()
		},
	}
}
