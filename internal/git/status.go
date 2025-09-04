package git

// StatusOptions specifies arguments for git status command.
type StatusOptions struct {
	Porcelain bool
}

// GetStatus retrieves the git status and returns it as a string.
func (g *GitCommands) GetStatus(options StatusOptions) (string, error) {
	args := []string{"status"}
	if options.Porcelain {
		args = append(args, "--porcelain")
	}
	cmd := ExecCommand("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), err
	}
	return string(output), nil
}
