package tui

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines the keybindings for the application.
type KeyMap struct {
	// miscellaneous keybindings
	Quit       key.Binding
	Escape     key.Binding
	ToggleHelp key.Binding

	// keybindings for changing theme
	SwitchTheme key.Binding

	// keybindings for navigation
	FocusNext  key.Binding
	FocusPrev  key.Binding
	FocusZero  key.Binding
	FocusOne   key.Binding
	FocusTwo   key.Binding
	FocusThree key.Binding
	FocusFour  key.Binding
	FocusFive  key.Binding
	FocusSix   key.Binding

	// Keybindings for FilesPanel
	StageItem key.Binding
	StageAll  key.Binding
	Discard   key.Binding
	Reset     key.Binding
	Stash     key.Binding
	StashAll  key.Binding
	Commit    key.Binding
}

// HelpSection is a struct to hold a title and keybindings for a help section.
type HelpSection struct {
	Title    string
	Bindings []key.Binding
}

// FullHelp returns a structured slice of HelpSection, which is used to build
// the full help view.
func (k KeyMap) FullHelp() []HelpSection {
	return []HelpSection{
		{
			Title: "Navigation",
			Bindings: []key.Binding{
				k.FocusNext, k.FocusPrev, k.FocusZero, k.FocusOne,
				k.FocusTwo, k.FocusThree, k.FocusFour, k.FocusFive,
				k.FocusSix,
			},
		},
		{
			Title: "Files",
			Bindings: []key.Binding{
				k.Commit, k.Stash, k.StashAll, k.StageItem,
				k.StageAll, k.Discard, k.Reset,
			},
		},
		{
			Title:    "Misc",
			Bindings: []key.Binding{k.SwitchTheme, k.ToggleHelp, k.Escape, k.Quit},
		},
	}
}

// ShortHelp returns a slice of key.Binding containing help for default keybindings.
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.ToggleHelp, k.Escape, k.Quit}
}

// HelpViewHelp returns a slice of key.Binding containing help for keybindings related to Help View.
func (k KeyMap) HelpViewHelp() []key.Binding {
	return []key.Binding{k.ToggleHelp, k.Escape, k.Quit}
}

// FilesPanelHelp returns a slice of key.Binding containing help for keybindings related to Files Panel.
func (k KeyMap) FilesPanelHelp() []key.Binding {
	help := []key.Binding{k.Commit, k.Stash, k.Discard, k.StageItem}
	return append(help, k.ShortHelp()...)
}

// DefaultKeyMap returns a set of default keybindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		// misc
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Escape: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("<esc>", "cancel"),
		),
		ToggleHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),

		// theme
		SwitchTheme: key.NewBinding(
			key.WithKeys("ctrl+t"),
			key.WithHelp("<c+t>", "switch theme"),
		),

		// navigation
		FocusNext: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "Focus Next Window"),
		),
		FocusPrev: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("<s+tab>", "Focus Previous Window"),
		),
		FocusZero: key.NewBinding(
			key.WithKeys("0"),
			key.WithHelp("0", "Focus Main Window"),
		),
		FocusOne: key.NewBinding(
			key.WithKeys("1"),
			key.WithHelp("1", "Focus Status Window"),
		),
		FocusTwo: key.NewBinding(
			key.WithKeys("2"),
			key.WithHelp("2", "Focus Files Window"),
		),
		FocusThree: key.NewBinding(
			key.WithKeys("3"),
			key.WithHelp("3", "Focus Branches Window"),
		),
		FocusFour: key.NewBinding(
			key.WithKeys("4"),
			key.WithHelp("4", "Focus Commits Window"),
		),
		FocusFive: key.NewBinding(
			key.WithKeys("5"),
			key.WithHelp("5", "Focus Stash Window"),
		),
		FocusSix: key.NewBinding(
			key.WithKeys("6"),
			key.WithHelp("6", "Focus Command log Window"),
		),

		// FilesPanel
		StageItem: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "Stage Item"),
		),
		StageAll: key.NewBinding(
			key.WithKeys("space"),
			key.WithHelp("<space>", "Stage All"),
		),
		Discard: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "Discard"),
		),
		Reset: key.NewBinding(
			key.WithKeys("D"),
			key.WithHelp("D", "Reset"),
		),
		Stash: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "Stash"),
		),
		StashAll: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "Stage all"),
		),
		Commit: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "Commit"),
		),
	}
}
