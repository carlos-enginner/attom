package utils

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"src/post_relay/config"
	"src/post_relay/models/environment"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/spf13/viper"
)

func LoadConfig() (environment.Config, error) {
	viper.SetConfigFile(config.FILE_ENVIRONMENT_APPLICATION)
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		return environment.Config{}, fmt.Errorf("erro ao ler o arquivo de configuração: %v", err)
	}

	var config environment.Config
	if err := viper.Unmarshal(&config); err != nil {
		return environment.Config{}, fmt.Errorf("erro ao mapear as configurações para a struct: %v", err)
	}

	return config, nil
}

func OnlyNumber(rawText string) string {
	re := regexp.MustCompile(`\D+`)
	newText := re.ReplaceAllString(rawText, "")
	newText = strings.TrimSpace(newText)

	return newText
}

func OnlyText(rawText string) string {
	re := regexp.MustCompile(`^\d+\s*-`)
	newText := re.ReplaceAllString(rawText, "")
	newText = strings.TrimSpace(newText)

	return newText
}

func Contains(value string, slice []string) bool {

	if reflect.TypeOf(slice).Kind() != reflect.Slice {
		fmt.Println("Erro: O argumento passado não é um slice!")
		return false
	}

	for _, item := range slice {
		if strings.EqualFold(strings.TrimSpace(item), value) {
			return true
		}
	}

	return false
}

func ToString(value float64) string {
	return fmt.Sprintf("%.0f", value)
}

func ToUpperCase(s string) string {
	return strings.ToUpper(s)
}

func ContainsWord(texto, palavra string) bool {
	re := regexp.MustCompile(`\b` + regexp.QuoteMeta(palavra) + `\b`)
	return re.MatchString(texto)
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

	return latest.GreaterThan(current)
}

func ExtractVersionFromURL(url string) (string, error) {
	re := regexp.MustCompile(`v(\d+\.\d+\.\d+)`)

	matches := re.FindStringSubmatch(url)

	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", fmt.Errorf("versão não encontrada na URL")
}
