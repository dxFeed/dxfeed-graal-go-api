package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// gorelease is the pinned gorelease that "release api" runs with go run.
const gorelease = "golang.org/x/exp/cmd/gorelease@v0.0.0-20261118215729-2791e4028fd4"

// goreleaseFunc runs gorelease in a directory with arguments and returns its report.
type goreleaseFunc func(dir string, args ...string) (string, error)

// runGorelease runs gorelease in dir and returns its report, also when it fails. gorelease compares the commit HEAD,
// refuses a working tree with any change (untracked files too), and needs cgo set up as for a build.
func runGorelease(dir string, args ...string) (string, error) {
	cmd := exec.Command("go", append([]string{"run", gorelease}, args...)...)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String(), fmt.Errorf("gorelease %s: %w: %s", strings.Join(args, " "), err,
			strings.TrimSpace(stderr.String()))
	}
	return stdout.String(), nil
}

// apiCheck compares the exported API with the latest release tag before version, or with the latest release tag if
// version is empty, and returns the report of gorelease. gorelease itself accepts any change in v0, so an
// incompatible change that it reports needs:
//   - with a version (a release): a new major version, or for v0 at least a new minor version;
//   - without a version, on a branch: a fragment that the branch changes since mainRef, marked [BREAKING]; the
//     incompatible changes that mainRef already has do not count;
//   - without a version, on mainRef itself: an unreleased fragment marked [BREAKING].
func apiCheck(root, version, mainRef string, gorelease goreleaseFunc) (string, error) {
	var target [3]int
	if version != "" {
		var ok bool
		if target, ok = parseVersion(version); !ok || !strings.HasPrefix(version, "v") {
			return "", fmt.Errorf("version %q is not vX.Y.Z", version)
		}
	}
	tags, err := git(root, "tag", "--list", "v*")
	if err != nil {
		return "", err
	}
	base, baseVersion := "", [3]int{}
	for _, tag := range strings.Fields(tags) {
		v, ok := parseVersion(tag)
		if !ok || (version != "" && !less(v, target)) || (base != "" && !less(baseVersion, v)) {
			continue
		}
		base, baseVersion = tag, v
	}
	if base == "" {
		return "", errors.New("no vX.Y.Z release tag to compare with; fetch the tags")
	}

	// The version is checked below, not by gorelease: it refuses a version that the Go proxy already has, which a
	// re-run of a release would hit.
	report, err := compare(root, base, gorelease)
	report = fmt.Sprintf("Compared with %s.\n\n%s", base, report)
	if err != nil {
		return report, err
	}
	incompatible := incompatibleChanges(report)
	if len(incompatible) == 0 {
		return report, nil
	}
	if version != "" {
		if target[0] > baseVersion[0] || (baseVersion[0] == 0 && target[1] > baseVersion[1]) {
			return report, nil
		}
		need := "a new major version"
		if baseVersion[0] == 0 {
			need = "at least a new minor version"
		}
		return report, fmt.Errorf("%s has incompatible changes since %s, which need %s:\n  %s", version, base, need,
			strings.Join(incompatible, "\n  "))
	}
	return report, branchCheck(root, base, mainRef, incompatible, gorelease)
}

