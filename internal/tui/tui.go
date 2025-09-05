package tui

import (
	"log"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fsnotify/fsnotify"
	"github.com/gitxtui/gitx/internal/git"
)

// App is the main application struct.
type App struct {
	program *tea.Program
}

// NewApp initializes a new TUI application.
func NewApp() *App {
	m := initialModel()
	return &App{
		program: tea.NewProgram(
			m,
			tea.WithoutCatchPanics(),
			tea.WithAltScreen(),
			tea.WithMouseAllMotion(),
		),
	}
}

// Run starts the TUI application and the file watcher.
func (a *App) Run() error {
	go a.watchGitDir()
	// program.Run() returns the final model and an error. We only need the error.
	_, err := a.program.Run()
	return err
}

// watchGitDir starts a file watcher on the .git directory and sends a message on change.
func (a *App) watchGitDir() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("error creating file watcher: %v", err)
		return
	}
	defer func() {
		if err := watcher.Close(); err != nil {
			log.Printf("error closing file watcher: %v", err)
		}
	}()

	gc := git.NewGitCommands()
	gitDir, err := gc.GetGitRepoPath()
	if err != nil {
		return
	}

	repoRoot := filepath.Dir(gitDir)

	watchPaths := []string{
		repoRoot,
		gitDir,
		filepath.Join(gitDir, "HEAD"),
		filepath.Join(gitDir, "index"),
		filepath.Join(gitDir, "refs"),
	}

	for _, path := range watchPaths {
		if err := watcher.Add(path); err != nil {
			// ignore errors for paths that might not exist yet
			log.Printf("error watching path %s: %v", path, err.Error())
		}
	}

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	var needsUpdate bool

	for {
		select {
		case _, ok := <-watcher.Events: // We don't need to inspect the event
			if !ok {
				return
			}
			needsUpdate = true // Set to true on ANY event
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("file watcher error: %v", err)
		case <-ticker.C:
			if needsUpdate {
				a.program.Send(fileWatcherMsg{})
				needsUpdate = false
			}
		}
	}
}
