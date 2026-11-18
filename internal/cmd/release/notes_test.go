package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const testNotes = `# Release notes

Intro.

## v0.4.0

* **\[MDAPI-1]\[GO]** First.

## 0.3.0

* Old.
`

func writeFiles(t *testing.T, root string, files map[string]string) {
	t.Helper()
	for name, text := range files {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestCheckFragment(t *testing.T) {
	tests := []struct {
		name, text, err string
	}{
		{"valid", "* **\\[MDAPI-7]\\[GO]** Fixed.\r\n  More text.\r\n  * Nested.\r\n", ""},
		{"empty", "\n", "empty"},
		{"no key", "* Fixed.\n", "line 1"},
		{"other ticket", "* **\\[MDAPI-8]\\[GO]** Fixed.\n", "names MDAPI-8"},
		{"two items", "* **\\[MDAPI-7]\\[GO]** Fixed.\n* **\\[MDAPI-7]\\[GO]** Again.\n", "line 2 is not indented"},
		{"blank line", "* **\\[MDAPI-7]\\[GO]** Fixed.\n\n  * Nested.\n", "line 2 is blank"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text, err := checkFragment(7, tt.text)
			if tt.err == "" {
				if err != nil {
					t.Fatal(err)
				}
				if strings.Contains(text, "\r") || strings.HasSuffix(text, "\n") {
					t.Fatalf("text must have LF endings and no trailing newline: %q", text)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.err) {
				t.Fatalf("error %v, want one containing %q", err, tt.err)
			}
		})
	}
}

func TestReadFragmentsReportsEveryProblem(t *testing.T) {
	dir := t.TempDir()
	writeFiles(t, dir, map[string]string{
		"README.md":    "Docs, not a fragment.",
		"MDAPI-2.md":   "* **\\[MDAPI-2]\\[GO]** Two.\n",
		"MDAPI-10.md":  "* **\\[MDAPI-10]\\[GO]** Ten.\n",
		"MDAPI-3.md":   "* Three.\n",
		"notes.txt":    "stray",
		"MDAPI-02.md":  "* **\\[MDAPI-2]\\[GO]** Two again.\n",
		"MDAPI-4.md/x": "A directory named like a fragment.",
	})
	fragments, err := readFragments(dir)
	if err == nil {
		t.Fatal("bad files must be reported")
	}
	for _, problem := range []string{"MDAPI-3.md: line 1", "notes.txt: a fragment", "MDAPI-02.md: a fragment", "MDAPI-4.md: unexpected directory"} {
		if !strings.Contains(err.Error(), problem) {
			t.Errorf("%q must be reported, got %v", problem, err)
		}
	}
	if len(fragments) != 2 || fragments[0].key != 2 || fragments[1].key != 10 {
		t.Fatalf("valid fragments must be sorted by ticket number, got %+v", fragments)
	}
}

func TestVersionSection(t *testing.T) {
	section, err := versionSection(testNotes, "v0.4.0")
	if err != nil {
		t.Fatal(err)
	}
	if section != "* **\\[MDAPI-1]\\[GO]** First.\n" {
		t.Fatalf("unexpected section %q", section)
	}
	if section, err = versionSection(testNotes, "0.3.0"); err != nil || section != "* Old.\n" {
		t.Fatalf("last section: %q, %v", section, err)
	}
	if _, err = versionSection(testNotes, "v9.9.9"); err == nil {
		t.Fatal("a missing version must be an error")
	}
	if _, err = versionSection("## v1.0.0\n\n## v0.9.0\n* x\n", "v1.0.0"); err == nil {
		t.Fatal("an empty section must be an error")
	}
}

func TestAddRelease(t *testing.T) {
	fragments := []fragment{{key: 2, text: "* **\\[MDAPI-2]\\[GO]** Two."}, {key: 5, text: "* **\\[MDAPI-5]\\[GO]** Five."}}
	updated, err := addRelease(testNotes, "v0.5.0", fragments)
	if err != nil {
		t.Fatal(err)
	}
	want := "Intro.\n\n## v0.5.0\n\n* **\\[MDAPI-2]\\[GO]** Two.\n* **\\[MDAPI-5]\\[GO]** Five.\n\n## v0.4.0\n"
	if !strings.Contains(updated, want) {
		t.Fatalf("new section not above the latest one:\n%s", updated)
	}
	for _, version := range []string{"v0.4.0", "v0.3.9", "0.5.0", "v0.5"} {
		if _, err := addRelease(testNotes, version, fragments); err == nil {
			t.Errorf("version %s must be rejected", version)
		}
	}
	if _, err := addRelease(testNotes, "v0.5.0", nil); err == nil {
		t.Error("a release without fragments must be rejected")
	}
	if _, err := addRelease("# Notes\n\n## Unreleased\n", "v0.5.0", fragments); err == nil {
		t.Error("an unparseable latest heading must be rejected")
	}
	updated, err = addRelease("# Notes\n", "v0.1.0", fragments[:1])
	if err != nil || !strings.HasSuffix(updated, "## v0.1.0\n\n* **\\[MDAPI-2]\\[GO]** Two.\n") {
		t.Errorf("first release: %q, %v", updated, err)
	}
	crlf := strings.ReplaceAll(testNotes, "\n", "\r\n")
	updated, err = addRelease(crlf, "v0.5.0", fragments)
	if err != nil || strings.Contains(updated, "\r") || !strings.Contains(updated, want) {
		t.Errorf("CRLF notes must be rewritten with LF: %q, %v", updated, err)
	}
}

func TestRunNotesAndCheck(t *testing.T) {
	root := t.TempDir()
	writeFiles(t, root, map[string]string{
		notesFile:                    testNotes,
		fragmentsDir + "/MDAPI-2.md": "* **\\[MDAPI-2]\\[GO]** Two.\n",
	})
	var out bytes.Buffer
	if err := run(root, []string{"notes", "v0.4.0"}, &out); err != nil || out.String() != "* **\\[MDAPI-1]\\[GO]** First.\n" {
		t.Fatalf("notes: %q, %v", out.String(), err)
	}
	out.Reset()
	if err := run(root, []string{"check"}, &out); err != nil || !strings.Contains(out.String(), "1 fragments OK") {
		t.Fatalf("check: %q, %v", out.String(), err)
	}
	if err := run(root, []string{"publish"}, &out); err == nil {
		t.Fatal("an unknown command must be an error")
	}
}

// gitRepo creates a repository with a release notes file and one fragment, isolated from the user's git config.
func gitRepo(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}
	t.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("GIT_AUTHOR_NAME", "test")
	t.Setenv("GIT_AUTHOR_EMAIL", "test@example.com")
	t.Setenv("GIT_COMMITTER_NAME", "test")
	t.Setenv("GIT_COMMITTER_EMAIL", "test@example.com")
	root := t.TempDir()
	writeFiles(t, root, map[string]string{
		notesFile:                    testNotes,
		fragmentsDir + "/README.md":  "Docs.\n",
		fragmentsDir + "/MDAPI-2.md": "* **\\[MDAPI-2]\\[GO]** Two.\n",
	})
	for _, args := range [][]string{{"init", "--quiet"}, {"add", "."}, {"commit", "--quiet", "-m", "init"}} {
		if _, err := git(root, args...); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestPrepareDryRunChangesNothing(t *testing.T) {
	root := gitRepo(t)
	var out bytes.Buffer
	if err := run(root, []string{"prepare", "v0.5.0", "--dry-run"}, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "## v0.5.0") || !strings.Contains(out.String(), "release-notes/MDAPI-2.md") {
		t.Fatalf("dry run output:\n%s", out.String())
	}
	if status, _ := git(root, "status", "--porcelain"); status != "" {
		t.Fatalf("dry run changed the tree: %s", status)
	}
	if tags, _ := git(root, "tag"); tags != "" {
		t.Fatalf("dry run created tags: %s", tags)
	}
}

func TestPrepareRefusesDirtyTree(t *testing.T) {
	root := gitRepo(t)
	writeFiles(t, root, map[string]string{"stray.txt": "x"})
	if err := run(root, []string{"prepare", "v0.5.0"}, &bytes.Buffer{}); err == nil || !strings.Contains(err.Error(), "working tree") {
		t.Fatalf("prepare on a dirty tree: %v", err)
	}
}

func TestPrepareRefusesExistingTag(t *testing.T) {
	root := gitRepo(t)
	if _, err := git(root, "tag", "v0.5.0"); err != nil {
		t.Fatal(err)
	}
	if err := run(root, []string{"prepare", "v0.5.0"}, &bytes.Buffer{}); err == nil || err.Error() != "tag v0.5.0 already exists" {
		t.Fatalf("prepare with an existing tag: %v", err)
	}
	if subject, _ := git(root, "log", "-1", "--format=%s"); subject != "init" {
		t.Fatalf("a refused release made commit %q", subject)
	}
	if status, _ := git(root, "status", "--porcelain"); status != "" {
		t.Fatalf("a refused release changed the tree: %s", status)
	}
}

func TestPrepareFailedCommitTellsHowToUndo(t *testing.T) {
	root := gitRepo(t)
	hook := filepath.Join(root, ".git", "hooks", "pre-commit")
	if err := os.WriteFile(hook, []byte("#!/bin/sh\nexit 1\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	err := run(root, []string{"prepare", "v0.5.0"}, &bytes.Buffer{})
	if err == nil || !strings.HasSuffix(err.Error(), "undo it with: git reset --hard HEAD") {
		t.Fatalf("prepare with a failing commit hook: %v", err)
	}
	if _, err := git(root, "reset", "--hard", "HEAD"); err != nil {
		t.Fatal(err)
	}
	if status, _ := git(root, "status", "--porcelain"); status != "" {
		t.Fatalf("the undo command left changes: %s", status)
	}
}

func TestPrepareCommitsAndTags(t *testing.T) {
	root := gitRepo(t)
	if err := run(root, []string{"prepare", "v0.5.0"}, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	if tag, _ := git(root, "describe", "--exact-match", "HEAD"); tag != "v0.5.0" {
		t.Fatalf("HEAD tag %q", tag)
	}
	if _, err := os.Stat(filepath.Join(root, fragmentsDir, "MDAPI-2.md")); !os.IsNotExist(err) {
		t.Fatal("the released fragment must be removed")
	}
	if _, err := os.Stat(filepath.Join(root, fragmentsDir, "README.md")); err != nil {
		t.Fatal("the fragment docs must stay")
	}
	notes, _ := os.ReadFile(filepath.Join(root, notesFile))
	if section, err := versionSection(string(notes), "v0.5.0"); err != nil || section != "* **\\[MDAPI-2]\\[GO]** Two.\n" {
		t.Fatalf("released section: %q, %v", section, err)
	}
	if status, _ := git(root, "status", "--porcelain"); status != "" {
		t.Fatalf("release left changes: %s", status)
	}
}
