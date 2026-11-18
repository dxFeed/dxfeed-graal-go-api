package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	notesFile    = "ReleaseNotes.md"
	fragmentsDir = "release-notes"
	fragmentDocs = "README.md"
)

var (
	fragmentName = regexp.MustCompile(`^MDAPI-([1-9]\d*)\.md$`)
	topItem      = regexp.MustCompile(`^\* \*\*\\\[MDAPI-(\d+)]\\\[GO](?:\\\[BREAKING])?\*\* \S`)
	versionRe    = regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)$`)
)

// fragment is one release-notes/MDAPI-NNN.md file: a single top-level item for its ticket.
type fragment struct {
	key  int
	file string
	text string
}

// readFragments reads and checks every fragment in dir, sorted by ticket number; it reports all problems at once.
func readFragments(dir string) ([]fragment, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var fragments []fragment
	var problems []error
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() {
			problems = append(problems, fmt.Errorf("%s: unexpected directory; fragments are not read from subdirectories", name))
			continue
		}
		if name == fragmentDocs {
			continue
		}
		m := fragmentName.FindStringSubmatch(name)
		if m == nil {
			problems = append(problems, fmt.Errorf("%s: a fragment must be named MDAPI-NNN.md", name))
			continue
		}
		key, _ := strconv.Atoi(m[1])
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			problems = append(problems, err)
			continue
		}
		text, err := checkFragment(key, string(data))
		if err != nil {
			problems = append(problems, fmt.Errorf("%s: %w", name, err))
			continue
		}
		fragments = append(fragments, fragment{key: key, file: name, text: text})
	}
	sort.Slice(fragments, func(i, j int) bool { return fragments[i].key < fragments[j].key })
	return fragments, errors.Join(problems...)
}

// checkFragment returns the fragment text with LF line endings and no trailing blank lines.
func checkFragment(key int, data string) (string, error) {
	text := strings.TrimRight(strings.ReplaceAll(data, "\r\n", "\n"), "\n ")
	if text == "" {
		return "", errors.New("empty fragment")
	}
	lines := strings.Split(text, "\n")
	m := topItem.FindStringSubmatch(lines[0])
	if m == nil {
		return "", fmt.Errorf(`line 1 must start with "* **\[MDAPI-%d]\[GO]** " or "* **\[MDAPI-%d]\[GO]\[BREAKING]** "`,
			key, key)
	}
	if m[1] != strconv.Itoa(key) {
		return "", fmt.Errorf("line 1 names MDAPI-%s, the file name MDAPI-%d", m[1], key)
	}
	for i, line := range lines[1:] {
		switch {
		case strings.TrimSpace(line) == "":
			return "", fmt.Errorf("line %d is blank; a blank line splits the item", i+2)
		case !strings.HasPrefix(line, "  "):
			return "", fmt.Errorf("line %d is not indented; a fragment has one top-level item, nest the others", i+2)
		}
	}
	return text, nil
}

// heading reports whether line is a version heading and returns its version.
func heading(line string) (string, bool) {
	if !strings.HasPrefix(line, "## ") {
		return "", false
	}
	return strings.TrimSpace(strings.TrimPrefix(line, "## ")), true
}

// versionSection returns the body of the "## <version>" section of the release notes.
func versionSection(notes, version string) (string, error) {
	lines := strings.Split(strings.ReplaceAll(notes, "\r\n", "\n"), "\n")
	start := -1
	for i, line := range lines {
		v, ok := heading(line)
		if !ok {
			continue
		}
		if start >= 0 {
			return trimSection(lines[start:i], version)
		}
		if v == version {
			start = i + 1
		}
	}
	if start < 0 {
		return "", fmt.Errorf("%s has no section %s", notesFile, version)
	}
	return trimSection(lines[start:], version)
}

func trimSection(lines []string, version string) (string, error) {
	body := strings.Trim(strings.Join(lines, "\n"), "\n ")
	if body == "" {
		return "", fmt.Errorf("section %s of %s is empty", version, notesFile)
	}
	return body + "\n", nil
}

// parseVersion accepts vX.Y.Z and, for the old tags, X.Y.Z.
func parseVersion(version string) ([3]int, bool) {
	m := versionRe.FindStringSubmatch(version)
	if m == nil {
		return [3]int{}, false
	}
	var v [3]int
	for i := range v {
		v[i], _ = strconv.Atoi(m[i+1])
	}
	return v, true
}

func less(a, b [3]int) bool {
	for i := range a {
		if a[i] != b[i] {
			return a[i] < b[i]
		}
	}
	return false
}

// addRelease puts a "## <version>" section with the fragments above the latest version of the release notes.
func addRelease(notes, version string, fragments []fragment) (string, error) {
	next, ok := parseVersion(version)
	if !ok || !strings.HasPrefix(version, "v") {
		return "", fmt.Errorf("version %q must look like v1.2.3", version)
	}
	if len(fragments) == 0 {
		return "", fmt.Errorf("no fragments in %s/ to release", fragmentsDir)
	}
	lines := strings.Split(strings.ReplaceAll(notes, "\r\n", "\n"), "\n")
	at := len(lines)
	for i, line := range lines {
		latest, ok := heading(line)
		if !ok {
			continue
		}
		prev, ok := parseVersion(latest)
		if !ok {
			return "", fmt.Errorf("cannot parse the latest version %q in %s", latest, notesFile)
		}
		if !less(prev, next) {
			return "", fmt.Errorf("%s is not above the latest released version %s", version, latest)
		}
		at = i
		break
	}
	items := make([]string, len(fragments))
	for i, f := range fragments {
		items[i] = f.text
	}
	section := []string{"## " + version, "", strings.Join(items, "\n"), ""}
	result := append(append(append([]string{}, lines[:at]...), section...), lines[at:]...)
	return strings.Join(result, "\n"), nil
}
