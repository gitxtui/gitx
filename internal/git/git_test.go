package git

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("failed to remove temp dir: %v", err)
		}
	}()

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change to temp dir: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Fatalf("failed to change back to original directory: %v", err)
		}
	}()

	g := NewGitCommands()
	repoPath := "test-repo"
	output, err := g.InitRepository(repoPath)
	if err != nil {
		t.Fatalf("InitRepository() failed: %v", err)
	}

	if !strings.Contains(output, "Initialized empty Git repository") {
		t.Errorf("expected success message, got: %s", output)
	}

	if _, err := os.Stat(filepath.Join(repoPath, ".git")); os.IsNotExist(err) {
		t.Errorf("expected .git directory to be created at %s", repoPath)
	}
}

func TestGitCommands_CloneRepository(t *testing.T) {
	// 1. Arrange: Set up a source "remote" repository
	remotePath, cleanup := setupRemoteRepo(t)
	defer cleanup()

	// Set up a destination directory for the clone
	localPath, err := os.MkdirTemp("", "git-local-")
	if err != nil {
		t.Fatalf("failed to create local repo dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(localPath); err != nil {
			t.Logf("failed to remove local repo dir: %v", err)
		}
	}()

	// 2. Act: Perform the clone
	g := NewGitCommands()
	output, err := g.CloneRepository(remotePath, localPath)
	if err != nil {
		t.Fatalf("CloneRepository() failed: %v", err)
	}

	// 3. Assert: Verify the results
	if !strings.Contains(output, "Successfully cloned repository") {
		t.Errorf("expected success message, got: %s", output)
	}
	if _, err := os.Stat(filepath.Join(localPath, ".git")); os.IsNotExist(err) {
		t.Error("expected .git directory to exist in cloned repo")
	}
	if _, err := os.Stat(filepath.Join(localPath, "testfile.txt")); os.IsNotExist(err) {
		t.Error("expected cloned file to exist in cloned repo")
	}

	// Test failure case
	_, err = g.CloneRepository("invalid-url", "")
	if err == nil {
		t.Error("CloneRepository() with invalid URL should have failed, but did not")
	}
}

