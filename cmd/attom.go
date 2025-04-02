package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "attom",
}

func init() {
	rootCmd.AddCommand(ApplicationInitCmd())
	rootCmd.AddCommand(ServiceRemove())
	rootCmd.AddCommand(ServiceInstall())
	rootCmd.AddCommand(DatabaseNotificationEnableCmd())
	rootCmd.AddCommand(DatabaseNotificationListenCmd())
	rootCmd.AddCommand(ApplicationGetVersion())
	rootCmd.AddCommand(PanelNewRegister())
}

// Execute executa o comando raiz
func Execute() error {
	return rootCmd.Execute()
}
