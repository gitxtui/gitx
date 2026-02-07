package git

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// UnsafeDirectories is a list of system directories that should not be initialized as git repositories.
var UnsafeDirectories = []string{"/", "/home", "/tmp"}

// WarnIfUnsafe checks if the given path is potentially unsafe for git initialization.
// If unsafe, it prints a warning and prompts the user for confirmation.
// Returns true if the user confirmed or the path is safe, false otherwise.
func WarnIfUnsafe(path string) (bool, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false, fmt.Errorf("failed to resolve path: %w", err)
	}

	// Check against unsafe system directories
	for _, unsafe := range UnsafeDirectories {
		if absPath == unsafe {
			return promptConfirmation(absPath, "system root directory")
		}
	}

	// Check against home directory
	homeDir, err := os.UserHomeDir()
	if err == nil && absPath == homeDir {
		return promptConfirmation(absPath, "home directory")
	}

	return true, nil // Safe to proceed
}

// promptConfirmation displays a warning and asks for user confirmation.
func promptConfirmation(path, reason string) (bool, error) {
	fmt.Printf("\n⚠️  WARNING: You are about to initialize a git repository in a %s:\n", reason)
	fmt.Printf("   Path: %s\n\n", path)
	fmt.Printf("This may not be what you intended. Initializing git here could cause issues.\n")
	fmt.Printf("Continue? [y/N]: ")

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read user input: %w", err)
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y", nil
}

// InitRepository initializes a new Git repository in the specified path.
func (g *GitCommands) InitRepository(path string) (string, error) {
	if path == "" {
		path = "."
	}
	args := []string{"init", path}

	output, _, err := g.executeCommand(args...)
	if err != nil {
		return string(output), err
	}

	absPath, _ := filepath.Abs(path)
	return fmt.Sprintf("Initialized empty Git repository in %s", absPath), nil
}
