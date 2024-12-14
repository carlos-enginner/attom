package cmd

import (
	"errors"
	"log"

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

func PanelNewRegister() *cobra.Command {
	return &cobra.Command{
		Use:   "new_panel",
		Short: "New Panel Register",
		Run: func(cmd *cobra.Command, args []string) {
			// m := &model{
			// 	flavors: []item{
			// 		{"orange", false},
			// 		{"strawberry", false},
			// 		{"watermelon", false},
			// 		{"apple", false},
			// 	},
			// 	adds: []item{
			// 		{"candies", false},
			// 		{"mixed colors", false},
			// 		{"strange glass shape", false},
			// 	},
			// }
			// p := tea.NewProgram(m)

			// if _, err := p.Run(); err != nil {
			// 	panic(fmt.Sprintf("failed to run program: %v", err))
			// }

			form := huh.NewForm(
				huh.NewGroup(
					// Ask the user for a base burger and toppings.
					huh.NewSelect[string]().
						Title("Choose your burger").
						Options(
							huh.NewOption("Charmburger Classic", "classic"),
							huh.NewOption("Chickwich", "chickwich"),
							huh.NewOption("Fishburger", "fishburger"),
							huh.NewOption("Charmpossible™ Burger", "charmpossible"),
						).
						Value(&burger), // store the chosen option in the "burger" variable

					// Let the user select multiple toppings.
					huh.NewMultiSelect[string]().
						Title("Toppings").
						Options(
							huh.NewOption("Lettuce", "lettuce").Selected(true),
							huh.NewOption("Tomatoes", "tomatoes").Selected(true),
							huh.NewOption("Jalapeños", "jalapeños"),
							huh.NewOption("Cheese", "cheese"),
							huh.NewOption("Vegan Cheese", "vegan cheese"),
							huh.NewOption("Nutella", "nutella"),
						).
						Limit(4). // there’s a 4 topping limit!
						Value(&toppings),

					// Option values in selects and multi selects can be any type you
					// want. We’ve been recording strings above, but here we’ll store
					// answers as integers. Note the generic "[int]" directive below.
					huh.NewSelect[int]().
						Title("How much Charm Sauce do you want?").
						Options(
							huh.NewOption("None", 0),
							huh.NewOption("A little", 1),
							huh.NewOption("A lot", 2),
						).
						Value(&sauceLevel),
				),

				// Gather some final details about the order.
				huh.NewGroup(
					huh.NewInput().
						Title("What’s your name?").
						Value(&name).
						// Validating fields is easy. The form will mark erroneous fields
						// and display error messages accordingly.
						Validate(func(str string) error {
							if str == "Frank" {
								return errors.New("Sorry, we don’t serve customers named Frank.")
							}
							return nil
						}),

					huh.NewText().
						Title("Special Instructions").
						CharLimit(400).
						Value(&instructions),

					huh.NewConfirm().
						Title("Would you like 15% off?").
						Value(&discount),
				),
			)

			err := form.Run()
			if err != nil {
				log.Fatal(err)
			}

		},
	}
}
