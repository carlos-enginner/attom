package win64

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"src/post_relay/internal/logger"
)

//go:embed assets/nssm.exe
var nssmData []byte

const NSSM_EXECUTABLE_TITLE = "nssm.exe"
const WINDOWS_SERVICE_NAME = "AttomSvc"

func NssmExtractApp() (string, error) {

	logger := logger.GetLogger()

	execDir, err := os.Getwd()
	if err != nil {
		logger.Errorf("erro ao obter o diretório atual: %s", err)
		return "", fmt.Errorf("erro ao obter o diretório atual: %w", err)
	}

	nssmDir := filepath.Join(execDir, ".nssm")

	// Verifica se a pasta já existe antes de tentar criá-la
	if _, err := os.Stat(nssmDir); os.IsNotExist(err) {
		err = os.MkdirAll(nssmDir, 0755)
		if err != nil {
			logger.Errorf("erro ao criar a pasta '.nssm': %s", err)
			return "", fmt.Errorf("erro ao criar a pasta '.nssm': %w", err)
		}
	}

	cmd := exec.Command("attrib", "+h", nssmDir)
	err = cmd.Run()
	if err != nil {
		logger.Errorf("erro ao tornar a pasta '.nssm' oculta: %s", err)
		return "", fmt.Errorf("erro ao tornar a pasta '.nssm' oculta: %w", err)
	}

	// Agora cria o arquivo nssm.exe
	nssmPath := filepath.Join(nssmDir, NSSM_EXECUTABLE_TITLE)
	err = os.WriteFile(nssmPath, nssmData, 0755)
	if err != nil {
		logger.Errorf("erro ao escrever o arquivo nssm.exe: %s", err)
		return "", fmt.Errorf("erro ao escrever o arquivo nssm.exe: %w", err)
	}

	return nssmPath, nil
}

func NssmInstallService() {

	logger := logger.GetLogger()

	nssmPath, err := NssmExtractApp()
	if err != nil {
		logger.Errorf("Error extracting nssm application: %s", err)
		fmt.Println("Error extracting nssm application:", err)
		return
	}

	execDir, err := os.Getwd()
	if err != nil {
		logger.Errorf("Error getting current directory: %s", err)
		fmt.Println("Error getting current directory:", err)
		return
	}

	executablePath := execDir + "\\attom.exe"

	cmdCreateService := exec.Command(nssmPath, "install", WINDOWS_SERVICE_NAME, executablePath)

	err = cmdCreateService.Run()
	if err != nil {
		logger.Errorf("Error creating service: %s", err)
		fmt.Println("Error creating service:", err)
		return
	}

	cmdSetDescription := exec.Command(nssmPath, "set", WINDOWS_SERVICE_NAME, "Description", "O serviço responsável por detectar e capturar eventos de atendimento no e-sus/PEC e envia-lós a um serviço externo de painel eletrônico")

	err = cmdSetDescription.Run()
	if err != nil {
		logger.Errorf("Error setting service description: %s", err)
		fmt.Println("Error setting service description:", err)
		return
	}

	fmt.Println("Service created successfully:", WINDOWS_SERVICE_NAME)
}

func NssmRemoveService() {

	logger := logger.GetLogger()
	nssmPath, err := NssmExtractApp()
	if err != nil {
		logger.Errorf("Error setting service description: %s", err)
		fmt.Println("Error extracting nssm application:", err)
		return
	}

	cmdRemoveService := exec.Command(nssmPath, "remove", WINDOWS_SERVICE_NAME, "confirm")

	err = cmdRemoveService.Run()
	if err != nil {
		logger.Errorf("Error removing service: %s", err)
		fmt.Println("Error removing service:", err)
		return
	}

	fmt.Println("Service removed successfully:", WINDOWS_SERVICE_NAME)
}

func NssmStartService() {

	logger := logger.GetLogger()
	nssmPath, err := NssmExtractApp()
	if err != nil {
		logger.Errorf("Error extracting nssm application: %s", err)
		fmt.Println("Error extracting nssm application:", err)
		return
	}

	startArgument := "start_service"
	cmdStartService := exec.Command(nssmPath, "start", WINDOWS_SERVICE_NAME, startArgument)

	err = cmdStartService.Run()
	if err != nil {
		logger.Errorf("Error starting service: %s", err)
		fmt.Println("Error starting service:", err)
		return
	}
}
