package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

// View is the main render function for the application.
func (m Model) View() string {
	if m.showHelp {
		return m.renderHelpView()
	}
	return m.renderMainView()
}

// renderMainView renders the primary user interface.
func (m Model) renderMainView() string {
	if m.width == 0 || m.height == 0 || len(m.panelHeights) == 0 {
		return "Initializing..."
	}

	leftSectionWidth := int(float64(m.width) * 0.3)
	rightSectionWidth := m.width - leftSectionWidth
	leftpanels := []Panel{StatusPanel, FilesPanel, BranchesPanel, CommitsPanel, StashPanel}
	rightpanels := []Panel{MainPanel, SecondaryPanel}
	titles := map[Panel]string{
		MainPanel: "Main", StatusPanel: "Status", FilesPanel: "Files",
		BranchesPanel: "Branches", CommitsPanel: "Commits", StashPanel: "Stash", SecondaryPanel: "Secondary",
	}

	leftColumn := m.renderPanelColumn(leftpanels, titles, leftSectionWidth)
	rightColumn := m.renderPanelColumn(rightpanels, titles, rightSectionWidth)

	content := lipgloss.JoinHorizontal(lipgloss.Top, leftColumn, rightColumn)
	helpBar := m.renderHelpBar()
	finalView := lipgloss.JoinVertical(lipgloss.Bottom, content, helpBar)
	zone.Scan(finalView)
	return finalView
}

// renderPanelColumn renders a vertical stack of panels.
func (m Model) renderPanelColumn(panels []Panel, titles map[Panel]string, width int) string {
	var renderedPanels []string
	for _, panel := range panels {
		height := m.panelHeights[panel]
		title := titles[panel]
		renderedPanels = append(renderedPanels, m.renderPanel(title, width, height, panel))
	}
	return lipgloss.JoinVertical(lipgloss.Left, renderedPanels...)
}

// renderPanel is the single source of truth for styling panel content.
func (m Model) renderPanel(title string, width, height int, panel Panel) string {
	var borderStyle BorderStyle
	var titleStyle lipgloss.Style
	isFocused := m.focusedPanel == panel

	if isFocused {
		borderStyle = m.theme.ActiveBorder
		titleStyle = m.theme.ActiveTitle
	} else {
		borderStyle = m.theme.InactiveBorder
		titleStyle = m.theme.InactiveTitle
	}

	formattedTitle := fmt.Sprintf("[%d] %s", int(panel), title)
	p := m.panels[panel]
	content := p.content
	contentWidth := width - 2

	if panel == FilesPanel || panel == BranchesPanel || panel == CommitsPanel || panel == StashPanel {
		var builder strings.Builder
		for i, line := range p.lines {
			lineID := fmt.Sprintf("%s-line-%d", panel.ID(), i)
			var finalLine string // Use a single variable for the final output

			if i == p.cursor && isFocused {
				// --- THE CORRECTED LOGIC ---
				// 1. Clean the raw data string.
				cleanLine := strings.ReplaceAll(line, "\t", "  ")

				// 2. Create the selection style WITH the full width.
				selectionStyle := m.theme.SelectedLine.Width(contentWidth)

				// 3. Render the final line. This string is now correctly padded.
				finalLine = selectionStyle.Render(cleanLine)

			} else {
				// For unselected lines, parse, style, and then apply MaxWidth to truncate if needed.
				styledLine := styleUnselectedLine(line, panel, m.theme)
				finalLine = lipgloss.NewStyle().MaxWidth(contentWidth).Render(styledLine)
			}

			// Write the final, correctly styled/padded line to the builder.
			builder.WriteString(zone.Mark(lineID, finalLine))
			builder.WriteRune('\n')
		}
		content = strings.TrimRight(builder.String(), "\n")
	}
	p.viewport.SetContent(content)

	isScrollable := !p.viewport.AtTop() || !p.viewport.AtBottom()
	showScrollbar := isScrollable
	if panel == StashPanel || panel == SecondaryPanel {
		showScrollbar = isScrollable && isFocused
	}

	box := renderBox(
		formattedTitle, titleStyle, borderStyle, p.viewport,
		m.theme.ScrollbarThumb, width, height, showScrollbar,
	)
	return zone.Mark(panel.ID(), box)
}

