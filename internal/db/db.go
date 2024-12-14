package db

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"src/post_relay/config"
	"src/post_relay/internal/dispatch"
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

type Unidade struct {
	NuCnes      string
	NomeUnidade string
}

type Tipo struct {
	Codigo    int64
	Descricao string
}

func GetUnidades(conn *pgx.Conn) ([]Unidade, error) {
	// Definindo a consulta SQL para buscar as unidades de saúde
	sql := `SELECT tus.nu_cnes, 
		UPPER(tus.no_unidade_saude_filtro) AS nome_unidade
		FROM tb_unidade_saude tus
		ORDER BY tus.no_unidade_saude;`

	// Executando a consulta e obtendo os resultados
	rows, err := conn.Query(context.Background(), sql)
	if err != nil {
		// Logando o erro caso a consulta falhe
		logger.GetLogger().Errorf("could not execute SQL query: %v", err)
		return nil, fmt.Errorf("could not execute SQL query: %v", err)
	}
	defer rows.Close()

	// Criando um slice para armazenar as unidades
	var unidades []Unidade

	// Iterando pelas linhas retornadas e populando o slice
	for rows.Next() {
		var unidade Unidade
		err := rows.Scan(&unidade.NuCnes, &unidade.NomeUnidade)
		if err != nil {
			// Logando erro de scan
			logger.GetLogger().Errorf("could not scan row: %v", err)
			return nil, fmt.Errorf("could not scan row: %v", err)
		}
		// Adicionando a unidade ao slice
		unidades = append(unidades, unidade)
	}

	// Verificando se houve erro durante a iteração das linhas
	if err := rows.Err(); err != nil {
		// Logando o erro de iteração
		logger.GetLogger().Errorf("error during rows iteration: %v", err)
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}

	return unidades, nil
}

func GetTipos(conn *pgx.Conn) ([]Tipo, error) {

	// Definindo a consulta SQL para buscar as unidades de saúde
	sql := `select
	t.co_tipo_atend_prof codigo,
	t.no_tipo_atend_prof descricao
from
	tb_tipo_atend_prof t
order by
	co_tipo_atend_prof`

	// Executando a consulta e obtendo os resultados
	rows, err := conn.Query(context.Background(), sql)
	if err != nil {
		// Logando o erro caso a consulta falhe
		logger.GetLogger().Errorf("could not execute SQL query: %v", err)
		return nil, fmt.Errorf("could not execute SQL query: %v", err)
	}
	defer rows.Close()

	// Criando um slice para armazenar as unidades
	var tipos []Tipo

	// Iterando pelas linhas retornadas e populando o slice
	for rows.Next() {
		var tipo Tipo
		err := rows.Scan(&tipo.Codigo, &tipo.Descricao)
		if err != nil {
			// Logando erro de scan
			logger.GetLogger().Errorf("could not scan row: %v", err)
			return nil, fmt.Errorf("could not scan row: %v", err)
		}
		// Adicionando a unidade ao slice
		tipos = append(tipos, tipo)
	}

	// Verificando se houve erro durante a iteração das linhas
	if err := rows.Err(); err != nil {
		// Logando o erro de iteração
		logger.GetLogger().Errorf("error during rows iteration: %v", err)
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}

	return tipos, nil
}

func StartNotifications() {
	conn, err := Connect()
	// Conectar ao banco de dados
	if err != nil {
		log.Fatal("Error connecting to database:", err)
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
			payload, err := dispatch.MakePayload(notification.Payload)
			if err != nil {
				log.Errorln("Error:", err)
				return
			}

			if err := dispatch.SendMessage(payload); err != nil {
				log.Errorln("Error sending to API:", err)
			} else {
				log.Info("Notification sent to API successfully!")
			}
		} else {
			log.Warning("Query without returning records. Check params date your query")
		}
	}
}
