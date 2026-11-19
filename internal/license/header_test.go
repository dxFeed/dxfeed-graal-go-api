// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

// Package license checks that every Go file of the repository starts with the license header. A blank line follows
// the header, so it is not a part of the package documentation.
package license

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const header = "// Copyright (c) 2026 Devexperts LLC.\n// SPDX-License-Identifier: MPL-2.0\n\n"

// TestGoFilesHaveLicenseHeader walks the repository, leaving out .git and the local directories .work and .claude (the
// worktrees of other branches), and fails for each Go file without the header.
func TestGoFilesHaveLicenseHeader(t *testing.T) {
	root := filepath.Join("..", "..")
	files := 0
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if name := entry.Name(); name == ".git" || name == ".work" || name == ".claude" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		files++
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if !bytes.HasPrefix(bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n")), []byte(header)) {
			t.Errorf("%s: no license header; start the file with\n%s", filepath.ToSlash(path), header)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if files == 0 {
		t.Fatal("no Go files found")
	}
}
