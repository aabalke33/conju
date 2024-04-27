package main

import (
	"conju/utils"
	"fmt"
	"strconv"
	"time"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type GameModel struct {
    verbs       []map[string]string
    language    string
    tense       string
    timer       timer.Model
    round       tea.Model
    count       int
    loaded      bool
    completed   bool
}

func initialGameModel(game Game) *GameModel {

    timeout := time.Duration(game.duration) * time.Minute
    timer := timer.NewWithInterval(timeout, time.Second)

    verbs := utils.QueryData(game.language, game.tense)
    verb, pov, pronoun := setupRound(verbs)
    round := initialRoundModel(verb, pov, pronoun)

    model := GameModel{
        verbs: verbs,
        language: game.language,
        tense: game.tense,
        timer: timer, 
        round: round,
        count: 0,
    }

    return &model
}

func (m GameModel) Init() tea.Cmd {
    return nil
}

func (m GameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    var cmds []tea.Cmd

    m.timer, cmd = m.timer.Update(msg)
    cmds = append(cmds, cmd)


    switch msg := msg.(type) {
    case timer.TickMsg:
        m.timer, cmd = m.timer.Update(msg)
        return m, cmd

    case timer.StartStopMsg:
        m.timer, cmd = m.timer.Update(msg)
        return m, cmd

    case tea.KeyMsg:
        switch key := msg.String(); key {
        case "ctrl+c":
            return m, tea.Quit

        case "tab":
            verb, pov, pronoun := setupRound(m.verbs)
            m.round = initialRoundModel(verb, pov, pronoun)
            return m, cmd
        }
	}

    if m.timer.Timedout() {
        m.completed = true
        return m, cmd
    }

    newRound, newCmd := m.round.Update(msg)
    roundModel, ok := newRound.(RoundModel)
    if !ok {
        panic("Round Model Assertion Failed")
    }

    if roundModel.pass {
        m.count++
        verb, pov, pronoun := setupRound(m.verbs)
        m.round = *initialRoundModel(verb, pov, pronoun)
        cmds = append(cmds, cmd)
        return m, tea.Batch(cmds...)
    }

    m.round = roundModel
    cmd = newCmd
    cmds = append(cmds, cmd)
    return m, tea.Batch(cmds...)
}

func (m GameModel) View() string {
    mainContent := (
        "Conju - Language Conjugation App\n")

    applyStyling := func(childElement string) (formatted string) {
        return lipgloss.NewStyle().
            Width(40).Height(20).
            Border(lipgloss.RoundedBorder()).
            BorderForeground(lipgloss.Color("8")).
            Render(childElement)
    }

    output := fmt.Sprintf("%s\n%s\n%s\n%s\nCount %s",
        m.language, m.tense, m.timer.View(), m.round.View(), strconv.Itoa(m.count))

    return mainContent + applyStyling(output)
}

func setupRound(verbs []map[string]string) (verb map[string]string, pov, pronoun string) {

    var povs []string

    keepPronouns := map[string]bool{
        "first_single": true,
        "first_plural": true,
        "second_single": true,
        "second_plural": false,
        "third_single": true,
        "third_plural": true,
    }

    for pronoun, keepPronoun := range keepPronouns {
        if keepPronoun { povs = append(povs, pronoun) }
    }

    verb, pov, pronoun = utils.ChooseVerb(verbs, keepPronouns)
    return
}
