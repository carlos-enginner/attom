package win64

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed assets/nssm.exe
var nssmData []byte

const NSSM_EXECUTABLE_TITLE = "nssm.exe"
const WINDOWS_SERVICE_NAME = "AttomSvc"

func NssmExtractApp() (string, error) {
	execDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %w", err)
	}

	nssmDir := filepath.Join(execDir, ".nssm")

	err = os.MkdirAll(nssmDir, 0755)
	if err != nil {
		return "", fmt.Errorf("erro ao criar a pasta 'nssm': %w", err)
	}

	cmd := exec.Command("attrib", "+h", nssmDir)
	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("erro ao tornar a pasta '.nssm' oculta: %w", err)

	}

	err = os.MkdirAll(nssmDir, 0755)
	if err != nil {
		return "", fmt.Errorf("erro ao criar a pasta '.nssm': %w", err)
	}

	nssmPath := filepath.Join(nssmDir, NSSM_EXECUTABLE_TITLE)
	err = os.WriteFile(nssmPath, nssmData, 0755)
	if err != nil {
		return "", fmt.Errorf("erro ao escrever o arquivo nssm.exe: %w", err)
	}

	return nssmPath, nil
}

func NssmInstallService() {

	nssmPath, err := NssmExtractApp()
	if err != nil {
		fmt.Println("Error extracting nssm application:", err)
		return
	}

	execDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	executablePath := execDir + "\\attom.exe"

	cmdCreateService := exec.Command(nssmPath, "install", WINDOWS_SERVICE_NAME, executablePath)

	err = cmdCreateService.Run()
	if err != nil {
		fmt.Println("Error creating service:", err)
		return
	}

	cmdSetDescription := exec.Command(nssmPath, "set", WINDOWS_SERVICE_NAME, "Description", "O serviço responsável por detectar e capturar eventos de atendimento no e-sus/PEC e envia-lós a um serviço externo de painel eletrônico")

	err = cmdSetDescription.Run()
	if err != nil {
		fmt.Println("Error setting service description:", err)
		return
	}

	fmt.Println("Service created successfully:", WINDOWS_SERVICE_NAME)
}

func NssmRemoveService() {

	nssmPath, err := NssmExtractApp()
	if err != nil {
		fmt.Println("Error extracting nssm application:", err)
		return
	}

	cmdRemoveService := exec.Command(nssmPath, "remove", WINDOWS_SERVICE_NAME, "confirm")

	err = cmdRemoveService.Run()
	if err != nil {
		fmt.Println("Error removing service:", err)
		return
	}

	fmt.Println("Service removed successfully:", WINDOWS_SERVICE_NAME)
}

func NssmStartService() {

	nssmPath, err := NssmExtractApp()
	if err != nil {
		fmt.Println("Error extracting nssm application:", err)
		return
	}

	startArgument := "start_service"

	cmdStartService := exec.Command(nssmPath, "start", WINDOWS_SERVICE_NAME, startArgument)

	err = cmdStartService.Run()
	if err != nil {
		fmt.Println("Error starting service:", err)
		return
	}
}
