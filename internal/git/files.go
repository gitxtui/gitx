package git

import (
	"fmt"
	"os/exec"
)

// ListFiles shows information about files in the index and the working tree.
func (g *GitCommands) ListFiles() (string, error) {
	cmd := exec.Command("git", "ls-files")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("failed to list files: %v", err)
	}

	return string(output), nil
}

// BlameFile shows what revision and author last modified each line of a file.
func (g *GitCommands) BlameFile(filePath string) (string, error) {
	if filePath == "" {
		return "", fmt.Errorf("file path is required")
	}

	cmd := exec.Command("git", "blame", filePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("failed to blame file: %v", err)
	}

	return string(output), nil
}
