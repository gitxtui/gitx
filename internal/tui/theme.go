package tui

import "github.com/charmbracelet/lipgloss"

type Palette struct {
	Black, Red, Green, Yellow, Blue, Magenta, Cyan, White,
	BrightBlack, BrightRed, BrightGreen, BrightYellow, BrightBlue, BrightMagenta, BrightCyan, BrightWhite,
	Bg, Fg string
}

// Palettes holds all the available color palettes.
var Palettes = map[string]Palette{
	"GitHub Dark": {
		// Normal
		Black:   "#24292E",
		Red:     "#ff7b72",
		Green:   "#3fb950",
		Yellow:  "#d29922",
		Blue:    "#58a6ff",
		Magenta: "#bc8cff",
		Cyan:    "#39c5cf",
		White:   "#b1bac4",

		// Bright
		BrightBlack:   "#6e7681",
		BrightRed:     "#ffa198",
		BrightGreen:   "#56d364",
		BrightYellow:  "#e3b341",
		BrightBlue:    "#79c0ff",
		BrightMagenta: "#d2a8ff",
		BrightCyan:    "#56d4dd",
		BrightWhite:   "#f0f6fc",

		// Special
		Bg: "#0d1117",
		Fg: "#c9d1d9",
	},
	"Gruvbox": {
		// Normal
		Black:   "#282828",
		Red:     "#cc241d",
		Green:   "#98971a",
		Yellow:  "#d79921",
		Blue:    "#458588",
		Magenta: "#b16286",
		Cyan:    "#689d6a",
		White:   "#a89984",

		// Bright
		BrightBlack:   "#928374",
		BrightRed:     "#fb4934",
		BrightGreen:   "#b8bb26",
		BrightYellow:  "#fabd2f",
		BrightBlue:    "#83a598",
		BrightMagenta: "#d3869b",
		BrightCyan:    "#8ec07c",
		BrightWhite:   "#ebdbb2",

		// Special
		Bg: "#282828",
		Fg: "#ebdbb2",
	},
}

// Theme represents the styles for different components of the UI.
type Theme struct {
	ActivePanel   lipgloss.Style
	InactivePanel lipgloss.Style
	ActiveTitle   lipgloss.Style
	InactiveTitle lipgloss.Style
	NormalText    lipgloss.Style
	HelpTitle     lipgloss.Style
	HelpButton    lipgloss.Style
}

// NewThemeFromPalette creates a Theme from a Palette.
func NewThemeFromPalette(p Palette) Theme {
	return Theme{
		ActivePanel: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(p.BrightCyan)),
		InactivePanel: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(p.BrightBlack)),
		ActiveTitle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(p.Bg)).
			Background(lipgloss.Color(p.BrightCyan)),
		InactiveTitle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(p.Fg)).
			Background(lipgloss.Color(p.Black)),
		NormalText: lipgloss.NewStyle().
			Foreground(lipgloss.Color(p.Fg)),
		HelpTitle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(p.Yellow)).
			Bold(true),
		HelpButton: lipgloss.NewStyle().
			Foreground(lipgloss.Color(p.Bg)).
			Background(lipgloss.Color(p.Green)).
			Margin(0, 1),
	}
}

// Themes holds all the available themes, generated from palettes.
var Themes = map[string]Theme{}

func init() {
	for name, p := range Palettes {
		Themes[name] = NewThemeFromPalette(p)
	}
}

// ThemeNames returns a slice of the available theme names.
func ThemeNames() []string {
	names := make([]string, 0, len(Palettes))
	for name := range Palettes {
		names = append(names, name)
	}
	return names
}
