package samples_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

var (
	// runLine matches a command line of a sample with its continuation lines, which end with a backslash.
	runLine  = regexp.MustCompile(`go run \./cmd/samples/([\w-]+)((?:[^\n]*\\\n)*[^\n]*)`)
	flagName = regexp.MustCompile(`(?:^|\s)-([a-zA-Z][\w-]*)`)
	// usageFlag matches a flag in the usage that the flag package prints: two spaces, a dash and the name.
	usageFlag = regexp.MustCompile(`(?m)^  -([\w-]+)`)
)

// TestSamplesUsage checks the command lines of the samples in their doc comments, in the README of the samples and in
// the documents of the repository: the sample exists and each flag of the line is in the usage that -h prints.
func TestSamplesUsage(t *testing.T) {
	dirs, err := filepath.Glob("*/main.go")
	if err != nil {
		t.Fatal(err)
	}
	names := make([]string, 0, len(dirs))
	for _, dir := range dirs {
		names = append(names, filepath.Base(filepath.Dir(dir)))
	}
	programs := buildAll(t, t.TempDir(), names...)
	flags := make(map[string]map[string]bool, len(names))
	for _, name := range names {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		out, err := exec.CommandContext(ctx, programs[name], "-h").CombinedOutput()
		cancel()
		if err != nil {
			t.Errorf("%s -h: %v\n%s", name, err, out)
			continue
		}
		flags[name] = map[string]bool{"h": true}
		for _, match := range usageFlag.FindAllStringSubmatch(string(out), -1) {
			flags[name][match[1]] = true
		}
	}

	files := append(dirs, "README.md", "../../README.md", "../../CONTRIBUTING.md", "../../HOW_TO_DEVELOP.md")
	docs, err := filepath.Glob("../../docs/*.md")
	if err != nil {
		t.Fatal(err)
	}
	lines := 0
	for _, file := range append(files, docs...) {
		source, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		text := strings.ReplaceAll(string(source), "\r\n", "\n")
		for _, match := range runLine.FindAllStringSubmatch(text, -1) {
			lines++
			name, args := match[1], strings.ReplaceAll(match[2], "//", "")
			known, ok := flags[name]
			if !ok {
				t.Errorf("%s: %q: no sample %s", file, match[0], name)
				continue
			}
			for _, flag := range flagName.FindAllStringSubmatch(args, -1) {
				if !known[flag[1]] {
					t.Errorf("%s: %q: %s has no flag -%s", file, match[0], name, flag[1])
				}
			}
		}
	}
	if lines == 0 {
		t.Fatal("no command lines of the samples found")
	}
}
