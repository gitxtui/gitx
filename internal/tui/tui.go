package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// App is the main application struct.
type App struct {
	program *tea.Program
}

// NewApp initializes a new TUI application.
func NewApp() *App {
	model := initialModel()
	// Use WithAltScreen to have a dedicated screen for the TUI.
	program := tea.NewProgram(model, tea.WithAltScreen())
	return &App{program: program}
}

// Run starts the TUI application.
func (a *App) Run() error {
	_, err := a.program.Run()
	return err
}
