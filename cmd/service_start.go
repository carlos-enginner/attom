package cmd

import (
	"fmt"
	"log"
	"src/post_relay/internal/win64"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"
)

var logger service.Logger

func ServiceStart() *cobra.Command {
	return &cobra.Command{
		Use:   "start_service",
		Short: "Inicia o serviço no Windows",
		Run: func(cmd *cobra.Command, args []string) {
			svc, err := win64.NewService()
			if err != nil {
				log.Fatal("Erro ao criar o serviço:", err)
			}

			logger, err = svc.Logger(nil)
			if err != nil {
				log.Fatal(err)
			}

			err = svc.Run()
			if err != nil {
				log.Fatalf("Erro ao iniciar o serviço: %v", err)
			}

			fmt.Println("Serviço iniciado com sucesso!")
		},
	}
}
