package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type KeyMap interface {
	ShortHelp() []key.Binding
	FullHelp() [][]key.Binding
}

type Styles struct {
	Ellipsis       lipgloss.Style
	ShortKey       lipgloss.Style
	ShortDesc      lipgloss.Style
	ShortSeparator lipgloss.Style
	FullKey        lipgloss.Style
	FullDesc       lipgloss.Style
	FullSeparator  lipgloss.Style
}

type HelpModel struct {
	Width          int
	ShowAll        bool
	ShortSeparator string
	FullSeparator  string
	Ellipsis       string
	Styles         Styles
}

func NewHelpModel() HelpModel {
	sharedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	return HelpModel{
		//ShortSeparator: " • ",
		ShortSeparator: " ",
		FullSeparator:  " ",
		//FullSeparator:  "    ",
		Ellipsis: "…",
		Styles: Styles{
			ShortKey:       sharedStyle,
			ShortDesc:      sharedStyle,
			ShortSeparator: sharedStyle,
			Ellipsis:       sharedStyle.Copy(),
			FullKey:        sharedStyle.Copy(),
			FullDesc:       sharedStyle.Copy(),
			FullSeparator:  sharedStyle.Copy(),
		},
	}
}

func (m HelpModel) Update(_ tea.Msg) (HelpModel, tea.Cmd) {
	return m, nil
}

func (m HelpModel) View(k KeyMap) string {
	if m.ShowAll {
		return m.FullHelpView(k.FullHelp())
	}
	return m.ShortHelpView(k.ShortHelp())
}

func (m HelpModel) ShortHelpView(bindings []key.Binding) string {
	if len(bindings) == 0 {
		return ""
	}

	var b strings.Builder
	var totalWidth int
	var separator = m.Styles.ShortSeparator.Inline(true).Render(m.ShortSeparator)

	for i, kb := range bindings {
		if !kb.Enabled() {
			continue
		}

		var sep string
		if totalWidth > 0 && i < len(bindings) {
			sep = separator
		}

		str := sep +
			m.Styles.ShortKey.Inline(true).Render(kb.Help().Key) + " " +
			m.Styles.ShortDesc.Inline(true).Render(kb.Help().Desc)

		w := lipgloss.Width(str)

		if m.Width > 0 && totalWidth+w > m.Width {
			tail := " " + m.Styles.Ellipsis.Inline(true).Render(m.Ellipsis)
			tailWidth := lipgloss.Width(tail)

			if totalWidth+tailWidth < m.Width {
				b.WriteString(tail)
			}

			break
		}

		totalWidth += w
		b.WriteString(str)
	}

	return b.String()
}

func (m HelpModel) FullHelpView(groups [][]key.Binding) string {
	if len(groups) == 0 {
		return ""
	}

	var (
		out []string

		totalWidth int
		sep        = m.Styles.FullSeparator.Render(m.FullSeparator)
		sepWidth   = lipgloss.Width(sep)
	)

	for i, group := range groups {
		if group == nil || !shouldRenderColumn(group) {
			continue
		}

		var (
			keys         []string
			descriptions []string
		)

		for _, kb := range group {
			if !kb.Enabled() {
				continue
			}
			keys = append(keys, kb.Help().Key)
			descriptions = append(descriptions, kb.Help().Desc)
		}

		col := lipgloss.JoinHorizontal(lipgloss.Top,
			m.Styles.FullKey.Render(strings.Join(keys, "\n")),
			m.Styles.FullKey.Render(" "),
			m.Styles.FullDesc.Render(strings.Join(descriptions, "\n")),
		)

		totalWidth += lipgloss.Width(col)
		if m.Width > 0 && totalWidth > m.Width {
			break
		}

		out = append(out, col)

		if i < len(group)-1 {
			totalWidth += sepWidth
			if m.Width > 0 && totalWidth > m.Width {
				break
			}
		}

		out = append(out, sep)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, out...)
}

func shouldRenderColumn(b []key.Binding) (ok bool) {
	for _, v := range b {
		if v.Enabled() {
			return true
		}
	}
	return false
}

type keyMap struct {
	Enter key.Binding
	Tab   key.Binding
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Help  key.Binding
	Quit  key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Help,
		k.Quit,
	}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Enter, k.Tab},
		{k.Help, k.Quit},
	}
}

var settingKeys = keyMap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "confirm"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
var confirmKeys = keyMap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "confirm"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
var gameKeys = keyMap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "confirm"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "skip"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}

var performanceKeys = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
