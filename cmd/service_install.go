package cmd

import (
	"src/post_relay/internal/win64"

	"github.com/spf13/cobra"
)

func ServiceInstall() *cobra.Command {
	return &cobra.Command{
		Use:   "install_service",
		Short: "Installs the service in Windows Services for automatic execution",
		Run: func(cmd *cobra.Command, args []string) {
			win64.NssmInstallService()
		},
	}
}
