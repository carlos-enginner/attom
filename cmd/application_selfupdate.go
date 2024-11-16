package cmd

import (
	"fmt"
	"log"
	selfupdate "src/post_relay/internal/self-update"

	"github.com/spf13/cobra"
)

func ApplicationSelfUpdate() *cobra.Command {
	return &cobra.Command{
		Use:   "self-update",
		Short: "Auto atualização da aplicação",
		Run: func(cmd *cobra.Command, args []string) {
			owner := "carlos-enginner"
			repo := "attom"
			token := "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" // Substitua pelo seu token de acesso pessoal
			targetPath := "./push-relay.new"                                                                         // Caminho onde o novo binário será salvo

			// Baixar a última release
			err := selfupdate.DownloadLatestRelease(owner, repo, token, targetPath)
			if err != nil {
				log.Fatalf("Erro ao baixar a release: %v", err)
			}
			fmt.Println("Binário baixado com sucesso!")

			// Substituir o binário atual pelo novo
			err = selfupdate.ReplaceExecutable(targetPath)
			if err != nil {
				log.Fatalf("Erro ao substituir o binário: %v", err)
			}

			fmt.Println("Self-update concluído com sucesso!")
		},
	}
}
