package git

import (
	"os/exec"
)

// ExecCommand is a variable that holds the exec.Command function
// This allows it to be mocked in tests
var ExecCommand = exec.Command

// GitCommands provides an interface to execute Git commands.
type GitCommands struct{}

// NewGitCommands creates a new instance of GitCommands.
func NewGitCommands() *GitCommands {
	return &GitCommands{}
}
