package cmd

import (
	"src/post_relay/internal/win64"

	"github.com/spf13/cobra"
)

func ServiceRemove() *cobra.Command {
	return &cobra.Command{
		Use:   "service_remove",
		Short: "Remove the application on Windows Service",
		Run: func(cmd *cobra.Command, args []string) {
			win64.NssmRemoveService()
		},
	}
}
