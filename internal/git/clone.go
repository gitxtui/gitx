package git

import (
	"fmt"
)

// CloneRepository clones a repository from a given URL into a specified directory.
func (g *GitCommands) CloneRepository(repoURL, directory string) (string, error) {
	if repoURL == "" {
		return "", fmt.Errorf("repository URL is required")
	}

	args := []string{"clone", repoURL}
	if directory != "" {
		args = append(args, directory)
	}

	output, err := g.executeCommand(args...)
	if err != nil {
		return string(output), fmt.Errorf("failed to clone repository: %v", err)
	}

	return fmt.Sprintf("Successfully cloned repository: %s", repoURL), nil
}