func TestGitCommands_Status(t *testing.T) {
	_, cleanup := setupTestRepo(t)
	defer cleanup()

	g := NewGitCommands()

	// Test on a clean repo
	status, err := g.GetStatus(StatusOptions{Porcelain: false})
	if err != nil {
		t.Errorf("GetStatus() on clean repo failed: %v", err)
	}
	if !strings.Contains(status, "nothing to commit, working tree clean") {
		t.Errorf("expected clean status, got: %s", status)
	}

	// Test with a new file
	if err := os.WriteFile("new-file.txt", []byte("content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	status, err = g.GetStatus(StatusOptions{Porcelain: true})
	if err != nil {
		t.Errorf("GetStatus() with new file failed: %v", err)
	}
	if !strings.Contains(status, "?? new-file.txt") {
		t.Errorf("expected untracked file status, got: %s", status)
	}
}

func TestGitCommands_Log(t *testing.T) {
	_, cleanup := setupTestRepo(t)
	defer cleanup()

	g := NewGitCommands()
	commitMessage := "Initial commit for log test"
	createAndCommitFile(t, g, "log-test.txt", "content", commitMessage)

	log, err := g.ShowLog(LogOptions{})
	if err != nil {
		t.Errorf("ShowLog() failed: %v", err)
	}
	if !strings.Contains(log, commitMessage) {
		t.Errorf("expected log to contain commit message, got: %s", log)
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

	diff, err := g.ShowDiff(DiffOptions{})
	if err != nil {
		t.Errorf("ShowDiff() failed: %v", err)
	}
	if !strings.Contains(diff, "+modified") {
		t.Errorf("expected diff to show added line, got: %s", diff)
	}
}

func TestGitCommands_Commit(t *testing.T) {
	_, cleanup := setupTestRepo(t)
	defer cleanup()

	g := NewGitCommands()

	// Test empty commit message
	if _, err := g.Commit(CommitOptions{}); err == nil {
		t.Error("Commit() with empty message should fail")
	}

	// Test successful commit
	createAndCommitFile(t, g, "commit-test.txt", "content", "Successful commit")

	// Test amend
	if err := os.WriteFile("commit-test.txt", []byte("amended content"), 0644); err != nil {
		t.Fatalf("failed to amend test file: %v", err)
	}
	if _, err := g.AddFiles([]string{"commit-test.txt"}); err != nil {
		t.Fatalf("failed to add amended file: %v", err)
	}
	if _, err := g.Commit(CommitOptions{Amend: true, Message: "Amended commit"}); err != nil {
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
	if _, err := g.ManageBranch(BranchOptions{Create: true, Name: branchName}); err != nil {
		t.Fatalf("ManageBranch() create failed: %v", err)
	}

	// Checkout branch
	if _, err := g.Checkout(branchName); err != nil {
		t.Fatalf("Checkout() failed: %v", err)
	}

	// Switch back to main/master
	if _, err := g.Switch("master"); err != nil {
		t.Fatalf("Switch() failed: %v", err)
	}

	// Delete branch
	if _, err := g.ManageBranch(BranchOptions{Delete: true, Name: branchName}); err != nil {
		t.Fatalf("ManageBranch() delete failed: %v", err)
	}
}

func TestGitCommands_Merge(t *testing.T) {
	// Setup: Create a repo with two branches and diverging commits
	_, cleanup := setupTestRepo(t)
	defer cleanup()
	g := NewGitCommands()

	// Create and commit on master
	createAndCommitFile(t, g, "master.txt", "master content", "master commit")

	// Create feature branch and commit
	branchName := "feature"
	if _, err := g.ManageBranch(BranchOptions{Create: true, Name: branchName}); err != nil {
		t.Fatalf("failed to create branch: %v", err)
	}
	if _, err := g.Checkout(branchName); err != nil {
		t.Fatalf("failed to checkout branch: %v", err)
	}
	createAndCommitFile(t, g, "feature.txt", "feature content", "feature commit")

	// Switch back to master and make another commit
	if _, err := g.Checkout("master"); err != nil {
		t.Fatalf("failed to checkout master: %v", err)
	}
	createAndCommitFile(t, g, "master2.txt", "master2 content", "master2 commit")

	// Merge feature branch into master
	output, err := g.Merge(MergeOptions{BranchName: branchName})
	if err != nil {
		t.Fatalf("Merge() failed: %v\nOutput: %s", err, output)
	}
	if !strings.Contains(output, "Merge made by") && !strings.Contains(output, "Already up to date") && !strings.Contains(output, "Fast-forward") {
		t.Errorf("expected merge output, got: %s", output)
	}
}

func TestGitCommands_Rebase(t *testing.T) {
	// Setup: Create a repo with two branches and diverging commits
	_, cleanup := setupTestRepo(t)
	defer cleanup()
	g := NewGitCommands()

	// Create and commit on master
	createAndCommitFile(t, g, "master.txt", "master content", "master commit")

	// Create feature branch and commit
	branchName := "feature"
	if _, err := g.ManageBranch(BranchOptions{Create: true, Name: branchName}); err != nil {
		t.Fatalf("failed to create branch: %v", err)
	}
	if _, err := g.Checkout(branchName); err != nil {
		t.Fatalf("failed to checkout branch: %v", err)
	}
	createAndCommitFile(t, g, "feature.txt", "feature content", "feature commit")

	// Switch back to master and make another commit
	if _, err := g.Checkout("master"); err != nil {
		t.Fatalf("failed to checkout master: %v", err)
	}
	createAndCommitFile(t, g, "master2.txt", "master2 content", "master2 commit")

	// Switch to feature branch and rebase onto master
	if _, err := g.Checkout(branchName); err != nil {
		t.Fatalf("failed to checkout feature branch: %v", err)
	}
	output, err := g.Rebase(RebaseOptions{BranchName: "master"})
	if err != nil {
		t.Fatalf("Rebase() failed: %v\nOutput: %s", err, output)
	}
	if !strings.Contains(output, "Successfully rebased") && !strings.Contains(output, "Fast-forwarded") && !strings.Contains(output, "Applying") {
		t.Errorf("expected rebase output, got: %s", output)
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
	if _, err := g.AddFiles([]string{"new-file.txt"}); err != nil {
		t.Errorf("AddFiles() failed: %v", err)
	}

	// Test Reset
	if _, err := g.ResetFiles([]string{"new-file.txt"}); err != nil {
		t.Errorf("ResetFiles() failed: %v", err)
	}

	// Test Remove
	if _, err := g.RemoveFiles([]string{"file-ops.txt"}, false); err != nil {
		t.Errorf("RemoveFiles() failed: %v", err)
	}
	if _, err := os.Stat("file-ops.txt"); !os.IsNotExist(err) {
		t.Error("file should have been removed from working directory")
	}

	// Test Move
	createAndCommitFile(t, g, "source.txt", "move content", "Commit for move test")
	if _, err := g.MoveFile("source.txt", "destination.txt"); err != nil {
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
	if _, err := g.Stash(StashOptions{Push: true, Message: "test stash"}); err != nil {
		t.Fatalf("Stash() push failed: %v", err)
	}

	// Stash apply
	if _, err := g.Stash(StashOptions{Apply: true}); err != nil {
		t.Errorf("Stash() apply failed: %v", err)
	}
}
