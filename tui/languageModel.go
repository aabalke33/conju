package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	//"github.com/charmbracelet/bubbles/help"
	//"github.com/charmbracelet/bubbles/key"
)

type LanguageModel struct {
	title    string
	selected string
	options  list.Model
}

func initialLanguageModel() *LanguageModel {

	height := 10

	items := []list.Item{
		Item("Spanish"),
		Item("German"),
		Item("Japanese"),
		Item("Russian"),
		Item("Chinese"),
	}
	options := list.New(items, itemDelegate{}, 0, height)
	options.SetShowStatusBar(false)
	options.SetShowTitle(false)
	options.SetShowHelp(false)

	model := LanguageModel{
		title:    "Language",
		options:  options,
		selected: "",
	}
	return &model
}

func (m LanguageModel) Init() tea.Cmd {
	return nil
}

func (m LanguageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.options.SetWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			i, ok := m.options.SelectedItem().(Item)
			if ok {
				m.selected = string(i)
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.options, cmd = m.options.Update(msg)
	return m, cmd
}

func (m LanguageModel) View() string {

	titleStyle := func(title string) (formatted string) {
		return lipgloss.NewStyle().
			//Padding(0, 1).
			//Background(lipgloss.Color("6")).
			Foreground(lipgloss.Color("10")).
			Render(title)
	}

	return titleStyle(m.title) + m.options.View()
}