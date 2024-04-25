package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var m *MainModel

func main() {
    m = initialMainModel()
    p := tea.NewProgram(m, tea.WithAltScreen())
    if _, err := p.Run(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
