// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package readme

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/doc/comment"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// A code block of the documentation: a ```go block of a Markdown file or a code block of a doc comment of package pkg.
type codeBlock struct {
	where string // file:line
	text  string
}

// skipMarker before a ```go block of a Markdown file leaves the block out of the check; the comment gives the reason.
const skipMarker = "<!-- not compiled:"

// markdownGoBlocks returns the ```go blocks of the Markdown files that are not the text of a sample: the build compiles
// the samples.
func markdownGoBlocks(t *testing.T) []codeBlock {
	t.Helper()
	samples, err := filepath.Glob("../../cmd/samples/*/main.go")
	if err != nil {
		t.Fatal(err)
	}
	sampleTexts := make(map[string]bool, len(samples))
	for _, sample := range samples {
		sampleTexts[readSample(t, sample)] = true
	}
	var blocks []codeBlock
	for _, file := range markdownFiles(t) {
		text := readNormalized(t, file)
		for _, match := range goBlock.FindAllStringSubmatchIndex(text, -1) {
			body := text[match[2]:match[3]]
			before := strings.TrimRight(text[:match[0]], "\n")
			if sampleTexts[body] || strings.HasPrefix(before[strings.LastIndex(before, "\n")+1:], skipMarker) {
				continue
			}
			line := strings.Count(text[:match[2]], "\n") + 1
			blocks = append(blocks, codeBlock{fmt.Sprintf("%s:%d", relative(file), line), body})
		}
	}
	return blocks
}

