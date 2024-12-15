package utils

import (
	"fmt"
	"log"
	"regexp"
	"src/post_relay/config"
	"src/post_relay/internal/logger"
	"src/post_relay/models/environment"
	"strings"

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

func SaveConfig(cnes string, panel string, tipos string) (environment.Config, error) {
	// Configuração do Viper
	viper.SetConfigFile(config.FILE_ENVIRONMENT_APPLICATION)
	viper.SetConfigType("toml")

	logger.GetLogger().Info(cnes, panel, tipos)

	// Tenta ler o arquivo de configuração
	if err := viper.ReadInConfig(); err != nil {
		return environment.Config{}, fmt.Errorf("erro ao ler o arquivo de configuração: %v", err)
	}

	// Definir a struct de destino e fazer o mapeamento
	var config environment.Config
	if err := viper.Unmarshal(&config); err != nil {
		return environment.Config{}, fmt.Errorf("erro ao mapear as configurações para a struct: %v", err)
	}

	// divindo a string dos paneis
	panelInfo := strings.Split(panel, " - ")

	// Novo item para ser adicionado ao painel
	newPanel := map[string]interface{}{
		"cnes":        OnlyNumber(cnes),
		"description": "Novo Painel Registrado",
		"type":        []string{tipos},
		"queue": map[string]string{
			"panelUuid":  panelInfo[1],
			"sectorUuid": panelInfo[3],
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

func OnlyNumber(cnes string) string {
	re := regexp.MustCompile(`\D+`)
	cnesInfo := re.ReplaceAllString(cnes, "")

	return cnesInfo
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
