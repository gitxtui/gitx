package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// setupTestRepo creates a new temporary directory, initializes a git repository,
// and returns a cleanup function to be deferred.
func setupTestRepo(t *testing.T) (string, func()) {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "git-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	// Create a temporary home directory to isolate from global git config
	tempHome, err := os.MkdirTemp("", "git-home-")
	if err != nil {
		t.Fatalf("failed to create temp home dir: %v", err)
	}

	originalHome := os.Getenv("HOME")
	if err := os.Setenv("HOME", tempHome); err != nil {
		t.Fatalf("failed to set HOME env var: %v", err)
	}

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change to temp dir: %v", err)
	}

	g := NewGitCommands()
	if _, err := g.InitRepository(""); err != nil {
		t.Fatalf("failed to initialize git repository: %v", err)
	}

	// Configure git user for commits
	if err := runGitConfig(tempDir); err != nil {
		t.Fatalf("failed to set git config: %v", err)
	}

	// Create an initial commit to make it a "clean" repo with history.
	if err := os.WriteFile("initial.txt", []byte("initial content"), 0644); err != nil {
		t.Fatalf("failed to create initial file: %v", err)
	}
	if _, err := g.AddFiles([]string{"initial.txt"}); err != nil {
		t.Fatalf("failed to add initial file: %v", err)
	}
	if _, err := g.Commit(CommitOptions{Message: "Initial commit"}); err != nil {
		t.Fatalf("failed to create initial commit: %v", err)
	}

	cleanup := func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Fatalf("failed to change back to original directory: %v", err)
		}
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("failed to remove temp dir: %v", err)
		}
		if err := os.RemoveAll(tempHome); err != nil {
			t.Logf("failed to remove temp home dir: %v", err)
		}
		if err := os.Setenv("HOME", originalHome); err != nil {
			t.Logf("failed to restore HOME env var: %v", err)
		}
	}

	return tempDir, cleanup
}

// setupRemoteRepo creates a temporary directory, initializes a git repository in it,
// and commits a file. It returns the path to the repo and a cleanup function.
func setupRemoteRepo(t *testing.T) (string, func()) {
	t.Helper()
	remotePath, err := os.MkdirTemp("", "git-remote-")
	if err != nil {
		t.Fatalf("failed to create remote repo dir: %v", err)
	}

	// Initialize the remote repo
	cmd := exec.Command("git", "init")
	cmd.Dir = remotePath
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to init remote repo: %v", err)
	}

	// Create a file and commit it to the remote
	if err := runGitConfig(remotePath); err != nil {
		t.Fatalf("failed to set git config on remote: %v", err)
	}
	if err := os.WriteFile(filepath.Join(remotePath, "testfile.txt"), []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to create file in remote: %v", err)
	}
	cmd = exec.Command("git", "add", ".")
	cmd.Dir = remotePath
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to add file in remote: %v", err)
	}
	cmd = exec.Command("git", "commit", "-m", "initial commit")
	cmd.Dir = remotePath
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to commit in remote: %v", err)
	}

	cleanup := func() {
		if err := os.RemoveAll(remotePath); err != nil {
			t.Logf("failed to remove remote repo dir: %v", err)
		}
	}
	return remotePath, cleanup
}

// createAndCommitFile creates a file with content and commits it.
func createAndCommitFile(t *testing.T, g *GitCommands, filename, content, message string) {
	t.Helper()
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create test file %s: %v", filename, err)
	}
	if _, err := g.AddFiles([]string{filename}); err != nil {
		t.Fatalf("failed to add file %s: %v", filename, err)
	}
	if _, err := g.Commit(CommitOptions{Message: message}); err != nil {
		t.Fatalf("failed to commit file %s: %v", filename, err)
	}
}

// Helper function to set git config for tests
func runGitConfig(dir string) error {
	cmd := exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return err
	}
	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return err
	}
	// Disable GPG signing for commits
	cmd = exec.Command("git", "config", "commit.gpgsign", "false")
	cmd.Dir = dir
	return cmd.Run()
}
