package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gitxtui/gitx/internal/tui"
	zone "github.com/lrstanley/bubblezone"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("gitx version: %s\n", version)
		return
	}

	zone.NewGlobal()
	defer zone.Close()

	app := tui.NewApp()
	if err := app.Run(); err != nil {
		if !errors.Is(err, tea.ErrProgramKilled) {
			log.Fatalf("Error running application: %v", err)
		}
	}
	fmt.Println("Bye from gitx! :)")
}
