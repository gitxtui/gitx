package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gitxtui/gitx/internal/tui"
	zone "github.com/lrstanley/bubblezone"
)

var version = "dev"

func main() {
	if err := ensureGitRepo(); err != nil {
		fmt.Fprintln(os.Stderr, err) // print to stderr
		os.Exit(1)
  }
  
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("gitx version: %s\n", version)
		return
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

func ensureGitRepo() error {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error: not a git repository")
	}
	return nil
}
