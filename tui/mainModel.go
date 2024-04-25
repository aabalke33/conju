package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	//"github.com/charmbracelet/bubbles/help"
	//"github.com/charmbracelet/bubbles/key"
	//"github.com/charmbracelet/bubbles/list"
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
    selectedLanguage string
    tense tea.Model
    selectedTense string
    duration tea.Model
    SelectedDuration int
    loaded bool
    quitting bool
}

func initialMainModel() *MainModel {

    language := initialLanguageModel()
    tense := initialTenseModel()
    duration := initialDurationModel()

    model := MainModel{
        state: languageView,
        language: language,
        tense: tense,
        duration: duration,
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
        switch key := msg.String(); key {
        case "shift+tab":
            if m.state > languageView {
                m.state--
            }
            return m, nil
        case "tab":
            if m.state < durationView {
                m.state++
            }
            return m, nil
        case "q":
			m.quitting = true
			return m, tea.Quit
		}
	}

    switch m.state {
    case languageView:
        newLanguage, newCmd := m.language.Update(msg)
        newLanguageModel, ok := newLanguage.(LanguageModel)

        if !ok {
            panic("Language Model assertion failed")
        }

        if newLanguageModel.selected != m.selectedLanguage {
            m.selectedLanguage = newLanguageModel.selected
            m.state++
        }

        m.language = newLanguageModel
        cmd = newCmd

    case tenseView:
        newTense, newCmd := m.tense.Update(msg)
        newTenseModel, ok := newTense.(TenseModel)
        
        if !ok {
            panic("Tense Model assertion failed")
        }

        if newTenseModel.selected != m.selectedTense {
            m.selectedTense = newTenseModel.selected
            m.state++
        }

        m.tense = newTenseModel
        cmd = newCmd

    case durationView:
        return m, nil
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

    mainContent := (
        "Conju - Language Conjugation App" +
        "\n" +
        m.selectedLanguage +
        "\n" +
        m.selectedTense +
        "\n" +
        //string(m.selectedDuration) +
        "\n")

    applyStyling := func(childElement string) (formatted string) {
        return lipgloss.NewStyle().
            Width(40).Height(20).
            Border(lipgloss.RoundedBorder()).
            BorderForeground(lipgloss.Color("8")).
            Render(childElement)
    }

    switch m.state {
    case languageView:  return mainContent + applyStyling(m.language.View())
    case tenseView:     return mainContent + applyStyling(m.tense.View())
    case durationView:  return mainContent + applyStyling("Duration")
    }
    return ""
}
