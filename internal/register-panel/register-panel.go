package registerpanel

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"src/post_relay/config"
	"src/post_relay/internal/db"
	"src/post_relay/internal/logger"
	"src/post_relay/internal/utils"
	"src/post_relay/models/environment"
	"src/post_relay/models/panels"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var panelTypes = &panels.TypesPanels{}
var panelUnidades = &panels.Unidades{}
var panelActive = &panels.PainelActive{}
var activesResponse = &panels.ActivesResponse{}

func GetUnidades() ([]panels.Unidade, error) {

	log := logger.GetLogger()

	conn, err := db.Connect()
	if err != nil {
		log.Infof("Error connecting to database: %s", err)
	}
	defer conn.Close(context.Background())

	config, err := utils.LoadConfig()
	if err != nil {
		return nil, err
	}

	sql := fmt.Sprintf(`select
			tus.nu_cnes,
			upper(tus.no_unidade_saude_filtro) nome_unidade,
			tl.no_localidade  nome_municipio
		from
			tb_unidade_saude tus
		inner join tb_localidade tl on
			tus.co_localidade_endereco = tl.co_localidade
		where
			tl.co_ibge in('%s')
		order by
			tus.no_unidade_saude;`, config.API.IBGE)

	rows, err := conn.Query(context.Background(), sql)
	if err != nil {
		logger.GetLogger().Errorf("could not execute SQL query: %v", err)
		return nil, fmt.Errorf("could not execute SQL query: %v", err)
	}
	defer rows.Close()

	var listaUnidades []panels.Unidade

	for rows.Next() {
		var unidade panels.Unidade
		err := rows.Scan(&unidade.NuCnes, &unidade.NomeUnidade, &unidade.NomeMunicipio)
		if err != nil {
			logger.GetLogger().Errorf("could not scan row: %v", err)
			return nil, fmt.Errorf("could not scan row: %v", err)
		}
		listaUnidades = append(listaUnidades, unidade)
	}

	if err := rows.Err(); err != nil {
		logger.GetLogger().Errorf("error during rows iteration: %v", err)
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}

	panelUnidades.SetItems(listaUnidades)

	return listaUnidades, nil
}

func GetTipos() ([]panels.TypeItem, error) {

	log := logger.GetLogger()
	conn, err := db.Connect()
	if err != nil {
		log.Infof("Error connecting to database: %s", err)
	}
	defer conn.Close(context.Background())

	sql := `select
	t.co_tipo_atend_prof codigo,
	t.no_tipo_atend_prof descricao
from
	tb_tipo_atend_prof t
order by
	co_tipo_atend_prof`

	rows, err := conn.Query(context.Background(), sql)
	if err != nil {
		logger.GetLogger().Errorf("could not execute SQL query: %v", err)
		return nil, fmt.Errorf("could not execute SQL query: %v", err)
	}
	defer rows.Close()

	var tipos []panels.TypeItem

	for rows.Next() {
		var tipo panels.TypeItem
		err := rows.Scan(&tipo.Codigo, &tipo.Descricao)
		if err != nil {
			logger.GetLogger().Errorf("could not scan row: %v", err)
			return nil, fmt.Errorf("could not scan row: %v", err)
		}
		tipos = append(tipos, tipo)
	}

	if err := rows.Err(); err != nil {
		logger.GetLogger().Errorf("error during rows iteration: %v", err)
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}

	panelTypes.SetItems(tipos)

	return tipos, nil
}

