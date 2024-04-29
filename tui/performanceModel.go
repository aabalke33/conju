package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PerformanceModel struct {
	game     Game
	count    int
	wpm      int
	quitting bool
	help     HelpModel
	keys     keyMap
}

func initialPerformanceModel(width int, game Game, count int) *PerformanceModel {

	wpm := int(float32(count) / float32(game.duration))
	help := NewHelpModel()
	help.Width = width

	model := PerformanceModel{
		game:  game,
		count: count,
		wpm:   wpm,
		help:  help,
		keys:  performanceKeys,
	}
	return &model
}

func (m PerformanceModel) Init() tea.Cmd {
	return nil
}

func (m PerformanceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch key := msg.String(); key {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m PerformanceModel) View() string {
	if m.quitting {
		return ""
	}

	applyStyling := func(childElement string) (formatted string) {
		return lipgloss.NewStyle().
			//BorderForeground(lipgloss.Color("8")).
			Render(childElement)
	}
	helpView := helpStyle.Render(m.help.View(m.keys))

	output := fmt.Sprintf(
		"Completed.\n%s - %s Test.\n%d Minutes\n%d Answered\n%d Per Minute",
		m.game.language, m.game.tense, m.game.duration, m.count, m.wpm)

	return applyStyling(output + "\n" + helpView)
}
