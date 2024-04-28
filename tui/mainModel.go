package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type mainSessionState int

const (
	settingView mainSessionState = iota
	gameView
	performanceView
)

type Game struct {
	language string
	tense    string
	duration int
}

type MainModel struct {
	width        int
	height       int
	state        mainSessionState
	setting      tea.Model
	gameSettings Game
	game         tea.Model
	count        int
	performance  tea.Model
	loaded       bool
	quitting     bool
}

func InitialMainModel() *MainModel {
	model := MainModel{
		state:   settingView,
		setting: NewSettingsModel(),
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
		m.width = msg.Width
		m.height = msg.Height
		var cmd tea.Cmd
		var cmds []tea.Cmd
		cmds = append(cmds, cmd)
		m.loaded = true
		return m, tea.Batch(cmds...)
	}

	switch m.state {
	case settingView:
		newSetting, newCmd := m.setting.Update(msg)
		settingModel, ok := newSetting.(SettingModel)

		if !ok {
			panic("Language Model assertion failed")
		}

		if settingModel.selectedConfirm {
			m.gameSettings = Game{
				language: settingModel.selectedLanguage,
				tense:    settingModel.selectedTense,
				duration: settingModel.selectedDuration,
			}
			m.game = *newGameModel(m.gameSettings)
			m.state = gameView
		}

		m.setting = settingModel
		cmd = newCmd

	case gameView:
		newGame, newCmd := m.game.Update(msg)
		gameModel, ok := newGame.(GameModel)

		if !ok {
			panic("Game Model assertion failed")
		}

		if gameModel.completed {
			m.count = gameModel.count
			m.performance = *initialPerformanceModel(m.gameSettings, m.count)
			m.state = performanceView
		}

		m.game = gameModel
		cmd = newCmd

	case performanceView:
		newPerformance, newCmd := m.performance.Update(msg)
		performanceModel, ok := newPerformance.(PerformanceModel)

		if !ok {
			panic("Performance Model assertion failed")
		}

		m.performance = performanceModel
		cmd = newCmd
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

	title := ("Conju - Language Conjugation App")
    header := lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(lipgloss.Color("8")).
			Render(title)

	buildView := func(div string) string {
		dialogBoxStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("6")).
			Padding(1, 2).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)

		return lipgloss.Place(m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			dialogBoxStyle.Render(div),
			lipgloss.WithWhitespaceChars("充电"),
			lipgloss.WithWhitespaceForeground(lipgloss.Color("8")),
		)

	}

	switch m.state {
	case settingView:
		return buildView(header + "\n" + m.setting.View())
	case gameView:
		return buildView(header + "\n" + m.game.View())
	case performanceView:
		return buildView(header + "\n" + m.performance.View())
	}
	return ""
}
