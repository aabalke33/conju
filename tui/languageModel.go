package tui

import (
	"conju/utils"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type LanguageModel struct {
	title     string
	selected  Item
	fileMap   map[string]string
	options   list.Model
	databases map[string]utils.Database
}

func initialLanguageModel(directory string) *LanguageModel {

	databases := utils.GetDatabases(directory)

	var items []list.Item
	fileMap := make(map[string]string)

	for _, database := range databases {

		item := Item(database.ProperName)
		items = append(items, item)

		fileMap[database.ProperName] = database.Basename
	}

	height := len(databases)

	options := list.New(items, itemDelegate{}, 0, height)
	options.SetShowStatusBar(false)
	options.SetShowTitle(false)
	options.SetShowHelp(false)
	options.SetShowPagination(false)
	options.SetShowFilter(false)

	model := LanguageModel{
		title:     "Language",
		options:   options,
		selected:  "",
		fileMap:   fileMap,
		databases: databases,
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
				m.selected = i
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
		return lipgloss.NewStyle().Render(title)
	}

	return titleStyle(m.title) + "\n" + m.options.View()
}
