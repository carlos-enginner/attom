package cmd

import (
	"src/post_relay/internal/win64"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"
)

var logger service.Logger

func ServiceStart() *cobra.Command {
	return &cobra.Command{
		Use:   "start_service",
		Short: "Starts the application on Windows Service",
		Run: func(cmd *cobra.Command, args []string) {
			win64.NssmStartService()
		},
	}
}
