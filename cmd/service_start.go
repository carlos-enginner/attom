package cmd

import (
	"src/post_relay/internal/db"
	"src/post_relay/internal/win64"

	"github.com/spf13/cobra"
)

func ServiceStart() *cobra.Command {
	return &cobra.Command{
		Use:   "start_service",
		Short: "Starts the application on Windows Service",
		Run: func(cmd *cobra.Command, args []string) {
			go db.StartNotifications()
			go win64.NssmStartService()

			select {}
		},
	}
}
