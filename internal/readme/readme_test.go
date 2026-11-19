// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

// Package readme checks the Markdown documentation: every ```go block of README.md must be the exact text of a sample
// under cmd/samples without its license header, which the build and go vet compile, and the links of the Markdown
// files must lead to existing files and headings.
package readme

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var goBlock = regexp.MustCompile("(?ms)^```go\n(.*?)^```$")

func readNormalized(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return strings.ReplaceAll(string(data), "\r\n", "\n")
}

// licenseHeader starts every Go file; README.md shows a sample without it.
const licenseHeader = "// Copyright (c) 2026 Devexperts LLC.\n// SPDX-License-Identifier: MPL-2.0\n\n"

// readSample returns the text of a sample without its license header, as README.md shows it.
func readSample(t *testing.T, path string) string {
	t.Helper()
	return strings.TrimPrefix(readNormalized(t, path), licenseHeader)
}

func TestReadmeGoBlocksAreSamples(t *testing.T) {
	samples, err := filepath.Glob("../../cmd/samples/*/main.go")
	if err != nil {
		t.Fatal(err)
	}
	sampleTexts := make(map[string]string, len(samples))
	for _, sample := range samples {
		sampleTexts[readSample(t, sample)] = sample
	}

	blocks := goBlock.FindAllStringSubmatch(readNormalized(t, "../../README.md"), -1)
	if len(blocks) == 0 {
		t.Fatal("README.md has no ```go blocks")
	}
	for i, block := range blocks {
		if _, ok := sampleTexts[block[1]]; !ok {
			first, _, _ := strings.Cut(block[1], "\n")
			t.Errorf("README.md ```go block %d (%q) is not the text of any cmd/samples/*/main.go; "+
				"copy the sample into the README or the README block into a sample", i+1, first)
		}
	}
}