// renderHelpView renders the help view.
func (m Model) renderHelpView() string {
	// For the help view, the scrollbar should always be visible if scrollable.
	showScrollbar := !m.helpViewport.AtTop() || !m.helpViewport.AtBottom()

	helpBox := renderBox(
		"Help",
		m.theme.ActiveTitle,
		m.theme.ActiveBorder,
		m.helpViewport,
		m.theme.ScrollbarThumb,
		m.helpViewport.Width,
		m.helpViewport.Height,
		showScrollbar,
	)

	centeredHelp := lipgloss.Place(m.width, m.height-1, lipgloss.Center, lipgloss.Center, helpBox)
	helpBar := m.renderHelpBar()
	return lipgloss.JoinVertical(lipgloss.Bottom, centeredHelp, helpBar)
}

// renderBox manually constructs a bordered box with a title and an integrated scrollbar.
func renderBox(title string, titleStyle lipgloss.Style, borderStyle BorderStyle, vp viewport.Model, thumbStyle lipgloss.Style, width, height int, showScrollbar bool) string {

	// 1. Get content and calculate internal dimensions.
	contentLines := strings.Split(vp.View(), "\n")
	contentWidth := width - 2   // Account for left/right borders.
	contentHeight := height - 2 // Account for top/bottom borders.
	if contentHeight < 0 {
		contentHeight = 0
	}

	// 2. Build the top border with the title embedded.
	var builder strings.Builder
	renderedTitle := titleStyle.Render(" " + title + " ")
	builder.WriteString(borderStyle.Style.Render(borderStyle.TopLeft))
	builder.WriteString(renderedTitle)
	remainingWidth := width - lipgloss.Width(renderedTitle) - 2
	if remainingWidth > 0 {
		builder.WriteString(borderStyle.Style.Render(strings.Repeat(borderStyle.Top, remainingWidth)))
	}
	builder.WriteString(borderStyle.Style.Render(borderStyle.TopRight))
	builder.WriteString("\n")

	// 3. Build the content rows with side borders and the scrollbar.
	thumbPosition := -1
	if showScrollbar {
		thumbPosition = int(float64(contentHeight-1) * vp.ScrollPercent())
	}

	for i := 0; i < contentHeight; i++ {
		builder.WriteString(borderStyle.Style.Render(borderStyle.Left))
		if i < len(contentLines) {
			builder.WriteString(lipgloss.NewStyle().MaxWidth(contentWidth).Render(contentLines[i]))
		} else {
			builder.WriteString(strings.Repeat(" ", contentWidth))
		}
		if thumbPosition == i {
			builder.WriteString(thumbStyle.Render(scrollThumb))
		} else {
			builder.WriteString(borderStyle.Style.Render(borderStyle.Right))
		}
		builder.WriteString("\n")
	}

	// 4. Build the bottom border.
	builder.WriteString(borderStyle.Style.Render(borderStyle.BottomLeft))
	builder.WriteString(borderStyle.Style.Render(strings.Repeat(borderStyle.Bottom, width-2)))
	builder.WriteString(borderStyle.Style.Render(borderStyle.BottomRight))

	return builder.String()
}

// generateHelpContent builds the formatted help string from the keymap.
func (m Model) generateHelpContent() string {
	helpSections := keys.FullHelp()
	var renderedSections []string
	for _, section := range helpSections {
		title := m.theme.HelpTitle.
			MarginLeft(9).
			Render(strings.Join([]string{"---", section.Title, "---"}, " "))
		bindings := m.renderHelpSection(section.Bindings)
		renderedSections = append(renderedSections, lipgloss.JoinVertical(lipgloss.Left, title, bindings))
	}
	return lipgloss.JoinVertical(lipgloss.Left, renderedSections...)
}

