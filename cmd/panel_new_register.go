package cmd

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var (
	burger       string
	toppings     []string
	sauceLevel   int
	name         string
	instructions string
	discount     bool
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
	var unidades string
	// var panels string
	var types string
	return Model{
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Key("units").
					Options(huh.NewOptions("United States", "Canada", "Mexico")...).
					Value(&unidades).
					Title("Unidades"),

				huh.NewSelect[string]().
					Height(8).
					Title("Paineis").
					Key("panels").
					OptionsFunc(func() []huh.Option[string] {
						opts := []string{
							"painel 1",
							"painel 2",
							"painel 3"}
						return huh.NewOptions(opts...)
					}, &unidades),

				huh.NewSelect[string]().
					Height(8).
					Title("Tipos").
					Key("types").
					OptionsFunc(func() []huh.Option[string] {
						opts := []string{
							"tipo1", "tipo2"}
						return huh.NewOptions(opts...)
					}, &types),

				huh.NewConfirm().
					Title("Confirma inclusão?"),
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
		units := m.form.GetString("units")
		panels := m.form.GetString("panels")
		types := m.form.GetString("types")

		return fmt.Sprintf("panel added %s - %s - %s. Press enter for exit", units, types, panels)
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
