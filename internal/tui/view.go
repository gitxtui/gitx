package tui

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

// ansiRegex is used to strip ANSI escape codes from strings.
var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// stripAnsi removes ANSI escape codes from a string.
func stripAnsi(str string) string {
	return ansiRegex.ReplaceAllString(str, "")
}

// parseIntFromString converts a string to an integer, returning an error if it fails.
func parseIntFromString(s string) (int, error) {
	n, err := strconv.Atoi(strings.TrimSpace(s))
	return n, err
}

// diffLineType represents the type of a diff line.
type diffLineType int

const (
	lineTypeContext diffLineType = iota
	lineTypeAdded
	lineTypeRemoved
	lineTypeHeader
	lineTypeFileHeader
	lineTypeHunkHeader
)

// diffRow represents a single row in the structured diff.
type diffRow struct {
	lineType   diffLineType
	oldLine    string
	newLine    string
	rawLine    string // Original line for fallback
	oldLineNum int    // Line number in old file (0 if not applicable)
	newLineNum int    // Line number in new file (0 if not applicable)
}

// parseDiffStructure transforms a unified diff into structured rows suitable for split-view rendering.
func parseDiffStructure(content string) []diffRow {
	lines := strings.Split(content, "\n")
	var rows []diffRow
	var oldLineNum, newLineNum int

	for _, line := range lines {
		if len(line) == 0 {
			rows = append(rows, diffRow{lineType: lineTypeContext, oldLine: "", newLine: "", rawLine: ""})
			continue
		}

		// Always strip ANSI codes first before any parsing logic
		cleanedLine := stripAnsi(line)

		if len(cleanedLine) == 0 {
			rows = append(rows, diffRow{lineType: lineTypeContext, oldLine: "", newLine: "", rawLine: line})
			continue
		}

		firstChar := cleanedLine[0]

		// File headers
		if strings.HasPrefix(cleanedLine, "diff --git") ||
			strings.HasPrefix(cleanedLine, "index ") {
			rows = append(rows, diffRow{lineType: lineTypeFileHeader, rawLine: line})
			continue
		}

		// --- and +++ headers (but not as content)
		if strings.HasPrefix(cleanedLine, "---") && len(cleanedLine) > 3 && cleanedLine[3] == ' ' {
			rows = append(rows, diffRow{lineType: lineTypeFileHeader, rawLine: line})
			continue
		}
		if strings.HasPrefix(cleanedLine, "+++") && len(cleanedLine) > 3 && cleanedLine[3] == ' ' {
			rows = append(rows, diffRow{lineType: lineTypeFileHeader, rawLine: line})
			continue
		}

		// Hunk headers - extract the starting line numbers
		if strings.HasPrefix(cleanedLine, "@@") {
			rows = append(rows, diffRow{lineType: lineTypeHunkHeader, rawLine: line})
			// Parse hunk header format: @@ -oldStart,oldCount +newStart,newCount @@
			hunkParts := strings.Fields(cleanedLine)
			if len(hunkParts) >= 2 {
				oldPart := strings.TrimPrefix(hunkParts[1], "-")
				oldStart := strings.Split(oldPart, ",")[0]
				if num, err := parseIntFromString(oldStart); err == nil {
					oldLineNum = num
				}
				if len(hunkParts) >= 3 {
					newPart := strings.TrimPrefix(hunkParts[2], "+")
					newStart := strings.Split(newPart, ",")[0]
					if num, err := parseIntFromString(newStart); err == nil {
						newLineNum = num
					}
				}
			}
			continue
		}

		// Newline marker
		if strings.HasPrefix(cleanedLine, "\\ No newline") {
			rows = append(rows, diffRow{lineType: lineTypeFileHeader, rawLine: line})
			continue
		}

		// Added lines
		if firstChar == '+' {
			contentWithoutPrefix := cleanedLine[1:]
			row := diffRow{lineType: lineTypeAdded, newLine: contentWithoutPrefix, rawLine: line, newLineNum: newLineNum}
			rows = append(rows, row)
			newLineNum++
			continue
		}

		// Removed lines
		if firstChar == '-' {
			contentWithoutPrefix := cleanedLine[1:]
			row := diffRow{lineType: lineTypeRemoved, oldLine: contentWithoutPrefix, rawLine: line, oldLineNum: oldLineNum}
			rows = append(rows, row)
			oldLineNum++
			continue
		}

		// Context lines (start with space)
		if firstChar == ' ' {
			contentWithoutSpace := cleanedLine[1:]
			row := diffRow{lineType: lineTypeContext, oldLine: contentWithoutSpace, newLine: contentWithoutSpace, rawLine: line, oldLineNum: oldLineNum, newLineNum: newLineNum}
			rows = append(rows, row)
			oldLineNum++
			newLineNum++
			continue
		}

		// Fallback
		rows = append(rows, diffRow{lineType: lineTypeContext, oldLine: cleanedLine, newLine: cleanedLine, rawLine: line})
	}

	return rows
}

