package tui

import (
	"github.com/gitxtui/gitx/internal/git"
	"github.com/rivo/tview"
)

// App represents the main application
type App struct {
	*tview.Application
}

// NewApp creates and configures the main tview application
func NewApp() *App {
	// Get git status
	status, err := git.GetStatus()
	if err != nil {
		// If git command fails, return error
		status = err.Error()
	}

	// Create TUI components
	textView := tview.NewTextView().
		SetText(status).
		SetScrollable(true)

	textView.SetBorder(true).SetTitle("Git Status")

	app := tview.NewApplication().
		SetRoot(textView, true).
		SetFocus(textView)

	return &App{app}
}
