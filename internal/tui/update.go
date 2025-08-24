package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gitxtui/gitx/internal/git"
	zone "github.com/lrstanley/bubblezone"
)

var keys = DefaultKeyMap()

// panelContentUpdatedMsg is a generic message used to signal that a panel's
// content has been updated.
type panelContentUpdatedMsg struct {
	panel   Panel
	content string
}

// fetchPanelContent is a generic command that fetches content for a given panel.
func fetchPanelContent(gc *git.GitCommands, panel Panel) tea.Cmd {
	return func() tea.Msg {
		var content string
		var err error

		switch panel {
		case StatusPanel:
			content, err = gc.GetStatus()
		case FilesPanel:
			content = "\nPLACEHOLDER DATA??\n1\n2\n cmd/gitx\n MM internal/tui\n file1.go\n M file2.txt\n A file3.md" // FIXME: Placeholder
		case BranchesPanel:
			content = "\nPLACEHOLDER DATA??\n1\n2\n3a\n4b\n  main\n* feature/new-ui\n test/add-test\n hotfix/bug-123" // FIXME: Placeholder
		case CommitsPanel:
			content = strings.Join([]string{
				"\nPLACEHOLDER DATA??\n1\n2\nf7875b4 (HEAD -> feature/new-ui) feat: add new panel layout",
				"a3e8b1c (origin/main, main) fix: correct scrolling logic",
				"c1d9f2e chore: update dependencies",
				"f7875b4 (HEAD -> feature/new-ui) feat: add new panel layout",
				"a3e8b1c (origin/main, main) fix: correct scrolling logic",
				"c1d9f2e chore: update dependencies",
			}, "\n") // FIXME: Placeholder
		case StashPanel:
			content = "PLACEHOLDER DATA??\n1\n2\n\n3\n4\n5\n6stash@{0}: WIP on feature/new-ui: 52f3a6b feat: add panels" // FIXME: Placeholder
		case MainPanel:
			content = "\nPLACEHOLDER DATA??\n1\n2\nThis is the main panel.\n\nSelect an item from another panel to see details here." // FIXME: Placeholder
		case SecondaryPanel:
			content = "PLACEHOLDER DATA??\n1\n2\nThis is the secondary panel." // FIXME: Placeholder
		}

		if err != nil {
			content = "Error: " + err.Error()
		}

		return panelContentUpdatedMsg{panel: panel, content: content}
	}
}

// Update is the central message handler for the application.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	oldFocus := m.focusedPanel

	switch msg := msg.(type) {
	case panelContentUpdatedMsg:
		m.panels[msg.panel].content = msg.content
		m.panels[msg.panel].viewport.SetContent(msg.content)
		return m, nil

	case tea.WindowSizeMsg:
		m, cmd = m.handleWindowSizeMsg(msg)
		cmds = append(cmds, cmd)

	case tea.MouseMsg:
		m, cmd = m.handleMouseMsg(msg)
		cmds = append(cmds, cmd)

	case tea.KeyMsg:
		m, cmd = m.handleKeyMsg(msg)
		cmds = append(cmds, cmd)
	}

	// If a message caused the focus to change, we need to recalculate the layout.
	if m.focusedPanel != oldFocus {
		if m.focusedPanel == StashPanel || m.focusedPanel == SecondaryPanel {
			// If the new panel is Stash or Secondary, scroll to top.
			m.panels[m.focusedPanel].viewport.GotoTop()
		}
		m = m.recalculateLayout()
	}

	return m, tea.Batch(cmds...)
}

// handleWindowSizeMsg recalculates the layout and resizes all viewports.
func (m Model) handleWindowSizeMsg(msg tea.WindowSizeMsg) (Model, tea.Cmd) {
	m.width = msg.Width
	m.height = msg.Height
	m.help.Width = msg.Width
	m.helpViewport.Width = int(float64(m.width) * 0.5)
	m.helpViewport.Height = int(float64(m.height) * 0.75)

	m = m.recalculateLayout()
	return m, nil
}

