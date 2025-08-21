package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// View renders the UI.
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	// --- Layout ---
	leftSectionRenderedWidth := int(float64(m.width) * 0.3)
	rightSectionRenderedWidth := m.width - leftSectionRenderedWidth

	// --- Left Section (5 panels) ---
	leftPanelTitles := []string{"Status", "Files", "Branches", "Commits", "Stash"}
	leftPanels := m.renderVerticalPanels(
		leftPanelTitles,
		leftSectionRenderedWidth,
		m.height-1,
		[]Panel{StatusPanel, FilesPanel, BranchesPanel, CommitsPanel, StashPanel},
	)

	// --- Right Section (2 panels) ---
	rightPanelTitles := []string{"Main", "Secondary"}
	rightPanels := m.renderVerticalPanels(
		rightPanelTitles,
		rightSectionRenderedWidth,
		m.height-1,
		[]Panel{MainPanel, SecondaryPanel},
	)

	// --- Final Layout ---
	content := lipgloss.JoinHorizontal(lipgloss.Top, leftPanels, rightPanels)
	helpBindings := m.panelShortHelp()
	helpPanel := m.help.ShortHelpView(helpBindings)

	return lipgloss.JoinVertical(lipgloss.Bottom, content, helpPanel)
}

// renderVerticalPanels renders a stack of vertical panels.
func (m Model) renderVerticalPanels(titles []string, width, height int, panelTypes []Panel) string {
	panelCount := len(titles)
	if panelCount == 0 {
		return ""
	}

	availableHeight := height
	panelHeight := availableHeight / panelCount
	lastPanelHeight := availableHeight - (panelHeight * (panelCount - 1))

	var panels []string
	for i, title := range titles {
		h := panelHeight
		if i == panelCount-1 {
			h = lastPanelHeight
		}
		panels = append(panels, m.renderPanel(title, width, h, panelTypes[i]))
	}
	return lipgloss.JoinVertical(lipgloss.Left, panels...)
}

// renderPanel renders a single panel with a title bar.
func (m Model) renderPanel(title string, width, height int, panelType Panel) string {
	var panelStyle lipgloss.Style
	var titleStyle lipgloss.Style

	if m.focusedPanel == panelType {
		panelStyle = m.theme.ActivePanel
		titleStyle = m.theme.ActiveTitle
	} else {
		panelStyle = m.theme.InactivePanel
		titleStyle = m.theme.InactiveTitle
	}

	// Set the width and height for the panel style
	panelStyle = panelStyle.Width(width - panelStyle.GetHorizontalBorderSize()).Height(height - panelStyle.GetVerticalBorderSize())

	// Create the title bar
	var formattedTitle string
	if panelType == SecondaryPanel {
		formattedTitle = title
	} else {
		formattedTitle = fmt.Sprintf("[%d] %s", int(panelType), title)
	}
	titleBar := titleStyle.Width(width - panelStyle.GetHorizontalBorderSize()).Render(" " + formattedTitle)

	// Placeholder for content
	content := m.theme.NormalText.Render(fmt.Sprintf("This is the %s panel.", title))
	contentHeight := height - panelStyle.GetVerticalBorderSize() - 1 // 1 for title bar

	// Combine title bar and content
	panelContent := lipgloss.JoinVertical(lipgloss.Left, titleBar, lipgloss.Place(
		width-panelStyle.GetHorizontalBorderSize(),
		contentHeight,
		lipgloss.Left,
		lipgloss.Top,
		content,
	))

	return panelStyle.Render(panelContent)
}
