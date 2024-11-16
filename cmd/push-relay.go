package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "post_relay",
	Short: "Post Relay Application",
	Long:  "This application listens for notifications from the database and forwards them to a webhook.",
}

func init() {
	rootCmd.AddCommand(ApplicationInitCmd())
	rootCmd.AddCommand(DatabaseNotificationEnableCmd())
	// rootCmd.AddCommand(DatabaseNotificationListenCmd())
	rootCmd.AddCommand(ApplicationInstall())
	rootCmd.AddCommand(ApplicationStart())
	rootCmd.AddCommand(ApplicationSelfUpdate())
	rootCmd.AddCommand(ApplicationGetVersion())
}

// Execute executa o comando raiz
func Execute() error {
	return rootCmd.Execute()
}

// to do:
// incluir a coluna do serviço na query que retorna o evento disparado; feito
// segmentar o código da notificação em um package o NotifyDatabase;
// sepaar os commands em comandos separados;
// adicionar o cobra para ter as opções:
// - adicionar a opção de instalar o notify_database;
// - adicionar a opção para instalar o software no serviço de start do windows;
// - adicionar a opção de auto_update para checagem das atualizações;
// adicionar o viper para uso das variaveis e configuração/leitura do json; feito
// - ajustar a lógica de direcionamento;
// resolver o erro ao disparar o payload "Error: toml: cannot unmarshal object into Go struct field Association.panels.IdLocalAtendimento of type string" - feito
// criar models (para centralizar) para o struct do template de config da aplicação; feito
// ter um lugar central para escrever e ler o json do template; feito
// tem um models do payload da api externa; feito
// criar o build para windows
// criar o CI/CD no github;
// centralizar o environment, err := utils.LoadConfig() em uma constante global