// branchCheck checks the incompatible changes since the release tag base on a branch or on mainRef itself.
func branchCheck(root, base, mainRef string, incompatible []string, gorelease goreleaseFunc) error {
	fragments, err := readFragments(filepath.Join(root, fragmentsDir))
	if err != nil {
		return err
	}
	head, err := git(root, "rev-parse", "HEAD")
	if err != nil {
		return err
	}
	mergeBase, err := git(root, "merge-base", "HEAD", mainRef)
	if err != nil {
		return fmt.Errorf("%w\nfetch %s to compare the branch with it", err, mainRef)
	}
	if mergeBase == head {
		for _, f := range fragments {
			if isBreaking(f) {
				return nil
			}
		}
		return fmt.Errorf("incompatible changes since %s, and no fragment in %s/ is marked [BREAKING]:\n  %s", base,
			fragmentsDir, strings.Join(incompatible, "\n  "))
	}

	// The incompatible changes that the branch adds to those of mainRef.
	dir, err := os.MkdirTemp("", "release-api-")
	if err != nil {
		return err
	}
	defer func() { _ = os.RemoveAll(dir) }()
	if _, err := git(root, "worktree", "add", "--quiet", "--detach", dir, mergeBase); err != nil {
		return err
	}
	defer func() { _, _ = git(root, "worktree", "remove", "--force", dir) }()
	mainReport, err := compare(dir, base, gorelease)
	if err != nil {
		return fmt.Errorf("on %s: %w", mainRef, err)
	}
	// An identifier that is incompatible on mainRef stays so, whatever the branch changes in it further.
	onMain := map[string]bool{}
	for _, change := range incompatibleChanges(mainReport) {
		onMain[identifier(change)] = true
	}
	var added []string
	for _, change := range incompatible {
		if !onMain[identifier(change)] {
			added = append(added, change)
		}
	}
	if len(added) == 0 {
		return nil
	}

	changed, err := git(root, "diff", "--name-only", mergeBase, "HEAD", "--", fragmentsDir)
	if err != nil {
		return err
	}
	changedFiles := map[string]bool{}
	for _, path := range strings.Fields(changed) {
		changedFiles[filepath.Base(path)] = true
	}
	for _, f := range fragments {
		if changedFiles[f.file] && isBreaking(f) {
			return nil
		}
	}
	return fmt.Errorf("incompatible changes, and no fragment that the branch changes in %s/ is marked [BREAKING]:"+
		"\n  %s\nmark the fragment of the change as breaking (see %s/%s), or keep the API compatible", fragmentsDir,
		strings.Join(added, "\n  "), fragmentsDir, fragmentDocs)
}

// compare runs gorelease in dir against the release tag base and returns its report. gorelease fails from v1 on when
// it finds an incompatible change, which apiCheck decides on itself, so only a report without a summary or with
// errors in loading the packages is an error.
func compare(dir, base string, gorelease goreleaseFunc) (string, error) {
	report, err := gorelease(dir, "-base="+base)
	if strings.Contains(report, "\n## errors in ") || strings.HasPrefix(report, "## errors in ") {
		return report, fmt.Errorf("gorelease could not load all the packages; see the report:\n%s", report)
	}
	if err != nil && !strings.Contains(report, "# summary") {
		return report, err
	}
	return report, nil
}

// identifier returns "<package>: <identifier>" of a change "<package>: <identifier>: <message>".
func identifier(change string) string {
	pkg, rest, _ := strings.Cut(change, ": ")
	id, _, _ := strings.Cut(rest, ":")
	return pkg + ": " + id
}

// isBreaking tells whether the top-level item of a fragment is marked [BREAKING].
func isBreaking(f fragment) bool {
	first, _, _ := strings.Cut(f.text, "\n")
	return strings.Contains(first, `\[GO]\[BREAKING]**`)
}

// incompatibleChanges returns the incompatible changes of a gorelease report as "<package>: <change>".
func incompatibleChanges(report string) []string {
	var changes []string
	pkg, inSection := "", false
	scanner := bufio.NewScanner(strings.NewReader(report))
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "# "):
			pkg, inSection = strings.TrimPrefix(line, "# "), false
		case strings.HasPrefix(line, "## "):
			inSection = strings.TrimSpace(strings.TrimPrefix(line, "## ")) == "incompatible changes"
		case inSection && strings.TrimSpace(line) != "":
			changes = append(changes, pkg+": "+strings.TrimSpace(line))
		}
	}
	return changes
}

// writeSummary writes the report and the result of the check in the Markdown of a GitHub job summary.
func writeSummary(w io.Writer, report string, err error) {
	result := "OK"
	if err != nil {
		result = "failed"
	}
	fmt.Fprintf(w, "### API compatibility: %s\n\n```text\n%s\n```\n", result, strings.TrimSpace(report))
	if err != nil {
		fmt.Fprintf(w, "\n```text\n%v\n```\n", err)
	}
}
