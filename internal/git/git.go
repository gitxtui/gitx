package git

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

// ExecCommand is a variable that holds the exec.Command function
// This allows it to be mocked in tests
var ExecCommand = exec.Command

type GitCommands struct{}

func NewGitCommands() *GitCommands {
	return &GitCommands{}
}

func (g *GitCommands) InitRepository(path string) error {
	if path == "" {
		path = "."
	}

	cmd := ExecCommand("git", "init", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to initialize repository: %v\nOutput: %s", err, output)
	}

	absPath, _ := filepath.Abs(path)
	fmt.Printf("Initialized empty Git repository in %s\n", absPath)
	return nil
}

func (g *GitCommands) CloneRepository(repoURL, directory string) error {
	if repoURL == "" {
		return fmt.Errorf("repository URL is required")
	}

	var cmd *exec.Cmd
	if directory != "" {
		cmd = exec.Command("git", "clone", repoURL, directory)
	} else {
		cmd = exec.Command("git", "clone", repoURL)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to clone repository: %v\nOutput: %s", err, output)
	}

	fmt.Printf("Successfully cloned repository: %s\n", repoURL)
	return nil
}

func (g *GitCommands) ShowStatus() error {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get status: %v", err)
	}

	if len(output) == 0 {
		fmt.Println("Working directory clean")
		return nil
	}

	cmd = exec.Command("git", "status")
	output, err = cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get detailed status: %v", err)
	}

	fmt.Print(string(output))
	return nil
}

func (g *GitCommands) ShowLog(options LogOptions) error {
	args := []string{"log"}

	if options.Oneline {
		args = append(args, "--oneline")
	}
	if options.Graph {
		args = append(args, "--graph")
	}
	if options.MaxCount > 0 {
		args = append(args, fmt.Sprintf("-%d", options.MaxCount))
	}

	cmd := exec.Command("git", args...)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get log: %v", err)
	}

	fmt.Print(string(output))
	return nil
}

func (g *GitCommands) ShowDiff(options DiffOptions) error {
	args := []string{"diff"}

	if options.Cached {
		args = append(args, "--cached")
	}
	if options.Stat {
		args = append(args, "--stat")
	}
	if options.Commit1 != "" {
		args = append(args, options.Commit1)
	}
	if options.Commit2 != "" {
		args = append(args, options.Commit2)
	}

	cmd := exec.Command("git", args...)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get diff: %v", err)
	}

	fmt.Print(string(output))
	return nil
}

func (g *GitCommands) ShowCommit(commitHash string) error {
	if commitHash == "" {
		commitHash = "HEAD"
	}

	cmd := exec.Command("git", "show", commitHash)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to show commit: %v", err)
	}

	fmt.Print(string(output))
	return nil
}

type LogOptions struct {
	Oneline  bool
	Graph    bool
	MaxCount int
}

type DiffOptions struct {
	Commit1 string
	Commit2 string
	Cached  bool
	Stat    bool
}

func GetStatus() (string, error) {
	cmd := exec.Command("git", "status")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), err
	}
	return string(output), nil
}

func (g *GitCommands) AddFiles(paths []string) error {
	if len(paths) == 0 {
		paths = []string{"."}
	}

	args := append([]string{"add"}, paths...)
	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to add files: %v\nOutput: %s", err, output)
	}

	fmt.Printf("Successfully added files to staging area\n")
	return nil
}

func (g *GitCommands) ResetFiles(paths []string) error {
	if len(paths) == 0 {
		return fmt.Errorf("at least one file path is required")
	}

	args := append([]string{"reset"}, paths...)
	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to unstage files: %v\nOutput: %s", err, output)
	}

	fmt.Printf("Successfully unstaged files\n")
	return nil
}

type CommitOptions struct {
	Message string
	Amend   bool
}

