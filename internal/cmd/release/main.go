// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

// Command release checks release note fragments, prints the notes of a version and prepares a release.
//
// Run it from the repository root:
//
//	go run ./internal/cmd/release check
//	go run ./internal/cmd/release api [v0.4.0]
//	go run ./internal/cmd/release notes v0.4.0
//	go run ./internal/cmd/release prepare v0.4.0 [--dry-run]
package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const usage = `usage:
  release check                      check release-notes/MDAPI-NNN.md fragments
  release api [<version>]            compare the API with the latest release tag (before <version>) with gorelease;
                                     an incompatible change needs a [BREAKING] fragment or, with <version>, a new
                                     major version (minor for v0)
  release notes <version>            print the section of <version> from ReleaseNotes.md
  release prepare <version> [--dry-run]
                                     move the fragments into ReleaseNotes.md, commit and tag <version>`

func main() {
	if err := run(".", os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "release:", err)
		os.Exit(1)
	}
}

func run(root string, args []string, out io.Writer) error {
	if len(args) == 0 {
		return errors.New(usage)
	}
	switch {
	case args[0] == "check" && len(args) == 1:
		fragments, err := readFragments(filepath.Join(root, fragmentsDir))
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "%d fragments OK\n", len(fragments))
		return nil
	case args[0] == "notes" && len(args) == 2:
		notes, err := os.ReadFile(filepath.Join(root, notesFile))
		if err != nil {
			return err
		}
		section, err := versionSection(string(notes), args[1])
		if err != nil {
			return err
		}
		_, err = io.WriteString(out, section)
		return err
	case args[0] == "api" && len(args) <= 2:
		version := ""
		if len(args) == 2 {
			version = args[1]
		}
		report, err := apiCheck(root, version, "origin/main", runGorelease)
		_, _ = io.WriteString(out, report)
		if path := os.Getenv("GITHUB_STEP_SUMMARY"); path != "" {
			f, ferr := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
			if ferr != nil {
				return errors.Join(err, ferr)
			}
			writeSummary(f, report, err)
			if ferr := f.Close(); ferr != nil {
				return errors.Join(err, ferr)
			}
		}
		return err
	case args[0] == "prepare" && len(args) == 2:
		return prepare(root, args[1], false, out)
	case args[0] == "prepare" && len(args) == 3 && args[2] == "--dry-run":
		return prepare(root, args[1], true, out)
	default:
		return errors.New(usage)
	}
}

// prepare moves the fragments under a new version heading, then commits and tags unless dryRun is set.
func prepare(root, version string, dryRun bool, out io.Writer) error {
	if !dryRun {
		status, err := git(root, "status", "--porcelain")
		if err != nil {
			return err
		}
		if status != "" {
			return errors.New("the working tree has changes; commit or stash them first")
		}
	}
	if _, err := git(root, "rev-parse", "-q", "--verify", "refs/tags/"+version); err == nil {
		return fmt.Errorf("tag %s already exists", version)
	}
	fragments, err := readFragments(filepath.Join(root, fragmentsDir))
	if err != nil {
		return err
	}
	notesPath := filepath.Join(root, notesFile)
	notes, err := os.ReadFile(notesPath)
	if err != nil {
		return err
	}
	updated, err := addRelease(string(notes), version, fragments)
	if err != nil {
		return err
	}
	section, err := versionSection(updated, version)
	if err != nil {
		return err
	}
	if dryRun {
		fmt.Fprintf(out, "## %s\n\n%s\nFragments to remove:\n", version, section)
		for _, f := range fragments {
			fmt.Fprintf(out, "  %s/%s\n", fragmentsDir, f.file)
		}
		return nil
	}

	if err := os.WriteFile(notesPath, []byte(updated), 0o644); err != nil {
		return err
	}
	rm := []string{"rm", "--quiet"}
	for _, f := range fragments {
		rm = append(rm, fragmentsDir+"/"+f.file)
	}
	steps := [][]string{
		rm,
		{"add", notesFile},
		{"commit", "--quiet", "-m", "Release " + version},
		{"tag", "-a", version, "-m", version},
	}
	// The tree was clean, so a failed step is undone by resetting to the commit the release started from.
	undo := "git reset --hard HEAD"
	for _, step := range steps {
		if _, err := git(root, step...); err != nil {
			return fmt.Errorf("%w\nthe release is incomplete; undo it with: %s", err, undo)
		}
		if step[0] == "commit" {
			undo = "git reset --hard HEAD~1"
		}
	}
	fmt.Fprintf(out, "Released %s locally with %d fragments. Check the commit, then push it and the tag:\n"+
		"  git push --atomic origin HEAD %s\n", version, len(fragments), version)
	return nil
}

func git(root string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = root
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return strings.TrimSpace(stdout.String()), nil
}
