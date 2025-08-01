package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gitxtui/gitx/internal/git/internal/git"
)

type Model struct {
	content string
	err     error
}

type App struct {
	program *tea.Program
}

func NewApp() *App {
	status, err := git.GetStatus()
	if err != nil {
		status = err.Error()
	}

	model := Model{
		content: status,
		err:     err,
	}

	program := tea.NewProgram(model, tea.WithAltScreen())

	return &App{program: program}
}

func (a *App) Run() error {
	_, err := a.program.Run()
	return err
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(1, 2).
		Width(80).
		Height(24)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("240")).
		Padding(0, 1).
		MarginBottom(1)

	title := titleStyle.Render("Git Status")
	content := borderStyle.Render(m.content)

	return fmt.Sprintf("%s\n%s\n\nPress 'q' or 'ctrl+c' to quit", title, content)
}
