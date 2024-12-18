package win64

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"src/post_relay/internal/logger"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
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
		logger.Errorf("erro ao obter o diretório atual: %s", err)
		fmt.Println("Erro ao obter o diretório atual:", err)
		return
	}

	executablePath := filepath.Join(execDir, "Attom.exe")
	batScript := fmt.Sprintf(`
@echo off
%s install AttomSvc "%s"
%s set AttomSvc Application "%s"
%s set AttomSvc AppDirectory "%s"
%s set AttomSvc AppParameters "start"
%s set AttomSvc Description "The service responsible for detecting and capturing service events in e-sus/PEC and sending them to an external electronic panel service"
%s set AttomSvc Start SERVICE_AUTO_START
`, nssmPath, executablePath, nssmPath, executablePath, nssmPath, execDir, nssmPath, nssmPath, nssmPath)

	filePath := filepath.Join(execDir, ".nssm", "nssm.bat")

	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		fmt.Println("Erro ao obter o caminho absoluto:", err)
		return
	}

	file, err := os.Create(absFilePath)
	if err != nil {
		fmt.Println("Erro ao criar o arquivo:", err)
		return
	}
	defer file.Close()

	writer := transform.NewWriter(file, charmap.ISO8859_1.NewEncoder())
	_, err = writer.Write([]byte(batScript))
	if err != nil {
		fmt.Println("Erro ao criar o arquivo:", err)
		return
	}

	cmdRunBat := exec.Command("cmd.exe", "/C", absFilePath)
	err = cmdRunBat.Run()
	if err != nil {
		logger.Errorf("Error running .bat script: %s", err)
		fmt.Println("Error running .bat script:", err)
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

	startArgument := "start"
	cmdStartService := exec.Command(nssmPath, "start", WINDOWS_SERVICE_NAME, startArgument)

	err = cmdStartService.Run()
	if err != nil {
		logger.Errorf("Error starting service: %s", err)
		fmt.Println("Error starting service:", err)
		return
	}
}
