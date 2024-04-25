package main

import (
	tea "github.com/charmbracelet/bubbletea"
	//"github.com/charmbracelet/lipgloss"
)

type mainSessionState int

const (
    settingView mainSessionState = iota
    gameView
)

type MainModel struct {
    state mainSessionState
    setting tea.Model
    confirmed bool
    //game tea.Model
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
    switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		var cmd tea.Cmd
		var cmds []tea.Cmd
		cmds = append(cmds, cmd)
		m.loaded = true
		return m, tea.Batch(cmds...)
    case tea.KeyMsg:
        switch key := msg.String(); key {
        case "q":
			m.quitting = true
			return m, tea.Quit
		}
	}

    switch m.state {
    case settingView:
        newSetting, newCmd := m.setting.Update(msg)
        newSettingModel, ok := newSetting.(SettingModel)

        if !ok {
            panic("Language Model assertion failed")
        }

        if newSettingModel.selectedConfirm != m.confirmed {
            m.confirmed = newSettingModel.selectedConfirm
            m.state++
        }

        m.setting = newSettingModel
        cmd = newCmd

    case gameView:
//        newTense, newCmd := m.tense.Update(msg)
//        newTenseModel, ok := newTense.(TenseModel)
//        
//        if !ok {
//            panic("Tense Model assertion failed")
//        }
//
//        if newTenseModel.selected != m.selectedTense {
//            m.selectedTense = newTenseModel.selected
//            m.state++
//        } else {
//            switch msg := msg.(type) {
//            case tea.KeyMsg:
//                switch key := msg.String(); key {
//                case "enter":
//                    m.state++
//        }
//
//        m.tense = newTenseModel
//        cmd = newCmd
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
    case settingView:  return m.setting.View()
    //case gameView:   return mainContent + applyStyling(m.confirm.View())
    }
    return ""
}