// calculateMaxLineNumber returns the maximum line number in the diff.
func calculateMaxLineNumber(rows []diffRow) int {
	max := 0
	for _, row := range rows {
		if row.oldLineNum > max {
			max = row.oldLineNum
		}
		if row.newLineNum > max {
			max = row.newLineNum
		}
	}
	return max
}

// renderSplitDiffView renders a GitHub-style split-view diff with line numbers and a separator.

func renderSplitDiffView(rows []diffRow, columnWidth int, theme Theme) string {
	if columnWidth < 20 {
		return ""
	}

	// Calculate line number width (minimum 4 chars).
	maxLineNum := calculateMaxLineNumber(rows)
	lineNumWidth := 4
	if maxLineNum > 0 {
		w := len(fmt.Sprintf("%d", maxLineNum)) + 1
		if w > lineNumWidth {
			lineNumWidth = w
		}
	}

	// Layout per half: [lineNum(lineNumWidth)] [space(1)] [content(contentColWidth)]
	// Total = lineNumWidth + 1 + contentColWidth = columnWidth
	// Reserve 1 char for spacing between line number and content
	contentColWidth := columnWidth - lineNumWidth - 1

	// Separator: exactly 1 character wide, same color as inactive border (matches panel borders), zero padding/margin.
	separatorColor := theme.InactiveBorder.Style.GetForeground()
	separatorStyle := lipgloss.NewStyle().
		Width(1).
		Foreground(separatorColor)
	sep := separatorStyle.Render("│")

	// Line number style — right-aligned, dimmed, NO extra padding.
	lineNumStyle := lipgloss.NewStyle().
		Width(lineNumWidth).
		Align(lipgloss.Right).
		Foreground(lipgloss.Color("8"))

	// renderLineNum returns a formatted, fixed-width line number string.
	renderLineNum := func(n int) string {
		if n > 0 {
			return lineNumStyle.Render(fmt.Sprintf("%d", n))
		}
		return lineNumStyle.Render("")
	}

	// Helper function to wrap text to fit within contentColWidth.
	// Tries to wrap on word boundaries; only splits inside a word when it alone
	// is longer than the available width.
	wrapText := func(text string, width int) []string {
		if width <= 0 {
			return []string{""}
		}

		// We already store diff content without ANSI codes in diffRow, but strip
		// again defensively in case future changes introduce styling here.
		cleaned := stripAnsi(text)
		runes := []rune(cleaned)
		if len(runes) <= width {
			return []string{cleaned}
		}

		words := strings.Fields(cleaned)
		if len(words) == 0 {
			return []string{""}
		}

		var (
			lines       []string
			currentLine []rune
			lineLen     int
		)

		flushLine := func() {
			if len(currentLine) > 0 {
				lines = append(lines, string(currentLine))
				currentLine = currentLine[:0]
				lineLen = 0
			}
		}

		for _, w := range words {
			wordRunes := []rune(w)
			wordLen := len(wordRunes)

			// If the word itself is longer than the width, we need to hard-wrap it.
			if wordLen > width {
				// First, flush any existing content on the current line.
				flushLine()

				for start := 0; start < wordLen; start += width {
					end := start + width
					if end > wordLen {
						end = wordLen
					}
					lines = append(lines, string(wordRunes[start:end]))
				}
				continue
			}

			// If this word (plus a space when needed) doesn't fit on the current line,
			// start a new line.
			additional := wordLen
			if lineLen > 0 {
				additional++ // space
			}
			if lineLen+additional > width {
				flushLine()
			}

			if lineLen > 0 {
				currentLine = append(currentLine, ' ')
				lineLen++
			}
			currentLine = append(currentLine, wordRunes...)
			lineLen += wordLen
		}

		flushLine()

		if len(lines) == 0 {
			return []string{""}
		}
		return lines
	}

	// Content styles: NO padding, fixed width with wrapping so the half-panel total is exact.
	// Background color is applied only; width controls the column, not padding.
	addedStyle := theme.GitStaged.
		Width(contentColWidth).
		MaxWidth(contentColWidth)

	removedStyle := theme.GitUnstaged.
		Width(contentColWidth).
		MaxWidth(contentColWidth)

	contextStyle := lipgloss.NewStyle().
		Width(contentColWidth).
		MaxWidth(contentColWidth)

	emptyStyle := lipgloss.NewStyle().
		Width(contentColWidth).
		MaxWidth(contentColWidth)

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Width(columnWidth).
		MaxWidth(columnWidth)

	// buildHalf assembles [lineNum][space][content] into a fixed-width columnWidth string.
	buildHalf := func(lineNum string, content string) string {
		// Add a space between line number and content
		spacer := " "
		return lipgloss.JoinHorizontal(lipgloss.Top, lineNum, spacer, content)
	}

	var renderedRows []string

	for _, row := range rows {
		switch row.lineType {

		case lineTypeFileHeader, lineTypeHunkHeader:
			// Span the full width (both halves + separator).
			// Left half holds the header text; right half is blank.
			left := headerStyle.Render(stripAnsi(row.rawLine))
			right := lipgloss.NewStyle().Width(columnWidth).Render("")
			gap := " "
			combined := lipgloss.JoinHorizontal(lipgloss.Top, left, gap, sep, gap, right)
			renderedRows = append(renderedRows, combined)

		case lineTypeRemoved:
			// Wrap the old line if needed
			wrappedOld := wrapText(row.oldLine, contentColWidth)
			for i, line := range wrappedOld {
				lineNum := renderLineNum(0)
				if i == 0 {
					lineNum = renderLineNum(row.oldLineNum)
				}
				left := buildHalf(lineNum, removedStyle.Render(line))
				right := buildHalf(renderLineNum(0), emptyStyle.Render(""))
				gap := " "
				renderedRows = append(renderedRows, lipgloss.JoinHorizontal(lipgloss.Top, left, gap, sep, gap, right))
			}

		case lineTypeAdded:
			// Wrap the new line if needed
			wrappedNew := wrapText(row.newLine, contentColWidth)
			for i, line := range wrappedNew {
				lineNum := renderLineNum(0)
				if i == 0 {
					lineNum = renderLineNum(row.newLineNum)
				}
				left := buildHalf(renderLineNum(0), emptyStyle.Render(""))
				right := buildHalf(lineNum, addedStyle.Render(line))
				gap := " "
				renderedRows = append(renderedRows, lipgloss.JoinHorizontal(lipgloss.Top, left, gap, sep, gap, right))
			}

		case lineTypeContext:
			// Wrap both old and new lines if needed
			wrappedOld := wrapText(row.oldLine, contentColWidth)
			wrappedNew := wrapText(row.newLine, contentColWidth)
			maxLines := len(wrappedOld)
			if len(wrappedNew) > maxLines {
				maxLines = len(wrappedNew)
			}
			for i := 0; i < maxLines; i++ {
				oldLineNum := renderLineNum(0)
				newLineNum := renderLineNum(0)
				if i == 0 {
					oldLineNum = renderLineNum(row.oldLineNum)
					newLineNum = renderLineNum(row.newLineNum)
				}
				oldLine := ""
				newLine := ""
				if i < len(wrappedOld) {
					oldLine = wrappedOld[i]
				}
				if i < len(wrappedNew) {
					newLine = wrappedNew[i]
				}
				left := buildHalf(oldLineNum, contextStyle.Render(oldLine))
				right := buildHalf(newLineNum, contextStyle.Render(newLine))
				gap := " "
				renderedRows = append(renderedRows, lipgloss.JoinHorizontal(lipgloss.Top, left, gap, sep, gap, right))
			}

		default:
			// Wrap both old and new lines if needed
			wrappedOld := wrapText(row.oldLine, contentColWidth)
			wrappedNew := wrapText(row.newLine, contentColWidth)
			maxLines := len(wrappedOld)
			if len(wrappedNew) > maxLines {
				maxLines = len(wrappedNew)
			}
			for i := 0; i < maxLines; i++ {
				oldLineNum := renderLineNum(0)
				newLineNum := renderLineNum(0)
				if i == 0 {
					oldLineNum = renderLineNum(row.oldLineNum)
					newLineNum = renderLineNum(row.newLineNum)
				}
				oldLine := ""
				newLine := ""
				if i < len(wrappedOld) {
					oldLine = wrappedOld[i]
				}
				if i < len(wrappedNew) {
					newLine = wrappedNew[i]
				}
				left := buildHalf(oldLineNum, contextStyle.Render(oldLine))
				right := buildHalf(newLineNum, contextStyle.Render(newLine))
				gap := " "
				renderedRows = append(renderedRows, lipgloss.JoinHorizontal(lipgloss.Top, left, gap, sep, gap, right))
			}
		}
	}

	return strings.Join(renderedRows, "\n")
}

