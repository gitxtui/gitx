package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// KeyMap stores keybindings by action name.
type KeyMap map[string]string

// HelpSection is a struct to hold a title and keybindings for a help section.
type HelpSection struct {
	Title    string
	Bindings []key.Binding
}

var keybindingDescriptions = map[string]string{
	"quit":              "quit",
	"escape":            "cancel",
	"toggle_help":       "toggle help",
	"switch_theme":      "switch theme",
	"focus_next":        "Focus Next Window",
	"focus_prev":        "Focus Previous Window",
	"focus_main":        "Focus Main Window",
	"focus_status":      "Focus Status Window",
	"focus_files":       "Focus Files Window",
	"focus_branches":    "Focus Branches Window",
	"focus_commits":     "Focus Commits Window",
	"focus_stash":       "Focus Stash Window",
	"focus_command_log": "Focus Command log Window",
	"up":                "up",
	"down":              "down",
	"stage_item":        "Stage Item",
	"stage_all":         "Stage All",
	"discard":           "Discard",
	"stash":             "Stash",
	"stash_all":         "Stash all",
	"commit":            "Commit",
	"checkout":          "Checkout",
	"new_branch":        "New Branch",
	"delete_branch":     "Delete",
	"rename_branch":     "Rename",
	"amend_commit":      "Amend",
	"revert":            "Revert",
	"reset_to_commit":   "Reset to Commit",
	"stash_apply":       "Apply",
	"stash_pop":         "Pop",
	"stash_drop":        "Drop",
}

func keySpec(keys ...string) string {
	return strings.Join(keys, ",")
}

// DefaultKeybindings returns default keybindings for each action.
func DefaultKeybindings() map[string]string {
	return map[string]string{
		"quit":              keySpec("q", "ctrl+c"),
		"escape":            keySpec("esc"),
		"toggle_help":       keySpec("?"),
		"switch_theme":      keySpec("ctrl+t"),
		"focus_next":        keySpec("tab"),
		"focus_prev":        keySpec("shift+tab"),
		"focus_main":        keySpec("0"),
		"focus_status":      keySpec("1"),
		"focus_files":       keySpec("2"),
		"focus_branches":    keySpec("3"),
		"focus_commits":     keySpec("4"),
		"focus_stash":       keySpec("5"),
		"focus_command_log": keySpec("6"),
		"up":                keySpec("k", "up"),
		"down":              keySpec("j", "down"),
		"stage_item":        keySpec("a"),
		"stage_all":         keySpec("space"),
		"discard":           keySpec("d"),
		"stash":             keySpec("s"),
		"stash_all":         keySpec("S"),
		"commit":            keySpec("c"),
		"checkout":          keySpec("enter"),
		"new_branch":        keySpec("n"),
		"delete_branch":     keySpec("d"),
		"rename_branch":     keySpec("r"),
		"amend_commit":      keySpec("A"),
		"revert":            keySpec("v"),
		"reset_to_commit":   keySpec("R"),
		"stash_apply":       keySpec("a"),
		"stash_pop":         keySpec("p"),
		"stash_drop":        keySpec("d"),
	}
}

// DefaultKeyMap returns default keybindings.
func DefaultKeyMap() KeyMap {
	defaults := DefaultKeybindings()
	result := make(KeyMap, len(defaults))
	for k, v := range defaults {
		result[k] = v
	}
	return result
}

// MergeKeybindings merges user overrides into defaults and ignores empty override values.
func MergeKeybindings(defaults, overrides map[string]string) map[string]string {
	result := make(map[string]string, len(defaults)+len(overrides))
	for k, v := range defaults {
		result[k] = v
	}
	for k, v := range overrides {
		if strings.TrimSpace(v) != "" {
			result[k] = v
		}
	}
	return result
}

// KeyMapFromConfig returns keybindings with user overrides applied on top of defaults.
func KeyMapFromConfig(overrides map[string]string) KeyMap {
	merged := MergeKeybindings(DefaultKeybindings(), overrides)
	result := make(KeyMap, len(merged))
	for k, v := range merged {
		result[k] = v
	}
	if alias, ok := result["open"]; ok && strings.TrimSpace(result["checkout"]) == "" {
		result["checkout"] = alias
	}
	return result
}

