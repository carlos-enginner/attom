package cmd

import (
	"fmt"
	"log"
	"os"
	"src/post_relay/models/environment"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

func newConfig() *environment.Config {
	return &environment.Config{
		Application: environment.Application{
			TimeoutConnection: 10,
		},
		API: environment.API{
			Endpoint: "https://example.com/api",
			Token:    "your-api-token",
			IBGE:     "1234567",
		},
		Database: environment.Database{
			Host:     "localhost",
			Port:     5432,
			User:     "user",
			Password: "password",
			DBName:   "example_db",
		},
		Panels: environment.Panels{
			Items: []environment.PanelItem{
				{
					Cnes:        "2382857",
					Description: "Painel Procedimentos",
					Cbos:        []string{},
					Queue: environment.Queue{
						PanelUuid: "eb6e9c6b-a196-42ed-8847-752da50bf95c", SectorUuid: "631da7f0-fe75-44fc-85c0-bafb56ab12d1",
					},
				},
				{
					Cnes:        "2382857",
					Description: "Painel Exames",
					Cbos:        []string{"2251", "2252"},
					Queue: environment.Queue{
						PanelUuid: "eb6e9c6b-a196-42ed-8847-752da50bf95c", SectorUuid: "631da7f0-fe75-44fc-85c0-bafb56ab12d1",
					},
				},
				{
					Cnes:        "2382857",
					Description: "Painel Vacinas",
					Cbos:        []string{"2251", "2252"},
					Queue: environment.Queue{
						PanelUuid: "eb6e9c6b-a196-42ed-8847-752da50bf95c", SectorUuid: "631da7f0-fe75-44fc-85c0-bafb56ab12d1",
					},
				},
				{
					Cnes:        "2382857",
					Description: "Painel Triagem",
					Cbos:        []string{"3222"},
					Queue: environment.Queue{
						PanelUuid: "eb6e9c6b-a196-42ed-8847-752da50bf95c", SectorUuid: "631da7f0-fe75-44fc-85c0-bafb56ab12d1",
					},
				},
				{
					Cnes:        "2382857",
					Description: "Painel Consultório Médico",
					Cbos:        []string{"2251", "2252", "2235", "2239", "2237"},
					Queue: environment.Queue{
						PanelUuid: "eb6e9c6b-a196-42ed-8847-752da50bf95c", SectorUuid: "631da7f0-fe75-44fc-85c0-bafb56ab12d1",
					},
				},
				{
					Cnes:        "2382857",
					Description: "Painel Consultório Odontológico",
					Cbos:        []string{"2232", "322425", "322415"},
					Queue: environment.Queue{
						PanelUuid: "eb6e9c6b-a196-42ed-8847-752da50bf95c", SectorUuid: "631da7f0-fe75-44fc-85c0-bafb56ab12d1",
					},
				},
			},
		},
	}
}

func ApplicationInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Generates a .config/environment.toml file",
		Long:  `This command asks for user input and generates a TOML configuration file with the provided settings.`,
		Run: func(cmd *cobra.Command, args []string) {

			// Criando a configuração
			config := newConfig()

			// Coletando dados da API
			fmt.Print("Enter the API endpoint: ")
			fmt.Scanln(&config.API.Endpoint)
			fmt.Print("Enter the IBGE code: ")
			fmt.Scanln(&config.API.IBGE)
			fmt.Print("Enter the token: ")
			fmt.Scanln(&config.API.Token)

			// Coletando dados do banco de dados
			fmt.Print("Enter the database host: ")
			fmt.Scanln(&config.Database.Host)
			fmt.Print("Enter the database port: ")
			fmt.Scanln(&config.Database.Port)
			fmt.Print("Enter the database user: ")
			fmt.Scanln(&config.Database.User)
			fmt.Print("Enter the database password: ")
			fmt.Scanln(&config.Database.Password)
			fmt.Print("Enter the database name: ")
			fmt.Scanln(&config.Database.DBName)

			// Pergunta se a aplicação está pronta
			fmt.Print("Is the application ready (true/false)? ")

			// Verificando e criando a pasta "config" se não existir
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

			fmt.Println("\nGerando o template environment.toml com valores padrão...")
			if err := toml.NewEncoder(file).Encode(config); err != nil {
				log.Fatalf("Erro ao escrever configuração no arquivo TOML: %v", err)
			}

			// Sucesso
			fmt.Printf("Configuration file 'config/environment.toml' generated successfully!\n")
		},
	}
}
