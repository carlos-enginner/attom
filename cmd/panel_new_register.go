package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"src/post_relay/internal/logger"
	registerpanel "src/post_relay/internal/register-panel"
	"src/post_relay/internal/utils"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	appStyle     = lipgloss.NewStyle().Margin(1, 2, 0, 2)
)

type model struct {
	spinner  spinner.Model
	timer    time.Time // Marca quando o spinner começou
	timeout  time.Duration
	quitting bool
	form     *huh.Form
}

func getUnidades() []string {

	unidades, err := registerpanel.GetUnidades()
	if err != nil {
		log.Fatalf("Error retrieving unidades: %v", err)
	}

	var options []string

	if len(unidades) == 0 {
		return append(options, "nenhum registro encontrado")
	}

	for _, unidade := range unidades {
		options = append(options, unidade.NuCnes+" - "+unidade.NomeUnidade)
	}

	return options
}

func getPaineis(cnes string) []string {

	cnes = "2569841"
	panels, err := registerpanel.GetPaineis(cnes)
	if err != nil {
		log.Fatalf("Error retrieving paneis: %v", err)
	}

	var options []string

	if len(panels.Obj) == 0 {
		return append(options, "nenhum registro encontrado")
	}

	for _, painel := range panels.Obj {
		for _, local := range painel.LocalAtendimento {
			options = append(options, fmt.Sprintf("%s - %s - %s - %s", painel.NomePainel, painel.IDPainel, local.Nome, local.ID))

		}
	}

	return options
}

func getTipos() []string {

	tipos, err := registerpanel.GetTipos()
	if err != nil {
		log.Fatalf("Error retrieving unidades: %v", err)
	}

	var options []string

	if len(tipos) == 0 {
		return append(options, "nenhum registro encontrado")
	} else {
		options = append(options, "0 - TODOS")
	}

	for _, tipo := range tipos {
		options = append(options, fmt.Sprintf("%d - %s", tipo.Codigo, tipo.Descricao))
	}

	return options
}

func registerPanel(cnes string, panel string, tipos string) {

	_, err := registerpanel.SavePanel(cnes, panel, tipos)
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
	s.Spinner = spinner.Points

	newConfirm := huh.NewConfirm().
		Key("btn_confirm").
		Validate(func(b bool) error {
			if painelSelected == "nenhum registro encontrado" {
				return errors.New("preenchimento incorreto, não é possível seguir, desculpe")
			}
			return nil
		}).
		TitleFunc(func() string {
			return "Confirma registro do painel: " + tipoSelected + "?"
		}, &tipoSelected).Affirmative("Sim").Negative("Não")

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Form register panel").
				Description("Here you can register panels dynamically"),
			huh.NewSelect[string]().
				Height(8).
				Title("1. Tipo do painel").
				Key("field_tipos").
				Value(&tipoSelected).
				OptionsFunc(func() []huh.Option[string] {
					var options = getTipos()
					return huh.NewOptions(options...)
				}, &tipoSelected),

			huh.NewSelect[string]().
				Height(8).
				Title("2. Unidade de saúde").
				Key("field_unidades").
				Value(&unidadeSelected).
				OptionsFunc(func() []huh.Option[string] {
					var options = getUnidades()
					return huh.NewOptions(options...)
				}, nil),

			huh.NewSelect[string]().
				Height(5).
				TitleFunc((func() string {
					label := "3. Paineis ativos na API para unidade:"
					if unidadeSelected != "" && unidadeSelected != "nenhum registro encontrado" {
						return fmt.Sprintf("%s: %s", label, unidadeSelected)
					} else {
						return label
					}
				}), &unidadeSelected).
				Key("field_paineis").
				Value(&painelSelected).
				OptionsFunc(func() []huh.Option[string] {
					options := []string{"nenhum registro encontrado"}
					cnes := utils.OnlyNumber(unidadeSelected)
					if cnes != "" {
						options = getPaineis(cnes)
					}
					return huh.NewOptions(options...)
				}, &unidadeSelected),

			newConfirm,
		),
	)

	return model{
		spinner: s,
		form:    form,
		timeout: 3 * time.Second,
		timer:   time.Now(),
	}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c", "q":
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
			s += "@ Loading form register panel @ \n\n" + m.spinner.View() + " preparing application data..."
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
