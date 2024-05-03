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
	povs       []string
	pronouns   map[string][]string
}

func newGameModel(selectedDb utils.Database, width int, game Game, config utils.Config) *GameModel {

	timer := setupTimer(game.duration)
	verbs, povs, pronouns := setupGame(selectedDb, config, game.tense)
	verb, pov, pronoun := utils.ChooseVerb(verbs, pronouns)

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
		povs:       povs,
		pronouns:   pronouns,
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
			verb, pov, pronoun := utils.ChooseVerb(m.verbs, m.pronouns)
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
	cmds = append(cmds, newCmd)

	if roundModel.pass {
		m.count++
		verb, pov, pronoun := utils.ChooseVerb(m.verbs, m.pronouns)
		m.round = *initialRoundModel(verb, pov, pronoun, m.config)
		return m, tea.Batch(cmds...)
	}

	m.round = roundModel
	return m, tea.Batch(cmds...)
}

func (m GameModel) View() string {

	var thirtySeconds int64 = 30_000_000

	timerStyled := func() (formatted string) {

		timeRemaining := m.timer.Timeout.Microseconds()
		if timeRemaining < thirtySeconds {
			return lipgloss.NewStyle().
				Foreground(lipgloss.Color("1")).
				Render(m.timer.View())
		}

		return m.timer.View()
	}

	helpView := helpStyle.Render(m.help.View(m.keys))

	output := fmt.Sprintf("%s\n%s\n%s\n%s\nCount %s",
		m.language, m.tense, timerStyled(), m.round.View(), strconv.Itoa(m.count))

	return output + "\n" + helpView
}

func setupTimer(duration int) timer.Model {
	timeout := time.Duration(duration) * time.Minute
	timer := timer.NewWithInterval(timeout, time.Second)
	return timer
}

func setupGame(selectedDb utils.Database, config utils.Config, tense string) (
	[]map[string]string, []string, map[string][]string) {

	verbs := selectedDb.QueryData(tense)
	defaultPronouns := config.Languages[selectedDb.LowerName].DefaultConjugations
	povs, pronouns := utils.ChoosePronouns(defaultPronouns, selectedDb, tense)
	return verbs, povs, pronouns
}
