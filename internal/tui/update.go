package tui

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gitxtui/gitx/internal/git"
	zone "github.com/lrstanley/bubblezone"
)

var keys = DefaultKeyMap()

// panelContentUpdatedMsg is sent when new content for a panel has been fetched.
type panelContentUpdatedMsg struct {
	panel   Panel
	content string
}

// mainContentUpdatedMsg is sent when the content for the main panel has been fetched.
type mainContentUpdatedMsg struct {
	content string
}

// lineClickedMsg is sent when a user clicks on a line in a selectable panel.
type lineClickedMsg struct {
	panel     Panel
	lineIndex int
}

// fileWatcherMsg is sent by the file watcher when the repository state changes.
type fileWatcherMsg struct{}

// errMsg is used to propagate errors back to the update loop.
type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

// Update is the main message handler for the TUI. It processes user input,
// window events, and application-specific messages.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.mode {
	case modeInput:
		return m.updateInput(msg)
	case modeConfirm:
		return m.updateConfirm(msg)
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd
	oldFocus := m.focusedPanel

	switch msg := msg.(type) {
	case errMsg:
		// You can improve this to show errors in the UI
		log.Printf("error: %v", msg)
		return m, nil

	case mainContentUpdatedMsg:
		m.panels[MainPanel].content = msg.content
		m.panels[MainPanel].viewport.SetContent(msg.content)
		return m, nil

	case panelContentUpdatedMsg:
		var selectedPath string
		// If the FilesPanel is being updated, try to find the path of the
		// currently selected item to preserve the cursor position after the refresh.
		if msg.panel == FilesPanel && m.panels[FilesPanel].cursor < len(m.panels[FilesPanel].lines) {
			line := m.panels[FilesPanel].lines[m.panels[FilesPanel].cursor]
			parts := strings.Split(line, "\t")
			if len(parts) == 4 {
				selectedPath = parts[3]
			}
		}

		oldCursor := m.panels[msg.panel].cursor

		if msg.panel == FilesPanel {
			root := BuildTree(msg.content)
			renderedTree := root.Render(m.theme)
			m.panels[FilesPanel].lines = renderedTree
			m.panels[FilesPanel].viewport.SetContent(strings.Join(renderedTree, "\n"))

			// Restore the cursor to the previously selected file path.
			newCursorPos := 0 // Default to top.
			if selectedPath != "" {
				for i, line := range renderedTree {
					parts := strings.Split(line, "\t")
					if len(parts) == 4 && parts[3] == selectedPath {
						newCursorPos = i
						break
					}
				}
			}
			m.panels[FilesPanel].cursor = newCursorPos
		} else {
			lines := strings.Split(msg.content, "\n")
			m.panels[msg.panel].lines = lines
			m.panels[msg.panel].viewport.SetContent(msg.content)
			m.panels[msg.panel].content = msg.content

			// Restore cursor by index for other, more stable panels.
			if oldCursor < len(lines) {
				m.panels[msg.panel].cursor = oldCursor
			} else if len(lines) > 0 {
				m.panels[msg.panel].cursor = len(lines) - 1
			} else {
				m.panels[msg.panel].cursor = 0
			}
		}
		return m, m.updateMainPanel()

	case fileWatcherMsg:
		// When the repository changes, trigger a content refresh for all panels.
		return m, tea.Batch(
			m.fetchPanelContent(StatusPanel),
			m.fetchPanelContent(FilesPanel),
			m.fetchPanelContent(BranchesPanel),
			m.fetchPanelContent(CommitsPanel),
			m.fetchPanelContent(StashPanel),
		)

	case lineClickedMsg:
		// Handle direct selection of a line via mouse click.
		if msg.lineIndex < len(m.panels[msg.panel].lines) {
			p := &m.panels[msg.panel]
			p.cursor = msg.lineIndex
			// Ensure the selected line is visible in the viewport.
			if p.cursor < p.viewport.YOffset {
				p.viewport.SetYOffset(p.cursor)
			}
			if p.cursor >= p.viewport.YOffset+p.viewport.Height {
				p.viewport.SetYOffset(p.cursor - p.viewport.Height + 1)
			}
		}
		m.activeSourcePanel = msg.panel
		m.panels[MainPanel].viewport.GotoTop()
		return m, m.updateMainPanel()

	case tea.WindowSizeMsg:
		m, cmd = m.handleWindowSizeMsg(msg)
		cmds = append(cmds, cmd)

	case tea.MouseMsg:
		m, cmd = m.handleMouseMsg(msg)
		cmds = append(cmds, cmd)

	case tea.KeyMsg:
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
			m.handleFocusKeys(msg)
			return m, nil
		}

		cmd = m.handlePanelKeys(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}
	}

	if m.focusedPanel != oldFocus {
		// When focus changes, reset scroll for the Stash and Secondary panels
		if m.focusedPanel == StashPanel || m.focusedPanel == SecondaryPanel {
			m.panels[m.focusedPanel].viewport.GotoTop()
		}

		// Update the active source panel and main panel content if the new focus is a source panel
		if m.focusedPanel != MainPanel && m.focusedPanel != SecondaryPanel {
			m.activeSourcePanel = m.focusedPanel
			m.panels[MainPanel].viewport.GotoTop() // Reset main panel scroll on source change
			cmd = m.updateMainPanel()
			cmds = append(cmds, cmd)
		}

		m = m.recalculateLayout()
	}

	// The original viewport update logic for scrolling
	m.panels[m.focusedPanel].viewport, cmd = m.panels[m.focusedPanel].viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// updateInput handles updates when in text input mode.
