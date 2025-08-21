package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

var keys = DefaultKeyMap()

// Update handles all incoming messages and updates the model accordingly.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// tea.WindowSizeMsg is sent when the terminal window is resized.
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width

	// tea.KeyMsg is sent when a key is pressed.
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.SwitchTheme):
			m.nextTheme()

		case key.Matches(msg, keys.FocusNext):
			m.nextPanel()

		case key.Matches(msg, keys.FocusPrev):
			m.prevPanel()

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
	return m, nil
}
