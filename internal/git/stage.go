package git

import (
	"fmt"
)

// AddFiles adds file contents to the index (staging area).
func (g *GitCommands) AddFiles(paths []string) (string, string, error) {
	if len(paths) == 0 {
		paths = []string{"."}
	}

	args := append([]string{"add"}, paths...)

	output, cmdStr, err := g.executeCommand(args...)
	if err != nil {
		return string(output), cmdStr, err
	}

	return string(output), cmdStr, nil
}

// ResetFiles resets the current HEAD to the specified state, unstaging files.
func (g *GitCommands) ResetFiles(paths []string) (string, string, error) {
	if len(paths) == 0 {
		return "", "", fmt.Errorf("at least one file path is required")
	}

	args := append([]string{"reset"}, paths...)

	output, cmdStr, err := g.executeCommand(args...)
	if err != nil {
		return string(output), cmdStr, err
	}

	return string(output), cmdStr, nil
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

	output, _, err := g.executeCommand(args...)
	if err != nil {
		return string(output), err
	}

	return string(output), nil
}

// MoveFile moves or renames a file, a directory, or a symlink.
func (g *GitCommands) MoveFile(source, destination string) (string, error) {
	if source == "" || destination == "" {
		return "", fmt.Errorf("source and destination paths are required")
	}

	args := []string{"mv", source, destination}

	output, _, err := g.executeCommand(args...)
	if err != nil {
		return string(output), err
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
func (g *GitCommands) Restore(options RestoreOptions) (string, string, error) {
	if len(options.Paths) == 0 {
		return "", "", fmt.Errorf("at least one file path is required")
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

	output, cmdStr, err := g.executeCommand(args...)
	if err != nil {
		return string(output), cmdStr, err
	}

	return string(output), cmdStr, nil
}

// Revert is used to record some new commits to reverse the effect of some earlier commits.
func (g *GitCommands) Revert(commitHash string) (string, string, error) {
	if commitHash == "" {
		return "", "", fmt.Errorf("commit hash is required")
	}

	args := []string{"revert", commitHash}

	output, cmdStr, err := g.executeCommand(args...)
	if err != nil {
		return string(output), cmdStr, err
	}

	return string(output), cmdStr, nil
}

// ResetToCommit resets the current HEAD to the specified commit.
func (g *GitCommands) ResetToCommit(commitHash string) (string, string, error) {
	if commitHash == "" {
		return "", "", fmt.Errorf("commit hash is required")
	}

	args := []string{"reset", "--hard", commitHash}

	output, cmdStr, err := g.executeCommand(args...)
	if err != nil {
		return string(output), cmdStr, err
	}

	return string(output), cmdStr, nil
}
