package main

import (
	"bytes"
	"errors"
	"path/filepath"
	"strings"
	"testing"
)

const compatibleReport = `# example.com/m/pkg/a
## compatible changes
F: added

# summary
Suggested version: v0.5.0
`

const incompatibleReport = `# example.com/m/pkg/a
## incompatible changes
G: removed
H: changed from func() to func(int)

## compatible changes
F: added

# example.com/m/pkg/b
## incompatible changes
T.X: removed

# summary
Suggested version: v0.5.0
`

// mainReport has one of the incompatible changes of incompatibleReport.
const mainReport = `# example.com/m/pkg/b
## incompatible changes
T.X: removed
`

func TestIncompatibleChanges(t *testing.T) {
	got := incompatibleChanges(incompatibleReport)
	want := []string{
		"example.com/m/pkg/a: G: removed",
		"example.com/m/pkg/a: H: changed from func() to func(int)",
		"example.com/m/pkg/b: T.X: removed",
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("got %q, want %q", got, want)
	}
	if got := incompatibleChanges(compatibleReport); len(got) != 0 {
		t.Errorf("a compatible report has incompatible changes %q", got)
	}
}

// fakeGorelease returns report in root and mainReport elsewhere (in the worktree of the main branch), and records
// the arguments of its last call in root.
func fakeGorelease(root, report, mainReport string, args *[]string) goreleaseFunc {
	return func(dir string, a ...string) (string, error) {
		if filepath.Clean(dir) != filepath.Clean(root) {
			return mainReport, nil
		}
		*args = a
		return report, nil
	}
}