// docCodeBlocks returns the code blocks of the doc comments of the packages under pkg: they show the API, while the
// code blocks of the commands under cmd and internal are command lines.
func docCodeBlocks(t *testing.T) []codeBlock {
	t.Helper()
	var blocks []codeBlock
	fset := token.NewFileSet()
	var parse comment.Parser
	err := filepath.WalkDir("../../pkg", func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return err
		}
		source, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		file, err := parser.ParseFile(fset, path, source, parser.ParseComments)
		if err != nil {
			return err
		}
		for _, group := range file.Comments {
			for _, block := range parse.Parse(group.Text()).Content {
				if code, ok := block.(*comment.Code); ok {
					blocks = append(blocks, codeBlock{fmt.Sprintf("%s:%d", relative(path),
						codeLine(fset, source, group, code.Text)), code.Text})
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return blocks
}

// codeLine returns the line of a code block of a comment: the first line of the comment that holds the first line of
// the code after the comment marker and a tab.
func codeLine(fset *token.FileSet, source []byte, group *ast.CommentGroup, code string) int {
	start := fset.Position(group.Pos())
	first, _, _ := strings.Cut(code, "\n")
	lines := strings.Split(strings.ReplaceAll(string(source[start.Offset:fset.Position(group.End()).Offset]), "\r\n",
		"\n"), "\n")
	for i, line := range lines {
		if strings.TrimPrefix(line, "//\t") == first {
			return start.Line + i
		}
	}
	return start.Line
}

func relative(path string) string {
	return filepath.ToSlash(strings.TrimPrefix(filepath.Clean(path), filepath.Clean("../..")+string(filepath.Separator)))
}

// TestCodeBlocksTypeCheck type-checks the Go code blocks of the Markdown files and of the doc comments against the API,
// so that a renamed or changed function breaks the test. A block that is not a whole file is a fragment: it goes into
// the body of a function that returns an error, the packages it names are imported, and its lines "..." are dropped.
// The variables that a fragment uses but does not declare (a subscription, a feed, a time) and the variables it does
// not use are left alone, so the calls on an undeclared variable are not checked; any other error fails, such as a missing method, a wrong argument or an unknown name of a
// package.
func TestCodeBlocksTypeCheck(t *testing.T) {
	blocks := append(markdownGoBlocks(t), docCodeBlocks(t)...)
	if len(blocks) == 0 {
		t.Fatal("no code blocks found")
	}
	packages := modulePackages(t)
	fset := token.NewFileSet()
	type pending struct {
		block codeBlock
		file  *ast.File
	}
	var files []pending
	imports := map[string]bool{}
	for _, block := range blocks {
		source, used := wrap(block.text, packages)
		file, err := parser.ParseFile(fset, block.where, source, 0)
		if err != nil {
			t.Errorf("%s: the block does not parse: %v\n%s", block.where, err, source)
			continue
		}
		for _, path := range used {
			imports[path] = true
		}
		files = append(files, pending{block, file})
	}

	imp := exportImporter(t, fset, imports)
	for _, f := range files {
		var errs []string
		config := types.Config{Importer: imp, Error: func(err error) {
			if e, ok := err.(types.Error); ok && !ignored(e.Msg) {
				errs = append(errs, e.Msg)
			}
		}}
		_, _ = config.Check("fragment", fset, []*ast.File{f.file}, nil)
		for _, msg := range errs {
			t.Errorf("%s: %s", f.block.where, msg)
		}
	}
}

var (
	undefinedLocal = regexp.MustCompile(`^undefined: [A-Za-z_]\w*$`)
	unusedName     = regexp.MustCompile(`declared and not used|is not used|imported and not used`)
	ellipsisLine   = regexp.MustCompile(`(?m)^\s*\.\.\.\s*$`)
	selector       = regexp.MustCompile(`\b([a-z][a-z0-9]*)\.[A-Za-z_]`)
)

// ignored reports the errors that come from a fragment rather than from the API: a variable it does not declare or
// does not use.
func ignored(msg string) bool {
	return undefinedLocal.MatchString(msg) || unusedName.MatchString(msg)
}

// standardPackages are the standard packages that the code blocks may name.
var standardPackages = map[string]string{
	"bytes": "bytes", "context": "context", "errors": "errors", "fmt": "fmt", "io": "io", "log": "log", "net": "net",
	"os": "os", "signal": "os/signal", "slog": "log/slog", "sort": "sort", "strings": "strings", "sync": "sync",
	"time": "time",
}

// wrap returns the source of a block as a file of package fragment and the import paths of the packages it names.
func wrap(text string, packages map[string]string) (string, []string) {
	text = ellipsisLine.ReplaceAllString(text, "")
	var used []string
	seen := map[string]bool{}
	for _, match := range selector.FindAllStringSubmatch(text, -1) {
		name := match[1]
		path, ok := packages[name]
		if !ok {
			path, ok = standardPackages[name]
		}
		if ok && !seen[name] {
			seen[name] = true
			used = append(used, path)
		}
	}
	sort.Strings(used)
	if strings.HasPrefix(strings.TrimSpace(text), "package ") {
		return text, used
	}
	var b strings.Builder
	b.WriteString("package fragment\n\nimport (\n")
	for _, path := range used {
		fmt.Fprintf(&b, "\t%q\n", path)
	}
	b.WriteString(")\n\nfunc _() error {\n")
	b.WriteString(text)
	b.WriteString("\n\treturn nil\n}\n")
	return b.String(), used
}

// modulePackages maps the names of the packages under pkg to their import paths.
func modulePackages(t *testing.T) map[string]string {
	t.Helper()
	packages := map[string]string{}
	err := filepath.WalkDir("../../pkg", func(path string, entry os.DirEntry, err error) error {
		if err != nil || !entry.IsDir() {
			return err
		}
		if goFiles, _ := filepath.Glob(filepath.Join(path, "*.go")); len(goFiles) > 0 {
			packages[filepath.Base(path)] = "github.com/dxfeed/dxfeed-graal-go-api/" + relative(path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return packages
}

// exportImporter returns an importer of the export data that go list -export writes for the packages and their
// dependencies; it covers the packages that use cgo, which a source importer does not.
func exportImporter(t *testing.T, fset *token.FileSet, packages map[string]bool) types.Importer {
	t.Helper()
	args := []string{"list", "-export", "-deps", "-f", "{{.ImportPath}}={{.Export}}"}
	for path := range packages {
		args = append(args, path)
	}
	cmd := exec.Command("go", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("go list -export: %v\n%s", err, stderr.String())
	}
	exports := map[string]string{}
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		if path, export, ok := strings.Cut(scanner.Text(), "="); ok && export != "" {
			exports[path] = export
		}
	}
	return importer.ForCompiler(fset, "gc", func(path string) (io.ReadCloser, error) {
		export, ok := exports[path]
		if !ok {
			return nil, fmt.Errorf("no export data for %s", path)
		}
		return os.Open(export)
	})
}
