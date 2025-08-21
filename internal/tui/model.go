package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
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
	helpViewport viewport.Model
	helpContent  string
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
		helpViewport: viewport.New(0, 0),
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

// panelShortHelp returns a slice of key.Binding for the focused Panel.
func (m *Model) panelShortHelp() []key.Binding {
	switch m.focusedPanel {
	case FilesPanel:
		return keys.FilesPanelHelp()
	// TODO: Add cases for rest of the Panels
	default:
		return keys.ShortHelp()
	}
}