func (m Model) updateInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			val := m.textInput.Value()
			cmd = m.inputCallback(val)
			m.mode = modeNormal
			m.textInput.Reset()
			return m, cmd
		case tea.KeyEsc:
			m.mode = modeNormal
			m.textInput.Reset()
			return m, nil
		}
	}
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// updateConfirm handles updates when in confirmation mode.
func (m Model) updateConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch strings.ToLower(msg.String()) {
		case "y":
			cmd = m.confirmCallback(true)
			m.mode = modeNormal
			return m, cmd
		case "n", "esc":
			cmd = m.confirmCallback(false)
			m.mode = modeNormal
			return m, cmd
		}
	}
	return m, nil
}

// fetchPanelContent returns a command that fetches the content for a specific panel.
func (m Model) fetchPanelContent(panel Panel) tea.Cmd {
	return func() tea.Msg {
		var content, repoName, branchName string
		var err error
		switch panel {
		case StatusPanel:
			repoName, branchName, err = m.git.GetRepoInfo()
			if err == nil {
				repo := m.theme.BranchCurrent.Render(repoName)
				branch := m.theme.BranchCurrent.Render(branchName)
				content = fmt.Sprintf("%s → %s", repo, branch)
			}
		case FilesPanel:
			content, err = m.git.GetStatus(git.StatusOptions{Porcelain: true})
		case BranchesPanel:
			var branchList []*git.Branch
			branchList, err = m.git.GetBranches()
			if err == nil {
				var builder strings.Builder
				for _, b := range branchList {
					name := b.Name
					if b.IsCurrent {
						name = fmt.Sprintf("(*) → %s", b.Name)
					}
					line := fmt.Sprintf("%s\t%s", b.LastCommit, name)
					builder.WriteString(line + "\n")
				}
				content = strings.TrimSpace(builder.String())
			}
		case CommitsPanel:
			var logs []git.CommitLog
			logs, err = m.git.GetCommitLogsGraph()
			if err == nil {
				var builder strings.Builder
				for _, log := range logs {
					var line string
					if log.SHA != "" {
						line = fmt.Sprintf("%s\t%s\t%s\t%s", log.Graph, log.SHA, log.AuthorInitials, log.Subject)
					} else {
						line = log.Graph
					}
					builder.WriteString(line + "\n")
				}
				content = strings.TrimSpace(builder.String())
			}
		case StashPanel:
			var stashList []*git.Stash
			stashList, err = m.git.GetStashes()
			if err == nil {
				if len(stashList) == 0 {
					content = "No stashed changes."
				} else {
					var builder strings.Builder
					for _, s := range stashList {
						line := fmt.Sprintf("%s\t%s: %s", s.Name, s.Branch, s.Message)
						builder.WriteString(line + "\n")
					}
					content = strings.TrimSpace(builder.String())
				}
			}
		case SecondaryPanel:
			url := m.theme.Hyperlink.Render(githubMainPage)
			content = strings.Join([]string{
				"\t--- Feature in development! ---",
				"\n\t* This panel will contain all the command logs and history for of TUI app.",
				fmt.Sprintf("\t* visit for more details: %s.", url),
			}, "\n")
		}

		if err != nil {
			content = "Error: " + err.Error()
		}
		return panelContentUpdatedMsg{panel: panel, content: content}
	}
}

