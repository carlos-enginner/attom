package dispatchpanel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"src/post_relay/internal/associations"
	"src/post_relay/internal/logger"
	"src/post_relay/internal/utils"
	"src/post_relay/models/panels"
	"time"

	"github.com/sirupsen/logrus"
)

func MakePayload(notificationPayload string) (panels.APIPayload, error) {

	var notificationParsed map[string]interface{}
	if err := json.Unmarshal([]byte(notificationPayload), &notificationParsed); err != nil {
		log.Fatal("Error parsing notification payload:", err)
	}

	logger.GetLogger().Info("Dispach.MakePayload.start")

	cidadao, ok := notificationParsed["cidadao"].(string)
	if !ok {
		cidadao = "Unknown"
	}

	cnes, ok := notificationParsed["cnes"].(string)
	if !ok {
		cnes = ""
	}

	localChamada, ok := notificationParsed["local_chamada"].(string)
	if !ok {
		localChamada = ""
	}

	environment, err := utils.LoadConfig()
	if err != nil {
		return panels.APIPayload{}, err
	}

	painelQueue, err := associations.LoadPainel(environment.Panels, cnes, localChamada)
	if err != nil {
		return panels.APIPayload{}, err
	}

	payload := panels.APIPayload{
		Context:            notificationParsed,
		NomePaciente:       cidadao,
		IdPainel:           painelQueue.PanelUuid,
		IdLocalAtendimento: painelQueue.SectorUuid,
	}

	logger.GetLogger().Infof("Dispatch.MakePayload.data - [nome_paciente: %s id_painel: %s id_local_atendimento: %s])", payload.NomePaciente, payload.IdPainel, payload.IdLocalAtendimento)

	return payload, nil
}

func SendMessage(payload panels.APIPayload) error {

	logger.GetLogger().Info("Dispatch.SendMessage.start")

	apiConfig, err := utils.LoadConfig()
	if err != nil {
		logger.GetLogger().Errorf("erro ao carregar configuração do webhook: %v", err)
		return fmt.Errorf("erro ao carregar configuração do webhook: %v", err)
	}

	data, err := json.Marshal(payload)
	if err != nil {
		logger.GetLogger().Errorf("erro ao serializar payload: %v", err)
		return fmt.Errorf("erro ao serializar payload: %v", err)
	}

	var trace *httptrace.ClientTrace
	if apiConfig.Application.HttpDebug {
		trace = &httptrace.ClientTrace{
			DNSStart: func(dnsInfo httptrace.DNSStartInfo) {
				logger.GetLogger().Infof("Iniciando consulta DNS para %v\n", dnsInfo.Host)
			},
			DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
				logger.GetLogger().Infof("Consulta DNS concluída: %v\n", dnsInfo.Addrs)
			},
			ConnectStart: func(network, addr string) {
				logger.GetLogger().Infof("Iniciando conexão para %v\n", addr)
			},
			ConnectDone: func(network, addr string, err error) {
				if err != nil {
					logger.GetLogger().Infof("Erro na conexão para %v: %v\n", addr, err)
				} else {
					logger.GetLogger().Infof("Conexão estabelecida para %v\n", addr)
				}
			},
			GotConn: func(connInfo httptrace.GotConnInfo) {
				if connInfo.Reused {
					logger.GetLogger().Infof("Conexão reutilizada: %v\n", connInfo.Reused)
				} else {
					logger.GetLogger().Infof("Nova conexão estabelecida\n")
				}
			},
			// Quando a requisição foi enviada
			WroteRequest: func(reqInfo httptrace.WroteRequestInfo) {
				logger.GetLogger().Infof("Requisição enviada: %v\n", reqInfo)
			},
		}
	}

	req, err := http.NewRequest("POST", apiConfig.API.Endpoint, bytes.NewBuffer(data))
	if trace != nil {
		req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	}
	if err != nil {
		logger.GetLogger().Errorf("erro ao criar requisição: %v", err)
		return fmt.Errorf("erro ao criar requisição: %v", err)
	}

	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {apiConfig.API.Token},
		"ibge":          {apiConfig.API.IBGE},
	}

	timeoutConnection := apiConfig.Application.TimeoutConnection

	client := &http.Client{
		Timeout: timeoutConnection * time.Second,
	}

	logger.GetLogger().Infof("Dispatch.SendMessage.data - [ibge: %s token: %s end_point: %s]", apiConfig.API.IBGE, apiConfig.API.Token, apiConfig.API.Endpoint)

	resp, err := client.Do(req)
	if err != nil {
		if err, ok := err.(*url.Error); ok && err.Timeout() {
			logger.GetLogger().WithFields(logrus.Fields{
				"error": err,
				"type":  "timeout",
			}).Error("Timeout ao tentar conectar com a API")
		} else {
			logger.GetLogger().WithFields(logrus.Fields{
				"error": err,
				"type":  "connection",
			}).Error("Erro ao enviar requisição para API")
		}
		return fmt.Errorf("erro ao enviar requisição para API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logger.GetLogger().Errorf("resposta não-200 recebida da API: %s, corpo: %s", resp.Status, string(body))
		return fmt.Errorf("resposta não-200 recebida da API: %s, corpo: %s", resp.Status, string(body))
	}

	return nil
}
