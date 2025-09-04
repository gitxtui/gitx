package git

import (
	"fmt"
	"strings"
	"unicode"
)

// CommitLog represents a single line in the commit history, which could be a
// commit or just part of the graph.
type CommitLog struct {
	Graph          string
	SHA            string
	AuthorInitials string
	Subject        string
}

// GetCommitLogsGraph fetches the git log with a graph and returns a slice of CommitLog structs.
func (g *GitCommands) GetCommitLogsGraph() ([]CommitLog, error) {
	// We use a custom format with a unique delimiter "<COMMIT>" to reliably parse the output.
	format := "<COMMIT>%h|%an|%s"
	options := LogOptions{
		Graph:  true,
		Format: format,
		Color:  "never",
	}

	output, err := g.ShowLog(options)
	if err != nil {
		return nil, err
	}

	return parseCommitLogs(strings.TrimSpace(output)), nil
}

// parseCommitLogs processes the raw git log string into a slice of CommitLog structs.
func parseCommitLogs(output string) []CommitLog {
	var logs []CommitLog
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if strings.Contains(line, "<COMMIT>") {
			// This line represents a commit.
			parts := strings.SplitN(line, "<COMMIT>", 2)
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
			// This line is purely for drawing the graph.
			logs = append(logs, CommitLog{Graph: line})
		}
	}

	return logs
}

// getInitials extracts the first two letters from a name for display.
// It handles single names, multiple names, and empty strings.
func getInitials(name string) string {
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		return ""
	}

	parts := strings.Fields(name)
	if len(parts) > 1 {
		// For "John Doe", return "JD"
		return strings.ToUpper(string(parts[0][0]) + string(parts[len(parts)-1][0]))
	}

	// For "John", return "JO"
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

	return "" // Should not happen if name is not empty
}

// LogOptions specifies the options for the git log command.
type LogOptions struct {
	Oneline  bool
	Graph    bool
	Decorate bool
	MaxCount int
	Format   string
	Color    string
}

// ShowLog returns the commit logs.
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
	if options.Decorate {
		args = append(args, "--decorate")
	}
	if options.MaxCount > 0 {
		args = append(args, fmt.Sprintf("-%d", options.MaxCount))
	}
	if options.Color != "" {
		args = append(args, fmt.Sprintf("--color=%s", options.Color))
	}

	cmd := ExecCommand("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("failed to get log: %v", err)
	}

	return string(output), nil
}
