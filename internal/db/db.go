package db

import (
	"context"
	_ "embed"
	"fmt"
	"net/url"
	"src/post_relay/config"
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
		return nil, err
	}

	return conn, nil
}

func ListenForNotifications(conn *pgx.Conn) error {
	_, err := conn.Exec(context.Background(), "LISTEN status_change")
	if err != nil {
		return err
	}

	fmt.Println("Listening for notifications...")
	return nil
}

func EnableNotify(conn *pgx.Conn) error {
	sql := string(queryNotificationSystem)

	_, err := conn.Exec(context.Background(), sql)
	if err != nil {
		return fmt.Errorf("could not execute SQL query: %v", err)
	}

	fmt.Printf("Executed SQL file: %s\n", config.NOTIFY_STRUCTURED)
	return nil
}
