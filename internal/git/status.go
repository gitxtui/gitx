package git

import (
	"os/exec"
)

// GetStatus retrieves the git status and returns it as a string.
func (g *GitCommands) GetStatus() (string, error) {
	cmd := exec.Command("git", "status")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), err
	}
	return string(output), nil
}
