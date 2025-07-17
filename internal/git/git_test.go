package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change to temp dir: %v", err)
	}

	g := NewGitCommands()
	if err := g.InitRepository(""); err != nil {
		t.Fatalf("failed to initialize git repository: %v", err)
	}

	// Configure git user for commits
	if err := runGitConfig(tempDir); err != nil {
		t.Fatalf("failed to set git config: %v", err)
	}

	cleanup := func() {
		os.Chdir(originalDir)
		os.RemoveAll(tempDir)
	}

	return tempDir, cleanup
}

// createAndCommitFile creates a file with content and commits it.
func createAndCommitFile(t *testing.T, g *GitCommands, filename, content, message string) {
	t.Helper()
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create test file %s: %v", filename, err)
	}
	if err := g.AddFiles([]string{filename}); err != nil {
		t.Fatalf("failed to add file %s: %v", filename, err)
	}
	if err := g.Commit(CommitOptions{Message: message}); err != nil {
		t.Fatalf("failed to commit file %s: %v", filename, err)
	}
}

func TestNewGitCommands(t *testing.T) {
	if g := NewGitCommands(); g == nil {
		t.Error("NewGitCommands() returned nil")
	}
}

func TestGitCommands_InitRepository(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "git-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)

	g := NewGitCommands()
	repoPath := "test-repo"
	if err := g.InitRepository(repoPath); err != nil {
		t.Fatalf("InitRepository() failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(repoPath, ".git")); os.IsNotExist(err) {
		t.Errorf("expected .git directory to be created at %s", repoPath)
	}
}

func TestGitCommands_CloneRepository(t *testing.T) {
	g := NewGitCommands()
	err := g.CloneRepository("invalid-url", "")
	if err == nil {
		t.Error("CloneRepository() with invalid URL should have failed, but did not")
	}
	if !strings.Contains(err.Error(), "failed to clone repository") {
		t.Errorf("expected clone error, got: %v", err)
	}
}

func TestGitCommands_Status(t *testing.T) {
	_, cleanup := setupTestRepo(t)
	defer cleanup()

	g := NewGitCommands()

	// Test on a clean repo
	if err := g.ShowStatus(); err != nil {
		t.Errorf("ShowStatus() on clean repo failed: %v", err)
	}

	// Test with a new file
	if err := os.WriteFile("new-file.txt", []byte("content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	if err := g.ShowStatus(); err != nil {
		t.Errorf("ShowStatus() with new file failed: %v", err)
	}
}

func TestGitCommands_Log(t *testing.T) {
	_, cleanup := setupTestRepo(t)
	defer cleanup()

	g := NewGitCommands()
	createAndCommitFile(t, g, "log-test.txt", "content", "Initial commit for log test")

	if err := g.ShowLog(LogOptions{}); err != nil {
		t.Errorf("ShowLog() failed: %v", err)
	}
}

func TestGitCommands_Diff(t *testing.T) {
	_, cleanup := setupTestRepo(t)
	defer cleanup()

	g := NewGitCommands()
	createAndCommitFile(t, g, "diff-test.txt", "initial", "Initial commit for diff test")

	// Modify the file to create a diff
	if err := os.WriteFile("diff-test.txt", []byte("modified"), 0644); err != nil {
		t.Fatalf("failed to modify test file: %v", err)
	}

	if err := g.ShowDiff(DiffOptions{}); err != nil {
		t.Errorf("ShowDiff() failed: %v", err)
	}
}

func TestGitCommands_Commit(t *testing.T) {
	_, cleanup := setupTestRepo(t)
	defer cleanup()

	g := NewGitCommands()

	// Test empty commit message
	if err := g.Commit(CommitOptions{}); err == nil {
		t.Error("Commit() with empty message should fail")
	}

	// Test successful commit
	createAndCommitFile(t, g, "commit-test.txt", "content", "Successful commit")

	// Test amend
	if err := os.WriteFile("commit-test.txt", []byte("amended content"), 0644); err != nil {
		t.Fatalf("failed to amend test file: %v", err)
	}
	if err := g.AddFiles([]string{"commit-test.txt"}); err != nil {
		t.Fatalf("failed to add amended file: %v", err)
	}
	if err := g.Commit(CommitOptions{Amend: true}); err != nil {
		t.Errorf("Commit() with amend failed: %v", err)
	}
}

func TestGitCommands_BranchAndCheckout(t *testing.T) {
	_, cleanup := setupTestRepo(t)
	defer cleanup()

	g := NewGitCommands()
	createAndCommitFile(t, g, "branch-test.txt", "content", "Initial commit for branch test")

	branchName := "feature-branch"

	// Create branch
	if err := g.ManageBranch(BranchOptions{Create: true, Name: branchName}); err != nil {
		t.Fatalf("ManageBranch() create failed: %v", err)
	}

	// Checkout branch
	if err := g.Checkout(branchName); err != nil {
		t.Fatalf("Checkout() failed: %v", err)
	}

	// Switch back to main/master
	if err := g.Switch("master"); err != nil {
		t.Fatalf("Switch() failed: %v", err)
	}

	// Delete branch
	if err := g.ManageBranch(BranchOptions{Delete: true, Name: branchName}); err != nil {
		t.Fatalf("ManageBranch() delete failed: %v", err)
	}
}

func TestGitCommands_FileOperations(t *testing.T) {
	_, cleanup := setupTestRepo(t)
	defer cleanup()

	g := NewGitCommands()
	createAndCommitFile(t, g, "file-ops.txt", "content", "Initial commit for file ops")

	// Test Add
	if err := os.WriteFile("new-file.txt", []byte("new"), 0644); err != nil {
		t.Fatalf("failed to create new file: %v", err)
	}
	if err := g.AddFiles([]string{"new-file.txt"}); err != nil {
		t.Errorf("AddFiles() failed: %v", err)
	}

	// Test Reset
	if err := g.ResetFiles([]string{"new-file.txt"}); err != nil {
		t.Errorf("ResetFiles() failed: %v", err)
	}

	// Test Remove
	if err := g.RemoveFiles([]string{"file-ops.txt"}, false); err != nil {
		t.Errorf("RemoveFiles() failed: %v", err)
	}
	if _, err := os.Stat("file-ops.txt"); !os.IsNotExist(err) {
		t.Error("file should have been removed from working directory")
	}

	// Test Move
	createAndCommitFile(t, g, "source.txt", "move content", "Commit for move test")
	if err := g.MoveFile("source.txt", "destination.txt"); err != nil {
		t.Errorf("MoveFile() failed: %v", err)
	}
	if _, err := os.Stat("destination.txt"); os.IsNotExist(err) {
		t.Error("destination file does not exist after move")
	}
}

func TestGitCommands_Stash(t *testing.T) {
	_, cleanup := setupTestRepo(t)
	defer cleanup()

	g := NewGitCommands()
	createAndCommitFile(t, g, "stash-test.txt", "content", "Initial commit for stash test")

	// Modify file to create something to stash
	if err := os.WriteFile("stash-test.txt", []byte("modified"), 0644); err != nil {
		t.Fatalf("failed to modify test file: %v", err)
	}

	// Stash push
	if err := g.Stash(StashOptions{Push: true, Message: "test stash"}); err != nil {
		t.Fatalf("Stash() push failed: %v", err)
	}

	// Stash apply
	if err := g.Stash(StashOptions{Apply: true}); err != nil {
		t.Errorf("Stash() apply failed: %v", err)
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
	return cmd.Run()
}
