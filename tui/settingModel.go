package tui

import (
	"conju/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sessionState int

const (
	languageView sessionState = iota
	kindView
	tenseView
	durationView
	confirmView
)

type SettingModel struct {
	state            sessionState
	language         tea.Model
	kind             tea.Model
	tense            tea.Model
	duration         tea.Model
	confirm          tea.Model
	selectedDb       utils.Database
	selectedLanguage string
	selectedKind     string
	selectedTense    string
	selectedDuration int
	selectedConfirm  bool
	help             HelpModel
	keys             keyMap
	quitting         bool
}

func NewSettingsModel(width int, config utils.Config) *SettingModel {

	language := initialLanguageModel(config.DatabaseDirectory)
	kind := initialKindModel()
	duration := initialDurationModel()
	help := NewHelpModel()
	help.Width = width

	model := SettingModel{
		state:    languageView,
		language: language,
		kind:     kind,
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
			m.state++
		}

		m.language = newLanguageModel
		cmd = newCmd

	case kindView:
		m.keys = settingKeys
		newKind, newCmd := m.kind.Update(msg)
		newKindModel, ok := newKind.(KindModel)

		if !ok {
			panic("Kind Model assertion failed")
		}

		if newKindModel.selected != m.selectedKind {
			m.selectedKind = newKindModel.selected
			m.tense = initialTenseModel(m.selectedDb)
			m.state++
		}

		m.kind = newKindModel
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
				m.selectedKind,
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
		confirmModel, ok := newConfirm.(ConfirmModel)

		if !ok {
			panic("Confirmation Model assertion failed")
		}

		if confirmModel.confirmed {
			m.selectedConfirm = confirmModel.confirmed
			m.state++
		}

		m.confirm = confirmModel
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
	case kindView:
		return applyStyling(m.kind.View())
	case tenseView:
		return applyStyling(m.tense.View())
	case durationView:
		return applyStyling(m.duration.View())
	case confirmView:
		return applyStyling(m.confirm.View())
	}
	return ""
}
