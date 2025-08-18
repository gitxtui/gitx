package git

import (
	"fmt"
	"os/exec"
)

// CloneRepository clones a repository from a given URL into a specified directory.
func (g *GitCommands) CloneRepository(repoURL, directory string) (string, error) {
	if repoURL == "" {
		return "", fmt.Errorf("repository URL is required")
	}

	var cmd *exec.Cmd
	if directory != "" {
		cmd = exec.Command("git", "clone", repoURL, directory)
	} else {
		cmd = exec.Command("git", "clone", repoURL)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("failed to clone repository: %v", err)
	}

	return fmt.Sprintf("Successfully cloned repository: %s", repoURL), nil
}