func (g *GitCommands) Commit(options CommitOptions) error {
	if options.Message == "" && !options.Amend {
		return fmt.Errorf("commit message is required unless amending")
	}

	args := []string{"commit"}

	if options.Amend {
		args = append(args, "--amend")
	}

	if options.Message != "" {
		args = append(args, "-m", options.Message)
	}

	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to commit changes: %v\nOutput: %s", err, output)
	}

	fmt.Print(string(output))
	return nil
}

type BranchOptions struct {
	Create bool
	Delete bool
	Name   string
}

func (g *GitCommands) ManageBranch(options BranchOptions) error {
	args := []string{"branch"}

	if options.Delete {
		if options.Name == "" {
			return fmt.Errorf("branch name is required for deletion")
		}
		args = append(args, "-d", options.Name)
	} else if options.Create {
		if options.Name == "" {
			return fmt.Errorf("branch name is required for creation")
		}
		args = append(args, options.Name)
	}

	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("branch operation failed: %v\nOutput: %s", err, output)
	}

	fmt.Print(string(output))
	return nil
}

func (g *GitCommands) Checkout(branchName string) error {
	if branchName == "" {
		return fmt.Errorf("branch name is required")
	}

	cmd := exec.Command("git", "checkout", branchName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to checkout branch: %v\nOutput: %s", err, output)
	}

	fmt.Printf("Switched to branch '%s'\n", branchName)
	return nil
}

func (g *GitCommands) Switch(branchName string) error {
	if branchName == "" {
		return fmt.Errorf("branch name is required")
	}

	cmd := exec.Command("git", "switch", branchName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to switch branch: %v\nOutput: %s", err, output)
	}

	fmt.Printf("Switched to branch '%s'\n", branchName)
	return nil
}

type MergeOptions struct {
	BranchName    string
	NoFastForward bool
	Message       string
}

func (g *GitCommands) Merge(options MergeOptions) error {
	if options.BranchName == "" {
		return fmt.Errorf("branch name is required")
	}

	args := []string{"merge"}

	if options.NoFastForward {
		args = append(args, "--no-ff")
	}

	if options.Message != "" {
		args = append(args, "-m", options.Message)
	}

	args = append(args, options.BranchName)

	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to merge branch: %v\nOutput: %s", err, output)
	}

	fmt.Print(string(output))
	return nil
}

type TagOptions struct {
	Create  bool
	Delete  bool
	Name    string
	Message string
	Commit  string
}

func (g *GitCommands) ManageTag(options TagOptions) error {
	args := []string{"tag"}

	if options.Delete {
		if options.Name == "" {
			return fmt.Errorf("tag name is required for deletion")
		}
		args = append(args, "-d", options.Name)
	} else if options.Create {
		if options.Name == "" {
			return fmt.Errorf("tag name is required for creation")
		}
		if options.Message != "" {
			args = append(args, "-m", options.Message)
		}
		args = append(args, options.Name)
		if options.Commit != "" {
			args = append(args, options.Commit)
		}
	}

	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("tag operation failed: %v\nOutput: %s", err, output)
	}

	fmt.Print(string(output))
	return nil
}

type RemoteOptions struct {
	Add     bool
	Remove  bool
	Name    string
	URL     string
	Verbose bool
}

func (g *GitCommands) ManageRemote(options RemoteOptions) error {
	args := []string{"remote"}

	if options.Verbose {
		args = append(args, "-v")
	}

	if options.Add {
		if options.Name == "" || options.URL == "" {
			return fmt.Errorf("remote name and URL are required for adding")
		}
		args = append(args, "add", options.Name, options.URL)
	} else if options.Remove {
		if options.Name == "" {
			return fmt.Errorf("remote name is required for removal")
		}
		args = append(args, "remove", options.Name)
	}

	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("remote operation failed: %v\nOutput: %s", err, output)
	}

	fmt.Print(string(output))
	return nil
}

func (g *GitCommands) Fetch(remote string, branch string) error {
	args := []string{"fetch"}

	if remote != "" {
		args = append(args, remote)
	}

	if branch != "" {
		args = append(args, branch)
	}

	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to fetch: %v\nOutput: %s", err, output)
	}

	fmt.Print(string(output))
	return nil
}

