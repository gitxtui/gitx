package main

import (
	"errors"
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gitxtui/gitx/internal/tui"
	zone "github.com/lrstanley/bubblezone"
)

func main() {
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
