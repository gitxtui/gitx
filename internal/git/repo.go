package git

import (
	"fmt"
	"path/filepath"
	"strings"
)

// GetRepoInfo returns the current repository and active branch name.
func (g *GitCommands) GetRepoInfo() (repoName string, branchName string, err error) {
	// Get the root dir of the repo.
	repoPath, _, err := g.executeCommand("rev-parse", "--show-toplevel")
	if err != nil {
		return "", "", fmt.Errorf("could not get repo path: %w", err)
	}
	repoPath = strings.TrimSpace(repoPath)
	repoName = filepath.Base(repoPath)

	// Get the current branch name.
	branchName, _, err = g.executeCommand("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", "", fmt.Errorf("could not get branch name: %w", err)
	}
	branchName = strings.TrimSpace(branchName)

	return repoName, branchName, nil
}

func (g *GitCommands) GetGitRepoPath() (repoPath string, err error) {
	repoPath, _, err = g.executeCommand("rev-parse", "--git-dir")
	if err != nil {
		return "", fmt.Errorf("could not get git dir path: %w", err)
	}
	repoPath = strings.TrimSpace(repoPath)
	return repoPath, nil
}

// GetUserName returns the user's name from the git config.
func (g *GitCommands) GetUserName() (string, error) {
	userName, _, err := g.executeCommand("config", "user.name")
	if err != nil {
		return "", fmt.Errorf("could not get user name: %w", err)
	}
	return strings.TrimSpace(userName), nil
}
