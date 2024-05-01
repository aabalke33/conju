package tui

import (
	"conju/utils"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type RoundModel struct {
	verb    map[string]string
	pov     string
	pronoun string
	pass    bool
	fail    bool
	input   textinput.Model
	config  utils.Config
}

func initialRoundModel(verb map[string]string, pov, pronoun string, config utils.Config) *RoundModel {

	input := textinput.New()
	input.Focus()
	input.CharLimit = 32
	input.Width = 20

	model := RoundModel{
		verb:    verb,
		pov:     pov,
		pronoun: pronoun,
		input:   input,
		config:  config,
	}

	return &model
}

func (m RoundModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m RoundModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key := msg.String(); key == "enter" {
			if match := m.input.Value() == m.verb[m.pov]; match {
				utils.PlayAudio("./utils/resources/pass.mp3", m.config)
				m.pass = true
			} else {
				utils.PlayAudio("./utils/resources/fail.mp3", m.config)
				m.fail = true
				return m, cmd
			}
		}
	}

	m.input, cmd = m.input.Update(msg)

	return m, cmd
}

func (m RoundModel) View() string {

	applyStyling := func(childElement string) (formatted string) {

		if m.fail {
			return lipgloss.NewStyle().
				Width(30).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("1")).
				Render(childElement)
		}

		return lipgloss.NewStyle().
			Width(30).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("8")).
			Render(childElement)
	}

	output := fmt.Sprintf("\n%s %s\n", m.pronoun, m.verb["infinitive"])
	return output + applyStyling(m.input.View())
}
