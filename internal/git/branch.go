package git

import (
	"fmt"
	"strings"
)

// Branch represents a git branch with its metadata.
type Branch struct {
	Name       string
	IsCurrent  bool
	LastCommit string
}

// GetBranches fetches all local branches, their last commit time, and sorts them.
func (g *GitCommands) GetBranches() ([]*Branch, error) {
	// This format gives us: <relative_commit_date> <tab> <branch_name> <tab> <is_current_indicator>
	format := "%(committerdate:relative)\t%(refname:short)\t%(HEAD)"
	args := []string{"for-each-ref", "--sort=-committerdate", "refs/heads/", fmt.Sprintf("--format=%s", format)}

	cmd := ExecCommand("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return []*Branch{}, nil // No branches found
	}

	var branches []*Branch
	var currentBranch *Branch

	for _, line := range lines {
		parts := strings.Split(line, "\t")
		if len(parts) < 2 {
			continue
		}

		isCurrent := len(parts) == 3 && parts[2] == "*"
		branch := &Branch{
			Name:       parts[1],
			IsCurrent:  isCurrent,
			LastCommit: formatRelativeDate(parts[0]),
		}

		if isCurrent {
			currentBranch = branch
		} else {
			branches = append(branches, branch)
		}
	}

	// Prepend the current branch to the list to ensure it's always at the top.
	if currentBranch != nil {
		branches = append([]*Branch{currentBranch}, branches...)
	}

	return branches, nil
}

// formatRelativeDate converts git's "X units ago" to a shorter format.
func formatRelativeDate(dateStr string) string {
	parts := strings.Split(dateStr, " ")
	if len(parts) < 2 {
		return dateStr // Return original if format is unexpected
	}

	val := parts[0]
	unit := parts[1]

	switch {
	case strings.HasPrefix(unit, "second"):
		return val + "s"
	case strings.HasPrefix(unit, "minute"):
		return val + "m"
	case strings.HasPrefix(unit, "hour"):
		return val + "h"
	case strings.HasPrefix(unit, "day"):
		return val + "d"
	case strings.HasPrefix(unit, "week"):
		return val + "w"
	case strings.HasPrefix(unit, "month"):
		return val + "M"
	case strings.HasPrefix(unit, "year"):
		return val + "y"
	default:
		return dateStr
	}
}

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

	cmd := ExecCommand("git", args...)
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

	cmd := ExecCommand("git", "checkout", branchName)
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

	cmd := ExecCommand("git", "switch", branchName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("failed to switch branch: %v", err)
	}

	return string(output), nil
}