// Função para fazer a requisição HTTP e processar a resposta
func fetchPanels(endpoint string, apiConfig environment.Config) (panels.ActivesResponse, error) {
	logger.GetLogger().Infof("GetPaineis.endpoint %v", endpoint)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		logger.GetLogger().Errorf("Erro ao criar requisição: %v", err)
		return panels.ActivesResponse{}, err
	}

	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {apiConfig.API.Token},
		"ibge":          {apiConfig.API.IBGE},
	}

	client := &http.Client{
		Timeout: 20 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		if urlErr, ok := err.(*url.Error); ok && urlErr.Timeout() {
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
		return panels.ActivesResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.GetLogger().Errorf("Erro ao ler a resposta: %v", err)
		return panels.ActivesResponse{}, err
	}

	logger.GetLogger().Infof("GetPaineis.return.API: %v", string(body))

	if resp.StatusCode != http.StatusOK {
		logger.GetLogger().Errorf("Erro: Status code %d", resp.StatusCode)
		return panels.ActivesResponse{}, fmt.Errorf("erro na API, status code %d", resp.StatusCode)
	}

	var apiResp panels.ActivesResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		logger.GetLogger().Errorf("Erro ao deserializar o JSON: %v", err)
		return panels.ActivesResponse{}, err
	}

	if apiResp.Error {
		logger.GetLogger().Errorf("Erro na resposta da API: %s", apiResp.Msg)
		return panels.ActivesResponse{}, fmt.Errorf("erro da API: %s", apiResp.Msg)
	}

	return apiResp, nil
}

func GetPaineisOld(unidade string) (panels.ActivesResponse, error) {

	time.Sleep(1 * time.Second)

	apiConfig, err := utils.LoadConfig()
	if err != nil {
		logger.GetLogger().Errorf("erro ao carregar configuração do webhook: %v", err)
	}

	logger.GetLogger().Infof("GetPaineis.unidade: %v", string(unidade))

	// URL da API
	// cnes := utils.OnlyNumber(unidade)
	cnes := utils.OnlyText(unidade)

	endpoint := fmt.Sprintf(`%s/estabelecimentos/%s/paineis`, apiConfig.API.Endpoint, cnes)

	if cnes == "TODAS" {
		for _, unidade := range panelUnidades.GetAll() {
			fmt.Println(unidade.NuCnes)
		}

	} else {
		// cnes := utils.OnlyNumber(unidade)
	}

	logger.GetLogger().Infof("GetPaineis.endpoint %v", string(endpoint))

	// heads
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		logger.GetLogger().Infof("Erro ao criar requisição: %v", err)
	}

	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {apiConfig.API.Token},
		"ibge":          {apiConfig.API.IBGE},
	}

	client := &http.Client{
		Timeout: 20 * time.Second,
	}
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
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.GetLogger().Infof("Erro ao ler a resposta: %v", err)
	}

	logger.GetLogger().Infof("GetPaineis.return.API: %v", string(body))

	if resp.StatusCode != http.StatusOK {
		logger.GetLogger().Infof("Erro: Status code %d", resp.StatusCode)
	}

	var apiResp panels.ActivesResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		logger.GetLogger().Infof("Erro ao deserializar o JSON: %v", err)
	}

	if apiResp.Error {
		logger.GetLogger().Infof("Erro na resposta da API: %s", apiResp.Msg)
	}

	return apiResp, nil
}

func GetPaineis(unidade string) (panels.ActivesResponse, error) {
	time.Sleep(1 * time.Second)

	apiConfig, err := utils.LoadConfig()
	if err != nil {
		logger.GetLogger().Errorf("Erro ao carregar configuração do webhook: %v", err)
		return panels.ActivesResponse{}, err
	}

	logger.GetLogger().Infof("GetPaineis.unidade: %v", unidade)

	cnes := utils.OnlyNumber(unidade)
	endpoint := fmt.Sprintf(`%s/estabelecimentos/%s/paineis`, apiConfig.API.Endpoint, cnes)

	if panelActive.RequestGetAll(unidade) {
		var allPanels []panels.PainelActive

		for _, unidade := range panelUnidades.GetAll() {
			endpoint := fmt.Sprintf(`%s/estabelecimentos/%s/paineis`, apiConfig.API.Endpoint, unidade.NuCnes)

			panelsResp, err := fetchPanels(endpoint, apiConfig)
			if err != nil {
				logger.GetLogger().Errorf("Erro ao buscar painéis para unidade %s: %v", unidade.NuCnes, err)
				continue
			}

			for i := range panelsResp.Obj {
				panelsResp.Obj[i].NuCnes = unidade.NuCnes
			}

			allPanels = append(allPanels, panelsResp.Obj...)
		}

		activesResponse.SetItems(allPanels)

		return panels.ActivesResponse{
			Error: false,
			Msg:   "Consulta concluída",
			Obj:   allPanels,
		}, nil
	}

	panelsResp, err := fetchPanels(endpoint, apiConfig)
	if err != nil {
		logger.GetLogger().Errorf("Erro ao buscar painéis para unidade %s: %v", cnes, err)
	}

	for i := range panelsResp.Obj {
		panelsResp.Obj[i].NuCnes = cnes
	}

	return panelsResp, err
}

