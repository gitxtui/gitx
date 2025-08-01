package main

import (
	"log"

	"github.com/gitxtui/gitx/internal/git/internal/tui"
)

func main() {
	app := tui.NewApp()
	if err := app.Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}
