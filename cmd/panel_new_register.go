package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"src/post_relay/internal/db"
	"src/post_relay/internal/logger"
	"src/post_relay/internal/utils"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// type Model struct {
// 	form *huh.Form // huh.Form is just a tea.Model
// }

// func NewModel() Model {
// 	var unidadeSelected string
// 	var loadingUnidades bool = true
// 	var tipoSelected string
// 	return Model{
// 		form: huh.NewForm(
// 			huh.NewGroup(

// 				huh.NewSelect[string]().
// 					Height(8).
// 					Title("Unidades").
// 					Key("field_unidades").
// 					Value(&unidadeSelected).
// 					OptionsFunc(func() []huh.Option[string] {
// 						return huh.NewOptions(
// 							"unidade 1",
// 							"unidade 2",
// 							"unidade 3",
// 						)
// 					}, &loadingUnidades),

// 				huh.NewSelect[string]().
// 					Height(8).
// 					Title("Paineis").
// 					Key("field_paineis").
// 					OptionsFunc(func() []huh.Option[string] {
// 						var options []string
// 						switch unidadeSelected {
// 						case "unidade 1":
// 							options = []string{"Opção 1.1", "Opção 1.2", "Opção 1.3"}
// 						case "unidade 2":
// 							options = []string{"Opção 2.1", "Opção 2.2", "Opção 2.3"}
// 						case "unidade 3":
// 							options = []string{"Opção 3.1", "Opção 3.2", "Opção 3.3"}
// 						default:
// 							options = []string{}
// 						}
// 						return huh.NewOptions(options...)
// 					}, &unidadeSelected),

// 				huh.NewSelect[string]().
// 					Height(8).
// 					Title("Tipos").
// 					Key("field_tipos").
// 					OptionsFunc(func() []huh.Option[string] {
// 						var options = []string{"Opção 3.1", "Opção 3.2", "Opção 3.3"}
// 						return huh.NewOptions(options...)
// 					}, &tipoSelected),

// 				huh.NewConfirm().
// 					Title("Confirma registro?"),
// 			),
// 		),
// 	}
// }

// func (m Model) Init() tea.Cmd {
// 	return m.form.Init()
// }

// func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		switch msg.String() {
// 		case "esc", "ctrl+c", "q":
// 			return m, tea.Quit
// 		}
// 	}

// 	var cmds []tea.Cmd

// 	// process the form
// 	form, cmd := m.form.Update(msg)
// 	if f, ok := form.(*huh.Form); ok {
// 		m.form = f
// 		cmds = append(cmds, cmd)
// 	}

// 	if m.form.State == huh.StateCompleted {
// 		switch msg := msg.(type) {
// 		case tea.KeyMsg:
// 			switch msg.String() {
// 			case "enter":
// 				return m, tea.Quit
// 			}
// 		}
// 	}

// 	return m, tea.Batch(cmds...)
// }

// func (m Model) View() string {
// 	if m.form.State == huh.StateCompleted {
// 		unidades := m.form.GetString("field_unidades")
// 		paineis := m.form.GetString("field_paineis")
// 		tipos := m.form.GetString("field_tipos")

// 		return fmt.Sprintf("panel added %s - %s - %s. Press enter for exit", unidades, paineis, tipos)
// 	}
// 	return m.form.View()
// }

var (
	spinnerStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Margin(1, 0)
	dotStyle      = helpStyle.UnsetMargins()
	durationStyle = dotStyle
	appStyle      = lipgloss.NewStyle().Margin(1, 2, 0, 2)
)

type model struct {
	spinner  spinner.Model
	timer    time.Time // Marca quando o spinner começou
	timeout  time.Duration
	quitting bool
	form     *huh.Form
}
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

func getUnidades() []string {
	conn, err := db.Connect()
	// Conectar ao banco de dados
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer conn.Close(context.Background())

	unidades, err := db.GetUnidades(conn)
	if err != nil {
		log.Fatalf("Error retrieving unidades: %v", err)
	}

	var options []string
	for _, unidade := range unidades {
		options = append(options, unidade.NuCnes+" - "+unidade.NomeUnidade)
	}

	return options
}

func getPaineis(unidade string) []string {

	time.Sleep(2 * time.Second)

	apiConfig, err := utils.LoadConfig()
	if err != nil {
		logger.GetLogger().Errorf("erro ao carregar configuração do webhook: %v", err)
	}

	// URL da API
	cnes := utils.OnlyNumber(unidade)
	endpoint := fmt.Sprintf(`http://painel.icsgo.com.br:7001/ws/v1/estabelecimentos/%s/paineis`, cnes)

	// Cabeçalhos necessários
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Fatalf("Erro ao criar requisição: %v", err)
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
		log.Fatalf("Erro ao ler a resposta: %v", err)
	}

	// Verificando o status da resposta
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Erro: Status code %d", resp.StatusCode)
	}

	// Mapear a resposta para a estrutura Go
	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		log.Fatalf("Erro ao deserializar o JSON: %v", err)
	}

	// Verificando se houve erro no campo 'error' da resposta
	if apiResp.Error {
		log.Fatalf("Erro na resposta da API: %s", apiResp.Msg)
	}

	var options []string
	for _, painel := range apiResp.Obj {
		for _, local := range painel.LocalAtendimento {
			options = append(options, fmt.Sprintf("%s - %s - %s - %s", painel.NomePainel, painel.IDPainel, local.Nome, local.ID))

		}
	}

	if len(options) == 0 {
		options = append(options, "nenhum painel encontrado")
	}

	return options
}

