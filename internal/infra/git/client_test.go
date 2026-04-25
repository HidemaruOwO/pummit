package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
)

func TestStagedFiles(t *testing.T) {
	repo := initRepo(t)
	path := filepath.Join(repo, "file.txt")
	if err := os.WriteFile(path, []byte("hello"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	runGit(t, repo, "add", "file.txt")

	client := NewClient(repo)
	got, err := client.StagedFiles()
	if err != nil {
		t.Fatalf("StagedFiles returned error: %v", err)
	}

	want := []string{"file.txt"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("StagedFiles() = %v, want %v", got, want)
	}
}

func TestCurrentBranch(t *testing.T) {
	repo := initRepo(t)
	runGit(t, repo, "checkout", "-b", "feature/test")

	client := NewClient(repo)
	got, err := client.CurrentBranch()
	if err != nil {
		t.Fatalf("CurrentBranch returned error: %v", err)
	}

	if got != "feature/test" {
		t.Fatalf("CurrentBranch() = %q, want %q", got, "feature/test")
	}
}

func initRepo(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}

	dir := t.TempDir()
	runGit(t, dir, "init", "--initial-branch=main")
	runGit(t, dir, "config", "user.email", "test@example.com")
	runGit(t, dir, "config", "user.name", "Test")
	return dir
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, string(out))
	}
}
