package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sessionState int

const (
    languageView sessionState = iota
    tenseView
    durationView
    confirmView
)

type SettingModel struct {
    state sessionState
    language tea.Model
    tense tea.Model
    duration tea.Model
    confirm tea.Model
    selectedLanguage string
    selectedTense string
    selectedDuration int
    selectedConfirm bool
}

func initialSettingModel() *SettingModel {

    language := initialLanguageModel()
    tense := initialTenseModel()
    duration := initialDurationModel()
    confirm := initialConfirmModel()

    model := SettingModel{
        state: languageView,
        language: language,
        tense: tense,
        duration: duration,
        confirm: confirm,
    }

    return &model
}

func (m SettingModel) Init() tea.Cmd {
    return nil
}

func (m SettingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    var cmds []tea.Cmd
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch key := msg.String(); key {
        case "shift+tab":
            if m.state > languageView {
                m.state--
            }
            return m, nil
        case "tab":
            if m.state < confirmView {
                m.state++
            }
            return m, nil
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
        newDuration, newCmd := m.duration.Update(msg)
        newDurationModel, ok := newDuration.(DurationModel)

        if !ok {
            panic("Duration Model assertion failed")
        }

        if newDurationModel.value != m.selectedDuration {
            m.selectedDuration = newDurationModel.value
            m.state++
        }

        m.duration = newDurationModel
        cmd = newCmd

    case confirmView:
        newConfirm, newCmd := m.confirm.Update(msg)
        newConfirmModel, ok := newConfirm.(ConfirmModel)

        if !ok {
            panic("Confirmation Model assertion failed")
        }

        if newConfirmModel.confirmed != m.selectedConfirm {
            m.selectedConfirm = newConfirmModel.confirmed
            m.state++
        }

        //May be able to move these
        if newConfirmModel.language != m.selectedLanguage {
            newConfirmModel.language = m.selectedLanguage
        }
        if newConfirmModel.tense != m.selectedTense {
            newConfirmModel.tense = m.selectedTense
        }
        if newConfirmModel.duration != m.selectedDuration {
            newConfirmModel.duration = m.selectedDuration
        }

        m.confirm = newConfirmModel
        cmd = newCmd
    }


    cmds = append(cmds, cmd)
    return m, tea.Batch(cmds...)
}

func (m SettingModel) View() string {
    mainContent := (
        "Conju - Language Conjugation App\n")

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
    case durationView:  return mainContent + applyStyling(m.duration.View())
    case confirmView:   return mainContent + applyStyling(m.confirm.View())
    }
    return ""
}
