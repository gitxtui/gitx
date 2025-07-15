package git

import (
	"os/exec"
)

// GetStatus executes `git status` and returns its output as a string.
func GetStatus() (string, error) {
	cmd := exec.Command("git", "status")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), err
	}
	return string(output), nil
}
