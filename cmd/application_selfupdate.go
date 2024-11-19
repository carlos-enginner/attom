package cmd

import (
	selfupdate "src/post_relay/internal/self-update"

	"github.com/spf13/cobra"
)

func ApplicationSelfUpdate() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Serviço de atualização da aplicação",
		Run: func(cmd *cobra.Command, args []string) {
			selfupdate.CheckAndUpdateVersion()
		},
	}
}