// recalculateLayout is the single source of truth for panel sizes.
func (m Model) recalculateLayout() Model {
	if m.width == 0 || m.height == 0 {
		return m
	}

	contentHeight := m.height - 1
	m.panelHeights = make([]int, totalPanels)
	expandedHeight := int(float64(contentHeight) * 0.3)

	// --- Right Column ---
	if m.focusedPanel == SecondaryPanel {
		m.panelHeights[SecondaryPanel] = expandedHeight
		m.panelHeights[MainPanel] = contentHeight - expandedHeight
	} else {
		m.panelHeights[SecondaryPanel] = 3 // Default collapsed size
		m.panelHeights[MainPanel] = contentHeight - 3
	}

	// --- Left Column ---
	m.panelHeights[StatusPanel] = 3 // Always fixed
	remainingHeight := contentHeight - m.panelHeights[StatusPanel]
	flexiblePanels := []Panel{FilesPanel, BranchesPanel, CommitsPanel, StashPanel}
	expandedPanel := StashPanel // The only expandable panel on the left

	if m.focusedPanel == expandedPanel {
		m.panelHeights[expandedPanel] = expandedHeight
	} else {
		m.panelHeights[expandedPanel] = 3 // Default collapsed size
	}

	// Distribute remaining height among the other flexible panels
	otherPanelsCount := len(flexiblePanels) - 1
	otherPanelHeight := (remainingHeight - m.panelHeights[expandedPanel]) / otherPanelsCount

	for _, pType := range flexiblePanels {
		if pType != expandedPanel {
			m.panelHeights[pType] = otherPanelHeight
		}
	}
	// Give any remainder to the last non-expanded panel to fill space
	m.panelHeights[CommitsPanel] += (remainingHeight - m.panelHeights[expandedPanel]) % otherPanelsCount

	return m.updateViewportSizes()
}

// updateViewportSizes applies the calculated heights from the model to the viewports.
func (m Model) updateViewportSizes() Model {
	horizontalBorderWidth := m.theme.ActiveBorder.Style.GetHorizontalBorderSize()
	titleBarHeight := 2 // Top and bottom border

	rightSectionWidth := m.width - int(float64(m.width)*0.3)
	rightContentWidth := rightSectionWidth - horizontalBorderWidth
	m.panels[MainPanel].viewport.Width = rightContentWidth
	m.panels[MainPanel].viewport.Height = m.panelHeights[MainPanel] - titleBarHeight
	m.panels[SecondaryPanel].viewport.Width = rightContentWidth
	m.panels[SecondaryPanel].viewport.Height = m.panelHeights[SecondaryPanel] - titleBarHeight

	leftSectionWidth := int(float64(m.width) * 0.3)
	leftContentWidth := leftSectionWidth - horizontalBorderWidth
	leftPanels := []Panel{StatusPanel, FilesPanel, BranchesPanel, CommitsPanel, StashPanel}
	for _, panel := range leftPanels {
		m.panels[panel].viewport.Width = leftContentWidth
		m.panels[panel].viewport.Height = m.panelHeights[panel] - titleBarHeight
	}
	return m
}

// handleMouseMsg handles all mouse events.
func (m Model) handleMouseMsg(msg tea.MouseMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.showHelp {
		if msg.Action == tea.MouseActionRelease && msg.Button == tea.MouseButtonLeft && zone.Get("help-button").InBounds(msg) {
			m.toggleHelp()
		} else {
			m.helpViewport, cmd = m.helpViewport.Update(msg)
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)
	}

	if msg.Action == tea.MouseActionRelease && msg.Button == tea.MouseButtonLeft {
		if zone.Get("help-button").InBounds(msg) {
			m.toggleHelp()
			return m, nil
		}
	}

	for i := range m.panels {
		panel := Panel(i)
		if zone.Get(panel.ID()).InBounds(msg) {
			if msg.Action == tea.MouseActionRelease && msg.Button == tea.MouseButtonLeft {
				m.focusedPanel = panel
			}
			m.panels[i].viewport, cmd = m.panels[i].viewport.Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}
	}
	return m, nil
}

// handleKeyMsg handles all keyboard events.
func (m Model) handleKeyMsg(msg tea.KeyMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.showHelp {
		m.helpViewport, cmd = m.helpViewport.Update(msg)
		cmds = append(cmds, cmd)
		switch {
		case key.Matches(msg, keys.Quit), key.Matches(msg, keys.ToggleHelp), key.Matches(msg, keys.Escape):
			m.showHelp = false
		case key.Matches(msg, keys.SwitchTheme):
			m.nextTheme()
			m.styleHelpViewContent()
		}
		return m, tea.Batch(cmds...)
	}

	switch {
	case key.Matches(msg, keys.Quit):
		return m, tea.Quit

	case key.Matches(msg, keys.ToggleHelp):
		m.toggleHelp()

	case key.Matches(msg, keys.SwitchTheme):
		m.nextTheme()

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

	case key.Matches(msg, keys.FocusSix):
		m.focusedPanel = SecondaryPanel
	}
	return m, nil
}

// toggleHelp toggles the visibility of the help view and prepares its content.
func (m *Model) toggleHelp() {
	m.showHelp = !m.showHelp
	if m.showHelp {
		m.styleHelpViewContent()
	}
}

// styleHelpViewContent refreshes the styles of the Help View content.
func (m *Model) styleHelpViewContent() {
	m.helpContent = m.generateHelpContent()
	m.helpViewport.SetContent(m.helpContent)
	m.helpViewport.Style = lipgloss.NewStyle()
	m.helpViewport.GotoTop()
}