// renderHelpSection formats keybindings into a two-column layout.
func (m Model) renderHelpSection(bindings []key.Binding) string {
	var helpText string
	keyStyle := m.theme.HelpKey.Width(12).Align(lipgloss.Right).MarginRight(1)
	descStyle := lipgloss.NewStyle()
	for _, kb := range bindings {
		key := kb.Help().Key
		desc := kb.Help().Desc
		line := lipgloss.JoinHorizontal(lipgloss.Left, keyStyle.Render(key), descStyle.Render(desc))
		helpText += line + "\n"
	}
	return helpText
}

// renderHelpBar creates the help bar view.
func (m Model) renderHelpBar() string {
	var helpBindings []key.Binding
	if !m.showHelp {
		helpBindings = m.panelShortHelp()
	} else {
		helpBindings = keys.ShortHelp()
	}
	shortHelp := m.help.ShortHelpView(helpBindings)
	helpButton := m.theme.HelpButton.Render(" help:? ")
	markedButton := zone.Mark("help-button", helpButton)
	return lipgloss.JoinHorizontal(lipgloss.Left, shortHelp, markedButton)
}

// styleUnselectedLine parses a raw data line and applies panel-specific styling.
func styleUnselectedLine(line string, panel Panel, theme Theme) string {
	switch panel {
	case FilesPanel:
		parts := strings.Split(line, "\t")
		// Directory: "prefix+connector▼", "", "name"
		// File:      "prefix+connector", "status", "name"
		if len(parts) < 3 {
			return line
		}
		prefix, status, path := parts[0], parts[1], parts[2]
		if status == "" { // It's a directory
			return fmt.Sprintf("%s  %s", prefix, path)
		}
		// It's a file
		styledStatus := styleStatus(status, theme)
		return fmt.Sprintf("%s %s %s", prefix, styledStatus, path)
	case BranchesPanel:
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) != 2 {
			return line
		}
		date, name := parts[0], parts[1]
		styledDate := theme.BranchDate.Render(date)
		styledName := theme.NormalText.Render(name)
		if strings.Contains(name, "(*)") {
			styledName = theme.BranchCurrent.Render(name)
		}
		return lipgloss.JoinHorizontal(lipgloss.Left, styledDate, " ", styledName)
	case CommitsPanel:
		parts := strings.SplitN(line, "\t", 4)
		if len(parts) != 4 {
			return line // Just a graph line
		}
		graph, sha, author, subject := parts[0], parts[1], parts[2], parts[3]
		styledSHA := theme.CommitSHA.Render(sha)
		styledAuthor := theme.CommitAuthor.Render(author)
		if strings.HasPrefix(strings.ToLower(subject), "merge") {
			styledAuthor = theme.CommitMerge.Render(author)
		}
		return fmt.Sprintf("%s %s %2s %s", graph, styledSHA, styledAuthor, subject)
	case StashPanel:
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) != 2 {
			return line
		}
		name, message := parts[0], parts[1]
		styledName := theme.StashName.Render(name)
		styledMessage := theme.StashMessage.Render(message)
		return lipgloss.JoinHorizontal(lipgloss.Left, styledName, " ", styledMessage)
	}
	return line
}

// styleStatus takes a 2-character git status and returns a styled string.
func styleStatus(status string, theme Theme) string {
	if len(status) < 2 {
		return "  "
	}
	if status == "??" {
		return theme.GitUntracked.Render(status)
	}
	indexChar := status[0]
	workTreeChar := status[1]
	if indexChar == 'U' || workTreeChar == 'U' || (indexChar == 'A' && workTreeChar == 'A') || (indexChar == 'D' && workTreeChar == 'D') {
		return theme.GitConflicted.Render(status)
	}
	styledIndex := styleChar(indexChar, theme.GitStaged)
	styledWorkTree := styleChar(workTreeChar, theme.GitUnstaged)
	return styledIndex + styledWorkTree
}

func styleChar(char byte, style lipgloss.Style) string {
	if char == ' ' || char == '?' {
		return " "
	}
	return style.Render(string(char))
}