// updateMainPanel returns a command that fetches the content for the main panel
// based on the currently active source panel.
func (m *Model) updateMainPanel() tea.Cmd {
	return func() tea.Msg {
		var content string
		var err error
		switch m.activeSourcePanel {
		case StatusPanel:
			userName, _ := m.git.GetUserName()
			url := m.theme.Hyperlink.Render(docsPage)
			msgHeading := m.theme.WelcomeHeading.Render(asciiArt)
			msgBody := fmt.Sprintf(welcomeMsg, m.theme.UserName.Render(userName), url)
			content = fmt.Sprintf(msgHeading, m.theme.WelcomeMsg.Render(msgBody))
		case FilesPanel:
			if m.panels[FilesPanel].cursor < len(m.panels[FilesPanel].lines) {
				line := m.panels[FilesPanel].lines[m.panels[FilesPanel].cursor]
				parts := strings.Split(line, "\t")

				if len(parts) == 4 {
					status := parts[1]
					path := parts[3] // Always use the full path from the 4th column

					if path != "" {
						if status == "" { // It's a directory
							content, err = m.git.ShowDiff(git.DiffOptions{Color: true, Commit1: "HEAD", Commit2: path})
						} else { // It's a file
							stagedChanges := status[0] != ' ' && status[0] != '?'
							unstagedChanges := status[1] != ' '

							if stagedChanges {
								content, err = m.git.ShowDiff(git.DiffOptions{Color: true, Cached: true, Commit1: path})
							} else if unstagedChanges {
								content, err = m.git.ShowDiff(git.DiffOptions{Color: true, Commit1: path})
							} else if status == "??" {
								content = "Untracked file: Stage to see content as a diff."
							}
						}
					}
				}
			}
		case BranchesPanel:
			if m.panels[BranchesPanel].cursor < len(m.panels[BranchesPanel].lines) {
				line := m.panels[BranchesPanel].lines[m.panels[BranchesPanel].cursor]
				parts := strings.Split(line, "\t")
				if len(parts) > 1 {
					branchName := strings.TrimSpace(strings.TrimPrefix(parts[1], "(*) → "))
					content, err = m.git.ShowLog(git.LogOptions{Graph: true, Color: "always", Branch: branchName})
				}
			}
		case CommitsPanel:
			if m.panels[CommitsPanel].cursor < len(m.panels[CommitsPanel].lines) {
				line := m.panels[CommitsPanel].lines[m.panels[CommitsPanel].cursor]
				parts := strings.Split(line, "\t")
				if len(parts) >= 2 {
					sha := parts[1]
					content, err = m.git.ShowCommit(sha)
				}
			}
		case StashPanel:
			if len(m.panels[StashPanel].lines) == 1 && m.panels[StashPanel].lines[0] == "No stashed changes." {
				content = "No stashed changes."
			} else if m.panels[StashPanel].cursor < len(m.panels[StashPanel].lines) {
				line := m.panels[StashPanel].lines[m.panels[StashPanel].cursor]
				parts := strings.SplitN(line, "\t", 2)
				if len(parts) > 0 {
					stashID := parts[0]
					content, err = m.git.Stash(git.StashOptions{Show: true, StashID: stashID})
				}
			}
		}

		if err != nil {
			content = "Error: " + err.Error()
		}
		if content == "" {
			content = "Select an item to see details."
		}
		return mainContentUpdatedMsg{content: content}
	}
}

