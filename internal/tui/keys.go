package tui

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines the keybindings for the application.
type KeyMap struct {
	// miscellaneous keybindings
	Quit key.Binding
	Help key.Binding

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

	// Keybindings for FilesPanel
	StageItem key.Binding
	StageAll  key.Binding
}

// FullHelp returns a nested slice of key.Binding containing
// help for all keybindings
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		// Navigation Help
		{k.FocusNext, k.FocusPrev, k.FocusZero, k.FocusOne},
		{k.FocusTwo, k.FocusThree, k.FocusFour, k.FocusFive},

		// FilesPanel help
		{k.StageItem},
		{k.StageAll},

		// Misc commands help
		{k.SwitchTheme, k.Help, k.Quit},
	}
}

// ShortHelp returns a slice of key.Binding containing help for default keybindings
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.FocusNext, k.Help, k.Quit}
}

// FilesPanelHelp returns a slice of key.Binding containing help for keybindings related to Files Panel
func (k KeyMap) FilesPanelHelp() []key.Binding {
	help := []key.Binding{k.StageItem, k.StageAll}
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
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),

		// theme
		SwitchTheme: key.NewBinding(
			key.WithKeys("ctrl+t"),
			key.WithHelp("ctrl+t", "switch theme"),
		),

		// navigation
		FocusNext: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "Focus Next Window"),
		),
		FocusPrev: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "Focus Previous Window"),
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

		// FilesPanel
		StageItem: key.NewBinding(
			key.WithKeys("space"),
			key.WithHelp("space", "Stage/Unstage Item"),
		),
		StageAll: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "Stage/Unstage All"),
		),
	}
}
