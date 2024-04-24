package main

import (
	//"github.com/charmbracelet/bubbles/help"
	"fmt"

	//"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"

	//"github.com/charmbracelet/bubbles/key"
	//"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	//"github.com/charmbracelet/lipgloss"
)

type TenseModel struct {
    title string
    options string
    //options list.Model
}

func initialTenseModel() *TenseModel {
    model := TenseModel{
        title: "Tense",
        options: "Preterite",
    }
//
//    width := 20
//    height := 20
//
//    items := []list.Item{
//        item
//
//    }
//
//    model := LanguageModel{
//        options: list.New(items, list.NewDefaultDelegate(), width, height),
//    }
    return &model
}

func (m TenseModel) Init() tea.Cmd {
    return nil
}

func (m TenseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
		switch msg.String() {
        case "enter":
            return m, nil
		}
	}

    return m, nil
}

func (m TenseModel) View() string {
    string := fmt.Sprintf("%s\n%s", m.title, m.options)
    return lipgloss.NewStyle().
        Width(20).Height(4).
        Border(lipgloss.RoundedBorder()).
        Render(string)
}
