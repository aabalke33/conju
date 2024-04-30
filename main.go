package main

import (
	 "conju/tui"
	 "fmt"
	 "os"
	
	 tea "github.com/charmbracelet/bubbletea"
)

var m *tui.MainModel

func main() {


	m = tui.InitialMainModel()
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