// renderAdaptiveDiffView returns the appropriately formatted diff based on viewport width and view mode preference.
func renderAdaptiveDiffView(content string, width int, theme Theme, forcedViewMode *bool) string {
	if !strings.Contains(content, "diff --git") && !strings.Contains(content, "@@") {
		return content
	}

	const splitViewThreshold = 80
	const minSplitViewWidth = 40

	useSplitView := false
	if forcedViewMode != nil {
		// When forced mode is set, respect it if width allows
		if *forcedViewMode {
			// Split view requested - use it if width is sufficient
			useSplitView = width >= minSplitViewWidth
		} else {
			// Unified view explicitly requested - always use unified
			useSplitView = false
		}
	} else {
		// Auto mode - use split view only if width is sufficient
		useSplitView = width >= splitViewThreshold && width > 60
	}

	if useSplitView && width >= minSplitViewWidth {
		// Layout: [left(columnWidth)] [space] [separator] [space] [right(columnWidth)]
		// Total width = 2*columnWidth + 3
		columnWidth := (width - 3) / 2
		rows := parseDiffStructure(content)
		splitView := renderSplitDiffView(rows, columnWidth, theme)
		if splitView != "" {
			return splitView
		}
	}

	return styleDiffContent(content, theme)
}

// styleDiffContent applies visual highlighting to diff lines for better readability.
func styleDiffContent(content string, theme Theme) string {
	if !strings.Contains(content, "diff --git") && !strings.Contains(content, "@@") {
		return content
	}

	lines := strings.Split(content, "\n")

	headerStyle := lipgloss.NewStyle().Bold(true)
	addedStyle := theme.GitStaged
	removedStyle := theme.GitUnstaged

	var result []string
	for i, line := range lines {
		if len(line) == 0 {
			result = append(result, line)
			continue
		}

		cleanedLine := stripAnsi(line)

		if len(cleanedLine) == 0 {
			result = append(result, line)
			continue
		}

		firstChar := cleanedLine[0]

		if strings.HasPrefix(cleanedLine, "@@") && i > 0 && result[len(result)-1] != "" {
			result = append(result, "")
		}

		if strings.HasPrefix(cleanedLine, "diff --git") ||
			strings.HasPrefix(cleanedLine, "index ") ||
			strings.HasPrefix(cleanedLine, "---") ||
			strings.HasPrefix(cleanedLine, "+++") ||
			strings.HasPrefix(cleanedLine, "@@") {
			result = append(result, headerStyle.Render(line))
		} else if firstChar == '+' && !strings.HasPrefix(cleanedLine, "+++") {
			result = append(result, addedStyle.Render(line))
		} else if firstChar == '-' && !strings.HasPrefix(cleanedLine, "---") {
			result = append(result, removedStyle.Render(line))
		} else if firstChar == '\\' {
			result = append(result, headerStyle.Render(line))
		} else {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

// View is the main render function for the application.
func (m Model) View() string {
	var finalView string
	if m.showHelp {
		finalView = m.renderHelpView()
	} else {
		finalView = m.renderMainView()
	}

	if m.mode != modeNormal {
		var popup string
		switch m.mode {
		case modeInput:
			popup = m.renderInputPopup()
		case modeConfirm:
			popup = m.renderConfirmPopup()
		case modeCommit:
			popup = m.renderCommitPopup()
		}
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, popup)
	}

	return finalView
}

// renderInputPopup creates the view for the text input pop-up.
func (m Model) renderInputPopup() string {
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		m.theme.ActiveTitle.Render(" "+m.promptTitle+" "),
		m.textInput.View(),
		m.theme.InactiveTitle.Render(" (Enter to confirm, Esc to cancel) "),
	)

	return lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.ActiveBorder.Style.GetForeground()).
		Render(content)
}

