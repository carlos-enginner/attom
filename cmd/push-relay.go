package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "attom",
}

func init() {
	rootCmd.AddCommand(ApplicationInitCmd())
	rootCmd.AddCommand(DatabaseNotificationEnableCmd())
	rootCmd.AddCommand(DatabaseNotificationListenCmd())
	rootCmd.AddCommand(ServiceInstall())
	rootCmd.AddCommand(ServiceStart())
	rootCmd.AddCommand(ApplicationSelfUpdate())
	rootCmd.AddCommand(ApplicationGetVersion())
}

// Execute executa o comando raiz
func Execute() error {
	return rootCmd.Execute()
}