type PullOptions struct {
	Remote string
	Branch string
	Rebase bool
}

func (g *GitCommands) Pull(options PullOptions) error {
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

	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to pull: %v\nOutput: %s", err, output)
	}

	fmt.Print(string(output))
	return nil
}

type PushOptions struct {
	Remote      string
	Branch      string
	Force       bool
	SetUpstream bool
	Tags        bool
}

func (g *GitCommands) Push(options PushOptions) error {
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

	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to push: %v\nOutput: %s", err, output)
	}

	fmt.Print(string(output))
	return nil
}

func (g *GitCommands) RemoveFiles(paths []string, cached bool) error {
	if len(paths) == 0 {
		return fmt.Errorf("at least one file path is required")
	}

	args := []string{"rm"}

	if cached {
		args = append(args, "--cached")
	}

	args = append(args, paths...)

	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to remove files: %v\nOutput: %s", err, output)
	}

	fmt.Printf("Successfully removed files\n")
	return nil
}

func (g *GitCommands) MoveFile(source, destination string) error {
	if source == "" || destination == "" {
		return fmt.Errorf("source and destination paths are required")
	}

	cmd := exec.Command("git", "mv", source, destination)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to move file: %v\nOutput: %s", err, output)
	}

	fmt.Printf("Successfully moved %s to %s\n", source, destination)
	return nil
}

type RestoreOptions struct {
	Paths      []string
	Source     string
	Staged     bool
	WorkingDir bool
}

func (g *GitCommands) Restore(options RestoreOptions) error {
	if len(options.Paths) == 0 {
		return fmt.Errorf("at least one file path is required")
	}

	args := []string{"restore"}

	if options.Staged {
		args = append(args, "--staged")
	}

	if options.WorkingDir {
		args = append(args, "--worktree")
	}

	if options.Source != "" {
		args = append(args, "--source", options.Source)
	}

	args = append(args, options.Paths...)

	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to restore files: %v\nOutput: %s", err, output)
	}

	fmt.Printf("Successfully restored files\n")
	return nil
}

func (g *GitCommands) Revert(commitHash string) error {
	if commitHash == "" {
		return fmt.Errorf("commit hash is required")
	}

	cmd := exec.Command("git", "revert", commitHash)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to revert commit: %v\nOutput: %s", err, output)
	}

	fmt.Print(string(output))
	return nil
}

type StashOptions struct {
	Push    bool
	Pop     bool
	Apply   bool
	List    bool
	Show    bool
	Drop    bool
	Message string
	StashID string
}

func (g *GitCommands) Stash(options StashOptions) error {
	if !options.Push && !options.Pop && !options.Apply && !options.List && !options.Show && !options.Drop {
		options.Push = true
	}

	var args []string

	if options.Push {
		args = []string{"stash", "push"}
		if options.Message != "" {
			args = append(args, "-m", options.Message)
		}
	} else if options.Pop {
		args = []string{"stash", "pop"}
		if options.StashID != "" {
			args = append(args, options.StashID)
		}
	} else if options.Apply {
		args = []string{"stash", "apply"}
		if options.StashID != "" {
			args = append(args, options.StashID)
		}
	} else if options.List {
		args = []string{"stash", "list"}
	} else if options.Show {
		args = []string{"stash", "show"}
		if options.StashID != "" {
			args = append(args, options.StashID)
		}
	} else if options.Drop {
		args = []string{"stash", "drop"}
		if options.StashID != "" {
			args = append(args, options.StashID)
		}
	}

	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("stash operation failed: %v\nOutput: %s", err, output)
	}

	fmt.Print(string(output))
	return nil
}

func (g *GitCommands) ListFiles() error {
	cmd := exec.Command("git", "ls-files")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to list files: %v", err)
	}

	fmt.Print(string(output))
	return nil
}

func (g *GitCommands) BlameFile(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("file path is required")
	}

	cmd := exec.Command("git", "blame", filePath)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to blame file: %v", err)
	}

	fmt.Print(string(output))
	return nil
}
