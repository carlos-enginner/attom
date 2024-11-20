package cmd

import (
	"fmt"
	"log"
	"src/post_relay/internal/win64"

	"github.com/spf13/cobra"
)

func ServiceInstall() *cobra.Command {
	return &cobra.Command{
		Use:   "install_service",
		Short: "Instala o serviço no Windows para execução automática",
		Run: func(cmd *cobra.Command, args []string) {
			// Cria o serviço
			svc, err := win64.NewService()
			if err != nil {
				log.Fatal("Erro ao criar o serviço:", err)
			}

			// Instala o serviço no Windows
			err = svc.Install()
			if err != nil {
				log.Fatalf("Erro ao instalar o serviço: %v", err)
			}

			fmt.Println("Serviço instalado com sucesso!")
		},
	}
}
