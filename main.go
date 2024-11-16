package main

import (
	"fmt"
	"os"
	"src/post_relay/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println("Error executing command:", err)
		os.Exit(1)
	}
}

// ////////////////////////////
// // SELFUPDATE USANDO GO-UPDATE
// ///////////////////////////

// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"os"
// )

// func replaceExecutable(newBin string) error {
// 	// Obter o caminho do binário atual
// 	currentPath, err := os.Executable()
// 	if err != nil {
// 		return fmt.Errorf("erro ao obter o caminho do executável atual: %w", err)
// 	}

// 	// Salvar o nome do binário antigo para backup (opcional)
// 	backupPath := currentPath + ".bak"

// 	// Renomear o binário atual para um backup
// 	err = os.Rename(currentPath, backupPath)
// 	if err != nil {
// 		return fmt.Errorf("erro ao renomear o binário atual: %w", err)
// 	}

// 	// Renomear o novo binário para o nome do binário atual
// 	err = os.Rename(newBin, currentPath)
// 	if err != nil {
// 		// Se falhar, restaura o binário antigo
// 		os.Rename(backupPath, currentPath)
// 		return fmt.Errorf("erro ao substituir o binário atual: %w", err)
// 	}

// 	// Excluir o backup, se necessário
// 	os.Remove(backupPath)

// 	// Garantir permissões executáveis no novo binário
// 	err = os.Chmod(currentPath, 0755)
// 	if err != nil {
// 		return fmt.Errorf("erro ao definir permissões no novo binário: %w", err)
// 	}

// 	return nil
// }

// func main() {
// 	// Defina o repositório e a versão desejada
// 	owner := "carlos-enginner"
// 	repo := "attom"
// 	tag := "latest" // Substitua pela versão desejada

// 	// Substitua com o seu token de acesso pessoal
// 	token := "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"

// 	// URL para o arquivo binário da release (modifique para corresponder ao seu arquivo)
// 	// Aqui é necessário usar a API do GitHub para encontrar o arquivo binário
// 	assetURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/%s", owner, repo, tag)

// 	// Fazendo a requisição GET para obter os detalhes da release
// 	req, err := http.NewRequest("GET", assetURL, nil)
// 	if err != nil {
// 		log.Fatalf("Erro ao criar requisição: %v", err)
// 	}

// 	// Cabeçalho de autenticação
// 	req.Header.Set("Authorization", "token "+token)
// 	req.Header.Set("Accept", "application/vnd.github.v3+json")

// 	// Fazendo a requisição
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		log.Fatalf("Erro ao fazer a requisição para o GitHub: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	// Lendo a resposta
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatalf("Erro ao ler a resposta: %v", err)
// 	}

// 	// Verificar se a requisição foi bem-sucedida
// 	if resp.StatusCode != http.StatusOK {
// 		log.Fatalf("Erro ao obter os dados da release: %v", string(body))
// 	}

// 	// A resposta da API do GitHub inclui os assets. Precisamos extrair o link para o binário.
// 	// Aqui estamos assumindo que a primeira asset encontrada é o binário que precisamos.
// 	var releaseData struct {
// 		Assets []struct {
// 			Id                 uint64 `json:"id"`
// 			Url                string `json:"url"`
// 			BrowserDownloadURL string `json:"browser_download_url"`
// 			Name               string `json:"name"`
// 		} `json:"assets"`
// 	}

// 	// Parse do JSON da resposta
// 	err = json.Unmarshal(body, &releaseData)
// 	if err != nil {
// 		log.Fatalf("Erro ao parsear a resposta JSON: %v", err)
// 	}

// 	const linuxName = "attom_1.0.0_linux_amd64"

// 	// Encontrar o binário correto para o sistema operacional
// 	var releaseID uint64
// 	for _, asset := range releaseData.Assets {
// 		if asset.Name == linuxName {
// 			releaseID = asset.Id
// 			break
// 		}
// 	}

// 	// Se não encontrou o binário correto
// 	if releaseID == 0 {
// 		log.Fatalf("Binário adequado não encontrado nas assets da release.")
// 	}

// 	downloadURL := fmt.Sprintf("https://%s:@api.github.com/repos/%s/%s/releases/assets/%d", token, owner, repo, releaseID)

// 	// Baixando o binário mais recente

// 	fmt.Println(downloadURL)

// 	req2, err2 := http.NewRequest("GET", downloadURL, nil)
// 	if err2 != nil {
// 		log.Fatalf("erro ao criar requisição de download: %v", err2)
// 	}
// 	req2.Header.Set("Accept", "application/octet-stream")

// 	// Enviar a requisição para obter o arquivo
// 	resp2, err2 := client.Do(req2)
// 	if err2 != nil {
// 		log.Fatalf("erro ao baixar o arquivo: %v", err)
// 	}
// 	defer resp2.Body.Close()

// 	// salvando o novo arquivo

// 	// Criar um arquivo para o binário
// 	out, err := os.Create("./tmp/push-relay.new")
// 	if err != nil {
// 		log.Fatalf("erro ao criar o arquivo: %v", err)
// 	}
// 	defer out.Close()

// 	// Copiar os dados do corpo da resposta para o arquivo
// 	_, err = io.Copy(out, resp2.Body)
// 	if err != nil {
// 		log.Fatalf("erro ao salvar o arquivo: %v", err)
// 	}

// 	fmt.Println("Baixado com sucesso!")

// 	// Aplicando o self-update
// 	// Substituir o binário atual pelo novo
// 	const newBinName = "./tmp/push-relay.new"
// 	err = replaceExecutable(newBinName)
// 	if err != nil {
// 		log.Fatalf("Erro ao substituir o binário: %v", err)
// 	}

// 	// err = update.Apply(resp2.Body, update.Options{
// 	// 	TargetPath:  "./tmp/push-relay",
// 	// 	OldSavePath: "./tmp/push-relay.old",
// 	// })
// 	// if err != nil {
// 	// 	log.Fatalf("Erro ao aplicar o update: %v", err)
// 	// }

// 	fmt.Println("Self-update concluído com sucesso!")
// }
