package dispatch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"src/post_relay/internal/associations"
	"src/post_relay/internal/utils"
	"src/post_relay/models/panels"
)

func MakePayload(receivedJSON map[string]interface{}) (panels.APIPayload, error) {
	patientName, ok := receivedJSON["cidadao"].(string)
	if !ok {
		patientName = "Unknown"
	}

	idCbo, ok := receivedJSON["prof_cbo_nu"].(string)
	if !ok {
		idCbo = ""
	}

	cnes, ok := receivedJSON["cnes"].(string)
	if !ok {
		cnes = ""
	}

	idServico := utils.ToString(receivedJSON["co_tipo_servico"].(float64))
	if !ok {
		idServico = ""
	}

	environment, err := utils.LoadConfig()
	if err != nil {
		return panels.APIPayload{}, err
	}

	painelQueue, err := associations.LoadPainel(environment.Panels, cnes, idServico, idCbo)
	if err != nil {
		return panels.APIPayload{}, err
	}

	return panels.APIPayload{
		Context:            receivedJSON,
		NomePaciente:       patientName,
		IdPainel:           painelQueue.PanelUuid,
		IdLocalAtendimento: painelQueue.SectorUuid,
	}, nil
}

func SendMessage(payload panels.APIPayload) error {
	apiConfig, err := utils.LoadConfig()
	if err != nil {
		return fmt.Errorf("erro ao carregar configuração do webhook: %v", err)
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("erro ao serializar payload: %v", err)
	}

	req, err := http.NewRequest("POST", apiConfig.API.Endpoint, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("erro ao criar requisição: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("ibge", apiConfig.API.IBGE)
	req.Header.Set("token", apiConfig.API.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao enviar requisição para API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("resposta não-200 recebida da API: %s, corpo: %s", resp.Status, string(body))
	}

	return nil
}
