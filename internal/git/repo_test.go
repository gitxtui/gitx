package git

import (
	"strings"
	"testing"
)

func TestGetRepoInfo(t *testing.T) {
	g := NewGitCommands()

	// This test will only pass if run inside a git repository.
	repoName, branchName, err := g.GetRepoInfo()

	if err != nil {
		t.Fatalf("GetRepoInfo() returned an error: %v", err)
	}

	if repoName == "" {
		t.Error("Expected repoName to not be empty")
	}

	if branchName == "" {
		t.Error("Expected branchName to not be empty")
	}

	// Verify with actual git commands
	expectedBranchBytes, err := ExecCommand("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		t.Fatalf("Failed to get branch name from git: %v", err)
	}
	expectedBranch := strings.TrimSpace(string(expectedBranchBytes))

	if branchName != expectedBranch {
		t.Errorf("got branch %q, want %q", branchName, expectedBranch)
	}
}
