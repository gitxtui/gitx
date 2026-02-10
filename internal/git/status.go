package git

import (
	"fmt"
)

// StatusOptions specifies arguments for git status command.
type StatusOptions struct {
	Porcelain bool
}

// GetStatus retrieves the git status and returns it as a string.
func (g *GitCommands) GetStatus(options StatusOptions) (string, error) {
	args := []string{"status"}
	if options.Porcelain {
		args = append(args, "--porcelain")
	}

	output, _, err := g.executeCommand(args...)
	if err != nil {
		return string(output), fmt.Errorf("git status failed: %w", err)
	}
	return string(output), nil
}
