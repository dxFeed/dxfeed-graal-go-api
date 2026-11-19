// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package readme

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
	"unicode"
)

// repoBlob is the prefix of the links to files of this repository on GitHub; they are checked as local paths.
var repoBlob = regexp.MustCompile(`^https://github\.com/dxFeed/dxfeed-graal-go-api/(?:blob|tree)/main/`)

var (
	inlineLink    = regexp.MustCompile(`!?\[[^\]]*\]\(\s*<?([^)\s>]+)>?(?:\s+(?:"[^"]*"|'[^']*'))?\s*\)`)
	htmlLink      = regexp.MustCompile(`\b(?:src|href)="([^"]+)"`)
	referenceLink = regexp.MustCompile(`(?m)^\s{0,3}\[[^\]]+\]:\s*<?(\S+?)>?(?:\s+.*)?$`)
	inlineCode    = regexp.MustCompile("`[^`]*`")
	heading       = regexp.MustCompile(`^#{1,6}\s+(.*?)\s*#*\s*$`)
	htmlAnchor    = regexp.MustCompile(`<a\s+(?:name|id)=["']([^"']+)["']`)
	lineAnchor    = regexp.MustCompile(`^L\d+(?:-L\d+)?$`)
)

// markdownFiles returns the Markdown files whose links and code blocks are checked: those of the repository root, of
// docs and of the samples. The fragments under release-notes are excluded: they are joined into the release notes of a
// release.
func markdownFiles(t *testing.T) []string {
	t.Helper()
	var files []string
	for _, pattern := range []string{"../../*.md", "../../docs/*.md", "../../cmd/samples/*.md"} {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			t.Fatal(err)
		}
		files = append(files, matches...)
	}
	if len(files) == 0 {
		t.Fatal("no Markdown files found")
	}
	return files
}

// prose returns the lines of a Markdown text outside the fenced code blocks, with the inline code removed.
func prose(text string) []string {
	var lines []string
	fence := ""
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if fence != "" {
			if strings.HasPrefix(trimmed, fence) {
				fence = ""
			}
			lines = append(lines, "")
			continue
		}
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			fence = trimmed[:3]
			lines = append(lines, "")
			continue
		}
		lines = append(lines, inlineCode.ReplaceAllString(line, ""))
	}
	return lines
}

// slug turns a heading into its anchor as GitHub does: lower case, punctuation dropped, spaces turned into hyphens.
func slug(title string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(title) {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_':
			b.WriteRune(r)
		case r == ' ':
			b.WriteByte('-')
		}
	}
	return b.String()
}

// anchors returns the anchors of a Markdown file: its headings, numbered as GitHub numbers repeated ones, and the
// names of its HTML anchors.
func anchors(text string) map[string]bool {
	result := make(map[string]bool)
	seen := make(map[string]int)
	fenced := prose(text)
	raw := strings.Split(text, "\n")
	for i, line := range fenced {
		if m := heading.FindStringSubmatch(line); m != nil {
			// The inline code of a heading is part of its anchor: take the heading from the raw line.
			title := heading.FindStringSubmatch(raw[i])[1]
			title = strings.NewReplacer("`", "", "*", "", "[", "", "]", "").Replace(title)
			anchor := slug(title)
			if n := seen[anchor]; n > 0 {
				result[fmt.Sprintf("%s-%d", anchor, n)] = true
			} else {
				result[anchor] = true
			}
			seen[anchor]++
		}
		for _, m := range htmlAnchor.FindAllStringSubmatch(line, -1) {
			result[m[1]] = true
		}
	}
	return result
}

// links returns the link targets of a Markdown text with their line numbers.
func links(text string) map[int][]string {
	result := make(map[int][]string)
	for i, line := range prose(text) {
		for _, m := range inlineLink.FindAllStringSubmatch(line, -1) {
			result[i+1] = append(result[i+1], m[1])
		}
		for _, m := range htmlLink.FindAllStringSubmatch(line, -1) {
			result[i+1] = append(result[i+1], m[1])
		}
		if m := referenceLink.FindStringSubmatch(line); m != nil {
			result[i+1] = append(result[i+1], m[1])
		}
	}
	return result
}

