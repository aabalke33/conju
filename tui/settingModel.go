package tui

import (
	"conju/utils"
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
	state            sessionState
	language         tea.Model
	tense            tea.Model
	duration         tea.Model
	confirm          tea.Model
	selectedLanguage string
	selectedDb       utils.Database
	selectedTense    string
	selectedDuration int
	selectedConfirm  bool
	help             HelpModel
	keys             keyMap
	quitting         bool
}

func NewSettingsModel(width int, config utils.Config) *SettingModel {

	language := initialLanguageModel(config.DatabaseDirectory)
	duration := initialDurationModel()
	help := NewHelpModel()
	help.Width = width

	model := SettingModel{
		state:    languageView,
		language: language,
		duration: duration,
		help:     help,
		keys:     settingKeys,
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
		case "?":
			m.help.ShowAll = !m.help.ShowAll
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
	}

	switch m.state {
	case languageView:
		m.keys = settingKeys
		newLanguage, newCmd := m.language.Update(msg)
		newLanguageModel, ok := newLanguage.(LanguageModel)

		if !ok {
			panic("Language Model assertion failed")
		}

		selectedItem := string(newLanguageModel.selected)

		if selectedItem != m.selectedLanguage {
			m.selectedLanguage = selectedItem
			m.selectedDb = newLanguageModel.databases[selectedItem]
			m.tense = initialTenseModel(m.selectedDb)
			m.state++
		}

		m.language = newLanguageModel
		cmd = newCmd

	case tenseView:
		m.keys = settingKeys
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
		m.keys = settingKeys
		newDuration, newCmd := m.duration.Update(msg)
		newDurationModel, ok := newDuration.(DurationModel)

		if !ok {
			panic("Duration Model assertion failed")
		}

		if newDurationModel.value != m.selectedDuration {
			m.selectedDuration = newDurationModel.value

			m.confirm = initialConfirmModel(
				m.selectedLanguage,
				m.selectedTense,
				m.selectedDuration,
			)

			m.state++
		}

		m.duration = newDurationModel
		cmd = newCmd

	case confirmView:
		m.keys = confirmKeys
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
	if m.quitting {
		return ""
	}

	applyStyling := func(childElement string) (formatted string) {
		helpView := helpStyle.Render(m.help.View(m.keys))
		return lipgloss.NewStyle().Render(childElement + "\n" + helpView)
	}

	switch m.state {
	case languageView:
		return applyStyling(m.language.View())
	case tenseView:
		return applyStyling(m.tense.View())
	case durationView:
		return applyStyling(m.duration.View())
	case confirmView:
		return applyStyling(m.confirm.View())
	}
	return ""
}
