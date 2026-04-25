package usecase

import (
	"errors"
	"testing"

	domainconfig "github.com/HidemaruOwO/pummit/internal/domain/config"
	domainemoji "github.com/HidemaruOwO/pummit/internal/domain/emoji"
)

type fakeConfigLoader struct {
	cfg domainconfig.Config
	err error
}

func (f fakeConfigLoader) Load() (domainconfig.Config, error) {
	return f.cfg, f.err
}

type fakeGitClient struct {
	staged        []string
	branch        string
	commitMessage string
	err           error
}

func (f *fakeGitClient) StagedFiles() ([]string, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.staged, nil
}

func (f *fakeGitClient) CurrentBranch() (string, error) {
	if f.err != nil {
		return "", f.err
	}
	return f.branch, nil
}

func (f *fakeGitClient) Commit(message string) error {
	if f.err != nil {
		return f.err
	}
	f.commitMessage = message
	return nil
}

type fakeCatalog struct {
	mapping map[string]domainemoji.Match
}

func (f fakeCatalog) Lookup(name string) (domainemoji.Match, bool) {
	match, ok := f.mapping[name]
	return match, ok
}

type fakeRemote struct {
	mapping map[string]domainemoji.Match
	err     error
}

func (f fakeRemote) Lookup(name string) (domainemoji.Match, bool, error) {
	if f.err != nil {
		return domainemoji.Match{}, false, f.err
	}
	match, ok := f.mapping[name]
	return match, ok, nil
}

func TestCommitExplicitUsesAliasShortcut(t *testing.T) {
	cfg := domainconfig.Default()
	git := &fakeGitClient{staged: []string{"file.txt"}}
	service := NewCommitService(fakeConfigLoader{cfg: cfg}, git, fakeCatalog{}, fakeRemote{})

	result, err := service.CommitExplicit("feat", "Add x")
	if err != nil {
		t.Fatalf("CommitExplicit returned error: %v", err)
	}

	if !result.Committed {
		t.Fatal("CommitExplicit did not commit")
	}

	want := "✨ Add x (file.txt)"
	if git.commitMessage != want {
		t.Fatalf("commit message = %q, want %q", git.commitMessage, want)
	}
}

func TestCommitAutoUsesBranchMapping(t *testing.T) {
	cfg := domainconfig.Default()
	git := &fakeGitClient{staged: []string{"file.txt"}, branch: "feature/test"}
	service := NewCommitService(fakeConfigLoader{cfg: cfg}, git, fakeCatalog{}, fakeRemote{})

	result, err := service.CommitAuto("Add x")
	if err != nil {
		t.Fatalf("CommitAuto returned error: %v", err)
	}

	if !result.Committed {
		t.Fatal("CommitAuto did not commit")
	}

	want := "✨ Add x (file.txt)"
	if git.commitMessage != want {
		t.Fatalf("commit message = %q, want %q", git.commitMessage, want)
	}
}

func TestCommitExplicitFallsBackToColonName(t *testing.T) {
	cfg := domainconfig.Default()
	git := &fakeGitClient{staged: []string{"file.txt"}}
	service := NewCommitService(fakeConfigLoader{cfg: cfg}, git, fakeCatalog{}, fakeRemote{err: errors.New("offline")})

	_, err := service.CommitExplicit("unknown", "Add x")
	if err != nil {
		t.Fatalf("CommitExplicit returned error: %v", err)
	}

	want := ":unknown: Add x (file.txt)"
	if git.commitMessage != want {
		t.Fatalf("commit message = %q, want %q", git.commitMessage, want)
	}
}

func TestCommitNoStagedFiles(t *testing.T) {
	cfg := domainconfig.Default()
	git := &fakeGitClient{staged: []string{}}
	service := NewCommitService(fakeConfigLoader{cfg: cfg}, git, fakeCatalog{}, fakeRemote{})

	result, err := service.CommitExplicit("sparkles", "Add x")
	if err != nil {
		t.Fatalf("CommitExplicit returned error: %v", err)
	}

	if result.Committed {
		t.Fatal("expected no commit when there are no staged files")
	}
}
