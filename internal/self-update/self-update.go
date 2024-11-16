package selfupdate

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// Função para substituir o binário atual pelo novo
func ReplaceExecutable(newBin string) error {
	// Obter o caminho do binário atual
	currentPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("erro ao obter o caminho do executável atual: %w", err)
	}

	// Salvar o nome do binário antigo para backup (opcional)
	backupPath := currentPath + ".bak"

	// Renomear o binário atual para um backup
	err = os.Rename(currentPath, backupPath)
	if err != nil {
		return fmt.Errorf("erro ao renomear o binário atual: %w", err)
	}

	// Renomear o novo binário para o nome do binário atual
	err = os.Rename(newBin, currentPath)
	if err != nil {
		// Se falhar, restaura o binário antigo
		os.Rename(backupPath, currentPath)
		return fmt.Errorf("erro ao substituir o binário atual: %w", err)
	}

	// Excluir o backup, se necessário
	os.Remove(backupPath)

	// Garantir permissões executáveis no novo binário
	err = os.Chmod(currentPath, 0755)
	if err != nil {
		return fmt.Errorf("erro ao definir permissões no novo binário: %w", err)
	}

	return nil
}

// Função para baixar a última release de um repositório GitHub
func DownloadLatestRelease(owner, repo, token, targetPath string) error {
	// Montar a URL da API do GitHub para obter a release
	assetURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)

	// Criar requisição para obter os dados da release
	req, err := http.NewRequest("GET", assetURL, nil)
	if err != nil {
		return fmt.Errorf("erro ao criar requisição para o GitHub: %w", err)
	}
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	// Fazer a requisição
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao fazer a requisição para o GitHub: %w", err)
	}
	defer resp.Body.Close()

	// Verificar se a requisição foi bem-sucedida
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("erro ao obter os dados da release: %v", string(body))
	}

	// Estrutura para armazenar os dados da release
	var releaseData struct {
		Assets []struct {
			Id                 uint64 `json:"id"`
			Url                string `json:"url"`
			BrowserDownloadURL string `json:"browser_download_url"`
			Name               string `json:"name"`
		} `json:"assets"`
	}

	// Parse da resposta JSON
	err = json.NewDecoder(resp.Body).Decode(&releaseData)
	if err != nil {
		return fmt.Errorf("erro ao parsear a resposta JSON: %w", err)
	}

	// Encontrar o binário correto na lista de assets
	var releaseID uint64
	for _, asset := range releaseData.Assets {
		// Você pode adaptar isso para o nome do binário que você precisa
		if asset.Name == "attom_1.0.0_linux_amd64" { // Ajuste para o nome do seu binário
			releaseID = asset.Id
			break
		}
	}

	// Se não encontrou o binário, retornar erro
	if releaseID == 0 {
		return fmt.Errorf("binário adequado não encontrado nas assets da release")
	}

	// Fazer o download do binário
	return downloadFile(releaseID, targetPath, owner, repo, token)
}

// Função para baixar o binário de uma URL
func downloadFile(releaseID uint64, filePath, owner string, repo string, token string) error {

	if releaseID == 0 {
		log.Fatalf("Binário adequado não encontrado nas assets da release.")
	}

	downloadURL := fmt.Sprintf("https://%s:@api.github.com/repos/%s/%s/releases/assets/%d", token, owner, repo, releaseID)

	// Criar o arquivo para salvar o binário
	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("erro ao criar o arquivo: %w", err)
	}
	defer out.Close()

	// Criar requisição para download
	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		return fmt.Errorf("erro ao criar requisição de download: %w", err)
	}
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/octet-stream")

	// Fazer a requisição
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao fazer a requisição de download: %w", err)
	}
	defer resp.Body.Close()

	// Salvar o conteúdo do corpo no arquivo
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("erro ao salvar o arquivo: %w", err)
	}

	return nil
}
