package git

import (
	"fmt"
)

// MergeOptions specifies the options for the git merge command.
type MergeOptions struct {
	BranchName    string
	NoFastForward bool
	Message       string
}

// Merge joins two or more development histories together.
func (g *GitCommands) Merge(options MergeOptions) (string, error) {
	if options.BranchName == "" {
		return "", fmt.Errorf("branch name is required")
	}

	args := []string{"merge"}

	if options.NoFastForward {
		args = append(args, "--no-ff")
	}

	if options.Message != "" {
		args = append(args, "-m", options.Message)
	}

	args = append(args, options.BranchName)

	output, _, err := g.executeCommand(args...)
	if err != nil {
		return string(output), fmt.Errorf("failed to merge branch: %v", err)
	}

	return string(output), nil
}

// RebaseOptions specifies the options for the git rebase command.
type RebaseOptions struct {
	BranchName  string
	Interactive bool
	Abort       bool
	Continue    bool
}

// Rebase integrates changes from another branch.
func (g *GitCommands) Rebase(options RebaseOptions) (string, error) {
	args := []string{"rebase"}

	if options.Interactive {
		args = append(args, "-i")
	}
	if options.Abort {
		args = append(args, "--abort")
	}
	if options.Continue {
		args = append(args, "--continue")
	}
	if options.BranchName != "" && !options.Abort && !options.Continue {
		args = append(args, options.BranchName)
	}

	output, _, err := g.executeCommand(args...)
	if err != nil {
		return string(output), fmt.Errorf("failed to rebase branch: %v", err)
	}

	return string(output), nil
}
