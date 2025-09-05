package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gitxtui/gitx/internal/git"
	zone "github.com/lrstanley/bubblezone"
)

var keys = DefaultKeyMap()

type panelContentUpdatedMsg struct {
	panel   Panel
	content string
}

type lineClickedMsg struct {
	panel     Panel
	lineIndex int
}

type fileWatcherMsg struct{}

func (m Model) fetchPanelContent(panel Panel) tea.Cmd {
	return func() tea.Msg {
		var content, repoName, branchName string
		var err error

		switch panel {
		case StatusPanel:
			// --- THE FIX ---
			// Apply styling here for the simple, non-selectable status panel
			repoName, branchName, err = m.git.GetRepoInfo()
			if err == nil {
				repo := m.theme.BranchCurrent.Render(repoName)
				branch := m.theme.BranchCurrent.Render(branchName)
				content = fmt.Sprintf("%s → %s", repo, branch)
			}
		case FilesPanel:
			content, err = m.git.GetStatus(git.StatusOptions{Porcelain: true})
		case BranchesPanel:
			branchList, err := m.git.GetBranches()
			if err != nil {
				content = "Error getting branches: " + err.Error()
				break
			}
			var builder strings.Builder
			for _, b := range branchList {
				name := b.Name
				if b.IsCurrent {
					name = fmt.Sprintf("(*) → %s", b.Name)
				}
				line := fmt.Sprintf("%s\t%s", b.LastCommit, name) // Use tab separator
				builder.WriteString(line + "\n")
			}
			content = strings.TrimSpace(builder.String())
		case CommitsPanel:
			logs, err := m.git.GetCommitLogsGraph()
			if err != nil {
				content = "Error getting commit logs: " + err.Error()
				break
			}
			var builder strings.Builder
			for _, log := range logs {
				var line string
				if log.SHA != "" {
					line = fmt.Sprintf("%s\t%s\t%s\t%s", log.Graph, log.SHA, log.AuthorInitials, log.Subject) // Use tab separator
				} else {
					line = log.Graph
				}
				builder.WriteString(line + "\n")
			}
			content = strings.TrimSpace(builder.String())
		case StashPanel:
			stashList, err := m.git.GetStashes()
			if err != nil {
				content = "Error getting stashes: " + err.Error()
				break
			}
			if len(stashList) == 0 {
				content = "No stashed changes."
				break
			}
			var builder strings.Builder
			for _, s := range stashList {
				// Create a tab-delimited string: "stash@{0}\tWIP on master: ..."
				line := fmt.Sprintf("%s\t%s: %s", s.Name, s.Branch, s.Message)
				builder.WriteString(line + "\n")
			}
			content = strings.TrimSpace(builder.String())
		case MainPanel, SecondaryPanel:
			content = "Loading..." // Or placeholder data
		}

		if err != nil {
			content = "Error: " + err.Error()
		}
		return panelContentUpdatedMsg{panel: panel, content: content}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	oldFocus := m.focusedPanel

	switch msg := msg.(type) {
	case panelContentUpdatedMsg:
		// --- START: INTELLIGENT CURSOR PRESERVATION ---
		var selectedPath string
		// If the updated panel is the FilesPanel and it's focused, get the path of the currently selected line.
		if msg.panel == FilesPanel && m.focusedPanel == FilesPanel && m.panels[FilesPanel].cursor < len(m.panels[FilesPanel].lines) {
			line := m.panels[FilesPanel].lines[m.panels[FilesPanel].cursor]
			parts := strings.Split(line, "\t")
			if len(parts) == 3 {
				selectedPath = parts[2] // The path is the third element
			}
		}
		// Preserve cursor index for other panels
		oldCursor := m.panels[msg.panel].cursor
		// --- END: INTELLIGENT CURSOR PRESERVATION ---

		if msg.panel == FilesPanel {
			root := BuildTree(msg.content)
			renderedTree := root.Render(m.theme)
			m.panels[FilesPanel].lines = renderedTree
			m.panels[FilesPanel].viewport.SetContent(strings.Join(renderedTree, "\n"))

			// --- START: RESTORE CURSOR BY PATH ---
			newCursorPos := 0 // Default to top
			if selectedPath != "" {
				// Find the new index of the previously selected path
				for i, line := range renderedTree {
					parts := strings.Split(line, "\t")
					if len(parts) == 3 && parts[2] == selectedPath {
						newCursorPos = i
						break
					}
				}
			}
			m.panels[FilesPanel].cursor = newCursorPos
			// --- END: RESTORE CURSOR BY PATH ---

		} else {
			lines := strings.Split(msg.content, "\n")
			m.panels[msg.panel].lines = lines
			m.panels[msg.panel].viewport.SetContent(msg.content)
			// --- THE FIX ---
			m.panels[msg.panel].content = msg.content // Add this line

			// Restore cursor for other panels
			if oldCursor < len(lines) {
				m.panels[msg.panel].cursor = oldCursor
			} else if len(lines) > 0 {
				m.panels[msg.panel].cursor = len(lines) - 1
			} else {
				m.panels[msg.panel].cursor = 0
			}
		}
		return m, nil

	case fileWatcherMsg:
		return m, tea.Batch(
			m.fetchPanelContent(StatusPanel),
			m.fetchPanelContent(FilesPanel),
			m.fetchPanelContent(BranchesPanel),
			m.fetchPanelContent(CommitsPanel),
			m.fetchPanelContent(StashPanel),
			m.fetchPanelContent(MainPanel),
			m.fetchPanelContent(SecondaryPanel),
		)

	case lineClickedMsg:
		if msg.lineIndex < len(m.panels[msg.panel].lines) {
			m.panels[msg.panel].cursor = msg.lineIndex
		}
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

	if m.focusedPanel != oldFocus {
		if m.focusedPanel == StashPanel || m.focusedPanel == SecondaryPanel {
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
	expandedHeight := int(float64(contentHeight) * 0.4)
	collapsedHeight := 3

	// --- Right Column ---
	if m.focusedPanel == SecondaryPanel {
		m.panelHeights[SecondaryPanel] = expandedHeight
		m.panelHeights[MainPanel] = contentHeight - expandedHeight
	} else {
		m.panelHeights[SecondaryPanel] = collapsedHeight
		m.panelHeights[MainPanel] = contentHeight - collapsedHeight
	}

	// --- Left Column ---
	m.panelHeights[StatusPanel] = 3
	remainingHeight := contentHeight - m.panelHeights[StatusPanel]

	if m.focusedPanel == StashPanel {
		m.panelHeights[StashPanel] = expandedHeight
	} else {
		m.panelHeights[StashPanel] = collapsedHeight
	}

	flexiblePanels := []Panel{FilesPanel, BranchesPanel, CommitsPanel}
	heightForFlex := remainingHeight - m.panelHeights[StashPanel]
	focusedFlexPanelFound := false

	for _, p := range flexiblePanels {
		if p == m.focusedPanel {
			focusedFlexPanelFound = true
			break
		}
	}

	if focusedFlexPanelFound {
		m.panelHeights[m.focusedPanel] = expandedHeight
		heightForOthers := heightForFlex - expandedHeight
		otherPanels := []Panel{}
		for _, p := range flexiblePanels {
			if p != m.focusedPanel {
				otherPanels = append(otherPanels, p)
			}
		}
		if len(otherPanels) > 0 {
			share := heightForOthers / len(otherPanels)
			for _, p := range otherPanels {
				m.panelHeights[p] = share
			}
			m.panelHeights[otherPanels[len(otherPanels)-1]] += heightForOthers % len(otherPanels)
		}
	} else {
		// Default distribution when none of the main flexible panels are focused.
		m.panelHeights[FilesPanel] = int(float64(heightForFlex) * 0.4)
		m.panelHeights[BranchesPanel] = int(float64(heightForFlex) * 0.3)
		m.panelHeights[CommitsPanel] = heightForFlex - m.panelHeights[FilesPanel] - m.panelHeights[BranchesPanel]
	}

	return m.updateViewportSizes()
}

// updateViewportSizes applies the calculated heights from the model to the viewports.
func (m Model) updateViewportSizes() Model {
	horizontalBorderWidth := m.theme.ActiveBorder.Style.GetHorizontalBorderSize()
	titleBarHeight := 2

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

// handleMouseMsg handles all mouse events
func (m Model) handleMouseMsg(msg tea.MouseMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.showHelp {
		if zone.Get("help-button").InBounds(msg) && msg.Action == tea.MouseActionRelease {
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

		// Check for clicks on panel lines first.
		for p := range m.panels {
			panel := Panel(p)
			// Only check selectable panels.
			if panel != FilesPanel && panel != BranchesPanel && panel != CommitsPanel && panel != StashPanel {
				continue
			}
			// Check each line in the panel.
			for i := 0; i < len(m.panels[panel].lines); i++ {
				lineID := fmt.Sprintf("%s-line-%d", panel.ID(), i)
				if zone.Get(lineID).InBounds(msg) {
					m.focusedPanel = panel
					return m, func() tea.Msg {
						return lineClickedMsg{panel: panel, lineIndex: i}
					}
				}
			}
		}

		// If no line was clicked, check for clicks on the panel itself to change focus.
		for i := range m.panels {
			if zone.Get(Panel(i).ID()).InBounds(msg) {
				m.focusedPanel = Panel(i)
				break
			}
		}
	}

	for i := range m.panels {
		panel := Panel(i)
		if zone.Get(panel.ID()).InBounds(msg) {
			m.panels[panel].viewport, cmd = m.panels[panel].viewport.Update(msg)
			cmds = append(cmds, cmd)
			break
		}
	}

	return m, tea.Batch(cmds...)
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

	// Global key handling that should take precedence over panel-specific logic.
	switch {
	case key.Matches(msg, keys.Quit):
		return m, tea.Quit

	case key.Matches(msg, keys.ToggleHelp):
		m.toggleHelp()
		return m, nil

	case key.Matches(msg, keys.SwitchTheme):
		m.nextTheme()
		return m, nil

	case key.Matches(msg, keys.FocusNext), key.Matches(msg, keys.FocusPrev),
		key.Matches(msg, keys.FocusZero), key.Matches(msg, keys.FocusOne),
		key.Matches(msg, keys.FocusTwo), key.Matches(msg, keys.FocusThree),
		key.Matches(msg, keys.FocusFour), key.Matches(msg, keys.FocusFive),
		key.Matches(msg, keys.FocusSix):
		switch {
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
		case key.Matches(msg, keys.FocusSix):
			m.focusedPanel = SecondaryPanel
		}
		return m, nil
	}

	// Panel-specific key handling for custom logic (like cursor movement).
	switch m.focusedPanel {
	case FilesPanel, BranchesPanel, CommitsPanel, StashPanel:
		p := &m.panels[m.focusedPanel]
		switch {
		case key.Matches(msg, keys.Up):
			if p.cursor > 0 {
				p.cursor--
				// Scroll viewport up if cursor is out of view
				if p.cursor < p.viewport.YOffset {
					p.viewport.SetYOffset(p.cursor)
				}
			}
			// We handled the key, so we return to prevent the default viewport scrolling.
			return m, nil
		case key.Matches(msg, keys.Down):
			if p.cursor < len(p.lines)-1 {
				p.cursor++
				// Scroll viewport down if cursor is out of view
				if p.cursor >= p.viewport.YOffset+p.viewport.Height {
					p.viewport.SetYOffset(p.cursor - p.viewport.Height + 1)
				}
			}
			// We handled the key, so we return to prevent the default viewport scrolling.
			return m, nil
		}
	}

	// Always pass the key message to the focused panel's viewport for scrolling.
	m.panels[m.focusedPanel].viewport, cmd = m.panels[m.focusedPanel].viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
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
