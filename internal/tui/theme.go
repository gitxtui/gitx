package tui

import (
	"sort"

	"github.com/charmbracelet/lipgloss"

	"os"

	"path/filepath"

	"github.com/BurntSushi/toml"

	"fmt"
)

// DefaultThemeName is the name of the default theme.
const DefaultThemeName = "GitHub Dark"

// Palette defines a set of colors for a theme.
type Palette struct {
	Black, Red, Green, Yellow, Blue, Magenta, Cyan, White,
	BrightBlack, BrightRed, BrightGreen, BrightYellow, BrightBlue, BrightMagenta, BrightCyan, BrightWhite,
	DarkBlack, DarkRed, DarkGreen, DarkYellow, DarkBlue, DarkMagenta, DarkCyan, DarkWhite,
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

		// Dark
		DarkBlack:   "#1b1f23",
		DarkRed:     "#d73a49",
		DarkGreen:   "#28a745",
		DarkYellow:  "#dbab09",
		DarkBlue:    "#2188ff",
		DarkMagenta: "#a041f5",
		DarkCyan:    "#12aab5",
		DarkWhite:   "#8b949e",

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

		// Dark
		DarkBlack:   "#1d2021",
		DarkRed:     "#9d0006",
		DarkGreen:   "#79740e",
		DarkYellow:  "#b57614",
		DarkBlue:    "#076678",
		DarkMagenta: "#8f3f71",
		DarkCyan:    "#427b58",
		DarkWhite:   "#928374",

		// Special
		Bg: "#282828",
		Fg: "#ebdbb2",
	},
}

// Theme represents the styles for different components of the UI.
type Theme struct {
	ActiveTitle    lipgloss.Style
	InactiveTitle  lipgloss.Style
	NormalText     lipgloss.Style
	HelpTitle      lipgloss.Style
	HelpKey        lipgloss.Style
	HelpButton     lipgloss.Style
	ScrollbarThumb lipgloss.Style
	SelectedLine   lipgloss.Style
	Hyperlink      lipgloss.Style
	WelcomeHeading lipgloss.Style
	WelcomeMsg     lipgloss.Style
	UserName       lipgloss.Style
	GitStaged      lipgloss.Style
	GitUnstaged    lipgloss.Style
	GitUntracked   lipgloss.Style
	GitConflicted  lipgloss.Style
	BranchCurrent  lipgloss.Style
	BranchDate     lipgloss.Style
	CommitSHA      lipgloss.Style
	CommitAuthor   lipgloss.Style
	CommitMerge    lipgloss.Style
	GraphEdge      lipgloss.Style
	GraphNode      lipgloss.Style
	GraphColors    []lipgloss.Style
	StashName      lipgloss.Style
	StashMessage   lipgloss.Style
	ActiveBorder   BorderStyle
	InactiveBorder BorderStyle
	Tree           TreeStyle
}

// BorderStyle defines the characters and styles for a panel's border.
type BorderStyle struct {
	Top         string
	Bottom      string
	Left        string
	Right       string
	TopLeft     string
	TopRight    string
	BottomLeft  string
	BottomRight string
	Style       lipgloss.Style
}

// TreeStyle defines the characters used to render the file tree.
type TreeStyle struct {
	Connector, ConnectorLast, Prefix, PrefixLast string
}

//config.toml
type themeConfig struct{
	Theme string `toml:"theme"` 
}

// custom_theme.toml
type ThemeFile struct{
	Fg	string	`toml:"fg"`
	Bg	string	`toml:"bg"`
	Normal map[string]string `toml:"normal"`
	Bright map[string]string `toml:"bright"`
	Dark map[string]string `toml:"dark"`
}

// Themes holds all the available themes, generated from palettes.
var Themes = map[string]Theme{}

func init() {
	for name, p := range Palettes {
		Themes[name] = NewThemeFromPalette(p)
	}
}

