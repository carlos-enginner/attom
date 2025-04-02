package cmd

import (
	"context"
	"fmt"
	"src/post_relay/internal/db"
	"src/post_relay/internal/logger"

	"github.com/spf13/cobra"
)

func DatabaseNotificationEnableCmd() *cobra.Command {
	log := logger.GetLogger()

	return &cobra.Command{
		Use:   "prepare_database",
		Short: "Prepares the database for asynchronous notifications.",
		Run: func(cmd *cobra.Command, args []string) {
			conn, err := db.Connect()
			// Conectar ao banco de dados
			if err != nil {
				log.Infof("Error connecting to database:", err)
			}
			defer conn.Close(context.Background())

			// Habilitar notificações
			if err := db.EnableNotify(conn); err != nil {
				log.Infof("Error enabling notifications:", err)
			}

			fmt.Println("Notifications enabled successfully.")
		},
	}
}