func getTipos() []string {
	conn, err := db.Connect()
	// Conectar ao banco de dados
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer conn.Close(context.Background())

	tipos, err := db.GetTipos(conn)
	if err != nil {
		log.Fatalf("Error retrieving unidades: %v", err)
	}

	var options []string
	for _, tipo := range tipos {
		options = append(options, tipo.Descricao)
	}

	return options
}

func registerPanel(cnes string, panel string, tipos string) {

	_, err := utils.SaveConfig(cnes, panel, tipos)
	if err != nil {
		logger.GetLogger().Errorf("erro ao carregar configuração do webhook: %v", err)
	}
}

func newModel() model {

	var unidadeSelected string
	var tipoSelected string
	var painelSelected string

	s := spinner.New()
	s.Style = spinnerStyle

	newConfirm := huh.NewConfirm().
		Key("btn_confirm").
		Validate(func(b bool) error {
			logger.GetLogger().Infof("btn_confirm: %s", painelSelected)
			if painelSelected == "nenhum painel encontrado" {
				return errors.New("não é possível seguir, desculpe")
			}
			return nil
		}).
		TitleFunc(func() string {
			return "Confirma registro do painel: " + tipoSelected + "?"
		}, &tipoSelected).Affirmative("Sim").Negative("Não")

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Height(8).
				Title("Tipo do painel").
				Key("field_tipos").
				Value(&tipoSelected).
				OptionsFunc(func() []huh.Option[string] {
					var options = getTipos()
					return huh.NewOptions(options...)
				}, &tipoSelected),

			huh.NewSelect[string]().
				Height(8).
				Title("Unidade de saúde").
				Key("field_unidades").
				Value(&unidadeSelected).
				OptionsFunc(func() []huh.Option[string] {
					var options = getUnidades()
					return huh.NewOptions(options...)
				}, nil),

			huh.NewSelect[string]().
				Height(5).
				TitleFunc((func() string {
					return "Paineis do Munícipio: " + unidadeSelected
				}), &unidadeSelected).
				Key("field_paineis").
				Value(&painelSelected).
				OptionsFunc(func() []huh.Option[string] {
					var options []string
					if unidadeSelected != "" {
						options = getPaineis(unidadeSelected)
					}
					return huh.NewOptions(options...)
				}, &unidadeSelected),

			newConfirm,
		),
	)

	return model{
		spinner: s,
		form:    form,
		timeout: 1 * time.Second,
		timer:   time.Now(),
	}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// switch msg := msg.(type) {
	// case tea.KeyMsg:
	// 	m.quitting = true
	// 	return m, tea.Quit
	// case resultMsg:
	// 	m.results = append(m.results[1:], msg)
	// 	return m, nil
	// case spinner.TickMsg:
	// 	var cmd tea.Cmd
	// 	m.spinner, cmd = m.spinner.Update(msg)
	// 	return m, cmd
	// case huh.Form:
	// 	var cmd tea.Cmd
	// 	form, cmd := m.form.Update(msg)
	// 	if f, ok := form.(*huh.Form); ok {
	// 		m.form = f
	// 	}
	// 	return m, cmd
	// default:
	// 	return m, nil
	// }

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c", "q":
			m.quitting = true
			return m, nil
		case "y", "n":
			m.quitting = true
			return m, tea.Quit
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		if time.Since(m.timer) >= m.timeout {
			m.quitting = true
		}
		return m, cmd
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	if m.form.State == huh.StateCompleted {
		btnConfirm := m.form.GetBool("btn_confirm")
		if btnConfirm {
			unidades := m.form.GetString("field_unidades")
			paineis := m.form.GetString("field_paineis")
			tipos := m.form.GetString("field_tipos")
			registerPanel(unidades, paineis, tipos)
		}
		return m, tea.Quit
	}

	return m, cmd

}

func (m model) View() string {
	var s string

	if m.quitting {
		return m.form.View()
	} else {
		if time.Since(m.timer) < m.timeout {
			s += m.spinner.View() + " Loading data..."
		}
	}

	return appStyle.Render(s)
}

func PanelNewRegister() *cobra.Command {
	return &cobra.Command{
		Use:   "register_panel",
		Short: "New Panel Register",
		Run: func(cmd *cobra.Command, args []string) {

			p := tea.NewProgram(newModel())
			if _, err := p.Run(); err != nil {
				fmt.Println("Error running program:", err)
				os.Exit(1)
			}
		},
	}
}
