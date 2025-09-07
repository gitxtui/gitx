package git

import (
	"fmt"
	"os/exec"
)

// CommitOptions specifies the options for the git commit command.
type CommitOptions struct {
	Message string
	Amend   bool
}

// Commit records changes to the repository.
func (g *GitCommands) Commit(options CommitOptions) (string, error) {
	if options.Message == "" && !options.Amend {
		return "", fmt.Errorf("commit message is required unless amending")
	}

	args := []string{"commit"}

	if options.Amend {
		args = append(args, "--amend")
	}

	if options.Message != "" {
		args = append(args, "-m", options.Message)
	}

	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("failed to commit changes: %v", err)
	}

	return string(output), nil
}

// ShowCommit shows the details of a specific commit.
func (g *GitCommands) ShowCommit(commitHash string) (string, error) {
	if commitHash == "" {
		commitHash = "HEAD"
	}

	cmd := exec.Command("git", "show", "--color=always", commitHash)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("failed to show commit: %v", err)
	}

	return string(output), nil
}
