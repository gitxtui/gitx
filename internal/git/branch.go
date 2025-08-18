package git

import (
	"fmt"
	"os/exec"
)

// BranchOptions specifies the options for managing branches.
type BranchOptions struct {
	Create bool
	Delete bool
	Name   string
}

// ManageBranch creates or deletes branches.
func (g *GitCommands) ManageBranch(options BranchOptions) (string, error) {
	args := []string{"branch"}

	if options.Delete {
		if options.Name == "" {
			return "", fmt.Errorf("branch name is required for deletion")
		}
		args = append(args, "-d", options.Name)
	} else if options.Create {
		if options.Name == "" {
			return "", fmt.Errorf("branch name is required for creation")
		}
		args = append(args, options.Name)
	}

	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("branch operation failed: %v", err)
	}

	return string(output), nil
}

// Checkout switches branches or restores working tree files.
func (g *GitCommands) Checkout(branchName string) (string, error) {
	if branchName == "" {
		return "", fmt.Errorf("branch name is required")
	}

	cmd := exec.Command("git", "checkout", branchName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("failed to checkout branch: %v", err)
	}

	return string(output), nil
}

// Switch switches to a specified branch.
func (g *GitCommands) Switch(branchName string) (string, error) {
	if branchName == "" {
		return "", fmt.Errorf("branch name is required")
	}

	cmd := exec.Command("git", "switch", branchName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("failed to switch branch: %v", err)
	}

	return string(output), nil
}
