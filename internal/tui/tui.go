package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gitxtui/gitx/internal/git"
)

type Panel int

const (
	StatusPanel Panel = iota
	FilesPanel
	BranchesPanel
	LogPanel
	DiffPanel
)

type Model struct {
	activePanel    Panel
	width          int
	height         int
	statusContent  string
	filesContent   string
	branchContent  string
	logContent     string
	diffContent    string
	selectedFile   int
	selectedBranch int
	selectedCommit int
	err            error
}

type App struct {
	program *tea.Program
}

func NewApp() *App {
	status, err := git.GetStatus()
	if err != nil {
		status = err.Error()
	}

	model := Model{
		activePanel:   StatusPanel,
		statusContent: status,
		filesContent:  "README.md\ngo.mod\ngo.sum\ninternal/\ncmd/",
		branchContent: "* master\n  develop\n  feature/new-ui\n  hotfix/bug-123",
		logContent:    "abc123 Initial commit\ndef456 Add git commands\nghi789 Implement TUI\njkl012 Update dependencies",
		diffContent:   "",
		err:           err,
	}

	program := tea.NewProgram(model, tea.WithAltScreen())
	return &App{program: program}
}

func (a *App) Run() error {
	_, err := a.program.Run()
	return err
}

func (m Model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "1":
			m.activePanel = StatusPanel
		case "2":
			m.activePanel = FilesPanel
		case "3":
			m.activePanel = BranchesPanel
		case "4":
			m.activePanel = LogPanel
		case "5":
			m.activePanel = DiffPanel

		case "tab":
			m.activePanel = Panel((int(m.activePanel) + 1) % 5)

		case "up", "k":
			switch m.activePanel {
			case FilesPanel:
				if m.selectedFile > 0 {
					m.selectedFile--
				}
			case BranchesPanel:
				if m.selectedBranch > 0 {
					m.selectedBranch--
				}
			case LogPanel:
				if m.selectedCommit > 0 {
					m.selectedCommit--
				}
			}

		case "down", "j":
			switch m.activePanel {
			case FilesPanel:
				files := strings.Split(m.filesContent, "\n")
				if m.selectedFile < len(files)-1 {
					m.selectedFile++
				}
			case BranchesPanel:
				branches := strings.Split(m.branchContent, "\n")
				if m.selectedBranch < len(branches)-1 {
					m.selectedBranch++
				}
			case LogPanel:
				commits := strings.Split(m.logContent, "\n")
				if m.selectedCommit < len(commits)-1 {
					m.selectedCommit++
				}
			}

		case "enter":
			switch m.activePanel {
			case FilesPanel:
				// Show diff for selected file
				m.diffContent = "diff --git a/selected_file b/selected_file\n--- a/selected_file\n+++ b/selected_file\n@@ -1,3 +1,4 @@\n line 1\n line 2\n+new line\n line 3"
				m.activePanel = DiffPanel
			case BranchesPanel:
				// Switch to selected branch
				branches := strings.Split(m.branchContent, "\n")
				if m.selectedBranch < len(branches) {
					selectedBranch := strings.TrimSpace(strings.TrimPrefix(branches[m.selectedBranch], "* "))
					m.statusContent = fmt.Sprintf("Switched to branch '%s'", selectedBranch)
				}
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	// Define styles
	activeStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("86")).
		Padding(0, 1)

	inactiveStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("86")).
		Padding(0, 1).
		MarginBottom(1)

	selectedStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("86")).
		Foreground(lipgloss.Color("0")).
		Bold(true)

	// Calculate panel dimensions
	panelWidth := (m.width - 6) / 2
	panelHeight := (m.height - 4) / 3

	// Create panels
	statusPanel := m.createPanel("Status", m.statusContent, StatusPanel, panelWidth, panelHeight, activeStyle, inactiveStyle, selectedStyle)
	filesPanel := m.createPanel("Files", m.filesContent, FilesPanel, panelWidth, panelHeight, activeStyle, inactiveStyle, selectedStyle)
	branchesPanel := m.createPanel("Branches", m.branchContent, BranchesPanel, panelWidth, panelHeight, activeStyle, inactiveStyle, selectedStyle)
	logPanel := m.createPanel("Log", m.logContent, LogPanel, panelWidth, panelHeight, activeStyle, inactiveStyle, selectedStyle)

	// Create diff panel (full width)
	diffPanelStyle := inactiveStyle
	if m.activePanel == DiffPanel {
		diffPanelStyle = activeStyle
	}
	diffTitle := titleStyle.Render("Diff")
	diffContent := m.diffContent
	if diffContent == "" {
		diffContent = "Select a file to view diff"
	}
	diffPanel := diffPanelStyle.
		Width(m.width - 4).
		Height(panelHeight - 8).
		Render(diffTitle + "\n" + diffContent)

	// Arrange panels in grid layout
	topRow := lipgloss.JoinHorizontal(lipgloss.Top, statusPanel, filesPanel)
	middleRow := lipgloss.JoinHorizontal(lipgloss.Top, branchesPanel, logPanel)

	// Create help text
	helpText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("Navigation: 1-5 (panels) | Tab (next panel) | ↑↓/jk (select) | Enter (action) | q/Ctrl+C (quit)")

	// Combine all elements
	result := lipgloss.JoinVertical(lipgloss.Left,
		topRow,
		middleRow,
		diffPanel,
		helpText,
	)

	return result
}

func (m Model) createPanel(title, content string, panelType Panel, width, height int, activeStyle, inactiveStyle, selectedStyle lipgloss.Style) string {
	style := inactiveStyle
	if m.activePanel == panelType {
		style = activeStyle
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("86")).
		Padding(0, 1).
		MarginBottom(1)

	if m.activePanel != panelType {
		titleStyle = titleStyle.Background(lipgloss.Color("240"))
	}

	formattedTitle := titleStyle.Render(title)

	// Handle selection highlighting
	lines := strings.Split(content, "\n")
	var formattedLines []string

	for i, line := range lines {
		shouldHighlight := false
		switch panelType {
		case FilesPanel:
			shouldHighlight = i == m.selectedFile && m.activePanel == FilesPanel
		case BranchesPanel:
			shouldHighlight = i == m.selectedBranch && m.activePanel == BranchesPanel
		case LogPanel:
			shouldHighlight = i == m.selectedCommit && m.activePanel == LogPanel
		}

		if shouldHighlight {
			formattedLines = append(formattedLines, selectedStyle.Render(line))
		} else {
			formattedLines = append(formattedLines, line)
		}
	}

	formattedContent := strings.Join(formattedLines, "\n")

	return style.
		Width(width).
		Height(height).
		Render(formattedTitle + "\n" + formattedContent)
}
