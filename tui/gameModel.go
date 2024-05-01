package tui

import (
	"conju/utils"
	"fmt"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strconv"
	"time"
)

type GameModel struct {
	verbs      []map[string]string
	language   string
	tense      string
	timer      timer.Model
	round      tea.Model
	count      int
	loaded     bool
	completed  bool
	help       HelpModel
	keys       keyMap
	selectedDb utils.Database
	config     utils.Config
}

func newGameModel(selectedDb utils.Database, width int, game Game, config utils.Config) *GameModel {

	timeout := time.Duration(game.duration) * time.Minute
	timer := timer.NewWithInterval(timeout, time.Second)

	verbs := selectedDb.QueryData(game.tense)
	verb, pov, pronoun := setupRound(verbs, selectedDb, config)
	round := initialRoundModel(verb, pov, pronoun, config)
	help := NewHelpModel()
	help.Width = width

	model := GameModel{
		verbs:      verbs,
		language:   game.language,
		tense:      game.tense,
		timer:      timer,
		round:      round,
		count:      0,
		help:       help,
		keys:       gameKeys,
		selectedDb: selectedDb,
		config:     config,
	}

	return &model
}

func (m GameModel) Init() tea.Cmd {
	return nil
}

func (m GameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	var cmds []tea.Cmd

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
			verb, pov, pronoun := setupRound(m.verbs, m.selectedDb, m.config)
			m.round = initialRoundModel(verb, pov, pronoun, m.config)
			return m, cmd
		case "?":
			m.help.ShowAll = !m.help.ShowAll
			return m, cmd
		}
	}

	if m.timer.Timedout() {
		m.completed = true
		return m, cmd
	}
	if !m.loaded {
		m.loaded = true
		cmds = append(cmds, m.timer.Init())
	}

	newRound, newCmd := m.round.Update(msg)
	roundModel, ok := newRound.(RoundModel)
	if !ok {
		panic("Round Model Assertion Failed")
	}

	if roundModel.pass {
		m.count++
		verb, pov, pronoun := setupRound(m.verbs, m.selectedDb, m.config)
		m.round = *initialRoundModel(verb, pov, pronoun, m.config)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	m.round = roundModel
	cmd = newCmd
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m GameModel) View() string {

	var thirtySeconds int64 = 30_000_000

	timerStyled := func() (formatted string) {
		if m.timer.Timeout.Microseconds() < thirtySeconds {
			return lipgloss.NewStyle().
				Foreground(lipgloss.Color("1")).
				Render(m.timer.View())
		}

		return m.timer.View()
	}

	helpView := helpStyle.Render(m.help.View(m.keys))

	applyStyling := func(childElement string) (formatted string) {

		return lipgloss.NewStyle().Render(childElement)
	}

	output := fmt.Sprintf("%s\n%s\n%s\n%s\nCount %s",
		m.language, m.tense, timerStyled(), m.round.View(), strconv.Itoa(m.count))

	return applyStyling(output + "\n" + helpView)
}

func setupRound(
	verbs []map[string]string,
	selectedDb utils.Database,
	config utils.Config) (
	verb map[string]string,
	pov, pronoun string) {

	var povs []string

	defaultPronouns := config.Languages[selectedDb.LowerName].DefaultConjugations

	for pronoun, addPronoun := range defaultPronouns {
		if addPronoun {
			povs = append(povs, pronoun)
		}
	}

	verb, pov, pronoun = utils.ChooseVerb(verbs, defaultPronouns, selectedDb)
	return
}
