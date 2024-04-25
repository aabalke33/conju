package main

import (
	//"fmt"
	"fmt"
	"strconv"

	//"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	//"github.com/charmbracelet/lipgloss"
	//"github.com/charmbracelet/bubbles/help"
	//"github.com/charmbracelet/bubbles/key"
)

type ConfirmModel struct {
    confirmed bool
    language string
    tense string
    duration int
}

func initialConfirmModel() *ConfirmModel {

    model := ConfirmModel{
        language: "",
        tense: "",
        duration: 0,
    }
    return &model
}

func (m ConfirmModel) Init() tea.Cmd {
    return nil
}

func (m ConfirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    switch msg := msg.(type) {
    case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
        case "enter":
            m.confirmed = !m.confirmed
            return m, nil
		}
	}


	//m.options, cmd = m.options.Update(msg)
    return m, cmd
}

func (m ConfirmModel) View() string {

//    titleStyle := func (title string) (formatted string) {
//        return lipgloss.NewStyle().
//            Padding(0, 1).
//            //Background(lipgloss.Color("6")).
//            Foreground(lipgloss.Color("10")).
//            Render(title)
//    }



    rows := [][]string{
        {"Language",    m.language},
        {"Tense",       m.tense},
        {"Duration",    strconv.Itoa(m.duration) + " minutes"},
    }

    t := table.New().Rows(rows...).Border(lipgloss.HiddenBorder())

    return fmt.Sprintf("Are these settings correct?\n%s\nPress enter to confirm.", t)
}
