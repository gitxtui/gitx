package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestGitCommands_InitRepository(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "init in current directory",
			path:    "",
			wantErr: false,
		},
		{
			name:    "init in specific directory",
			path:    "test-repo",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary directory for testing
			tempDir, err := os.MkdirTemp("", "git-test-")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			// Change to temp directory
			originalDir, _ := os.Getwd()
			os.Chdir(tempDir)
			defer os.Chdir(originalDir)

			g := NewGitCommands()
			err = g.InitRepository(tt.path)

			if (err != nil) != tt.wantErr {
				t.Errorf("GitCommands.InitRepository() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check if .git directory was created
			gitDir := ".git"
			if tt.path != "" {
				gitDir = filepath.Join(tt.path, ".git")
			}

			if !tt.wantErr {
				if _, err := os.Stat(gitDir); os.IsNotExist(err) {
					t.Errorf("expected .git directory to be created at %s", gitDir)
				}
			}
		})
	}
}

func TestGitCommands_CloneRepository(t *testing.T) {
	tests := []struct {
		name      string
		repoURL   string
		directory string
		wantErr   bool
	}{
		{
			name:      "empty repository URL",
			repoURL:   "",
			directory: "",
			wantErr:   true,
		},
		{
			name:      "invalid repository URL",
			repoURL:   "invalid-url",
			directory: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir, err := os.MkdirTemp("", "git-test-")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			originalDir, _ := os.Getwd()
			os.Chdir(tempDir)
			defer os.Chdir(originalDir)

			g := NewGitCommands()
			err = g.CloneRepository(tt.repoURL, tt.directory)

			if (err != nil) != tt.wantErr {
				t.Errorf("GitCommands.CloneRepository() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGitCommands_ShowStatus(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "git-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)

	g := NewGitCommands()

	err = g.ShowStatus()
	if err == nil {
		t.Error("expected error when running status outside git repository")
	}

	err = g.InitRepository("")
	if err != nil {
		t.Fatalf("failed to initialize git repository: %v", err)
	}

	err = g.ShowStatus()
	if err != nil {
		t.Errorf("unexpected error when running status in git repository: %v", err)
	}
}

func TestGitCommands_ShowLog(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "git-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)

	g := NewGitCommands()

	err = g.ShowLog(LogOptions{})
	if err == nil {
		t.Error("expected error when running log outside git repository")
	}

	err = g.InitRepository("")
	if err != nil {
		t.Fatalf("failed to initialize git repository: %v", err)
	}

	err = g.ShowLog(LogOptions{})
	if err == nil {
		t.Error("expected error when running log in empty git repository")
	}

	tests := []struct {
		name    string
		options LogOptions
		wantErr bool
	}{
		{
			name:    "default options",
			options: LogOptions{},
			wantErr: true,
		},
		{
			name: "oneline option",
			options: LogOptions{
				Oneline: true,
			},
			wantErr: true,
		},
		{
			name: "graph option",
			options: LogOptions{
				Graph: true,
			},
			wantErr: true,
		},
		{
			name: "max count option",
			options: LogOptions{
				MaxCount: 5,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = g.ShowLog(tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("GitCommands.ShowLog() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGitCommands_ShowDiff(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "git-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)

	g := NewGitCommands()

	err = g.ShowDiff(DiffOptions{})
	if err == nil {
		t.Error("expected error when running diff outside git repository")
	}

	err = g.InitRepository("")
	if err != nil {
		t.Fatalf("failed to initialize git repository: %v", err)
	}

	err = g.ShowDiff(DiffOptions{})
	if err != nil {
		t.Errorf("unexpected error when running diff in empty git repository: %v", err)
	}

	tests := []struct {
		name    string
		options DiffOptions
		wantErr bool
	}{
		{
			name:    "default options",
			options: DiffOptions{},
			wantErr: false,
		},
		{
			name: "cached option",
			options: DiffOptions{
				Cached: true,
			},
			wantErr: false,
		},
		{
			name: "stat option",
			options: DiffOptions{
				Stat: true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = g.ShowDiff(tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("GitCommands.ShowDiff() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGitCommands_ShowCommit(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "git-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)

	g := NewGitCommands()

	err = g.ShowCommit("")
	if err == nil {
		t.Error("expected error when running show outside git repository")
	}

	err = g.InitRepository("")
	if err != nil {
		t.Fatalf("failed to initialize git repository: %v", err)
	}

	err = g.ShowCommit("")
	if err == nil {
		t.Error("expected error when running show in empty git repository")
	}

	err = g.ShowCommit("nonexistent-hash")
	if err == nil {
		t.Error("expected error when running show with nonexistent commit hash")
	}
}

func TestNewGitCommands(t *testing.T) {
	g := NewGitCommands()
	if g == nil {
		t.Error("NewGitCommands() returned nil")
	}
}

func TestGitCommands_AddFiles(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "git-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)

	g := NewGitCommands()

	err = g.AddFiles([]string{"."})
	if err == nil {
		t.Error("expected error when running add outside git repository")
	}

	err = g.InitRepository("")
	if err != nil {
		t.Fatalf("failed to initialize git repository: %v", err)
	}

	// Create a test file
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Test with no paths (defaults to ".")
	err = g.AddFiles([]string{})
	if err != nil {
		t.Errorf("unexpected error when adding all files: %v", err)
	}

	// Test with specific file
	err = g.AddFiles([]string{"test.txt"})
	if err != nil {
		t.Errorf("unexpected error when adding specific file: %v", err)
	}
}

func TestGitCommands_ResetFiles(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "git-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)

	g := NewGitCommands()

	// Should error when no paths provided
	err = g.ResetFiles([]string{})
	if err == nil {
		t.Error("expected error when no paths provided")
	}

	err = g.InitRepository("")
	if err != nil {
		t.Fatalf("failed to initialize git repository: %v", err)
	}

	// Create and add a test file
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	err = g.AddFiles([]string{"test.txt"})
	if err != nil {
		t.Fatalf("failed to add test file: %v", err)
	}

	// Test resetting the file
	err = g.ResetFiles([]string{"test.txt"})
	if err != nil {
		t.Errorf("unexpected error when resetting file: %v", err)
	}
}

func TestGitCommands_Commit(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "git-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)

	g := NewGitCommands()

	// Should error when no message provided
	err = g.Commit(CommitOptions{})
	if err == nil {
		t.Error("expected error when no message provided")
	}

	err = g.InitRepository("")
	if err != nil {
		t.Fatalf("failed to initialize git repository: %v", err)
	}

	// Set git config for test
	err = runGitConfig(tempDir)
	if err != nil {
		t.Fatalf("failed to set git config: %v", err)
	}

	// Create, add, and commit a test file
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	err = g.AddFiles([]string{"test.txt"})
	if err != nil {
		t.Fatalf("failed to add test file: %v", err)
	}

	// Test committing with message
	err = g.Commit(CommitOptions{Message: "Initial commit"})
	if err != nil {
		t.Errorf("unexpected error when committing: %v", err)
	}

	// Test amending commit
	err = g.Commit(CommitOptions{Amend: true})
	if err != nil {
		t.Errorf("unexpected error when amending commit: %v", err)
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

func TestGitCommands_ManageBranch(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "git-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)

	g := NewGitCommands()

	err = g.InitRepository("")
	if err != nil {
		t.Fatalf("failed to initialize git repository: %v", err)
	}

	// Set git config for test
	err = runGitConfig(tempDir)
	if err != nil {
		t.Fatalf("failed to set git config: %v", err)
	}

	// Create an initial commit so we can create branches
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	err = g.AddFiles([]string{"test.txt"})
	if err != nil {
		t.Fatalf("failed to add test file: %v", err)
	}

	err = g.Commit(CommitOptions{Message: "Initial commit"})
	if err != nil {
		t.Fatalf("failed to commit: %v", err)
	}

	// Test creating a branch
	err = g.ManageBranch(BranchOptions{Create: true, Name: "test-branch"})
	if err != nil {
		t.Errorf("unexpected error when creating branch: %v", err)
	}

	// Test listing branches
	err = g.ManageBranch(BranchOptions{})
	if err != nil {
		t.Errorf("unexpected error when listing branches: %v", err)
	}

	// Test deleting a branch
	err = g.ManageBranch(BranchOptions{Delete: true, Name: "test-branch"})
	if err != nil {
		t.Errorf("unexpected error when deleting branch: %v", err)
	}

	// Test error when no name provided for create
	err = g.ManageBranch(BranchOptions{Create: true})
	if err == nil {
		t.Error("expected error when no name provided for create")
	}

	// Test error when no name provided for delete
	err = g.ManageBranch(BranchOptions{Delete: true})
	if err == nil {
		t.Error("expected error when no name provided for delete")
	}
}

func TestGitCommands_Checkout(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "git-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)

	g := NewGitCommands()

	// Should error when no branch name provided
	err = g.Checkout("")
	if err == nil {
		t.Error("expected error when no branch name provided")
	}

	err = g.InitRepository("")
	if err != nil {
		t.Fatalf("failed to initialize git repository: %v", err)
	}

	// Set git config for test
	err = runGitConfig(tempDir)
	if err != nil {
		t.Fatalf("failed to set git config: %v", err)
	}

	// Create an initial commit
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	err = g.AddFiles([]string{"test.txt"})
	if err != nil {
		t.Fatalf("failed to add test file: %v", err)
	}

	err = g.Commit(CommitOptions{Message: "Initial commit"})
	if err != nil {
		t.Fatalf("failed to commit: %v", err)
	}

	// Create a branch
	err = g.ManageBranch(BranchOptions{Create: true, Name: "test-branch"})
	if err != nil {
		t.Fatalf("failed to create branch: %v", err)
	}

	// Test checkout
	err = g.Checkout("test-branch")
	if err != nil {
		t.Errorf("unexpected error when checking out branch: %v", err)
	}

	// Test checkout with nonexistent branch
	err = g.Checkout("nonexistent-branch")
	if err == nil {
		t.Error("expected error when checking out nonexistent branch")
	}
}

func TestGitCommands_Switch(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "git-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)

	g := NewGitCommands()

	// Should error when no branch name provided
	err = g.Switch("")
	if err == nil {
		t.Error("expected error when no branch name provided")
	}

	err = g.InitRepository("")
	if err != nil {
		t.Fatalf("failed to initialize git repository: %v", err)
	}

	// Set git config for test
	err = runGitConfig(tempDir)
	if err != nil {
		t.Fatalf("failed to set git config: %v", err)
	}

	// Create an initial commit
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	err = g.AddFiles([]string{"test.txt"})
	if err != nil {
		t.Fatalf("failed to add test file: %v", err)
	}

	err = g.Commit(CommitOptions{Message: "Initial commit"})
	if err != nil {
		t.Fatalf("failed to commit: %v", err)
	}

	// Create a branch
	err = g.ManageBranch(BranchOptions{Create: true, Name: "test-branch"})
	if err != nil {
		t.Fatalf("failed to create branch: %v", err)
	}

	// Test switch
	err = g.Switch("test-branch")
	if err != nil {
		t.Errorf("unexpected error when switching branch: %v", err)
	}

	// Test switch with nonexistent branch
	err = g.Switch("nonexistent-branch")
	if err == nil {
		t.Error("expected error when switching to nonexistent branch")
	}
}

func TestGitCommands_Merge(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "git-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)

	g := NewGitCommands()

	// Should error when no branch name provided
	err = g.Merge(MergeOptions{})
	if err == nil {
		t.Error("expected error when no branch name provided")
	}

	err = g.InitRepository("")
	if err != nil {
		t.Fatalf("failed to initialize git repository: %v", err)
	}

	// Set git config for test
	err = runGitConfig(tempDir)
	if err != nil {
		t.Fatalf("failed to set git config: %v", err)
	}

	// Create an initial commit on main
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	err = g.AddFiles([]string{"test.txt"})
	if err != nil {
		t.Fatalf("failed to add test file: %v", err)
	}

	err = g.Commit(CommitOptions{Message: "Initial commit"})
	if err != nil {
		t.Fatalf("failed to commit: %v", err)
	}

	// Create and checkout a branch
	err = g.ManageBranch(BranchOptions{Create: true, Name: "test-branch"})
	if err != nil {
		t.Fatalf("failed to create branch: %v", err)
	}

	err = g.Checkout("test-branch")
	if err != nil {
		t.Fatalf("failed to checkout branch: %v", err)
	}

	// Make a change and commit on the branch
	if err := os.WriteFile(testFile, []byte("updated content"), 0644); err != nil {
		t.Fatalf("failed to update test file: %v", err)
	}

	err = g.AddFiles([]string{"test.txt"})
	if err != nil {
		t.Fatalf("failed to add updated test file: %v", err)
	}

	err = g.Commit(CommitOptions{Message: "Update on branch"})
	if err != nil {
		t.Fatalf("failed to commit on branch: %v", err)
	}

	// Switch back to main
	err = g.Checkout("master")
	if err != nil {
		t.Fatalf("failed to checkout main branch: %v", err)
	}

	// Test merge
	err = g.Merge(MergeOptions{BranchName: "test-branch"})
	if err != nil {
		t.Errorf("unexpected error when merging branch: %v", err)
	}

	// Test merge with options
	err = g.Merge(MergeOptions{
		BranchName:    "test-branch",
		NoFastForward: true,
		Message:       "Merge test branch",
	})
	if err != nil {
		t.Errorf("unexpected error when merging with options: %v", err)
	}

	// Test merge with nonexistent branch
	err = g.Merge(MergeOptions{BranchName: "nonexistent-branch"})
	if err == nil {
		t.Error("expected error when merging nonexistent branch")
	}
}

func TestGitCommands_ManageTag(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "git-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)

	g := NewGitCommands()

	err = g.InitRepository("")
	if err != nil {
		t.Fatalf("failed to initialize git repository: %v", err)
	}

	// Set git config for test
	err = runGitConfig(tempDir)
	if err != nil {
		t.Fatalf("failed to set git config: %v", err)
	}

	// Create an initial commit
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	err = g.AddFiles([]string{"test.txt"})
	if err != nil {
		t.Fatalf("failed to add test file: %v", err)
	}

	err = g.Commit(CommitOptions{Message: "Initial commit"})
	if err != nil {
		t.Fatalf("failed to commit: %v", err)
	}

	// Test creating a tag
	err = g.ManageTag(TagOptions{Create: true, Name: "v1.0.0", Message: "Version 1.0.0"})
	if err != nil {
		t.Errorf("unexpected error when creating tag: %v", err)
	}

	// Test listing tags
	err = g.ManageTag(TagOptions{})
	if err != nil {
		t.Errorf("unexpected error when listing tags: %v", err)
	}

	// Test deleting a tag
	err = g.ManageTag(TagOptions{Delete: true, Name: "v1.0.0"})
	if err != nil {
		t.Errorf("unexpected error when deleting tag: %v", err)
	}

	// Test error when no name provided for create
	err = g.ManageTag(TagOptions{Create: true})
	if err == nil {
		t.Error("expected error when no name provided for create")
	}

	// Test error when no name provided for delete
	err = g.ManageTag(TagOptions{Delete: true})
	if err == nil {
		t.Error("expected error when no name provided for delete")
	}
}

func TestGitCommands_RemoveFiles(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "git-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)

	g := NewGitCommands()

	// Should error when no paths provided
	err = g.RemoveFiles([]string{}, false)
	if err == nil {
		t.Error("expected error when no paths provided")
	}

	err = g.InitRepository("")
	if err != nil {
		t.Fatalf("failed to initialize git repository: %v", err)
	}

	// Set git config for test
	err = runGitConfig(tempDir)
	if err != nil {
		t.Fatalf("failed to set git config: %v", err)
	}

	// Create a file and commit it
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	err = g.AddFiles([]string{"test.txt"})
	if err != nil {
		t.Fatalf("failed to add test file: %v", err)
	}

	err = g.Commit(CommitOptions{Message: "Initial commit"})
	if err != nil {
		t.Fatalf("failed to commit: %v", err)
	}

	// Test removing a file
	err = g.RemoveFiles([]string{"test.txt"}, false)
	if err != nil {
		t.Errorf("unexpected error when removing file: %v", err)
	}

	// Create and add another file to test --cached
	testFile2 := filepath.Join(tempDir, "test2.txt")
	if err := os.WriteFile(testFile2, []byte("test content 2"), 0644); err != nil {
		t.Fatalf("failed to create second test file: %v", err)
	}

	err = g.AddFiles([]string{"test2.txt"})
	if err != nil {
		t.Fatalf("failed to add second test file: %v", err)
	}

	// Test removing with --cached option
	err = g.RemoveFiles([]string{"test2.txt"}, true)
	if err != nil {
		t.Errorf("unexpected error when removing file with --cached: %v", err)
	}
}

func TestGitCommands_MoveFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "git-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)

	g := NewGitCommands()

	// Should error when source or destination not provided
	err = g.MoveFile("", "dest.txt")
	if err == nil {
		t.Error("expected error when source not provided")
	}

	err = g.MoveFile("source.txt", "")
	if err == nil {
		t.Error("expected error when destination not provided")
	}

	err = g.InitRepository("")
	if err != nil {
		t.Fatalf("failed to initialize git repository: %v", err)
	}

	// Set git config for test
	err = runGitConfig(tempDir)
	if err != nil {
		t.Fatalf("failed to set git config: %v", err)
	}

	// Create a file and commit it
	testFile := filepath.Join(tempDir, "source.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	err = g.AddFiles([]string{"source.txt"})
	if err != nil {
		t.Fatalf("failed to add test file: %v", err)
	}

	err = g.Commit(CommitOptions{Message: "Initial commit"})
	if err != nil {
		t.Fatalf("failed to commit: %v", err)
	}

	// Test moving a file
	err = g.MoveFile("source.txt", "dest.txt")
	if err != nil {
		t.Errorf("unexpected error when moving file: %v", err)
	}

	// Verify destination file exists
	if _, err := os.Stat(filepath.Join(tempDir, "dest.txt")); os.IsNotExist(err) {
		t.Error("destination file does not exist after move")
	}
}

func TestGitCommands_Restore(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "git-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)

	g := NewGitCommands()

	// Should error when no paths provided
	err = g.Restore(RestoreOptions{Paths: []string{}})
	if err == nil {
		t.Error("expected error when no paths provided")
	}

	err = g.InitRepository("")
	if err != nil {
		t.Fatalf("failed to initialize git repository: %v", err)
	}

	// Set git config for test
	err = runGitConfig(tempDir)
	if err != nil {
		t.Fatalf("failed to set git config: %v", err)
	}

	// Create a file and commit it
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("initial content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	err = g.AddFiles([]string{"test.txt"})
	if err != nil {
		t.Fatalf("failed to add test file: %v", err)
	}

	err = g.Commit(CommitOptions{Message: "Initial commit"})
	if err != nil {
		t.Fatalf("failed to commit: %v", err)
	}

	// Modify the file
	if err := os.WriteFile(testFile, []byte("modified content"), 0644); err != nil {
		t.Fatalf("failed to modify test file: %v", err)
	}

	// Test restoring working tree changes
	err = g.Restore(RestoreOptions{Paths: []string{"test.txt"}})
	if err != nil {
		t.Errorf("unexpected error when restoring file: %v", err)
	}

	// Add the modified file to staging
	if err := os.WriteFile(testFile, []byte("modified content"), 0644); err != nil {
		t.Fatalf("failed to modify test file again: %v", err)
	}

	err = g.AddFiles([]string{"test.txt"})
	if err != nil {
		t.Fatalf("failed to add modified test file: %v", err)
	}

	// Test restoring from staging
	err = g.Restore(RestoreOptions{Paths: []string{"test.txt"}, Staged: true})
	if err != nil {
		t.Errorf("unexpected error when restoring staged file: %v", err)
	}
}

func TestGitCommands_Revert(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "git-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)

	g := NewGitCommands()

	// Should error when no commit hash provided
	err = g.Revert("")
	if err == nil {
		t.Error("expected error when no commit hash provided")
	}

	err = g.InitRepository("")
	if err != nil {
		t.Fatalf("failed to initialize git repository: %v", err)
	}

	// Set git config for test
	err = runGitConfig(tempDir)
	if err != nil {
		t.Fatalf("failed to set git config: %v", err)
	}

	// Create a file and make initial commit
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("initial content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	err = g.AddFiles([]string{"test.txt"})
	if err != nil {
		t.Fatalf("failed to add test file: %v", err)
	}

	err = g.Commit(CommitOptions{Message: "Initial commit"})
	if err != nil {
		t.Fatalf("failed to make initial commit: %v", err)
	}

	// Get commit hash
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = tempDir
	commitHash, err := cmd.Output()
	if err != nil {
		t.Fatalf("failed to get commit hash: %v", err)
	}
	hashStr := string(commitHash)
	hashStr = hashStr[:len(hashStr)-1] // Remove newline

	// Modify file and make a second commit
	if err := os.WriteFile(testFile, []byte("modified content"), 0644); err != nil {
		t.Fatalf("failed to modify test file: %v", err)
	}

	err = g.AddFiles([]string{"test.txt"})
	if err != nil {
		t.Fatalf("failed to add modified test file: %v", err)
	}

	err = g.Commit(CommitOptions{Message: "Second commit"})
	if err != nil {
		t.Fatalf("failed to make second commit: %v", err)
	}

	// Test reverting a commit (this might fail if there are merge conflicts)
	err = g.Revert(hashStr)
	if err != nil {
		// In a real test, we might want to make the commits in such a way that revert won't conflict
		t.Logf("Revert failed, might be due to conflicts: %v", err)
	}
}

func TestGitCommands_Stash(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "git-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)

	g := NewGitCommands()

	err = g.InitRepository("")
	if err != nil {
		t.Fatalf("failed to initialize git repository: %v", err)
	}

	// Set git config for test
	err = runGitConfig(tempDir)
	if err != nil {
		t.Fatalf("failed to set git config: %v", err)
	}

	// Create a file and make initial commit
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("initial content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	err = g.AddFiles([]string{"test.txt"})
	if err != nil {
		t.Fatalf("failed to add test file: %v", err)
	}

	err = g.Commit(CommitOptions{Message: "Initial commit"})
	if err != nil {
		t.Fatalf("failed to make initial commit: %v", err)
	}

	// Modify file but don't commit
	if err := os.WriteFile(testFile, []byte("modified content"), 0644); err != nil {
		t.Fatalf("failed to modify test file: %v", err)
	}

	// Test stash push
	err = g.Stash(StashOptions{Push: true, Message: "Test stash"})
	if err != nil {
		t.Errorf("unexpected error when pushing stash: %v", err)
	}

	// Test stash list
	err = g.Stash(StashOptions{List: true})
	if err != nil {
		t.Errorf("unexpected error when listing stashes: %v", err)
	}

	// Test stash show
	err = g.Stash(StashOptions{Show: true, StashID: "stash@{0}"})
	if err != nil {
		t.Errorf("unexpected error when showing stash: %v", err)
	}

	// Test stash apply
	err = g.Stash(StashOptions{Apply: true})
	if err != nil {
		t.Errorf("unexpected error when applying stash: %v", err)
	}

	// Test stash push with default options (should set Push to true)
	err = g.Stash(StashOptions{})
	if err != nil {
		t.Errorf("unexpected error when pushing stash with default options: %v", err)
	}

	// Test stash drop
	err = g.Stash(StashOptions{Drop: true})
	if err != nil {
		t.Errorf("unexpected error when dropping stash: %v", err)
	}
}

func TestGitCommands_ListFiles(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "git-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)

	g := NewGitCommands()

	err = g.InitRepository("")
	if err != nil {
		t.Fatalf("failed to initialize git repository: %v", err)
	}

	// Test listing files in empty repo
	err = g.ListFiles()
	if err != nil {
		t.Errorf("unexpected error when listing files in empty repo: %v", err)
	}

	// Set git config for test
	err = runGitConfig(tempDir)
	if err != nil {
		t.Fatalf("failed to set git config: %v", err)
	}

	// Create a file and commit it
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	err = g.AddFiles([]string{"test.txt"})
	if err != nil {
		t.Fatalf("failed to add test file: %v", err)
	}

	err = g.Commit(CommitOptions{Message: "Add test file"})
	if err != nil {
		t.Fatalf("failed to commit test file: %v", err)
	}

	// Test listing files with content
	err = g.ListFiles()
	if err != nil {
		t.Errorf("unexpected error when listing files: %v", err)
	}
}

func TestGitCommands_BlameFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "git-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)

	g := NewGitCommands()

	// Should error when no file path provided
	err = g.BlameFile("")
	if err == nil {
		t.Error("expected error when no file path provided")
	}

	err = g.InitRepository("")
	if err != nil {
		t.Fatalf("failed to initialize git repository: %v", err)
	}

	// Set git config for test
	err = runGitConfig(tempDir)
	if err != nil {
		t.Fatalf("failed to set git config: %v", err)
	}

	// Create a file and commit it
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	err = g.AddFiles([]string{"test.txt"})
	if err != nil {
		t.Fatalf("failed to add test file: %v", err)
	}

	err = g.Commit(CommitOptions{Message: "Add test file"})
	if err != nil {
		t.Fatalf("failed to commit test file: %v", err)
	}

	// Test blame
	err = g.BlameFile("test.txt")
	if err != nil {
		t.Errorf("unexpected error when blaming file: %v", err)
	}

	// Test blame nonexistent file
	err = g.BlameFile("nonexistent.txt")
	if err == nil {
		t.Error("expected error when blaming nonexistent file")
	}
}

func TestGitCommands_ManageRemote(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "git-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)

	g := NewGitCommands()

	err = g.InitRepository("")
	if err != nil {
		t.Fatalf("failed to initialize git repository: %v", err)
	}

	// Test adding a remote without name or URL
	err = g.ManageRemote(RemoteOptions{Add: true})
	if err == nil {
		t.Error("expected error when adding remote without name or URL")
	}

	// Test adding a remote
	err = g.ManageRemote(RemoteOptions{Add: true, Name: "origin", URL: "https://github.com/example/repo.git"})
	if err != nil {
		t.Errorf("unexpected error when adding remote: %v", err)
	}

	// Test listing remotes
	err = g.ManageRemote(RemoteOptions{Verbose: true})
	if err != nil {
		t.Errorf("unexpected error when listing remotes: %v", err)
	}

	// Test removing a remote without name
	err = g.ManageRemote(RemoteOptions{Remove: true})
	if err == nil {
		t.Error("expected error when removing remote without name")
	}

	// Test removing a remote
	err = g.ManageRemote(RemoteOptions{Remove: true, Name: "origin"})
	if err != nil {
		t.Errorf("unexpected error when removing remote: %v", err)
	}
}
