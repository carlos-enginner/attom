package db

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"src/post_relay/config"
	dispatchpanel "src/post_relay/internal/dispatch-panel"
	"src/post_relay/internal/logger"
	"src/post_relay/internal/utils"

	"github.com/jackc/pgx/v4"
)

//go:embed migrations/notification_system.sql
var queryNotificationSystem string

// Connect estabelece uma conexão com o banco de dados PostgreSQL
func Connect() (*pgx.Conn, error) {

	config, err := utils.LoadConfig()
	if err != nil {
		return nil, err
	}

	// Encode a senha
	encodedPassword := url.QueryEscape(config.Database.Password)

	// Crie a string de conexão
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", config.Database.User, encodedPassword, config.Database.Host, config.Database.Port, config.Database.DBName)

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		logger.GetLogger().Errorf("connection error: %v", err)
		return nil, err
	}

	return conn, nil
}

func ListenForNotifications(conn *pgx.Conn) error {
	_, err := conn.Exec(context.Background(), "LISTEN call_record")
	if err != nil {
		return err
	}

	logger.GetLogger().Info("Listening for notifications...")
	fmt.Println("Listening for notifications...")
	return nil
}

func EnableNotify(conn *pgx.Conn) error {
	sql := string(queryNotificationSystem)

	_, err := conn.Exec(context.Background(), sql)
	if err != nil {
		logger.GetLogger().Errorf("could not execute SQL query: %v", err)
		return fmt.Errorf("could not execute SQL query: %v", err)
	}

	logger.GetLogger().Infof("Executed SQL file: %s\n", config.NOTIFY_STRUCTURED)
	fmt.Printf("Executed SQL file: %s\n", config.NOTIFY_STRUCTURED)
	return nil
}

func StartNotifications() {
	conn, err := Connect()
	// Conectar ao banco de dados
	if err != nil {
		log.Fatal("Error connecting to database: StartNotifications", err)
	}
	defer conn.Close(context.Background())

	// Escutar notificações
	if err := ListenForNotifications(conn); err != nil {
		log.Fatal("Error listening for notifications:", err)
	}

	log := logger.GetLogger()

	for {
		// Esperar por notificações
		notification, err := conn.WaitForNotification(context.Background())
		if err != nil {
			log.Fatal("Error waiting for notification:", err)
		}

		if json.Valid([]byte(notification.Payload)) {

			// Enviar para API
			payload, err := dispatchpanel.MakePayload(notification.Payload)
			if err != nil {
				log.Errorln("Error:", err)
				return
			}

			if err := dispatchpanel.SendMessage(payload); err != nil {
				log.Errorln("Error sending to API:", err)
			} else {
				log.Info("Notification sent to API successfully!")
			}
		} else {
			log.Warning("Query without returning records. Check params date your query")
		}
	}
}