// NewThemeFromPalette creates a Theme from a given color Palette.
func NewThemeFromPalette(p Palette) Theme {
	return Theme{
		ActiveTitle:    lipgloss.NewStyle().Foreground(lipgloss.Color(p.Bg)).Background(lipgloss.Color(p.BrightCyan)),
		InactiveTitle:  lipgloss.NewStyle().Foreground(lipgloss.Color(p.Fg)).Background(lipgloss.Color(p.Black)),
		NormalText:     lipgloss.NewStyle().Foreground(lipgloss.Color(p.Fg)),
		HelpTitle:      lipgloss.NewStyle().Foreground(lipgloss.Color(p.Green)).Bold(true),
		HelpKey:        lipgloss.NewStyle().Foreground(lipgloss.Color(p.Yellow)),
		HelpButton:     lipgloss.NewStyle().Foreground(lipgloss.Color(p.Bg)).Background(lipgloss.Color(p.Green)).Margin(0, 1),
		ScrollbarThumb: lipgloss.NewStyle().Foreground(lipgloss.Color(p.BrightGreen)),
		SelectedLine:   lipgloss.NewStyle().Background(lipgloss.Color(p.DarkBlue)).Foreground(lipgloss.Color(p.BrightWhite)),
		Hyperlink:      lipgloss.NewStyle().Foreground(lipgloss.Color(p.BrightBlue)).Underline(true),
		WelcomeHeading: lipgloss.NewStyle().Foreground(lipgloss.Color(p.BrightCyan)),
		WelcomeMsg:     lipgloss.NewStyle().Foreground(lipgloss.Color(p.BrightGreen)),
		UserName:       lipgloss.NewStyle().Foreground(lipgloss.Color(p.Magenta)),
		GitStaged:      lipgloss.NewStyle().Foreground(lipgloss.Color(p.Green)),
		GitUnstaged:    lipgloss.NewStyle().Foreground(lipgloss.Color(p.Red)),
		GitUntracked:   lipgloss.NewStyle().Foreground(lipgloss.Color(p.BrightBlack)),
		GitConflicted:  lipgloss.NewStyle().Foreground(lipgloss.Color(p.BrightRed)).Bold(true),
		BranchCurrent:  lipgloss.NewStyle().Foreground(lipgloss.Color(p.Green)).Bold(true),
		BranchDate:     lipgloss.NewStyle().Foreground(lipgloss.Color(p.Yellow)),
		CommitSHA:      lipgloss.NewStyle().Foreground(lipgloss.Color(p.Yellow)),
		CommitAuthor:   lipgloss.NewStyle().Foreground(lipgloss.Color(p.Green)),
		CommitMerge:    lipgloss.NewStyle().Foreground(lipgloss.Color(p.Magenta)),
		GraphEdge:      lipgloss.NewStyle().Foreground(lipgloss.Color(p.BrightBlack)),
		GraphNode:      lipgloss.NewStyle().Foreground(lipgloss.Color(p.Green)),
		GraphColors: []lipgloss.Style{
			lipgloss.NewStyle().Foreground(lipgloss.Color(p.Green)),
			lipgloss.NewStyle().Foreground(lipgloss.Color(p.Yellow)),
			lipgloss.NewStyle().Foreground(lipgloss.Color(p.Blue)),
			lipgloss.NewStyle().Foreground(lipgloss.Color(p.Magenta)),
			lipgloss.NewStyle().Foreground(lipgloss.Color(p.Cyan)),
			lipgloss.NewStyle().Foreground(lipgloss.Color(p.BrightGreen)),
			lipgloss.NewStyle().Foreground(lipgloss.Color(p.BrightYellow)),
			lipgloss.NewStyle().Foreground(lipgloss.Color(p.BrightBlue)),
			lipgloss.NewStyle().Foreground(lipgloss.Color(p.BrightMagenta)),
			lipgloss.NewStyle().Foreground(lipgloss.Color(p.BrightCyan)),
		},
		StashName:    lipgloss.NewStyle().Foreground(lipgloss.Color(p.Yellow)),
		StashMessage: lipgloss.NewStyle().Foreground(lipgloss.Color(p.Fg)),
		ActiveBorder: BorderStyle{
			Top: borderTop, Bottom: borderBottom, Left: borderLeft, Right: borderRight,
			TopLeft: borderTopLeft, TopRight: borderTopRight, BottomLeft: borderBottomLeft, BottomRight: borderBottomRight,
			Style: lipgloss.NewStyle().Foreground(lipgloss.Color(p.BrightCyan)),
		},
		InactiveBorder: BorderStyle{
			Top: borderTop, Bottom: borderBottom, Left: borderLeft, Right: borderRight,
			TopLeft: borderTopLeft, TopRight: borderTopRight, BottomLeft: borderBottomLeft, BottomRight: borderBottomRight,
			Style: lipgloss.NewStyle().Foreground(lipgloss.Color(p.BrightBlack)),
		},
		Tree: TreeStyle{
			Connector:     treeConnector,
			ConnectorLast: treeConnectorLast,
			Prefix:        treePrefix,
			PrefixLast:    treePrefixLast,
		},
	}
}

// ThemeNames returns a slice of the available theme names.
func ThemeNames() []string {
	names := make([]string, 0, len(Palettes))
	for name := range Palettes {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func load_config() (*themeConfig, error){
	cfgPath := filepath.Join(os.Getenv("HOME"), ".config", "gitx", "config.toml")
	if _,err := os.Stat(cfgPath); os.IsNotExist(err) {
		return &themeConfig{Theme: DefaultThemeName}, nil //fallback
	}

	var cfg themeConfig
	if _, err := toml.DecodeFile(cfgPath, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func load_custom_theme(name string) (*Palette, error){
	themePath := filepath.Join(os.Getenv("HOME"), ".config", "gitx", "themes", name+".toml")
	if _,err := os.Stat(themePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("theme not found: %s", name)
	}
	
	var tf ThemeFile
	if _, err := toml.DecodeFile(themePath, &tf); err != nil {
		return nil, err
	}

	// Create a Palette from the ThemeFile
	p := Palette{
		Fg: tf.Fg,
		Bg: tf.Bg,
		Black: tf.Normal["Black"], Red: tf.Normal["Red"], Green: tf.Normal["Green"], Yellow: tf.Normal["Yellow"],
        Blue: tf.Normal["Blue"], Magenta: tf.Normal["Magenta"], Cyan: tf.Normal["Cyan"], White: tf.Normal["White"],

        BrightBlack: tf.Bright["Black"], BrightRed: tf.Bright["Red"], BrightGreen: tf.Bright["Green"], BrightYellow: tf.Bright["Yellow"],
        BrightBlue: tf.Bright["Blue"], BrightMagenta: tf.Bright["Magenta"], BrightCyan: tf.Bright["Cyan"], BrightWhite: tf.Bright["White"],

        DarkBlack: tf.Dark["Black"], DarkRed: tf.Dark["Red"], DarkGreen: tf.Dark["Green"], DarkYellow: tf.Dark["Yellow"],
        DarkBlue: tf.Dark["Blue"], DarkMagenta: tf.Dark["Magenta"], DarkCyan: tf.Dark["Cyan"], DarkWhite: tf.Dark["White"],

	}

	Palettes[name] = p // Add to Palettes map for future use
	Themes[name] = NewThemeFromPalette(p) // Add to Themes map

	return &p, nil
}
