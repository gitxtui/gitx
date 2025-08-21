package tui

import "github.com/charmbracelet/lipgloss"

// Theme represents the styles for different components of the UI.
type Theme struct {
	ActivePanel   lipgloss.Style
	InactivePanel lipgloss.Style
	ActiveTitle   lipgloss.Style
	InactiveTitle lipgloss.Style
	NormalText    lipgloss.Style
}

// Themes holds all the available themes.
var Themes = map[string]Theme{
	"Default": {
		ActivePanel: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#cba6f7")), // Mauve
		InactivePanel: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")), // Gray
		ActiveTitle: lipgloss.NewStyle().
			Background(lipgloss.Color("#cba6f7")). // Mauve
			Foreground(lipgloss.Color("#1e1e2e")), // Base
		InactiveTitle: lipgloss.NewStyle().
			Background(lipgloss.Color("240")).     // Gray
			Foreground(lipgloss.Color("#cad3f5")), // Text
		NormalText: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#cad3f5")), // Text
	},
	"Dracula": {
		ActivePanel: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#bd93f9")), // Purple
		InactivePanel: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#6272a4")), // Comment
		ActiveTitle: lipgloss.NewStyle().
			Background(lipgloss.Color("#bd93f9")). // Purple
			Foreground(lipgloss.Color("#282a36")), // Background
		InactiveTitle: lipgloss.NewStyle().
			Background(lipgloss.Color("#6272a4")). // Comment
			Foreground(lipgloss.Color("#f8f8f2")), // Foreground
		NormalText: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f8f8f2")), // Foreground
	},
	"Nord": {
		ActivePanel: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#88c0d0")), // Frost 3
		InactivePanel: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#4c566a")), // Polar Night 3
		ActiveTitle: lipgloss.NewStyle().
			Background(lipgloss.Color("#88c0d0")). // Frost 3
			Foreground(lipgloss.Color("#2e3440")), // Polar Night 1
		InactiveTitle: lipgloss.NewStyle().
			Background(lipgloss.Color("#4c566a")). // Polar Night 3
			Foreground(lipgloss.Color("#d8dee9")), // Snow Storm 1
		NormalText: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#d8dee9")), // Snow Storm 1
	},
}

// ThemeNames returns a slice of the available theme names.
func ThemeNames() []string {
	names := make([]string, 0, len(Themes))
	for name := range Themes {
		names = append(names, name)
	}
	return names
}
