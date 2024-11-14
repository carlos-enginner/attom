package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"src/post_relay/internal/db"
	"src/post_relay/internal/dispatch"

	"github.com/spf13/cobra"
)

func DatabaseNotificationListenCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "listen",
		Short: "Listen for notifications from the database",
		Run: func(cmd *cobra.Command, args []string) {
			conn, err := db.Connect()
			// Conectar ao banco de dados
			if err != nil {
				log.Fatal("Error connecting to database:", err)
			}
			defer conn.Close(context.Background())

			// Escutar notificações
			if err := db.ListenForNotifications(conn); err != nil {
				log.Fatal("Error listening for notifications:", err)
			}

			for {
				// Esperar por notificações
				notification, err := conn.WaitForNotification(context.Background())
				if err != nil {
					log.Fatal("Error waiting for notification:", err)
				}

				// Exibir notificação
				fmt.Printf("Received notification: %s\n", notification.Payload)

				// Parse JSON da notificação
				var notificationJSON map[string]interface{}
				if err := json.Unmarshal([]byte(notification.Payload), &notificationJSON); err != nil {
					log.Fatal("Error parsing notification payload:", err)
				}

				// Enviar para API
				payload, err := dispatch.MakePayload(notificationJSON)
				if err != nil {
					fmt.Println("Error:", err)
					return
				}

				if err := dispatch.SendMessage(payload); err != nil {
					log.Println("Error sending to API:", err)
				} else {
					fmt.Println("Notification sent to API successfully!")
				}
			}
		},
	}
}
