package tui

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines the keybindings for the application.
type KeyMap struct {
	Quit        key.Binding
	Help        key.Binding
	SwitchTheme key.Binding
	FocusNext   key.Binding
	FocusPrev   key.Binding
	FocusZero   key.Binding
	FocusOne    key.Binding
	FocusTwo    key.Binding
	FocusThree  key.Binding
	FocusFour   key.Binding
	FocusFive   key.Binding
}

// DefaultKeyMap returns a set of default keybindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		SwitchTheme: key.NewBinding(
			key.WithKeys("ctrl+t"),
			key.WithHelp("ctrl+t", "switch theme"),
		),
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
	}
}
