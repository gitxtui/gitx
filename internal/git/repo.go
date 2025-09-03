package git

import (
	"os/exec"
	"path/filepath"
	"strings"
)

// GetRepoInfo returns the current repository and active branch name.
func (g *GitCommands) GetRepoInfo() (repoName string, branchName string, err error) {
	// Get the root dir of the repo.
	repoPathBytes, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", "", err
	}
	repoPath := strings.TrimSpace(string(repoPathBytes))

	repoName = filepath.Base(repoPath)

	// Get the current branch name.
	repoBranchBytes, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return "", "", err
	}
	branchName = strings.TrimSpace(string(repoBranchBytes))

	return repoName, branchName, nil
}
