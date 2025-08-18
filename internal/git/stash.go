package git

import (
	"fmt"
	"os/exec"
)

// StashOptions specifies the options for the git stash command.
type StashOptions struct {
	Push    bool
	Pop     bool
	Apply   bool
	List    bool
	Show    bool
	Drop    bool
	Message string
	StashID string
}

// Stash saves your local modifications away and reverts the working directory to match the HEAD commit.
func (g *GitCommands) Stash(options StashOptions) (string, error) {
	if !options.Push && !options.Pop && !options.Apply && !options.List && !options.Show && !options.Drop {
		options.Push = true
	}

	var args []string

	if options.Push {
		args = []string{"stash", "push"}
		if options.Message != "" {
			args = append(args, "-m", options.Message)
		}
	} else if options.Pop {
		args = []string{"stash", "pop"}
		if options.StashID != "" {
			args = append(args, options.StashID)
		}
	} else if options.Apply {
		args = []string{"stash", "apply"}
		if options.StashID != "" {
			args = append(args, options.StashID)
		}
	} else if options.List {
		args = []string{"stash", "list"}
	} else if options.Show {
		args = []string{"stash", "show"}
		if options.StashID != "" {
			args = append(args, options.StashID)
		}
	} else if options.Drop {
		args = []string{"stash", "drop"}
		if options.StashID != "" {
			args = append(args, options.StashID)
		}
	}

	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("stash operation failed: %v", err)
	}

	return string(output), nil
}
