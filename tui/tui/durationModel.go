package tui

import (
	//"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	//"github.com/charmbracelet/bubbles/help"
	//"github.com/charmbracelet/bubbles/key"
)

type DurationModel struct {
	title string
	value int
	input textinput.Model
}

func initialDurationModel() *DurationModel {

	input := textinput.New()
	input.Placeholder = "Minutes"
	input.Focus()
	input.CharLimit = 2
	input.Width = 20

	model := DurationModel{
		title: "Duration",
		input: input,
		value: 0,
	}
	return &model
}

func (m DurationModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m DurationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "backspace":
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		case "enter":
			s := m.input.Value()

			if s == "" || s == "0" {
				return m, nil
			}

			i, err := strconv.Atoi(s)
			if err != nil {
				panic("Could not parse duration")
			}

			m.value = i
			return m, nil
		}
	}

	return m, cmd
}

func (m DurationModel) View() string {

	titleStyle := func(title string) (formatted string) {
		return lipgloss.NewStyle().
			Padding(0, 1).
			//Background(lipgloss.Color("6")).
			Foreground(lipgloss.Color("10")).
			Render(title)
	}

	return titleStyle(m.title) + "\n" + m.input.View()
}
