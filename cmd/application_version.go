package cmd

import (
	"github.com/spf13/cobra"
)

func ApplicationGetVersion() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version application",
		Run: func(cmd *cobra.Command, args []string) {
			ApplicationGetVersion().Version
		},
	}
}
