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
	width             int
	height            int
	panels            []panel
	panelHeights      []int
	focusedPanel      Panel
	activeSourcePanel Panel
	theme             Theme
	themeNames        []string
	themeIndex        int
	help              help.Model
	helpViewport      viewport.Model
	helpContent       string
	showHelp          bool
	git               *git.GitCommands
	repoName          string
	branchName        string
}

// initialModel creates the initial state of the application.
func initialModel() Model {
	themeNames := ThemeNames()
	gc := git.NewGitCommands()
	repoName, branchName, _ := gc.GetRepoInfo()
	initialContent := initialContentLoading

	// Create a slice to hold all UI panels.
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
		theme:             Themes[themeNames[0]],
		themeNames:        themeNames,
		themeIndex:        0,
		focusedPanel:      StatusPanel,
		activeSourcePanel: StatusPanel,
		help:              help.New(),
		helpViewport:      viewport.New(0, 0),
		showHelp:          false,
		git:               gc,
		repoName:          repoName,
		branchName:        branchName,
		panels:            panels,
	}
}

// Init is the first command that is run when the program starts.
func (m Model) Init() tea.Cmd {
	// fetch initial content for all panels.
	return tea.Batch(
		m.fetchPanelContent(StatusPanel),
		m.fetchPanelContent(FilesPanel),
		m.fetchPanelContent(BranchesPanel),
		m.fetchPanelContent(CommitsPanel),
		m.fetchPanelContent(StashPanel),
		m.fetchPanelContent(SecondaryPanel),
		m.updateMainPanel(),
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
	case BranchesPanel:
		return keys.BranchesPanelHelp()
	case CommitsPanel:
		return keys.CommitsPanelHelp()
	case StashPanel:
		return keys.StashPanelHelp()
	default:
		return keys.ShortHelp()
	}
}