// renderCommitPopup creates the view for the commit message pop-up.
func (m Model) renderCommitPopup() string {
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		m.theme.ActiveTitle.Render(" Commit Message "),
		m.textInput.View(),
		m.descriptionInput.View(),
		m.theme.InactiveTitle.Render(" (Tab to switch, Enter to save, Esc to cancel) "),
	)

	return lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.ActiveBorder.Style.GetForeground()).
		Render(content)
}

// renderConfirmPopup creates the view for the confirmation pop-up.
func (m Model) renderConfirmPopup() string {
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		m.confirmMessage,
		m.theme.InactiveTitle.Render(" (y/n) "),
	)

	return lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.ActiveBorder.Style.GetForeground()).
		Render(content)
}

// renderMainView renders the primary user interface with all panels.
func (m Model) renderMainView() string {
	if m.width == 0 || m.height == 0 || len(m.panelHeights) == 0 {
		return initialContentLoading
	}

	leftSectionWidth := int(float64(m.width) * leftPanelWidthRatio)
	rightSectionWidth := m.width - leftSectionWidth

	leftpanels := []Panel{StatusPanel, FilesPanel, BranchesPanel, CommitsPanel, StashPanel}
	rightpanels := []Panel{MainPanel, SecondaryPanel}

	titles := map[Panel]string{
		MainPanel: panelZero, StatusPanel: panelOne, FilesPanel: panelTwo,
		BranchesPanel: panelThree, CommitsPanel: panelFour, StashPanel: panelFive, SecondaryPanel: panelSix,
	}

	leftColumn := m.renderPanelColumn(leftpanels, titles, leftSectionWidth)
	rightColumn := m.renderPanelColumn(rightpanels, titles, rightSectionWidth)

	content := lipgloss.JoinHorizontal(lipgloss.Top, leftColumn, rightColumn)
	helpBar := m.renderHelpBar()

	finalView := lipgloss.JoinVertical(lipgloss.Bottom, content, helpBar)
	zone.Scan(finalView)
	return finalView
}