// apiRepo is gitRepo with release tags and the branch main0 at its commit; its fragment MDAPI-2 is not breaking.
func apiRepo(t *testing.T, tags ...string) string {
	t.Helper()
	root := gitRepo(t)
	for _, tag := range tags {
		if _, err := git(root, "tag", tag); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := git(root, "branch", "main0"); err != nil {
		t.Fatal(err)
	}
	return root
}

func commitFiles(t *testing.T, root string, files map[string]string) {
	t.Helper()
	writeFiles(t, root, files)
	for _, args := range [][]string{{"add", "."}, {"commit", "--quiet", "-m", "change"}} {
		if _, err := git(root, args...); err != nil {
			t.Fatal(err)
		}
	}
}

const breakingFragment = "* **\\[MDAPI-3]\\[GO]\\[BREAKING]** Three.\n  * **\\[BREAKING]** Removed G.\n"

func TestAPICheckOnMain(t *testing.T) {
	root := apiRepo(t, "0.3.0", "v0.4.0", "v0.10.0", "v0.9.1")
	var args []string
	report, err := apiCheck(root, "", "main0", fakeGorelease(root, compatibleReport, "", &args))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(args, " ") != "-base=v0.10.0" {
		t.Errorf("gorelease ran with %q, want the latest tag v0.10.0", args)
	}
	if !strings.HasPrefix(report, "Compared with v0.10.0.") || !strings.Contains(report, "F: added") {
		t.Errorf("report:\n%s", report)
	}

	_, err = apiCheck(root, "", "main0", fakeGorelease(root, incompatibleReport, "", &args))
	if err == nil || !strings.Contains(err.Error(), "[BREAKING]") || !strings.Contains(err.Error(), "T.X: removed") {
		t.Errorf("an incompatible change without a breaking fragment: %v", err)
	}

	// On the main branch, any unreleased breaking fragment allows the incompatible changes.
	writeFiles(t, root, map[string]string{fragmentsDir + "/MDAPI-3.md": breakingFragment})
	if _, err := apiCheck(root, "", "main0", fakeGorelease(root, incompatibleReport, "", &args)); err != nil {
		t.Errorf("an incompatible change with a breaking fragment: %v", err)
	}
}

func TestAPICheckBranch(t *testing.T) {
	root := apiRepo(t, "v0.4.0")
	// main0 has a breaking fragment, which does not cover the changes of the branch.
	commitFiles(t, root, map[string]string{fragmentsDir + "/MDAPI-2.md": "* **\\[MDAPI-2]\\[GO]\\[BREAKING]** Two.\n"})
	if _, err := git(root, "branch", "--force", "main0"); err != nil {
		t.Fatal(err)
	}
	commitFiles(t, root, map[string]string{"code.go": "package m\n"})

	var args []string
	if _, err := apiCheck(root, "", "main0", fakeGorelease(root, incompatibleReport, incompatibleReport, &args)); err != nil {
		t.Errorf("the incompatible changes of the main branch only: %v", err)
	}
	_, err := apiCheck(root, "", "main0", fakeGorelease(root, incompatibleReport, mainReport, &args))
	if err == nil || !strings.Contains(err.Error(), "G: removed") || strings.Contains(err.Error(), "T.X") {
		t.Errorf("new incompatible changes without a breaking fragment of the branch: %v", err)
	}
	if out, _ := git(root, "worktree", "list"); strings.Count(out, "\n") != 0 {
		t.Errorf("the worktree of the main branch is left:\n%s", out)
	}

	commitFiles(t, root, map[string]string{fragmentsDir + "/MDAPI-3.md": breakingFragment})
	if _, err := apiCheck(root, "", "main0", fakeGorelease(root, incompatibleReport, mainReport, &args)); err != nil {
		t.Errorf("new incompatible changes with a breaking fragment of the branch: %v", err)
	}

	if _, err := apiCheck(root, "", "no-such-branch", fakeGorelease(root, incompatibleReport, "", &args)); err == nil {
		t.Error("a missing main branch must be an error")
	}
}

func TestAPICheckRelease(t *testing.T) {
	root := apiRepo(t, "v0.4.0", "v0.4.1", "v0.5.0", "v1.2.0")
	tests := []struct {
		version, report, base string
		ok                    bool
	}{
		{"v0.4.2", compatibleReport, "v0.4.1", true},
		{"v0.4.2", incompatibleReport, "v0.4.1", false},
		{"v0.5.0", incompatibleReport, "v0.4.1", true}, // the tag of the version itself is not the base
		{"v0.6.0", incompatibleReport, "v0.5.0", true},
		{"v1.0.0", incompatibleReport, "v0.5.0", true},
		{"v1.3.0", compatibleReport, "v1.2.0", true},
		{"v1.3.0", incompatibleReport, "v1.2.0", false},
		{"v2.0.0", incompatibleReport, "v1.2.0", true},
	}
	for _, test := range tests {
		var args []string
		_, err := apiCheck(root, test.version, "main0", fakeGorelease(root, test.report, "", &args))
		if want := "-base=" + test.base; strings.Join(args, " ") != want {
			t.Errorf("%s: gorelease ran with %q, want %q", test.version, args, want)
		}
		if (err == nil) != test.ok {
			t.Errorf("%s with base %s: err = %v, want ok = %v", test.version, test.base, err, test.ok)
		}
	}
	if _, err := apiCheck(root, "0.6.0", "main0", fakeGorelease(root, compatibleReport, "", new([]string))); err == nil {
		t.Error("a version without v must be an error")
	}
	if _, err := apiCheck(root, "v0.4.0", "main0", fakeGorelease(root, compatibleReport, "", new([]string))); err == nil {
		t.Error("a version without a release tag before it must be an error")
	}
}

func TestAPICheckGoreleaseError(t *testing.T) {
	root := apiRepo(t, "v0.4.0", "v1.0.0")
	failing := func(string, ...string) (string, error) { return "partial", errors.New("gorelease failed") }
	report, err := apiCheck(root, "", "main0", failing)
	if err == nil || !strings.Contains(report, "partial") {
		t.Errorf("report %q, err %v", report, err)
	}

	// From v1 on, gorelease fails on an incompatible change; the report decides.
	v1 := func(string, ...string) (string, error) { return incompatibleReport, errors.New("exit status 1") }
	if _, err := apiCheck(root, "v2.0.0", "main0", v1); err != nil {
		t.Errorf("an incompatible change in a new major version: %v", err)
	}
	if _, err := apiCheck(root, "v1.1.0", "main0", v1); err == nil || !strings.Contains(err.Error(), "major") {
		t.Errorf("an incompatible change in a minor version of v1: %v", err)
	}

	baseErrors := "# example.com/m/pkg/a\n## errors in base version:\nno Go files\n\n# summary\nv1.1.0 is valid\n"
	loading := func(string, ...string) (string, error) { return baseErrors, nil }
	if _, err := apiCheck(root, "v1.1.0", "main0", loading); err == nil {
		t.Error("errors in loading the base version must be an error")
	}
}

func TestIdentifier(t *testing.T) {
	a := identifier("example.com/m/pkg/a: H: changed from func() to func(int)")
	b := identifier("example.com/m/pkg/a: H: changed from func() to func(string)")
	if a != b || a != "example.com/m/pkg/a: H" {
		t.Errorf("identifiers %q and %q", a, b)
	}
}

func TestWriteSummary(t *testing.T) {
	var b bytes.Buffer
	writeSummary(&b, "Compared with v0.4.0.\n", errors.New("incompatible"))
	if s := b.String(); !strings.Contains(s, "API compatibility: failed") || !strings.Contains(s, "incompatible") {
		t.Errorf("summary:\n%s", s)
	}
}
