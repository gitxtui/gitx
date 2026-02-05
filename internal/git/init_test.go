package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitRepository(t *testing.T) {
	// Skip if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not found in PATH; skipping integration test")
	}

	tmp := t.TempDir()
	g := &GitCommands{}

	msg, err := g.InitRepository(tmp)
	if err != nil {
		t.Fatalf("InitRepository failed: %v (msg=%q)", err, msg)
	}

	// Check that .git directory was created
	gitdir := filepath.Join(tmp, ".git")
	if _, err := os.Stat(gitdir); err != nil {
		t.Fatalf(".git directory missing: %v", err)
	}

	// Check that the message contains the absolute path
	abs, _ := filepath.Abs(tmp)
	if !strings.Contains(msg, abs) {
		t.Fatalf("unexpected message: %q does not contain %q", msg, abs)
	}
}

func TestInitRepositoryAlreadyExists(t *testing.T) {
	// Skip if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not found in PATH; skipping integration test")
	}

	tmp := t.TempDir()
	g := &GitCommands{}

	// Initialize repo first time
	_, err := g.InitRepository(tmp)
	if err != nil {
		t.Fatalf("Initial InitRepository failed: %v", err)
	}

	// Try to initialize again
	msg, err := g.InitRepository(tmp)
	if err != nil {
		// It might error or succeed depending on git behavior
		// Just verify we handle it gracefully
		t.Logf("Second InitRepository returned error (expected): %v", err)
	} else {
		t.Logf("Second InitRepository succeeded with msg: %q", msg)
	}
}
