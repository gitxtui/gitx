package git

import (
	"fmt"
	"os/exec"
)

// TagOptions specifies the options for managing tags.
type TagOptions struct {
	Create  bool
	Delete  bool
	Name    string
	Message string
	Commit  string
}

// ManageTag creates, lists, deletes or verifies a tag object signed with GPG.
func (g *GitCommands) ManageTag(options TagOptions) (string, error) {
	args := []string{"tag"}

	if options.Delete {
		if options.Name == "" {
			return "", fmt.Errorf("tag name is required for deletion")
		}
		args = append(args, "-d", options.Name)
	} else if options.Create {
		if options.Name == "" {
			return "", fmt.Errorf("tag name is required for creation")
		}
		if options.Message != "" {
			args = append(args, "-m", options.Message)
		}
		args = append(args, options.Name)
		if options.Commit != "" {
			args = append(args, options.Commit)
		}
	}

	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("tag operation failed: %v", err)
	}

	return string(output), nil
}