// handleWindowSizeMsg recalculates the layout and resizes all viewports on window resize.
func (m Model) handleWindowSizeMsg(msg tea.WindowSizeMsg) (Model, tea.Cmd) {
	m.width = msg.Width
	m.height = msg.Height
	m.help.Width = msg.Width
	m.helpViewport.Width = int(float64(m.width) * helpViewWidthRatio)
	m.helpViewport.Height = int(float64(m.height) * helpViewHeightRatio)

	m = m.recalculateLayout()
	return m, nil
}

// handleMouseMsg handles all mouse events, including clicks and scrolling.
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

		// Check for clicks on selectable lines first.
		for p := range m.panels {
			panel := Panel(p)
			if panel != FilesPanel && panel != BranchesPanel && panel != CommitsPanel && panel != StashPanel {
				continue
			}
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

	// Pass mouse events to the corresponding panel's viewport for scrolling.
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

// handlePanelKeys handles keybindings that are specific to the focused panel.
func (m *Model) handlePanelKeys(msg tea.KeyMsg) tea.Cmd {
	switch m.focusedPanel {
	case FilesPanel:
		return m.handleFilesPanelKeys(msg)
	case BranchesPanel:
		return m.handleBranchesPanelKeys(msg)
	case CommitsPanel:
		return m.handleCommitsPanelKeys(msg)
	case StashPanel:
		return m.handleStashPanelKeys(msg)
	}
	return nil
}

// handleCursorMovement is a helper to handle up/down cursor movement in selectable panels.
// It returns true if the key was handled.
func (m *Model) handleCursorMovement(msg tea.KeyMsg) (bool, tea.Cmd) {
	p := &m.panels[m.focusedPanel]
	itemSelected := false
	switch {
	case key.Matches(msg, keys.Up):
		if p.cursor > 0 {
			p.cursor--
			if p.cursor < p.viewport.YOffset {
				p.viewport.SetYOffset(p.cursor)
			}
			itemSelected = true
		}
	case key.Matches(msg, keys.Down):
		if p.cursor < len(p.lines)-1 {
			p.cursor++
			if p.cursor >= p.viewport.YOffset+p.viewport.Height {
				p.viewport.SetYOffset(p.cursor - p.viewport.Height + 1)
			}
			itemSelected = true
		}
	}
	if itemSelected {
		m.panels[MainPanel].viewport.GotoTop()
		return true, m.updateMainPanel()
	}
	return false, nil
}

func (m *Model) handleFilesPanelKeys(msg tea.KeyMsg) tea.Cmd {
	if handled, cmd := m.handleCursorMovement(msg); handled {
		return cmd
	}

	if m.panels[FilesPanel].cursor >= len(m.panels[FilesPanel].lines) {
		return nil
	}
	line := m.panels[FilesPanel].lines[m.panels[FilesPanel].cursor]
	parts := strings.Split(line, "\t")
	if len(parts) < 4 {
		return nil
	}
	status := parts[1]
	filePath := parts[3]

	switch {
	case key.Matches(msg, keys.Commit):
		m.mode = modeInput
		m.promptTitle = "Commit Message"
		m.textInput.Placeholder = "Enter commit message"
		m.textInput.Focus()
		m.inputCallback = func(message string) tea.Cmd {
			if message == "" {
				// Don't commit with an empty message
				return nil
			}
			return func() tea.Msg {
				_, err := m.git.Commit(git.CommitOptions{Message: message})
				if err != nil {
					return errMsg{err}
				}
				return fileWatcherMsg{}
			}
		}
		return nil
	case key.Matches(msg, keys.StageItem):
		// If the item is unstaged, stage it, and vice-versa.
		isStaged := len(status) > 0 && status[0] != ' ' && status[0] != '?'
		return func() tea.Msg {
			var err error
			if isStaged {
				_, err = m.git.ResetFiles([]string{filePath})
			} else {
				_, err = m.git.AddFiles([]string{filePath})
			}
			if err != nil {
				return errMsg{err}
			}
			return fileWatcherMsg{}
		}
	case key.Matches(msg, keys.StageAll):
		return func() tea.Msg {
			_, err := m.git.AddFiles([]string{"."})
			if err != nil {
				return errMsg{err}
			}
			return fileWatcherMsg{}
		}
	case key.Matches(msg, keys.Reset):
		return func() tea.Msg {
			_, err := m.git.ResetFiles([]string{filePath})
			if err != nil {
				return errMsg{err}
			}
			return fileWatcherMsg{}
		}
	case key.Matches(msg, keys.Discard):
		return func() tea.Msg {
			_, err := m.git.Restore(git.RestoreOptions{Paths: []string{filePath}, WorkingDir: true})
			if err != nil {
				return errMsg{err}
			}
			return fileWatcherMsg{}
		}
	}
	return nil
}

func (m *Model) handleBranchesPanelKeys(msg tea.KeyMsg) tea.Cmd {
	if handled, cmd := m.handleCursorMovement(msg); handled {
		return cmd
	}

	if m.panels[BranchesPanel].cursor >= len(m.panels[BranchesPanel].lines) {
		return nil
	}
	line := m.panels[BranchesPanel].lines[m.panels[BranchesPanel].cursor]
	parts := strings.Split(line, "\t")
	if len(parts) < 2 {
		return nil
	}
	branchName := strings.TrimSpace(strings.TrimPrefix(parts[1], "(*) → "))

	switch {
	case key.Matches(msg, keys.Checkout):
		return func() tea.Msg {
			_, err := m.git.Checkout(branchName)
			if err != nil {
				return errMsg{err}
			}
			return fileWatcherMsg{}
		}
	case key.Matches(msg, keys.DeleteBranch):
		m.mode = modeConfirm
		m.confirmMessage = fmt.Sprintf("Are you sure you want to delete branch '%s'?", branchName)
		m.confirmCallback = func(confirmed bool) tea.Cmd {
			if !confirmed {
				return nil
			}
			return func() tea.Msg {
				_, err := m.git.ManageBranch(git.BranchOptions{Delete: true, Name: branchName})
				if err != nil {
					return errMsg{err}
				}
				return fileWatcherMsg{}
			}
		}
		return nil
	}
	return nil
}

func (m *Model) handleCommitsPanelKeys(msg tea.KeyMsg) tea.Cmd {
	if handled, cmd := m.handleCursorMovement(msg); handled {
		return cmd
	}

	if m.panels[CommitsPanel].cursor >= len(m.panels[CommitsPanel].lines) {
		return nil
	}
	line := m.panels[CommitsPanel].lines[m.panels[CommitsPanel].cursor]
	parts := strings.Split(line, "\t")
	if len(parts) < 2 {
		return nil
	}
	sha := parts[1]

	switch {
	case key.Matches(msg, keys.Revert):
		return func() tea.Msg {
			_, err := m.git.Revert(sha)
			if err != nil {
				return errMsg{err}
			}
			return fileWatcherMsg{}
		}
	}
	return nil
}

func (m *Model) handleStashPanelKeys(msg tea.KeyMsg) tea.Cmd {
	if handled, cmd := m.handleCursorMovement(msg); handled {
		return cmd
	}

	if m.panels[StashPanel].cursor >= len(m.panels[StashPanel].lines) {
		return nil
	}
	line := m.panels[StashPanel].lines[m.panels[StashPanel].cursor]
	parts := strings.SplitN(line, "\t", 2)
	if len(parts) < 1 {
		return nil
	}
	stashID := parts[0]

	switch {
	case key.Matches(msg, keys.StashApply):
		return func() tea.Msg {
			_, err := m.git.Stash(git.StashOptions{Apply: true, StashID: stashID})
			if err != nil {
				return errMsg{err}
			}
			return fileWatcherMsg{}
		}
	case key.Matches(msg, keys.StashPop):
		return func() tea.Msg {
			_, err := m.git.Stash(git.StashOptions{Pop: true, StashID: stashID})
			if err != nil {
				return errMsg{err}
			}
			return fileWatcherMsg{}
		}
	case key.Matches(msg, keys.StashDrop):
		return func() tea.Msg {
			_, err := m.git.Stash(git.StashOptions{Drop: true, StashID: stashID})
			if err != nil {
				return errMsg{err}
			}
			return fileWatcherMsg{}
		}
	}
	return nil
}

// handleFocusKeys changes the focused panel based on keyboard shortcuts.
func (m *Model) handleFocusKeys(msg tea.KeyMsg) {
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
}

// recalculateLayout is the single source of truth for panel sizes and layout.
func (m Model) recalculateLayout() Model {
	if m.width == 0 || m.height == 0 {
		return m
	}

	contentHeight := m.height - 1 // Account for help bar
	m.panelHeights = make([]int, totalPanels)
	expandedHeight := int(float64(contentHeight) * expandedPanelHeightRatio)

	// Right Column Layout
	if m.focusedPanel == SecondaryPanel {
		m.panelHeights[SecondaryPanel] = expandedHeight
		m.panelHeights[MainPanel] = contentHeight - expandedHeight
	} else {
		m.panelHeights[SecondaryPanel] = collapsedPanelHeight
		m.panelHeights[MainPanel] = contentHeight - collapsedPanelHeight
	}

	// Left Column Layout
	m.panelHeights[StatusPanel] = statusPanelHeight
	remainingHeight := contentHeight - m.panelHeights[StatusPanel]

	if m.focusedPanel == StashPanel {
		m.panelHeights[StashPanel] = expandedHeight
	} else {
		m.panelHeights[StashPanel] = collapsedPanelHeight
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
		var otherPanels []Panel
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
			// Distribute remainder pixels to the last panel.
			m.panelHeights[otherPanels[len(otherPanels)-1]] += heightForOthers % len(otherPanels)
		}
	} else {
		// Default distribution when no flexible panels are focused.
		m.panelHeights[FilesPanel] = int(float64(heightForFlex) * 0.4)
		m.panelHeights[BranchesPanel] = int(float64(heightForFlex) * 0.3)
		m.panelHeights[CommitsPanel] = heightForFlex - m.panelHeights[FilesPanel] - m.panelHeights[BranchesPanel]
	}

	return m.updateViewportSizes()
}

