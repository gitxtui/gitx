package tui

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

// Model represents the state of the TUI.
type Model struct {
	width        int
	height       int
	theme        Theme
	themeNames   []string
	themeIndex   int
	focusedPanel Panel
	help         help.Model
	showHelp     bool
}

func initialModel() Model {
	themeNames := ThemeNames()
	return Model{
		theme:        Themes[themeNames[0]],
		themeNames:   themeNames,
		themeIndex:   0,
		focusedPanel: MainPanel,
		help:         help.New(),
		showHelp:     false,
	}
}

// Init is the first command that is run when the program starts.
func (m Model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

// nextTheme cycles to the next theme.
func (m *Model) nextTheme() {
	m.themeIndex = (m.themeIndex + 1) % len(m.themeNames)
	m.theme = Themes[m.themeNames[m.themeIndex]]
}
