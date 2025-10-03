package git

import (
	"fmt"
)

// ListFiles shows information about files in the index and the working tree.
func (g *GitCommands) ListFiles() (string, error) {
	args := []string{"ls-files"}

	output, err := g.executeCommand(args...)
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

	args := []string{"blame", filePath}
	output, err := g.executeCommand(args...)
	if err != nil {
		return string(output), fmt.Errorf("failed to blame file: %v", err)
	}

	return string(output), nil
}
