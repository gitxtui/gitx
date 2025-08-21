package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

// View is the main render function for the application, called by the Bubble Tea
// runtime. It delegates rendering to other functions based on the application's state.
func (m Model) View() string {
	if m.showHelp {
		return m.renderHelpView()
	}
	return m.renderMainView()
}

// renderMainView renders the primary user interface, consisting of multiple panels
// and a short help bar at the bottom.
func (m Model) renderMainView() string {
	// If the terminal size has not been determined yet, show a loading message.
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	// Calculate the widths for the main left and right sections of the UI.
	leftSectionRenderedWidth := int(float64(m.width) * 0.3)
	rightSectionRenderedWidth := m.width - leftSectionRenderedWidth

	// Render the stack of panels for the left section.
	leftPanelTitles := []string{"Status", "Files", "Branches", "Commits", "Stash"}
	leftPanels := m.renderVerticalPanels(
		leftPanelTitles,
		leftSectionRenderedWidth,
		m.height-1, // Subtract 1 for the help bar at the bottom.
		[]Panel{StatusPanel, FilesPanel, BranchesPanel, CommitsPanel, StashPanel},
	)

	// Render the stack of panels for the right section.
	rightPanelTitles := []string{"Main", "Secondary"}
	rightPanels := m.renderVerticalPanels(
		rightPanelTitles,
		rightSectionRenderedWidth,
		m.height-1, // Subtract 1 for the help bar at the bottom.
		[]Panel{MainPanel, SecondaryPanel},
	)

	// Assemble the final view by joining the sections and adding the help bar.
	content := lipgloss.JoinHorizontal(lipgloss.Top, leftPanels, rightPanels)
	helpBindings := m.panelShortHelp()
	helpBar := m.help.ShortHelpView(helpBindings)

	return lipgloss.JoinVertical(lipgloss.Bottom, content, helpBar)
}

// renderHelpView renders the full-screen help menu. It centers the
// pre-rendered help content within the terminal window.
func (m Model) renderHelpView() string {
	// The viewport's content and style are set in the Update function.
	// Here, we just call its View() method to get the rendered string.
	styledHelp := m.helpViewport.View()

	// Place the rendered help content in the center of the screen.
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, styledHelp)
}

// generateHelpContent builds the complete, formatted help string from the keymap.
// This content is then used to populate the help viewport.
func (m Model) generateHelpContent() string {
	// Define titles for different sections of the help menu.
	navTitle := m.theme.HelpTitle.Render("Navigation")
	filesTitle := m.theme.HelpTitle.Render("Files")
	miscTitle := m.theme.HelpTitle.Render("Misc")

	// Render each section using the keybindings defined in the keymap.
	navHelp := m.renderHelpSection([]key.Binding{keys.FocusNext, keys.FocusPrev, keys.FocusZero, keys.FocusOne, keys.FocusTwo, keys.FocusThree, keys.FocusFour, keys.FocusFive})
	filesHelp := m.renderHelpSection([]key.Binding{keys.StageItem, keys.StageAll})
	miscHelp := m.renderHelpSection([]key.Binding{keys.ToggleHelp, keys.Quit, keys.SwitchTheme})

	// Assemble the sections into a single string for display.
	navSection := lipgloss.JoinVertical(lipgloss.Left, navTitle, navHelp)
	filesSection := lipgloss.JoinVertical(lipgloss.Left, filesTitle, filesHelp)
	miscSection := lipgloss.JoinVertical(lipgloss.Left, miscTitle, miscHelp)

	return lipgloss.JoinVertical(lipgloss.Left, navSection, "", filesSection, "", miscSection)
}

// renderHelpSection formats a slice of keybindings into a two-column layout
// (key and description) for the help menu.
func (m Model) renderHelpSection(bindings []key.Binding) string {
	var helpText string

	// Define styles for the key and description columns.
	keyStyle := lipgloss.NewStyle().Width(10).Align(lipgloss.Right).MarginRight(2)
	descStyle := lipgloss.NewStyle()

	for _, kb := range bindings {
		key := kb.Help().Key
		desc := kb.Help().Desc
		line := lipgloss.JoinHorizontal(lipgloss.Left,
			keyStyle.Render(key),
			descStyle.Render(desc),
		)
		helpText += line + "\n"
	}
	return helpText
}

// renderVerticalPanels takes a list of titles and dimensions and renders them
// as a stack of panels, distributing the available height evenly.
func (m Model) renderVerticalPanels(titles []string, width, height int, panelTypes []Panel) string {
	panelCount := len(titles)
	if panelCount == 0 {
		return ""
	}

	// Calculate the height for each panel, giving any remainder to the last one.
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

// renderPanel renders a single panel with a title bar and placeholder content.
// It applies different styles based on whether the panel is currently focused.
func (m Model) renderPanel(title string, width, height int, panelType Panel) string {
	var panelStyle lipgloss.Style
	var titleStyle lipgloss.Style

	// Apply active or inactive theme styles based on focus.
	if m.focusedPanel == panelType {
		panelStyle = m.theme.ActivePanel
		titleStyle = m.theme.ActiveTitle
	} else {
		panelStyle = m.theme.InactivePanel
		titleStyle = m.theme.InactiveTitle
	}

	// Account for border size when setting panel dimensions.
	panelStyle = panelStyle.Width(width - panelStyle.GetHorizontalBorderSize()).Height(height - panelStyle.GetVerticalBorderSize())

	// Create the title bar with the panel number and name.
	var formattedTitle string
	if panelType == SecondaryPanel {
		formattedTitle = title
	} else {
		formattedTitle = fmt.Sprintf("[%d] %s", int(panelType), title)
	}
	titleBar := titleStyle.Width(width - panelStyle.GetHorizontalBorderSize()).Render(" " + formattedTitle)

	// Placeholder for the panel's actual content.
	content := m.theme.NormalText.Render(fmt.Sprintf("This is the %s panel.", title))
	contentHeight := height - panelStyle.GetVerticalBorderSize() - 1 // Subtract 1 for the title bar.

	// Combine the title bar and content area.
	panelContent := lipgloss.JoinVertical(lipgloss.Left, titleBar, lipgloss.Place(
		width-panelStyle.GetHorizontalBorderSize(),
		contentHeight,
		lipgloss.Left,
		lipgloss.Top,
		content,
	))

	return panelStyle.Render(panelContent)
}
