package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
)

// Panel is an enumeration of all the panels in the UI.
type Panel int

const (
	MainPanel Panel = iota
	StatusPanel
	FilesPanel
	BranchesPanel
	CommitsPanel
	StashPanel
	SecondaryPanel
	totalPanels // Used to iterate over all panels
)

// ID returns a string ID for the panel, used for mouse zone detection.
func (p Panel) ID() string {
	return fmt.Sprintf("panel-%d", p)
}

// panel represents the state of a single UI panel.
type panel struct {
	viewport viewport.Model
	content  string
}

// nextPanel shifts focus to the next Panel.
func (m *Model) nextPanel() {
	// Skips SecondaryPanel
	m.focusedPanel = (m.focusedPanel + 1) % (totalPanels - 1)
}

// prevPanel shifts focus to the previous Panel.
func (m *Model) prevPanel() {
	// skip SecondaryPanel
	if m.focusedPanel == 0 {
		m.focusedPanel = totalPanels - 2
		return
	}
	m.focusedPanel = m.focusedPanel - 1
}
