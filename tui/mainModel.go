package main

import (
	//"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type mainSessionState int

const (
    settingView mainSessionState = iota
    gameView
    performanceView
)

type Game struct {
    language string
    tense string
    duration int
}

type MainModel struct {
    state mainSessionState
    setting tea.Model
    gameSettings Game 
    game tea.Model
    count int
    performance tea.Model
    loaded bool
    quitting bool
}

func initialMainModel() *MainModel {

    setting := initialSettingModel()

    model := MainModel{
        state: settingView,
        setting: setting,
    }
    return &model
}

func (m MainModel) Init() tea.Cmd {
    return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    var cmds []tea.Cmd
    switch msg.(type) {
	case tea.WindowSizeMsg:
		var cmd tea.Cmd
		var cmds []tea.Cmd
		cmds = append(cmds, cmd)
		m.loaded = true
		return m, tea.Batch(cmds...)
	}

    switch m.state {
    case settingView:
        newSetting, newCmd := m.setting.Update(msg)
        newSettingModel, ok := newSetting.(SettingModel)

        if !ok {
            panic("Language Model assertion failed")
        }

        if newSettingModel.selectedConfirm {

            m.gameSettings = Game{
                language: newSettingModel.selectedLanguage,
                tense: newSettingModel.selectedTense,
                duration: newSettingModel.selectedDuration,
            }
            m.game = initialGameModel(m.gameSettings)
            m.state = gameView

            m.setting = newSettingModel
            cmd = newCmd
            return m, cmd
        }

        m.setting = newSettingModel
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
        cmds = append(cmds, cmd)
        return m, tea.Batch(cmds...)


    case performanceView:
        newPerformance, newCmd := m.performance.Update(msg)
        newPerformanceModel, ok := newPerformance.(PerformanceModel)
        
        if !ok {
            panic("Performance Model assertion failed")
        }

        m.performance = newPerformanceModel
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

    switch m.state {
        case settingView:       return m.setting.View()
        case gameView:          return m.game.View()
        case performanceView:   return m.performance.View()
    }
    return ""
}
