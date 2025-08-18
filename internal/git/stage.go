package git

import (
	"fmt"
	"os/exec"
)

// AddFiles adds file contents to the index (staging area).
func (g *GitCommands) AddFiles(paths []string) (string, error) {
	if len(paths) == 0 {
		paths = []string{"."}
	}

	args := append([]string{"add"}, paths...)
	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("failed to add files: %v", err)
	}

	return string(output), nil
}

// ResetFiles resets the current HEAD to the specified state, unstaging files.
func (g *GitCommands) ResetFiles(paths []string) (string, error) {
	if len(paths) == 0 {
		return "", fmt.Errorf("at least one file path is required")
	}

	args := append([]string{"reset"}, paths...)
	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("failed to unstage files: %v", err)
	}

	return string(output), nil
}

// RemoveFiles removes files from the working tree and from the index.
func (g *GitCommands) RemoveFiles(paths []string, cached bool) (string, error) {
	if len(paths) == 0 {
		return "", fmt.Errorf("at least one file path is required")
	}

	args := []string{"rm"}

	if cached {
		args = append(args, "--cached")
	}

	args = append(args, paths...)

	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("failed to remove files: %v", err)
	}

	return string(output), nil
}

// MoveFile moves or renames a file, a directory, or a symlink.
func (g *GitCommands) MoveFile(source, destination string) (string, error) {
	if source == "" || destination == "" {
		return "", fmt.Errorf("source and destination paths are required")
	}

	cmd := exec.Command("git", "mv", source, destination)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("failed to move file: %v", err)
	}

	return string(output), nil
}

// RestoreOptions specifies the options for the git restore command.
type RestoreOptions struct {
	Paths      []string
	Source     string
	Staged     bool
	WorkingDir bool
}

// Restore restores working tree files.
func (g *GitCommands) Restore(options RestoreOptions) (string, error) {
	if len(options.Paths) == 0 {
		return "", fmt.Errorf("at least one file path is required")
	}

	args := []string{"restore"}

	if options.Staged {
		args = append(args, "--staged")
	}

	if options.WorkingDir {
		args = append(args, "--worktree")
	}

	if options.Source != "" {
		args = append(args, "--source", options.Source)
	}

	args = append(args, options.Paths...)

	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("failed to restore files: %v", err)
	}

	return string(output), nil
}

// Revert is used to record some new commits to reverse the effect of some earlier commits.
func (g *GitCommands) Revert(commitHash string) (string, error) {
	if commitHash == "" {
		return "", fmt.Errorf("commit hash is required")
	}

	cmd := exec.Command("git", "revert", commitHash)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("failed to revert commit: %v", err)
	}

	return string(output), nil
}
