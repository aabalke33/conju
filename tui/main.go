package main

import (
	"fmt"
	"os"
	tea "github.com/charmbracelet/bubbletea"
)

//var m *MainModel
var m *GameModel

func main() {
    //m = initialMainModel()
    m = initialGameModel(Game{"spanish", "present", 3})
    p := tea.NewProgram(m, tea.WithAltScreen())
    if _, err := p.Run(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
