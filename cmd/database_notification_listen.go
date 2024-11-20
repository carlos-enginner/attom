package cmd

import (
	"src/post_relay/internal/db"

	"github.com/spf13/cobra"
)

func DatabaseNotificationListenCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Application start",
		Run: func(cmd *cobra.Command, args []string) {
			go db.StartNotifications()
			select {}
		},
	}
}
