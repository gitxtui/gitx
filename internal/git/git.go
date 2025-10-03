package git

import (
	"log"
	"os/exec"
	"strings"
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

// executeCommand centralizes the execution of all git commands and serves
// as a single point for logging. It takes a list of flags passed to the git
// command as arguments and returns 1. standard output, 2. the command string
// and 3. standard error
func (g *GitCommands) executeCommand(args ...string) (string, string, error) {
	cmdStr := "git " + strings.Join(args, " ")
	log.Printf("Executing command: %s", cmdStr)

	cmd := ExecCommand("git", args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("Error: %v, Output: %s", err, string(output))
		return string(output), cmdStr, err
	}

	return string(output), cmdStr, nil
}
