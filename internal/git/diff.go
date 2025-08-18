package git

import (
	"fmt"
	"os/exec"
)

// DiffOptions specifies the options for the git diff command.
type DiffOptions struct {
	Commit1 string
	Commit2 string
	Cached  bool
	Stat    bool
}

// ShowDiff shows changes between commits, commit and working tree, etc.
func (g *GitCommands) ShowDiff(options DiffOptions) (string, error) {
	args := []string{"diff"}

	if options.Cached {
		args = append(args, "--cached")
	}
	if options.Stat {
		args = append(args, "--stat")
	}
	if options.Commit1 != "" {
		args = append(args, options.Commit1)
	}
	if options.Commit2 != "" {
		args = append(args, options.Commit2)
	}

	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("failed to get diff: %v", err)
	}

	return string(output), nil
}
