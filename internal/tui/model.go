package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gitxtui/gitx/internal/git"
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
	git          *git.GitCommands
	panels       []panel
	panelHeights []int
}

// initialModel creates the initial state of the application.
func initialModel() Model {
	themeNames := ThemeNames()
	gc := git.NewGitCommands()
	initialContent := "Loading..."

	// Create a slice to hold all our panels.
	panels := make([]panel, totalPanels)
	for i := range panels {
		vp := viewport.New(0, 0)
		vp.SetContent(initialContent)
		panels[i] = panel{
			viewport: vp,
			content:  initialContent,
		}
	}

	return Model{
		theme:        Themes[themeNames[0]],
		themeNames:   themeNames,
		themeIndex:   0,
		focusedPanel: StatusPanel,
		help:         help.New(),
		helpViewport: viewport.New(0, 0),
		showHelp:     false,
		git:          gc,
		panels:       panels,
	}
}

// Init is the first command that is run when the program starts.
func (m Model) Init() tea.Cmd {
	// fetch initial content for all panels.
	return tea.Batch(
		fetchPanelContent(m.git, StatusPanel),
		fetchPanelContent(m.git, FilesPanel),
		fetchPanelContent(m.git, BranchesPanel),
		fetchPanelContent(m.git, CommitsPanel),
		fetchPanelContent(m.git, StashPanel),
		fetchPanelContent(m.git, MainPanel),
		fetchPanelContent(m.git, SecondaryPanel),
	)
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
