package cmd

import (
	"fmt"
	"log"
	"os"
	"src/post_relay/models/environment"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

func newConfig() *environment.Config {
	return &environment.Config{
		Application: environment.Application{
			TimeoutConnection: 10 * time.Second,
			HttpDebug:         false,
		},
		API: environment.API{
			Endpoint: "http://painel.icsgo.com.br:7001/ws/v1",
			Token:    "58a846d58a0670ea3c3ad54ec5069130",
			IBGE:     "1234567",
		},
		Database: environment.Database{
			Host:     "localhost",
			Port:     5433,
			User:     "postgres",
			Password: "password",
			DBName:   "esus",
		},
		Panels: environment.Panels{
			Items: []environment.PanelItem{},
		},
	}
}

func ApplicationInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Generates a .config/environment.toml file",
		Long:  `This command asks for user input and generates a TOML configuration file with the provided settings.`,
		Run: func(cmd *cobra.Command, args []string) {

			config := newConfig()

			var input string
			fmt.Printf("Enter the IBGE code (default: %s): ", config.API.IBGE)
			fmt.Scanln(&input)
			if strings.TrimSpace(input) != "" {
				config.API.IBGE = input
			}

			input = ""
			fmt.Printf("Enter the database user (default: %s): ", config.Database.User)
			fmt.Scanln(&input)
			if strings.TrimSpace(input) != "" {
				config.Database.User = input
			}

			input = ""
			fmt.Printf("Enter the database password (default: %s): ", config.Database.Password)
			fmt.Scanln(&input)
			if strings.TrimSpace(input) != "" {
				config.Database.Password = input
			}

			input = ""
			fmt.Printf("Enter the database port (default: %d): ", config.Database.Port)
			fmt.Scanln(&input)
			if strings.TrimSpace(input) != "" {
				port, err := strconv.Atoi(input)
				if err != nil {
					log.Printf("Invalid input for port, using default value %d\n", config.Database.Port)
				} else {
					config.Database.Port = port
				}
			}

			if _, err := os.Stat("config"); os.IsNotExist(err) {
				err := os.Mkdir("config", 0755)
				if err != nil {
					log.Fatalf("Error creating config folder: %v", err)
				}
			}

			// Criar o arquivo environment.toml
			file, err := os.Create("config/environment.toml")
			if err != nil {
				log.Fatalf("Erro ao criar o arquivo TOML: %v", err)
			}
			defer file.Close()

			if err := toml.NewEncoder(file).Encode(config); err != nil {
				log.Fatalf("Erro ao escrever configuração no arquivo TOML: %v", err)
			}

			// Sucesso
			fmt.Printf("Configuration file 'config/environment.toml' generated successfully!\n")
		},
	}
}
