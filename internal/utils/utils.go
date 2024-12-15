package utils

import (
	"fmt"
	"log"
	"regexp"
	"src/post_relay/config"
	"src/post_relay/models/environment"

	"github.com/Masterminds/semver"
	"github.com/spf13/viper"
)

func LoadConfig() (environment.Config, error) {
	// Configuração do Viper
	viper.SetConfigFile(config.FILE_ENVIRONMENT_APPLICATION)
	viper.SetConfigType("toml")

	// Tenta ler o arquivo de configuração
	if err := viper.ReadInConfig(); err != nil {
		return environment.Config{}, fmt.Errorf("erro ao ler o arquivo de configuração: %v", err)
	}

	// Definir a struct de destino e fazer o mapeamento
	var config environment.Config
	if err := viper.Unmarshal(&config); err != nil {
		return environment.Config{}, fmt.Errorf("erro ao mapear as configurações para a struct: %v", err)
	}

	return config, nil
}

func SaveConfig() (environment.Config, error) {
	// Configuração do Viper
	viper.SetConfigFile(config.FILE_ENVIRONMENT_APPLICATION)
	viper.SetConfigType("toml")

	// Tenta ler o arquivo de configuração
	if err := viper.ReadInConfig(); err != nil {
		return environment.Config{}, fmt.Errorf("erro ao ler o arquivo de configuração: %v", err)
	}

	// Definir a struct de destino e fazer o mapeamento
	var config environment.Config
	if err := viper.Unmarshal(&config); err != nil {
		return environment.Config{}, fmt.Errorf("erro ao mapear as configurações para a struct: %v", err)
	}

	// Novo item para ser adicionado ao painel
	newPanel := map[string]interface{}{
		"cnes":        "2382857",
		"description": "Novo Painel de Teste",
		"type": []string{
			"CONSULTA", "ESCUTA INICIAL", "CONSULTA ODONTOLÓGICA",
			"AVALIAÇÃO DE ELEGIBILIDADE", "ATENÇÃO DOMICILIAR",
			"ATENDIMENTO DE PROCEDIMENTO", "PRÉ-NATAL", "PUERPÉRIO",
			"PUERICULTURA", "VACINAÇÃO", "ZIKA / MICROCEFALIA", "OBSERVAÇÃO",
		},
		"queue": map[string]string{
			"panelUuid":  "14c3c97e-805d-4128-b3cd-84159eac25f3_test",
			"sectorUuid": "36efe62b-231c-421d-b446-3106908f97a7_test",
		},
	}

	// Adicionando o novo item ao campo panels.items
	existingPanels := viper.Get("panels.items").([]interface{})
	viper.Set("panels.items", append(existingPanels, newPanel))

	// // Salvando a configuração modificada
	if err := viper.WriteConfig(); err != nil {
		log.Fatalf("Erro ao salvar a configuração: %v", err)
	}

	return config, nil
}

func Contains(value string, slice []string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func Substr(s string, start, end int) string {
	if start < 0 {
		start = 0
	}
	if end > len(s) {
		end = len(s)
	}
	if start > end {
		return ""
	}
	return s[start:end]
}

func ToString(value float64) string {
	return fmt.Sprintf("%.0f", value)
}

func VersionIsGreaterThan(latestVersion string) bool {
	currentVersion := config.Version
	current, err := semver.NewVersion(currentVersion)
	if err != nil {
		log.Fatalf("Erro ao parse da versão atual: %v", err)
	}

	latest, err := semver.NewVersion(latestVersion)
	if err != nil {
		log.Fatalf("Erro ao parse da versão mais recente: %v", err)
	}

	// Compara as versões
	return latest.GreaterThan(current)
}

func ExtractVersionFromURL(url string) (string, error) {
	// Expressão regular para encontrar o padrão de versão (vX.Y.Z)
	re := regexp.MustCompile(`v(\d+\.\d+\.\d+)`)

	// Tenta encontrar a versão na URL
	matches := re.FindStringSubmatch(url)

	// Se encontrar a versão, retorna ela, caso contrário retorna erro
	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", fmt.Errorf("versão não encontrada na URL")
}
