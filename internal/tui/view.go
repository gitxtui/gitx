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

// renderMainView renders the primary user interface using pre-calculated panel heights.
func (m Model) renderMainView() string {
	if m.width == 0 || m.height == 0 || len(m.panelHeights) == 0 {
		return "Initializing..."
	}

	leftSectionWidth := int(float64(m.width) * 0.3)
	rightSectionWidth := m.width - leftSectionWidth

	// Define the panels for each column.
	leftpanels := []Panel{StatusPanel, FilesPanel, BranchesPanel, CommitsPanel, StashPanel}
	rightpanels := []Panel{MainPanel, SecondaryPanel}

	// Create a map of titles for easy lookup.
	titles := map[Panel]string{
		MainPanel:      "Main",
		StatusPanel:    "Status",
		FilesPanel:     "Files",
		BranchesPanel:  "Branches",
		CommitsPanel:   "Commits",
		StashPanel:     "Stash",
		SecondaryPanel: "Secondary",
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

// renderPanel is a convenience function that calls renderBox with the correct
// styles and content for a specific panel.
func (m Model) renderPanel(title string, width, height int, panel Panel) string {
	var borderStyle BorderStyle
	var titleStyle lipgloss.Style
	isFocused := m.focusedPanel == panel

	if m.focusedPanel == panel {
		borderStyle = m.theme.ActiveBorder
		titleStyle = m.theme.ActiveTitle
	} else {
		borderStyle = m.theme.InactiveBorder
		titleStyle = m.theme.InactiveTitle
	}

	formattedTitle := fmt.Sprintf("[%d] %s", int(panel), title)
	if panel == SecondaryPanel {
		formattedTitle = title
	}

	viewport := m.panels[panel].viewport
	isScrollable := !viewport.AtTop() || !viewport.AtBottom()
	showScrollbar := isScrollable

	// For Stash and Secondary panels, only show the scrollbar when focused.
	if panel == StashPanel || panel == SecondaryPanel {
		showScrollbar = isScrollable && isFocused
	}

	box := renderBox(
		formattedTitle,
		titleStyle,
		borderStyle,
		m.panels[panel].viewport,
		m.theme.ScrollbarThumb,
		width,
		height,
		showScrollbar,
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
