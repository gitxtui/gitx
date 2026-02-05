package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gitxtui/gitx/internal/git"
	gitxlog "github.com/gitxtui/gitx/internal/log"
	"github.com/gitxtui/gitx/internal/tui"
	zone "github.com/lrstanley/bubblezone"
)

var version = "dev"

func printHelp() {
	fmt.Println("gitx - A Git TUI Helper")
	fmt.Println()
	fmt.Println("Usage: gitx [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -v, --version    Show version information")
	fmt.Println("  -h, --help       Show this help message")
	fmt.Println("  -i, --init       Initialize a new Git repository")
	fmt.Println()
	fmt.Println("Run 'gitx' inside a Git repository to start the TUI.")
	fmt.Println("Or run 'gitx -i' to initialize a new Git repository in the current directory.")
}

func main() {
	logFile, err := gitxlog.SetupLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to set up logger: %v\n", err)
	}
	defer func() {
		if err := logFile.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close log file: %v\n", err)
		}
	}()

	// Parse flags
	shouldInit := false
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-v":
			fmt.Printf("gitx version: %s\n", version)
			return
		case "--help", "-h":
			printHelp()
			return
		case "--init", "-i":
			shouldInit = true
		}
	}

	// Ensure git repo exists (initialize if flag is set)
	if err := ensureGitRepo(shouldInit); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	zone.NewGlobal()
	defer zone.Close()

	app := tui.NewApp()

	if err := app.Run(); err != nil {
		if !errors.Is(err, tea.ErrProgramKilled) {
			log.Fatalf("error running application: %v", err)
		}
	}
	fmt.Println("Bye from gitx! :)")
}

func ensureGitRepo(shouldInit bool) error {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	if err := cmd.Run(); err == nil {
		return nil // Already inside a git repo
	}

	if !shouldInit {
		return fmt.Errorf("error: not a git repository\nrun gitx -i/--init to initialize a new git repository and open gitx")
	}

	// Initialize a new git repository
	g := &git.GitCommands{}
	_, err := g.InitRepository(".")
	if err != nil {
		return fmt.Errorf("failed to initialize git repository: %w", err)
	}

	fmt.Println("Initialized new git repository in current directory")
	return nil
}
