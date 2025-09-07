package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// Stash represents a single entry in the git stash list.
type Stash struct {
	Name    string
	Branch  string
	Message string
}

// GetStashes fetches all stashes and returns them as a slice of Stash structs.
func (g *GitCommands) GetStashes() ([]*Stash, error) {
	// Format: stash@{0}
	// Branch: On master
	// Message: WIP on master: 52f3a6b feat: add panels
	// We use a unique delimiter to reliably parse the multi-line output for each stash.
	format := "%gD%n%gs"
	cmd := ExecCommand("git", "stash", "list", fmt.Sprintf("--format=%s", format))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	rawStashes := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(rawStashes) == 1 && rawStashes[0] == "" {
		return []*Stash{}, nil // No stashes found
	}

	var stashes []*Stash
	for i, rawStash := range rawStashes {
		parts := strings.SplitN(rawStash, ": ", 2)
		if len(parts) < 2 {
			continue // Malformed entry
		}
		stashes = append(stashes, &Stash{
			Name:    fmt.Sprintf("stash@{%d}", i),
			Branch:  parts[0],
			Message: parts[1],
		})
	}
	return stashes, nil
}

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
		args = []string{"stash", "show", "--color=always"}
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
		// The command fails if there's no stash.
		if strings.Contains(string(output), "No stash entries found") || strings.Contains(string(output), "No stash found") {
			return "No stashes found.", nil
		}
		return string(output), fmt.Errorf("stash operation failed: %v", err)
	}

	return string(output), nil
}
