package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type KindModel struct {
	title    string
	selected string
	options  list.Model
}

func initialKindModel() *KindModel {

	kinds := []string{"Vocabulary", "Conjugation"}

	var items []list.Item

	for _, kind := range kinds {
		items = append(items, Item(kind))
	}

	height := len(items)

	options := list.New(items, itemDelegate{}, 0, height)
	options.SetShowStatusBar(false)
	options.SetShowTitle(false)
	options.SetShowHelp(false)
	options.SetShowPagination(false)
	options.SetShowFilter(false)

	model := KindModel{
		title:    "Kind",
		options:  options,
		selected: "",
	}
	return &model
}

func (m KindModel) Init() tea.Cmd {
	return nil
}

func (m KindModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m KindModel) View() string {
	titleStyle := func(title string) (formatted string) {
		return lipgloss.NewStyle().Render(title)
	}

	return titleStyle(m.title) + "\n" + m.options.View()
}
