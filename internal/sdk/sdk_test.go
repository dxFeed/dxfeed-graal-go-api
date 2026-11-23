// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	"fmt"
	"strings"
	"testing"
)

func TestCheckVersion(t *testing.T) {
	major, minor, _ := strings.Cut(Version, ".")
	minor, _, _ = strings.Cut(minor, ".")
	tests := []struct {
		version string
		ok      bool
	}{
		{Version, true},
		{Version + "-SNAPSHOT", true},
		{major + "." + minor + ".99", true},
		{major + ".0.0", false},
		{major + "." + minor + "0.0", false},
		{"9" + major + "." + minor + ".0", false},
		{"", false},
		{"3.16", false},
		{"3.17.0", false}, // the release before the pinned one: no I/O statistics of message connectors
		{"3.17.9", false},
		{"3.18.0", false}, // no build information of the library
		{"3.18.9", false},
		{"a.b.c", false},
		{"3.-16.0", false},
		{"+" + major + ".+" + minor + ".+0", false},
	}
	for _, tt := range tests {
		err := CheckVersion(tt.version)
		if (err == nil) != tt.ok {
			t.Errorf("CheckVersion(%q) = %v, want ok=%v", tt.version, err, tt.ok)
		}
		if err != nil && !strings.Contains(err.Error(), Version) {
			t.Errorf("CheckVersion(%q) error %q does not name the pinned version", tt.version, err)
		}
	}
}

func TestCheckVersionOlderPatch(t *testing.T) {
	v, err := parseVersion(Version)
	if err != nil {
		t.Fatal(err)
	}
	if v[2] == 0 {
		t.Skip("the pinned version has patch 0")
	}
	older := fmt.Sprintf("%d.%d.%d", v[0], v[1], v[2]-1)
	if CheckVersion(older) == nil {
		t.Errorf("CheckVersion(%q) accepted an older patch", older)
	}
}

func TestArchivesMatchVersion(t *testing.T) {
	for platform, a := range Archives {
		if !strings.Contains(a.Name, "-"+Version+"-") {
			t.Errorf("%s: archive %s is not of SDK %s", platform, a.Name, Version)
		}
		if len(a.SHA256) != 64 || strings.ToLower(a.SHA256) != a.SHA256 {
			t.Errorf("%s: checksum %q is not 64 lower-case hex digits", platform, a.SHA256)
		}
	}
}
