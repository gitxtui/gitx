package git

import (
	"fmt"
	"path/filepath"
)

// InitRepository initializes a new Git repository in the specified path.
func (g *GitCommands) InitRepository(path string) (string, error) {
	if path == "" {
		path = "."
	}

	cmd := ExecCommand("git", "init", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("failed to initialize repository: %v", err)
	}

	absPath, _ := filepath.Abs(path)
	return fmt.Sprintf("Initialized empty Git repository in %s", absPath), nil
}
