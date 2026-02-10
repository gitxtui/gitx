package git

import (
	"fmt"
)

// RemoteOptions specifies the options for managing remotes.
type RemoteOptions struct {
	Add     bool
	Remove  bool
	Name    string
	URL     string
	Verbose bool
}

// ManageRemote manages the set of repositories ("remotes") whose branches you track.
func (g *GitCommands) ManageRemote(options RemoteOptions) (string, error) {
	args := []string{"remote"}

	if options.Verbose {
		args = append(args, "-v")
	}

	if options.Add {
		if options.Name == "" || options.URL == "" {
			return "", fmt.Errorf("remote name and URL are required for adding")
		}
		args = append(args, "add", options.Name, options.URL)
	} else if options.Remove {
		if options.Name == "" {
			return "", fmt.Errorf("remote name is required for removal")
		}
		args = append(args, "remove", options.Name)
	}

	output, _, err := g.executeCommand(args...)
	if err != nil {
		return string(output), fmt.Errorf(
			"failed to manage git remote: %w",
			err,
		)
	}

	return string(output), nil
}

// Fetch downloads objects and refs from another repository.
func (g *GitCommands) Fetch(remote string, branch string) (string, error) {
	args := []string{"fetch"}

	if remote != "" {
		args = append(args, remote)
	}

	if branch != "" {
		args = append(args, branch)
	}

	output, _, err := g.executeCommand(args...)
	if err != nil {
		return string(output), fmt.Errorf(
			"failed to fetch from remote: %w",
			err,
		)
	}

	return string(output), nil
}

// PullOptions specifies the options for the git pull command.
type PullOptions struct {
	Remote string
	Branch string
	Rebase bool
}

// Pull fetches from and integrates with another repository or a local branch.
func (g *GitCommands) Pull(options PullOptions) (string, error) {
	args := []string{"pull"}

	if options.Rebase {
		args = append(args, "--rebase")
	}

	if options.Remote != "" {
		args = append(args, options.Remote)
	}

	if options.Branch != "" {
		args = append(args, options.Branch)
	}

	output, _, err := g.executeCommand(args...)
	if err != nil {
		return string(output), fmt.Errorf(
			"failed to pull repository: %w",
			err,
		)
	}

	return string(output), nil
}

// PushOptions specifies the options for the git push command.
type PushOptions struct {
	Remote      string
	Branch      string
	Force       bool
	SetUpstream bool
	Tags        bool
}

// Push updates remote refs along with associated objects.
func (g *GitCommands) Push(options PushOptions) (string, error) {
	args := []string{"push"}

	if options.Force {
		args = append(args, "--force")
	}

	if options.SetUpstream {
		args = append(args, "--set-upstream")
	}

	if options.Tags {
		args = append(args, "--tags")
	}

	if options.Remote != "" {
		args = append(args, options.Remote)
	}

	if options.Branch != "" {
		args = append(args, options.Branch)
	}

	output, _, err := g.executeCommand(args...)
	if err != nil {
		return string(output), fmt.Errorf(
			"failed to push to remote: %w",
			err,
		)
	}

	return string(output), nil
}
