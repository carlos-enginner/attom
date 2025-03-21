package cmd

import (
	"context"
	"fmt"
	"log"
	"src/post_relay/internal/db"

	"github.com/spf13/cobra"
)

func DatabaseNotificationEnableCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "prepare_database",
		Short: "Prepares the database for asynchronous notifications.",
		Run: func(cmd *cobra.Command, args []string) {
			conn, err := db.Connect()
			// Conectar ao banco de dados
			if err != nil {
				log.Fatal("Error connecting to database: DatabaseNotificationEnableCmd", err)
			}
			defer conn.Close(context.Background())

			// Habilitar notificações
			if err := db.EnableNotify(conn); err != nil {
				log.Fatal("Error enabling notifications:", err)
			}

			fmt.Println("Notifications enabled successfully.")
		},
	}
}
