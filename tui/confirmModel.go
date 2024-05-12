package tui

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type ConfirmModel struct {
	confirmed bool
	language  string
	kind      string
	tense     string
	duration  int
}

func initialConfirmModel(language, kind string, tense string, duration int) *ConfirmModel {
	model := ConfirmModel{
		language: language,
		kind:     kind,
		tense:    tense,
		duration: duration,
	}
	return &model
}

func (m ConfirmModel) Init() tea.Cmd {
	return nil
}

func (m ConfirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			m.confirmed = !m.confirmed
			return m, nil
		}
	}

	return m, cmd
}

func (m ConfirmModel) View() string {

	rows := [][]string{
		{"Language", m.language},
		{"Kind", m.kind},
		{"Tense", m.tense},
		{"Duration", strconv.Itoa(m.duration) + " minutes"},
	}

	t := table.New().Rows(rows...).Border(lipgloss.HiddenBorder())

	enter := lipgloss.NewStyle().
		Foreground(lipgloss.Color("2")).
		Italic(true).Bold(true).
		Render("enter")

	return fmt.Sprintf("Are these settings correct?\n%s\nPress %s to confirm.", t, enter)
}
