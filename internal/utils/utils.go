package utils

import (
	"fmt"
	"src/post_relay/config"
	"src/post_relay/models/environment"

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