// TestMarkdownLinks checks that the relative links of the Markdown files, and the links to files of this repository on
// GitHub, lead to existing files and headings. External links are not checked: the test runs offline.
func TestMarkdownLinks(t *testing.T) {
	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	anchorCache := make(map[string]map[string]bool)
	anchorsOf := func(path string) map[string]bool {
		if a, ok := anchorCache[path]; ok {
			return a
		}
		a := anchors(readNormalized(t, path))
		anchorCache[path] = a
		return a
	}

	for _, file := range markdownFiles(t) {
		name, _ := filepath.Rel(root, mustAbs(t, file))
		fileLinks := links(readNormalized(t, file))
		lines := make([]int, 0, len(fileLinks))
		for line := range fileLinks {
			lines = append(lines, line)
		}
		sort.Ints(lines)
		for _, line := range lines {
			for _, target := range fileLinks[line] {
				if strings.Contains(target, "://") && !repoBlob.MatchString(target) ||
					strings.HasPrefix(target, "mailto:") {
					continue
				}
				// The fragment and the query are cut from the target, not from the joined path, which may hold
				// '#' or '?' itself.
				rel, anchor, _ := strings.Cut(target, "#")
				rel, _, _ = strings.Cut(rel, "?")
				if decoded, err := url.PathUnescape(rel); err == nil {
					rel = decoded
				}
				var path string
				switch {
				case repoBlob.MatchString(rel):
					path = filepath.Join(root, filepath.FromSlash(repoBlob.ReplaceAllString(rel, "")))
				case rel == "":
					path = mustAbs(t, file)
				case strings.HasPrefix(rel, "/"):
					path = filepath.Join(root, filepath.FromSlash(rel))
				default:
					path = filepath.Join(filepath.Dir(mustAbs(t, file)), filepath.FromSlash(rel))
				}
				info, err := os.Stat(path)
				if err != nil {
					t.Errorf("%s:%d: link %q: no such file or directory", name, line, target)
					continue
				}
				if anchor == "" || info.IsDir() || !strings.HasSuffix(path, ".md") || lineAnchor.MatchString(anchor) {
					continue
				}
				if !anchorsOf(path)[anchor] {
					t.Errorf("%s:%d: link %q: no heading with the anchor #%s", name, line, target, anchor)
				}
			}
		}
	}
}

// TestMarkdownControlCharacters finds control characters in the Markdown files, such as a backspace that a script
// wrote for the "\b" of a Windows path.
func TestMarkdownControlCharacters(t *testing.T) {
	for _, file := range markdownFiles(t) {
		for i, line := range strings.Split(readNormalized(t, file), "\n") {
			for _, r := range line {
				if r < 0x20 && r != '\t' {
					t.Errorf("%s:%d: control character %U", strings.TrimPrefix(filepath.ToSlash(file), "../../"), i+1, r)
				}
			}
		}
	}
}

// TestSlug pins the anchor rules on headings of this repository.
func TestSlug(t *testing.T) {
	for title, want := range map[string]string{
		"How to connect to QD endpoint": "how-to-connect-to-qd-endpoint",
		"IPF & Schedule":                "ipf--schedule",
		"Subscriptions & Models":        "subscriptions--models",
		"Line endings":                  "line-endings",
		"Windows: the C compiler":       "windows-the-c-compiler",
	} {
		if got := slug(title); got != want {
			t.Errorf("slug(%q) = %q, want %q", title, got, want)
		}
	}
}

func mustAbs(t *testing.T, path string) string {
	t.Helper()
	abs, err := filepath.Abs(path)
	if err != nil {
		t.Fatal(err)
	}
	return abs
}
