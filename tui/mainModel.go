package main

import (
	//"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
	//"github.com/charmbracelet/bubbles/key"
	//"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	//"github.com/charmbracelet/lipgloss"
)

type sessionState int

const (
    languageView sessionState = iota
    tenseView
    durationView
)

type MainModel struct {
    state sessionState
    language tea.Model
    tense tea.Model
    duration tea.Model
    loaded bool
    quitting bool
}

func initialMainModel() *MainModel {

    language := initialLanguageModel()
    tense := initialTenseModel()

    model := MainModel{
        state: languageView,
        language: language,
        tense: tense,
    }
    return &model
}

func (m MainModel) Init() tea.Cmd {
    return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    var cmds []tea.Cmd
    switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		var cmd tea.Cmd
		var cmds []tea.Cmd
		cmds = append(cmds, cmd)
		m.loaded = true
		return m, tea.Batch(cmds...)
    case tea.KeyMsg:
		switch msg.String() {
        case "h", "shift+tab":
            if m.state > languageView {
                m.state--
            }
            return m, nil
        case "l", "tab":
            if m.state < durationView {
                m.state++
            }
            return m, nil
        case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
	}

    switch m.state {
    case languageView:  {
        _, newCmd := m.language.Update(msg)
        m.state++
        cmd = newCmd
    }
    case tenseView: {
        _, newCmd := m.tense.Update(msg)
        m.state++
        cmd = newCmd
    }
    case durationView:  return m, nil
    }
    cmds = append(cmds, cmd)
    return m, tea.Batch(cmds...)
}

func (m MainModel) View() string {

    if m.quitting {
        return ""
    }
    if !m.loaded {
        return "loading..."
    }

    switch m.state {
    //case languageView:  return m.language.View()
    case languageView:  return m.language.View()
    case tenseView:     return m.tense.View()
    case durationView:  return lipgloss.NewStyle().Width(20).Height(4).Border(lipgloss.RoundedBorder()).Render("Conju\nDuration")
    }
    return ""
}
