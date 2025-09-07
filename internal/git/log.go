// Package git provides a wrapper around common git commands.
package git

import (
	"fmt"
	"strings"
	"unicode"
)

// CommitLog represents a single entry in the git log graph.
// It can be a commit or a line representing the graph structure.
type CommitLog struct {
	Graph          string // The graph structure string.
	SHA            string // The abbreviated commit hash.
	AuthorInitials string // The initials of the commit author.
	Subject        string // The subject line of the commit message.
}

// LogOptions specifies the options for the git log command.
type LogOptions struct {
	Oneline  bool
	Graph    bool
	All      bool
	MaxCount int
	Format   string
	Color    string
	Branch   string
}

// GetCommitLogsGraph fetches the git log with a graph format and returns it as a
// slice of CommitLog structs.
func (g *GitCommands) GetCommitLogsGraph() ([]CommitLog, error) {
	// A custom format with a unique delimiter is used to reliably parse the output.
	format := "<COMMIT>%h|%an|%s"
	options := LogOptions{
		Graph:  true,
		Format: format,
		Color:  "never",
		All:    true,
	}

	output, err := g.ShowLog(options)
	if err != nil {
		return nil, err
	}
	return parseCommitLogs(strings.TrimSpace(output)), nil
}

// ShowLog executes the `git log` command with the given options and returns the raw output.
func (g *GitCommands) ShowLog(options LogOptions) (string, error) {
	args := []string{"log"}

	if options.Format != "" {
		args = append(args, fmt.Sprintf("--pretty=format:%s", options.Format))
	} else if options.Oneline {
		args = append(args, "--oneline")
	}

	if options.Graph {
		args = append(args, "--graph")
	}
	if options.All {
		args = append(args, "--all")
	}
	if options.MaxCount > 0 {
		args = append(args, fmt.Sprintf("-%d", options.MaxCount))
	}
	if options.Color != "" {
		args = append(args, fmt.Sprintf("--color=%s", options.Color))
	}
	if options.Branch != "" {
		args = append(args, options.Branch)
	}

	cmd := ExecCommand("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("failed to get log: %w", err)
	}

	return string(output), nil
}

// parseCommitLogs processes the raw git log string into a slice of CommitLog structs.
func parseCommitLogs(output string) []CommitLog {
	var logs []CommitLog
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		lineWithNodeReplaced := strings.ReplaceAll(line, "*", "○")

		if strings.Contains(lineWithNodeReplaced, "<COMMIT>") {
			parts := strings.SplitN(lineWithNodeReplaced, "<COMMIT>", 2)
			graph := parts[0]
			commitData := strings.SplitN(parts[1], "|", 3)

			if len(commitData) == 3 {
				logs = append(logs, CommitLog{
					Graph:          graph,
					SHA:            commitData[0],
					AuthorInitials: getInitials(commitData[1]),
					Subject:        commitData[2],
				})
			}
		} else {
			logs = append(logs, CommitLog{Graph: lineWithNodeReplaced})
		}
	}
	return logs
}

// getInitials extracts up to two initials from a name string for concise display.
func getInitials(name string) string {
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		return ""
	}

	parts := strings.Fields(name)
	if len(parts) > 1 {
		// For "John Doe", return "JD".
		return strings.ToUpper(string(parts[0][0]) + string(parts[len(parts)-1][0]))
	}

	// For a single name like "John", return "JO".
	var initials []rune
	for _, r := range name {
		if unicode.IsLetter(r) {
			initials = append(initials, unicode.ToUpper(r))
		}
		if len(initials) == 2 {
			break
		}
	}

	if len(initials) == 1 {
		return string(initials[0])
	}
	if len(initials) > 1 {
		return string(initials[0:2])
	}

	return ""
}
