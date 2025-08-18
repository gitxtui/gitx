package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the state of the TUI.
type Model struct {
	width  int
	height int
}

// App is the main application struct.
type App struct {
	program *tea.Program
}

// NewApp initializes a new TUI application.
func NewApp() *App {
	model := Model{}
	// We're using WithAltScreen to have a dedicated screen for the TUI.
	program := tea.NewProgram(model, tea.WithAltScreen())
	return &App{program: program}
}

// Run starts the TUI application.
func (a *App) Run() error {
	_, err := a.program.Run()
	return err
}

// Init is the first command that is run when the program starts.
func (m Model) Init() tea.Cmd {
	// tea.EnterAltScreen is a command that tells the terminal to enter the alternate screen buffer.
	return tea.EnterAltScreen
}

// Update handles all incoming messages and updates the model accordingly.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// tea.WindowSizeMsg is sent when the terminal window is resized.
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	// tea.KeyMsg is sent when a key is pressed.
	case tea.KeyMsg:
		switch msg.String() {
		// These keys will exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	// Return the updated model to the Bubble Tea runtime.
	return m, nil
}

// View renders the UI.
func (m Model) View() string {
	// If the window size has not been determined yet, show a loading message.
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	// --- Styles ---
	// Define a generic panel style.
	panelStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240"))

	// Get border sizes to calculate content dimensions accurately.
	horizontalBorderWidth := panelStyle.GetHorizontalBorderSize()
	verticalBorderHeight := panelStyle.GetVerticalBorderSize()

	// --- Layout ---
	// Calculate rendered widths for the two main vertical sections.
	leftSectionRenderedWidth := int(float64(m.width) * 0.3)
	rightSectionRenderedWidth := m.width - leftSectionRenderedWidth

	// Calculate content widths for panels inside sections.
	leftPanelContentWidth := leftSectionRenderedWidth - horizontalBorderWidth
	rightPanelContentWidth := rightSectionRenderedWidth - horizontalBorderWidth

	// --- Left Section (5 panels) ---
	// Calculate content heights for the 5 panels on the left.
	leftSectionAvailableContentHeight := m.height - (5 * verticalBorderHeight)
	leftPanelContentHeight := leftSectionAvailableContentHeight / 5
	leftPanelLastContentHeight := leftSectionAvailableContentHeight - (leftPanelContentHeight * 4)

	leftPanelTitles := []string{"Status", "Files", "Branches", "Commits", "Stash"}
	var leftPanels []string
	for i, title := range leftPanelTitles {
		h := leftPanelContentHeight
		if i == len(leftPanelTitles)-1 {
			h = leftPanelLastContentHeight
		}
		panel := panelStyle.
			Width(leftPanelContentWidth).
			Height(h).
			Render(fmt.Sprintf("-> %s", title))
		leftPanels = append(leftPanels, panel)
	}
	leftSection := lipgloss.JoinVertical(lipgloss.Left, leftPanels...)

	// --- Right Section (2 panels) ---
	// Calculate content heights for the 2 panels on the right.
	rightSectionAvailableContentHeight := m.height - (2 * verticalBorderHeight)
	rightPanelContentHeight := rightSectionAvailableContentHeight / 2
	rightPanelLastContentHeight := rightSectionAvailableContentHeight - rightPanelContentHeight

	rightPanelTitles := []string{"Main", "Secondary"}
	var rightPanels []string
	for i, title := range rightPanelTitles {
		h := rightPanelContentHeight
		if i == len(rightPanelTitles)-1 {
			h = rightPanelLastContentHeight
		}
		panel := panelStyle.
			Width(rightPanelContentWidth).
			Height(h).
			Render(fmt.Sprintf("-> %s", title))
		rightPanels = append(rightPanels, panel)
	}
	rightSection := lipgloss.JoinVertical(lipgloss.Left, rightPanels...)

	// --- Final Layout ---
	// Join the left and right sections horizontally.
	return lipgloss.JoinHorizontal(lipgloss.Top, leftSection, rightSection)
}
