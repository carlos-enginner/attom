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

func NssmExtractApp() {
	execDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %s", err)
		return
	}

	nssmDir := filepath.Join(execDir, ".nssm")

	nssmPath := filepath.Join(nssmDir, NSSM_EXECUTABLE_TITLE)

	if _, err := os.Stat(nssmPath); err == nil {
		fmt.Println("Arquivo 'nssm.exe' já existe em:", nssmPath)
		return
	}

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

	err = os.MkdirAll(nssmDir, 0755)
	if err != nil {
		fmt.Println("Erro ao criar a pasta '.nssm':", err)
		return
	}

	err = os.WriteFile(nssmPath, nssmData, 0755)
	if err != nil {
		fmt.Println("Erro ao escrever o arquivo nssm.exe:", err)
		return
	}
}

func NssmInstallService() {

	execDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	nssmPath := filepath.Join(execDir+"\\.nssm", NSSM_EXECUTABLE_TITLE)

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

	execDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	nssmPath := filepath.Join(execDir+"\\.nssm", NSSM_EXECUTABLE_TITLE)

	cmdRemoveService := exec.Command(nssmPath, "remove", WINDOWS_SERVICE_NAME, "confirm")

	err = cmdRemoveService.Run()
	if err != nil {
		fmt.Println("Error removing service:", err)
		return
	}

	fmt.Println("Service removed successfully:", WINDOWS_SERVICE_NAME)
}

func NssmStartService() {

	execDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	nssmPath := filepath.Join(execDir+"\\.nssm", NSSM_EXECUTABLE_TITLE)

	startArgument := "start_service"

	cmdStartService := exec.Command(nssmPath, "start", WINDOWS_SERVICE_NAME, startArgument)

	err = cmdStartService.Run()
	if err != nil {
		fmt.Println("Error starting service:", err)
		return
	}
}
