package win64

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	// Importando a funcionalidade de embed
	_ "embed"
)

//go:embed assets/nssm.exe
var nssmData []byte

const NSSM_EXECUTABLE_TITLE = "nssm.exe"

func NssmExtractApp() {
	// Obter o diretório atual de onde o programa está sendo executado
	execDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Erro ao obter o diretório atual: %s", err)
		return
	}

	// Caminho para a pasta "nssm" no diretório de execução
	nssmDir := filepath.Join(execDir, ".nssm")

	// Caminho completo para salvar o nssm.exe
	nssmPath := filepath.Join(nssmDir, NSSM_EXECUTABLE_TITLE)

	if _, err := os.Stat(nssmPath); err == nil {
		fmt.Println("Arquivo 'nssm.exe' já existe em:", nssmPath)
		return
	}

	// Criar a pasta "nssm" se ela não existir
	err = os.MkdirAll(nssmDir, 0755)
	if err != nil {
		fmt.Println("Erro ao criar a pasta 'nssm':", err)
		return
	}

	cmd := exec.Command("attrib", "+h", nssmDir)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Erro ao tornar a pasta '.nssm' oculta:", err)
		return
	}

	fmt.Println("Pasta 'nssm' criada com sucesso em:", nssmDir)

	// Criar a pasta oculta ".nssm"
	err = os.MkdirAll(nssmDir, 0755)
	if err != nil {
		fmt.Println("Erro ao criar a pasta '.nssm':", err)
		return
	}

	// Salvar o conteúdo embutido do nssm.exe na pasta oculta
	err = os.WriteFile(nssmPath, nssmData, 0755)
	if err != nil {
		fmt.Println("Erro ao escrever o arquivo nssm.exe:", err)
		return
	}

	fmt.Println("Arquivo nssm.exe extraído com sucesso para:", nssmPath)
}

func NssmInstallService() {

	execDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Erro ao obter o diretório atual:", err)
		return
	}

	nssmPath := filepath.Join(execDir+"\\.nssm", NSSM_EXECUTABLE_TITLE)

	// Usar o nssm.exe para criar o serviço no Windows
	// Exemplo de comando para criar um serviço (ajuste conforme necessário)
	serviceName := "AttomSvc"
	executablePath := execDir + "\\attom.exe"

	// Comando nssm.exe para criar o serviço
	cmdCreateService := exec.Command(nssmPath, "install", serviceName, executablePath)

	// Executar o comando
	err = cmdCreateService.Run()
	if err != nil {
		fmt.Println("Erro ao criar o serviço:", err)
		return
	}

	// Comando para adicionar a descrição ao serviço
	cmdSetDescription := exec.Command(nssmPath, "set", serviceName, "Description", "O serviço responsável por detectar e capturar eventos de atendimento no e-sus/PEC e envia-lós a um serviço externo de painel eletrônico")

	// Executar o comando para adicionar a descrição
	err = cmdSetDescription.Run()
	if err != nil {
		fmt.Println("Erro ao definir a descrição do serviço:", err)
		return
	}

	fmt.Println("Serviço criado com sucesso:", serviceName)
}

func NssmRemoveService() {
	serviceName := "AttomSvc"

	execDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Erro ao obter o diretório atual:", err)
		return
	}

	nssmPath := filepath.Join(execDir+"\\.nssm", NSSM_EXECUTABLE_TITLE)

	// Comando para remover o serviço
	cmdRemoveService := exec.Command(nssmPath, "remove", serviceName, "confirm")

	// Executar o comando para remover o serviço
	err = cmdRemoveService.Run()
	if err != nil {
		fmt.Println("Erro ao remover o serviço:", err)
		return
	}

	fmt.Println("Serviço removido com sucesso:", serviceName)
}

func NssmStartService() {

	execDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Erro ao obter o diretório atual:", err)
		return
	}

	nssmPath := filepath.Join(execDir+"\\.nssm", NSSM_EXECUTABLE_TITLE)

	// Nome do serviço a ser iniciado
	serviceName := "AttomSvc"

	// Argumento para o comando start (pode ser "start" ou outro argumento necessário)
	startArgument := "start"

	// Comando para iniciar o serviço
	cmdStartService := exec.Command(nssmPath, "start", serviceName, startArgument)

	// Executar o comando para iniciar o serviço
	err = cmdStartService.Run()
	if err != nil {
		fmt.Println("Erro ao iniciar o serviço:", err)
		return
	}
}