// updateViewportSizes applies the calculated dimensions from the model to the viewports.
func (m Model) updateViewportSizes() Model {
	leftSectionWidth := int(float64(m.width) * leftPanelWidthRatio)
	rightSectionWidth := m.width - leftSectionWidth
	rightContentWidth := rightSectionWidth - borderWidth
	m.panels[MainPanel].viewport.Width = rightContentWidth
	m.panels[MainPanel].viewport.Height = m.panelHeights[MainPanel] - titleBarHeight
	m.panels[SecondaryPanel].viewport.Width = rightContentWidth
	m.panels[SecondaryPanel].viewport.Height = m.panelHeights[SecondaryPanel] - titleBarHeight

	leftContentWidth := leftSectionWidth - borderWidth
	leftPanels := []Panel{StatusPanel, FilesPanel, BranchesPanel, CommitsPanel, StashPanel}
	for _, panel := range leftPanels {
		m.panels[panel].viewport.Width = leftContentWidth
		m.panels[panel].viewport.Height = m.panelHeights[panel] - titleBarHeight
	}
	return m
}

// toggleHelp toggles the visibility of the help view.
func (m *Model) toggleHelp() {
	m.showHelp = !m.showHelp
	if m.showHelp {
		m.styleHelpViewContent()
	}
}

// styleHelpViewContent prepares and styles the content for the help view.
func (m *Model) styleHelpViewContent() {
	m.helpContent = m.generateHelpContent()
	m.helpViewport.SetContent(m.helpContent)
	m.helpViewport.GotoTop()
}
