package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// keys is a package-level variable that holds the application's keybindings.
var keys = DefaultKeyMap()

// Update is the central message handler for the application. It's called by the
// Bubble Tea runtime when a message is received. It's responsible for updating
// the model's state based on the message and returning any commands to execute.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	// Handle terminal window resize events.
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width
		// Recalculate the dimensions of the help viewport.
		m.helpViewport.Width = int(float64(m.width) * 0.5)
		m.helpViewport.Height = int(float64(m.height) * 0.75)

	// Handle keyboard input.
	case tea.KeyMsg:
		// If the help view is currently visible, handle its specific keybindings.
		if m.showHelp {
			// Allow the viewport to handle scrolling with arrow keys.
			m.helpViewport, cmd = m.helpViewport.Update(msg)
			cmds = append(cmds, cmd)

			// Check for keys that close the help view.
			switch {
			case key.Matches(msg, keys.Quit), key.Matches(msg, keys.ToggleHelp), key.Matches(msg, keys.Escape):
				m.showHelp = false
			}
			return m, tea.Batch(cmds...)
		}

		// Handle keybindings for the main application view.
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.ToggleHelp):
			m.showHelp = true
			// Generate and style the help content when the view is opened.
			m.helpContent = m.generateHelpContent()
			m.helpViewport.SetContent(m.helpContent)
			m.helpViewport.Style = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(m.theme.ActivePanel.GetBorderTopForeground()).
				Padding(1, 2)
			m.helpViewport.GotoTop()

		case key.Matches(msg, keys.SwitchTheme):
			m.nextTheme()
			// Regenerate help content to apply new theme colors.
			m.helpContent = m.generateHelpContent()
			m.helpViewport.SetContent(m.helpContent)
			m.helpViewport.Style = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(m.theme.ActivePanel.GetBorderTopForeground()).
				Padding(1, 2)

		// Handle panel focus navigation.
		case key.Matches(msg, keys.FocusNext):
			m.nextPanel()

		case key.Matches(msg, keys.FocusPrev):
			m.prevPanel()

		// Handle direct panel focus via number keys.
		case key.Matches(msg, keys.FocusZero):
			m.focusedPanel = MainPanel

		case key.Matches(msg, keys.FocusOne):
			m.focusedPanel = StatusPanel

		case key.Matches(msg, keys.FocusTwo):
			m.focusedPanel = FilesPanel

		case key.Matches(msg, keys.FocusThree):
			m.focusedPanel = BranchesPanel

		case key.Matches(msg, keys.FocusFour):
			m.focusedPanel = CommitsPanel

		case key.Matches(msg, keys.FocusFive):
			m.focusedPanel = StashPanel
		}
		// Return the updated model to the Bubble Tea runtime.
		return m, nil
	}

	// Batch and return any commands that were generated.
	return m, tea.Batch(cmds...)
}
