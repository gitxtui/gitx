package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gitxtui/gitx/internal/git"
)

// appMode defines the different operational modes of the TUI.
type appMode int

const (
	modeNormal appMode = iota
	modeInput
	modeConfirm
	modeCommit
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
	// New fields for pop-ups
	mode             appMode
	promptTitle      string
	confirmMessage   string
	textInput        textinput.Model
	descriptionInput textarea.Model
	inputCallback    func(string) tea.Cmd
	commitCallback   func(title, description string) tea.Cmd
	confirmCallback  func(bool) tea.Cmd
	// New fields for command history
	CommandHistory []string
	// Diff view mode: nil = auto (respects threshold), true = split, false = unified
	forcedDiffViewMode *bool
}

// initialModel creates the initial state of the application.
func initialModel() Model {
	themeNames := ThemeNames() //built-in themes load
	cfg, _ := load_config()

	var selectedThemeName string
	if t, ok := Themes[cfg.Theme]; ok {
		selectedThemeName = cfg.Theme
		_ = t // to avoid unused variable warning
	} else {
		if _, err := load_custom_theme(cfg.Theme); err == nil {
			selectedThemeName = cfg.Theme
		} else {
			//fallback
			selectedThemeName = themeNames[0]
		}
	}

	themeNames = ThemeNames() // reload

	gc := git.NewGitCommands()
	repoName, branchName, _ := gc.GetRepoInfo()
	initialContent := initialContentLoading

	panels := make([]panel, totalPanels)
	for i := range panels {
		vp := viewport.New(0, 0)
		vp.SetContent(initialContent)
		panels[i] = panel{
			viewport: vp,
			content:  initialContent,
		}
	}

	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 80

	ta := textarea.New()
	ta.Placeholder = "Enter commit description"
	ta.SetWidth(80)
	ta.SetHeight(5)

	historyVP := viewport.New(0, 0)
	historyVP.SetContent("Command history will appear here...")

	return Model{
		theme:              Themes[selectedThemeName],
		themeNames:         themeNames,
		themeIndex:         indexOf(themeNames, selectedThemeName),
		focusedPanel:       StatusPanel,
		activeSourcePanel:  StatusPanel,
		help:               help.New(),
		helpViewport:       viewport.New(0, 0),
		showHelp:           false,
		git:                gc,
		repoName:           repoName,
		branchName:         branchName,
		panels:             panels,
		mode:               modeNormal,
		textInput:          ti,
		descriptionInput:   ta,
		CommandHistory:     []string{},
		forcedDiffViewMode: nil,
	}
}

func indexOf(arr []string, val string) int {
	for i, s := range arr {
		if s == val {
			return i
		}
	}
	return 0
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

// toggleDiffView cycles through diff view modes: auto -> split -> unified -> auto
func (m *Model) toggleDiffView() {
	if m.forcedDiffViewMode == nil {
		// Currently auto - switch to split
		trueVal := true
		m.forcedDiffViewMode = &trueVal
	} else if *m.forcedDiffViewMode {
		// Currently split - switch to unified
		falseVal := false
		m.forcedDiffViewMode = &falseVal
	} else {
		// Currently unified - switch to auto
		m.forcedDiffViewMode = nil
	}
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
