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
	kind    string
	pov     string
	pronoun string
	pass    bool
	fail    bool
	input   textinput.Model
	config  utils.Config
}

func initialRoundModel(
	verb map[string]string,
	pov,
	pronoun string,
	config utils.Config,
	kind string,
) *RoundModel {

	input := textinput.New()
	input.Focus()
	input.CharLimit = 32
	input.Width = 20

	model := RoundModel{
		verb:    verb,
		kind:    kind,
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

	handlePass := func() {
		utils.PlayAudio("./utils/resources/pass.mp3", m.config)
		m.pass = true
	}
	handleFail := func() {
		utils.PlayAudio("./utils/resources/fail.mp3", m.config)
		m.fail = true
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key := msg.String(); key == "enter" {
			switch m.kind {
			case "Vocabulary":
				if match := m.input.Value() == m.verb["infinitive"]; match {
					handlePass()
				} else {
					handleFail()
					return m, cmd
				}
			default:
				if match := m.input.Value() == m.verb[m.pov]; match {
					handlePass()
				} else {
					handleFail()
					return m, cmd
				}
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

	switch m.kind {
	case "Vocabulary":
		output := fmt.Sprintf("\n%s\n", m.verb["meaning"])
		return output + applyStyling(m.input.View())
	default:
		output := fmt.Sprintf("\n%s %s\n", m.pronoun, m.verb["infinitive"])
		return output + applyStyling(m.input.View())
	}
}