func SavePanel(cnes string, panel string, nomePanel string) (environment.Config, error) {
	viper.SetConfigFile(config.FILE_ENVIRONMENT_APPLICATION)
	viper.SetConfigType("toml")

	logger.GetLogger().Info(cnes, panel, nomePanel)

	if err := viper.ReadInConfig(); err != nil {
		return environment.Config{}, fmt.Errorf("erro ao ler o arquivo de configuração: %v", err)
	}

	var config environment.Config
	if err := viper.Unmarshal(&config); err != nil {
		return environment.Config{}, fmt.Errorf("erro ao mapear as configurações para a struct: %v", err)
	}

	cnes = utils.OnlyNumber(cnes)

	painelsActives := activesResponse.GetPanelsActives()

	for _, panelActive := range painelsActives {
		sectors := panelActive.LocalAtendimento

		if cnes == panelActive.NuCnes {
			for _, sector := range sectors {

				if nomePanel == "0 - TODOS" {
					newPanel := map[string]interface{}{
						"cnes":        panelActive.NuCnes,
						"description": fmt.Sprintf("Painel %s adicionado", utils.ToUpperCase(sector.Nome)),
						"type":        []string{utils.ToUpperCase(sector.Nome)},
						"queue": map[string]string{
							"panelUuid":  panelActive.IDPainel,
							"sectorUuid": sector.ID,
						},
					}

					existingPanels := viper.Get("panels.items").([]interface{})
					viper.Set("panels.items", append(existingPanels, newPanel))
				} else {
					if utils.ToUpperCase(sector.Nome) == nomePanel {
						newPanel := map[string]interface{}{
							"cnes":        panelActive.NuCnes,
							"description": fmt.Sprintf("Painel %s adicionado", utils.ToUpperCase(sector.Nome)),
							"type":        []string{utils.ToUpperCase(sector.Nome)},
							"queue": map[string]string{
								"panelUuid":  panelActive.IDPainel,
								"sectorUuid": sector.ID,
							},
						}

						existingPanels := viper.Get("panels.items").([]interface{})
						viper.Set("panels.items", append(existingPanels, newPanel))
					}
				}
			}
		}

		if panelActive.RequestGetAll(cnes) {
			for _, sector := range sectors {
				newPanel := map[string]interface{}{
					"cnes":        panelActive.NuCnes,
					"description": fmt.Sprintf("Painel %s adicionado", utils.ToUpperCase(sector.Nome)),
					"type":        []string{utils.ToUpperCase(sector.Nome)},
					"queue": map[string]string{
						"panelUuid":  panelActive.IDPainel,
						"sectorUuid": sector.ID,
					},
				}

				existingPanels := viper.Get("panels.items").([]interface{})
				viper.Set("panels.items", append(existingPanels, newPanel))
			}
		}

	}

	if err := viper.WriteConfig(); err != nil {
		log.Fatalf("Erro ao salvar a configuração: %v", err)
	}

	return config, nil
}
