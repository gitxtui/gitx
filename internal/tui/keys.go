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
	Up         key.Binding
	Down       key.Binding

	// Keybindings for FilesPanel
	StageItem key.Binding
	StageAll  key.Binding
	Discard   key.Binding
	Stash     key.Binding
	StashAll  key.Binding
	Commit    key.Binding

	// Keybindings for BranchesPanel
	Checkout     key.Binding
	NewBranch    key.Binding
	DeleteBranch key.Binding
	RenameBranch key.Binding

	// Keybindings for CommitsPanel
	AmendCommit   key.Binding
	Revert        key.Binding
	ResetToCommit key.Binding

	// Keybindings for StashPanel
	StashApply key.Binding
	StashPop   key.Binding
	StashDrop  key.Binding
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
				k.FocusSix, k.Up, k.Down,
			},
		},
		{
			Title: "Files",
			Bindings: []key.Binding{
				k.Commit, k.Stash, k.StashAll, k.StageItem,
				k.StageAll, k.Discard,
			},
		},
		{
			Title:    "Branches",
			Bindings: []key.Binding{k.Checkout, k.NewBranch, k.DeleteBranch, k.RenameBranch},
		},
		{
			Title:    "Commits",
			Bindings: []key.Binding{k.AmendCommit, k.Revert, k.ResetToCommit},
		},
		{
			Title:    "Stash",
			Bindings: []key.Binding{k.StashApply, k.StashPop, k.StashDrop},
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

// BranchesPanelHelp returns a slice of key.Binding for the Branches Panel help bar.
func (k KeyMap) BranchesPanelHelp() []key.Binding {
	help := []key.Binding{k.Checkout, k.NewBranch, k.DeleteBranch}
	return append(help, k.ShortHelp()...)
}

// CommitsPanelHelp returns a slice of key.Binding for the Commits Panel help bar.
func (k KeyMap) CommitsPanelHelp() []key.Binding {
	help := []key.Binding{k.AmendCommit, k.Revert, k.ResetToCommit}
	return append(help, k.ShortHelp()...)
}

// StashPanelHelp returns a slice of key.Binding for the Stash Panel help bar.
func (k KeyMap) StashPanelHelp() []key.Binding {
	help := []key.Binding{k.StashApply, k.StashPop, k.StashDrop}
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
		Up: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("k/↑", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("j/↓", "down"),
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
		Stash: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "Stash"),
		),
		StashAll: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "Stash all"),
		),
		Commit: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "Commit"),
		),

		Checkout: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "Checkout"),
		),
		NewBranch: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "New Branch"),
		),
		DeleteBranch: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "Delete"),
		),
		RenameBranch: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "Rename"),
		),

		AmendCommit: key.NewBinding(
			key.WithKeys("A"),
			key.WithHelp("A", "Amend"),
		),
		Revert: key.NewBinding(
			key.WithKeys("v"),
			key.WithHelp("v", "Revert"),
		),
		ResetToCommit: key.NewBinding(
			key.WithKeys("R"),
			key.WithHelp("R", "Reset to Commit"),
		),

		StashApply: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "Apply"),
		),
		StashPop: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "Pop"),
		),
		StashDrop: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "Drop"),
		),
	}
}