// renderPanelColumn renders a vertical stack of panels for one column.
func (m Model) renderPanelColumn(panels []Panel, titles map[Panel]string, width int) string {
	var renderedPanels []string
	for _, panel := range panels {
		height := m.panelHeights[panel]
		title := titles[panel]
		renderedPanels = append(renderedPanels, m.renderPanel(title, width, height, panel))
	}
	return lipgloss.JoinVertical(lipgloss.Left, renderedPanels...)
}

// renderPanel renders a single panel with its border, title, and content.
func (m Model) renderPanel(title string, width, height int, panel Panel) string {
	isFocused := m.focusedPanel == panel
	borderStyle := m.theme.InactiveBorder
	titleStyle := m.theme.InactiveTitle
	if isFocused {
		borderStyle = m.theme.ActiveBorder
		titleStyle = m.theme.ActiveTitle
	}

	formattedTitle := fmt.Sprintf("[%d] %s", int(panel), title)
	p := m.panels[panel]

	if panel == FilesPanel || panel == BranchesPanel || panel == CommitsPanel || panel == StashPanel {
		if len(p.lines) > 0 {
			formattedTitle = fmt.Sprintf("[%d] %s (%d/%d)", int(panel), title, p.cursor+1, len(p.lines))
		}
	}

	content := p.content
	contentWidth := width - borderWidth

	if panel == FilesPanel || panel == BranchesPanel || panel == CommitsPanel || panel == StashPanel {
		var builder strings.Builder
		for i, line := range p.lines {
			lineID := fmt.Sprintf("%s-line-%d", panel.ID(), i)
			var finalLine string

			if i == p.cursor && isFocused {
				var cleanLine string
				if panel == FilesPanel {
					parts := strings.Split(line, "\t")
					if len(parts) >= 3 {
						cleanLine = fmt.Sprintf("%s %s %s", parts[0], parts[1], parts[2])
					} else {
						cleanLine = line
					}
				} else {
					cleanLine = stripAnsi(line)
				}

				cleanLine = strings.ReplaceAll(cleanLine, "\t", "  ")
				selectionStyle := m.theme.SelectedLine.Width(contentWidth)
				finalLine = selectionStyle.Render(cleanLine)
			} else {
				styledLine := styleUnselectedLine(line, panel, m.theme)
				finalLine = lipgloss.NewStyle().MaxWidth(contentWidth).Render(styledLine)
			}

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

// renderHelpView renders the full-screen help view.
func (m Model) renderHelpView() string {
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

// renderHelpBar creates the help bar displayed at the bottom of the screen.
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

// renderBox manually constructs a bordered box with a title and an integrated scrollbar.
func renderBox(title string, titleStyle lipgloss.Style, borderStyle BorderStyle, vp viewport.Model, thumbStyle lipgloss.Style, width, height int, showScrollbar bool) string {
	contentLines := strings.Split(vp.View(), "\n")
	contentWidth := width - borderWidth
	contentHeight := height - titleBarHeight
	if contentHeight < 0 {
		contentHeight = 0
	}

	var builder strings.Builder
	renderedTitle := titleStyle.Render(" " + title + " ")
	builder.WriteString(borderStyle.Style.Render(borderStyle.TopLeft))
	builder.WriteString(renderedTitle)
	remainingWidth := width - lipgloss.Width(renderedTitle) - 2
	if remainingWidth > 0 {
		builder.WriteString(borderStyle.Style.Render(strings.Repeat(borderStyle.Top, remainingWidth)))
	}
	builder.WriteString(borderStyle.Style.Render(borderStyle.TopRight))
	builder.WriteRune('\n')

	var thumbPosition = -1
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
			builder.WriteString(thumbStyle.Render(scrollThumbChar))
		} else {
			builder.WriteString(borderStyle.Style.Render(borderStyle.Right))
		}
		builder.WriteRune('\n')
	}

	builder.WriteString(borderStyle.Style.Render(borderStyle.BottomLeft))
	builder.WriteString(borderStyle.Style.Render(strings.Repeat(borderStyle.Bottom, width-2)))
	builder.WriteString(borderStyle.Style.Render(borderStyle.BottomRight))

	return builder.String()
}

// generateHelpContent builds the formatted help string from the application's keymap.
func (m Model) generateHelpContent() string {
	helpSections := keys.FullHelp()
	var renderedSections []string
	for _, section := range helpSections {
		title := m.theme.HelpTitle.
			MarginLeft(helpTitleMargin).
			Render(strings.Join([]string{"---", section.Title, "---"}, " "))
		bindings := m.renderHelpSection(section.Bindings)
		renderedSections = append(renderedSections, lipgloss.JoinVertical(lipgloss.Left, title, bindings))
	}
	return lipgloss.JoinVertical(lipgloss.Left, renderedSections...)
}

// renderHelpSection formats a set of keybindings into a two-column layout.
func (m Model) renderHelpSection(bindings []key.Binding) string {
	var helpText string
	keyStyle := m.theme.HelpKey.Width(helpKeyWidth).Align(lipgloss.Right).MarginRight(helpDescMargin)
	descStyle := lipgloss.NewStyle()
	for _, kb := range bindings {
		key := kb.Help().Key
		desc := kb.Help().Desc
		line := lipgloss.JoinHorizontal(lipgloss.Left, keyStyle.Render(key), descStyle.Render(desc))
		helpText += line + "\n"
	}
	return helpText
}

// styleUnselectedLine parses a raw data line and applies panel-specific styling.
func styleUnselectedLine(line string, panel Panel, theme Theme) string {
	switch panel {
	case FilesPanel:
		parts := strings.Split(line, "\t")
		if len(parts) < 3 {
			return line
		}
		prefix, status, path := parts[0], parts[1], parts[2]

		var styledStatus string
		if status == "" {
			styledStatus = "  "
		} else {
			styledStatus = styleStatus(status, theme)
		}
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
			return strings.ReplaceAll(line, "○", theme.GraphNode.Render("○"))
		}
		graph, sha, author, subject := parts[0], parts[1], parts[2], parts[3]

		styledGraph := strings.ReplaceAll(graph, "○", theme.GraphNode.Render("○"))
		styledSHA := theme.CommitSHA.Render(sha)
		styledAuthor := theme.CommitAuthor.Render(author)
		if strings.HasPrefix(strings.ToLower(subject), "merge") {
			styledAuthor = theme.CommitMerge.Render(author)
		}

		final := lipgloss.JoinHorizontal(lipgloss.Left, styledSHA, " ", styledAuthor, " ", subject)
		return fmt.Sprintf("%s %s", styledGraph, final)
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

// styleStatus takes a 2-character git status code and returns a styled string.
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

// styleChar styles a single character of a status code.
func styleChar(char byte, style lipgloss.Style) string {
	if char == ' ' || char == '?' {
		return " "
	}
	return style.Render(string(char))
}
