package git

import (
	"fmt"
	"os/exec"
)

// LogOptions specifies the options for the git log command.
type LogOptions struct {
	Oneline  bool
	Graph    bool
	MaxCount int
}

// ShowLog displays the commit logs.
func (g *GitCommands) ShowLog(options LogOptions) (string, error) {
	args := []string{"log"}

	if options.Oneline {
		args = append(args, "--oneline")
	}
	if options.Graph {
		args = append(args, "--graph")
	}
	if options.MaxCount > 0 {
		args = append(args, fmt.Sprintf("-%d", options.MaxCount))
	}

	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("failed to get log: %v", err)
	}

	return string(output), nil
}
