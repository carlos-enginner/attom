package cmd

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

// type model struct {
// 	flavors, adds []item
// 	list, item    int
// }

// type item struct {
// 	text    string
// 	checked bool
// }

// func (m *model) Init() tea.Cmd {
// 	return nil
// }

// func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch typed := msg.(type) {
// 	case tea.KeyMsg:
// 		return m, m.handleKeyMsg(typed)
// 	}
// 	return m, nil
// }

// func (m *model) handleKeyMsg(msg tea.KeyMsg) tea.Cmd {
// 	switch msg.String() {
// 	case "esc", "ctrl+c":
// 		return tea.Quit
// 	case " ", "enter":
// 		switch m.list {
// 		case 0:
// 			m.flavors[m.item].checked = !m.flavors[m.item].checked
// 		case 1:
// 			m.adds[m.item].checked = !m.adds[m.item].checked
// 		}
// 	case "up":
// 		if m.item > 0 {
// 			m.item--
// 		} else if m.list > 0 {
// 			m.list--
// 			m.item = len(m.flavors) - 1
// 		}
// 	case "down":
// 		switch m.list {
// 		case 0:
// 			if m.item+1 < len(m.flavors) {
// 				m.item++
// 			} else {
// 				m.list++
// 				m.item = 0
// 			}
// 		case 1:
// 			if m.item+1 < len(m.adds) {
// 				m.item++
// 			}
// 		}
// 	}
// 	return nil
// }

// func (m *model) View() string {
// 	curFlavor, curAdd := -1, -1
// 	switch m.list {
// 	case 0:
// 		curFlavor = m.item
// 	case 1:
// 		curAdd = m.item
// 	}
// 	return m.renderList("choose two flavors", m.flavors, curFlavor) +
// 		"\n" +
// 		m.renderList("select adds", m.adds, curAdd)
// }

// func (m *model) renderList(header string, items []item, selected int) string {
// 	out := "~ " + header + ":\n"
// 	for i, item := range items {
// 		sel := " "
// 		if i == selected {
// 			sel = ">"
// 		}
// 		check := " "
// 		if items[i].checked {
// 			check = "✓"
// 		}
// 		out += fmt.Sprintf("%s [%s] %s\n", sel, check, item.text)
// 	}
// 	return out
// }

type Model struct {
	form *huh.Form // huh.Form is just a tea.Model
}

func NewModel() Model {
	var unidadeSelected string
	// var paineis string
	var tipoSelected string
	return Model{
		form: huh.NewForm(
			huh.NewGroup(

				huh.NewSelect[string]().
					Height(8).
					Title("Unidades").
					Key("field_unidades").
					Value(&unidadeSelected).
					OptionsFunc(func() []huh.Option[string] {
						return huh.NewOptions(
							"unidade 1",
							"unidade 2",
							"unidade 3",
						)
					}, nil),

				huh.NewSelect[string]().
					Height(8).
					Title("Paineis").
					Key("field_paineis").
					OptionsFunc(func() []huh.Option[string] {
						var options []string
						switch unidadeSelected {
						case "unidade 1":
							options = []string{"Opção 1.1", "Opção 1.2", "Opção 1.3"}
						case "unidade 2":
							options = []string{"Opção 2.1", "Opção 2.2", "Opção 2.3"}
						case "unidade 3":
							options = []string{"Opção 3.1", "Opção 3.2", "Opção 3.3"}
						default:
							options = []string{}
						}
						return huh.NewOptions(options...)
					}, &unidadeSelected),

				huh.NewSelect[string]().
					Height(8).
					Title("Tipos").
					Key("field_tipos").
					OptionsFunc(func() []huh.Option[string] {
						var options = []string{"Opção 3.1", "Opção 3.2", "Opção 3.3"}
						return huh.NewOptions(options...)
					}, &tipoSelected),

				huh.NewConfirm().
					Title("Confirma registro?"),
			),
		),
	}
}

func (m Model) Init() tea.Cmd {
	return m.form.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	var cmds []tea.Cmd

	// process the form
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	if m.form.State == huh.StateCompleted {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				return m, tea.Quit
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.form.State == huh.StateCompleted {
		unidades := m.form.GetString("field_unidades")
		paineis := m.form.GetString("field_paineis")
		tipos := m.form.GetString("field_tipos")

		return fmt.Sprintf("panel added %s - %s - %s. Press enter for exit", unidades, paineis, tipos)
	}
	return m.form.View()
}

func PanelNewRegister() *cobra.Command {
	return &cobra.Command{
		Use:   "new_panel",
		Short: "New Panel Register",
		Run: func(cmd *cobra.Command, args []string) {

			// var unidades string
			// // var panels string
			// var types string

			// form := huh.NewForm(
			// 	huh.NewGroup(
			// 		huh.NewSelect[string]().
			// 			Options(huh.NewOptions("United States", "Canada", "Mexico")...).
			// 			Value(&unidades).
			// 			Title("Unidades"),

			// 		huh.NewSelect[string]().
			// 			Height(8).
			// 			Title("Paineis").
			// 			OptionsFunc(func() []huh.Option[string] {
			// 				opts := []string{
			// 					"painel 1",
			// 					"painel 2",
			// 					"painel 3"}
			// 				return huh.NewOptions(opts...)
			// 			}, &unidades),

			// 		huh.NewSelect[string]().
			// 			Height(8).
			// 			Title("Tipos").
			// 			OptionsFunc(func() []huh.Option[string] {
			// 				opts := []string{
			// 					"tipo1", "tipo2"}
			// 				return huh.NewOptions(opts...)
			// 			}, &types),

			// 		huh.NewConfirm().
			// 			Title("Confirma inclusão?"),
			// 	),
			// )

			p := tea.NewProgram(NewModel())
			_, err := p.Run()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Done")
		},
	}
}