func parseConfiguredKeys(configured string) ([]string, bool) {
	parts := strings.Split(configured, ",")
	keys := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			keys = append(keys, trimmed)
		}
	}
	return keys, len(keys) > 0
}

func helpLabel(keys []string) string {
	return strings.Join(keys, "/")
}

func (k KeyMap) binding(action string) key.Binding {
	spec := strings.TrimSpace(k[action])
	if spec == "" {
		spec = DefaultKeybindings()[action]
	}
	resolvedKeys, ok := parseConfiguredKeys(spec)
	if !ok {
		resolvedKeys, _ = parseConfiguredKeys(DefaultKeybindings()[action])
	}
	desc := keybindingDescriptions[action]
	return key.NewBinding(
		key.WithKeys(resolvedKeys...),
		key.WithHelp(helpLabel(resolvedKeys), desc),
	)
}

// Matches reports whether the message matches the configured key spec.
func Matches(msg tea.KeyMsg, spec string) bool {
	resolvedKeys, ok := parseConfiguredKeys(spec)
	if !ok {
		return false
	}
	return key.Matches(msg, key.NewBinding(key.WithKeys(resolvedKeys...)))
}

func (k KeyMap) bindings(actions ...string) []key.Binding {
	result := make([]key.Binding, 0, len(actions))
	for _, action := range actions {
		result = append(result, k.binding(action))
	}
	return result
}

// FullHelp returns a structured slice of HelpSection, which is used to build
// the full help view.
func (k KeyMap) FullHelp() []HelpSection {
	return []HelpSection{
		{Title: "Navigation", Bindings: k.bindings(
			"focus_next", "focus_prev", "focus_main", "focus_status",
			"focus_files", "focus_branches", "focus_commits", "focus_stash",
			"focus_command_log", "up", "down",
		)},
		{Title: "Files", Bindings: k.bindings("commit", "stash", "stash_all", "stage_item", "stage_all", "discard")},
		{Title: "Branches", Bindings: k.bindings("checkout", "new_branch", "delete_branch", "rename_branch")},
		{Title: "Commits", Bindings: k.bindings("amend_commit", "revert", "reset_to_commit")},
		{Title: "Stash", Bindings: k.bindings("stash_apply", "stash_pop", "stash_drop")},
		{Title: "Misc", Bindings: k.bindings("switch_theme", "toggle_help", "escape", "quit")},
	}
}

// ShortHelp returns a slice of key.Binding containing help for default keybindings.
func (k KeyMap) ShortHelp() []key.Binding {
	return k.bindings("toggle_help", "escape", "quit")
}

// HelpViewHelp returns a slice of key.Binding containing help for keybindings related to Help View.
func (k KeyMap) HelpViewHelp() []key.Binding {
	return k.ShortHelp()
}

// FilesPanelHelp returns a slice of key.Binding containing help for keybindings related to Files Panel.
func (k KeyMap) FilesPanelHelp() []key.Binding {
	help := k.bindings("commit", "stash", "discard", "stage_item")
	return append(help, k.ShortHelp()...)
}

// BranchesPanelHelp returns a slice of key.Binding for the Branches Panel help bar.
func (k KeyMap) BranchesPanelHelp() []key.Binding {
	help := k.bindings("checkout", "new_branch", "delete_branch")
	return append(help, k.ShortHelp()...)
}

// CommitsPanelHelp returns a slice of key.Binding for the Commits Panel help bar.
func (k KeyMap) CommitsPanelHelp() []key.Binding {
	help := k.bindings("amend_commit", "revert", "reset_to_commit")
	return append(help, k.ShortHelp()...)
}

// StashPanelHelp returns a slice of key.Binding for the Stash Panel help bar.
func (k KeyMap) StashPanelHelp() []key.Binding {
	help := k.bindings("stash_apply", "stash_pop", "stash_drop")
	return append(help, k.ShortHelp()...)
}
