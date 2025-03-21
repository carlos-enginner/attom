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
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type LocalAtendimento struct {
	ID   string `json:"id"`
	Nome string `json:"nome"`
}

type Painel struct {
	Descricao        string             `json:"descricao"`
	IDPainel         string             `json:"idPainel"`
	NomePainel       string             `json:"nomePainel"`
	DuracaoChamada   int                `json:"duracaoChamada"`
	LocalAtendimento []LocalAtendimento `json:"localAtendimento"`
}

type APIResponse struct {
	Error bool     `json:"error"`
	Msg   string   `json:"msg"`
	Obj   []Painel `json:"obj"`
}

type Unidade struct {
	NuCnes        string
	NomeUnidade   string
	NomeMunicipio string
}

type Tipo struct {
	Codigo    int64
	Descricao string
}

func GetUnidades() ([]Unidade, error) {

	conn, err := db.Connect()
	// Conectar ao banco de dados
	if err != nil {
		log.Fatal("Error connecting to database: - GetUnidades", err)
	}
	defer conn.Close(context.Background())

	config, err := utils.LoadConfig()
	if err != nil {
		return nil, err
	}

	// Definindo a consulta SQL para buscar as unidades de saúde
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

	// Executando a consulta e obtendo os resultados
	rows, err := conn.Query(context.Background(), sql)
	if err != nil {
		// Logando o erro caso a consulta falhe
		logger.GetLogger().Errorf("could not execute SQL query: %v", err)
		return nil, fmt.Errorf("could not execute SQL query: %v", err)
	}
	defer rows.Close()

	// Criando um slice para armazenar as unidades
	var unidades []Unidade

	// Iterando pelas linhas retornadas e populando o slice
	for rows.Next() {
		var unidade Unidade
		err := rows.Scan(&unidade.NuCnes, &unidade.NomeUnidade, &unidade.NomeMunicipio)
		if err != nil {
			// Logando erro de scan
			logger.GetLogger().Errorf("could not scan row: %v", err)
			return nil, fmt.Errorf("could not scan row: %v", err)
		}
		// Adicionando a unidade ao slice
		unidades = append(unidades, unidade)
	}

	// Verificando se houve erro durante a iteração das linhas
	if err := rows.Err(); err != nil {
		// Logando o erro de iteração
		logger.GetLogger().Errorf("error during rows iteration: %v", err)
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}

	return unidades, nil
}

func GetTipos() ([]Tipo, error) {

	conn, err := db.Connect()
	// Conectar ao banco de dados
	if err != nil {
		log.Fatal("Error connecting to database - GetTipos:", err)
	}
	defer conn.Close(context.Background())

	// Definindo a consulta SQL para buscar as unidades de saúde
	sql := `select
	t.co_tipo_atend_prof codigo,
	t.no_tipo_atend_prof descricao
from
	tb_tipo_atend_prof t
order by
	co_tipo_atend_prof`

	// Executando a consulta e obtendo os resultados
	rows, err := conn.Query(context.Background(), sql)
	if err != nil {
		// Logando o erro caso a consulta falhe
		logger.GetLogger().Errorf("could not execute SQL query: %v", err)
		return nil, fmt.Errorf("could not execute SQL query: %v", err)
	}
	defer rows.Close()

	// Criando um slice para armazenar as unidades
	var tipos []Tipo

	// Iterando pelas linhas retornadas e populando o slice
	for rows.Next() {
		var tipo Tipo
		err := rows.Scan(&tipo.Codigo, &tipo.Descricao)
		if err != nil {
			// Logando erro de scan
			logger.GetLogger().Errorf("could not scan row: %v", err)
			return nil, fmt.Errorf("could not scan row: %v", err)
		}
		// Adicionando a unidade ao slice
		tipos = append(tipos, tipo)

		environment.PainelTypes = append(environment.PainelTypes, tipo.Descricao)
	}

	// Verificando se houve erro durante a iteração das linhas
	if err := rows.Err(); err != nil {
		// Logando o erro de iteração
		logger.GetLogger().Errorf("error during rows iteration: %v", err)
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}

	return tipos, nil
}

func GetPaineis(unidade string) (APIResponse, error) {

	time.Sleep(1 * time.Second)

	apiConfig, err := utils.LoadConfig()
	if err != nil {
		logger.GetLogger().Errorf("erro ao carregar configuração do webhook: %v", err)
	}

	logger.GetLogger().Infof("GetPaineis.unidade: %v", string(unidade))

	// URL da API
	cnes := utils.OnlyNumber(unidade)
	endpoint := fmt.Sprintf(`%s/estabelecimentos/%s/paineis`, apiConfig.API.Endpoint, cnes)

	logger.GetLogger().Infof("GetPaineis.endpoint %v", string(endpoint))

	// Cabeçalhos necessários
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

	// Lendo a resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.GetLogger().Infof("Erro ao ler a resposta: %v", err)
	}

	logger.GetLogger().Infof("GetPaineis.return.API: %v", string(body))

	// Verificando o status da resposta
	if resp.StatusCode != http.StatusOK {
		logger.GetLogger().Infof("Erro: Status code %d", resp.StatusCode)
	}

	// Mapear a resposta para a estrutura Go
	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		logger.GetLogger().Infof("Erro ao deserializar o JSON: %v", err)
	}

	// Verificando se houve erro no campo 'error' da resposta
	if apiResp.Error {
		logger.GetLogger().Infof("Erro na resposta da API: %s", apiResp.Msg)
	}

	return apiResp, nil
}

func SavePanel(cnes string, panel string, tipo string) (environment.Config, error) {
	viper.SetConfigFile(config.FILE_ENVIRONMENT_APPLICATION)
	viper.SetConfigType("toml")

	logger.GetLogger().Info(cnes, panel, tipo)

	if err := viper.ReadInConfig(); err != nil {
		return environment.Config{}, fmt.Errorf("erro ao ler o arquivo de configuração: %v", err)
	}

	var config environment.Config
	if err := viper.Unmarshal(&config); err != nil {
		return environment.Config{}, fmt.Errorf("erro ao mapear as configurações para a struct: %v", err)
	}

	// formatting texts
	panelInfo := strings.Split(panel, " - ")
	tipo = utils.OnlyText(tipo)

	// por default pega o painel informado, se não pega o geral
	descriptionPanel := fmt.Sprintf("Painel %s registrado", tipo)
	typePanel := []string{tipo}

	if tipo == "TODOS" {
		descriptionPanel = "Painel UNIFICADO registrado"
		typePanel = environment.PainelTypes
	}

	newPanel := map[string]interface{}{
		"cnes":        utils.OnlyNumber(cnes),
		"description": descriptionPanel,
		"type":        typePanel,
		"queue": map[string]string{
			"panelUuid":  panelInfo[1],
			"sectorUuid": panelInfo[3],
		},
	}

	existingPanels := viper.Get("panels.items").([]interface{})
	viper.Set("panels.items", append(existingPanels, newPanel))

	if err := viper.WriteConfig(); err != nil {
		log.Fatalf("Erro ao salvar a configuração: %v", err)
	}

	return config, nil
}
