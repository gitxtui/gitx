package tui

// Panel represents a section of the UI.
type Panel int

const (
	MainPanel Panel = iota
	StatusPanel
	FilesPanel
	BranchesPanel
	CommitsPanel
	StashPanel
	SecondaryPanel
	totalPanels
)

// nextPanel shifts focus to the next Panel.
func (m *Model) nextPanel() {
	m.focusedPanel = (m.focusedPanel + 1) % totalPanels
	// skip SecondaryPanel
	if m.focusedPanel == SecondaryPanel {
		m.nextPanel()
	}
}

// prevPanel shifts focus to the previous Panel.
func (m *Model) prevPanel() {
	m.focusedPanel = (m.focusedPanel - 1 + totalPanels) % totalPanels
	// skip SecondaryPanel
	if m.focusedPanel == SecondaryPanel {
		m.prevPanel()
	}
}
