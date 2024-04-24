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

type LanguageModel struct {
    title string
    options string
    //options list.Model
}

func initialLanguageModel() *LanguageModel {
    model := LanguageModel{
        title: "Language",
        options: "Spanish",
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

func (m LanguageModel) Init() tea.Cmd {
    return nil
}

func (m LanguageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
		switch msg.String() {
        case "enter":
            return m, nil
		}
	}

    return m, nil
}

func (m LanguageModel) View() string {
    string := fmt.Sprintf("%s\n%s", m.title, m.options)
    return lipgloss.NewStyle().
        Width(20).Height(4).
        Border(lipgloss.RoundedBorder()).
        Render(string)
}
