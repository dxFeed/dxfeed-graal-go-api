// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package api_test

import (
	"regexp"
	"testing"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
)

func TestGetBuildInfo(t *testing.T) {
	info, err := api.GetBuildInfo()
	if err != nil {
		t.Fatal(err)
	}
	// The QDS version starts with major.minor; the runtime and its collectors are always known on the desktop
	// platforms. The texts are for display, so only their presence is checked.
	if !regexp.MustCompile(`^\d+\.\d+`).MatchString(info.QDVersion) {
		t.Errorf("QDVersion = %q", info.QDVersion)
	}
	if info.Vendor == "" || info.GC == "" {
		t.Errorf("vendor or GC is empty: %+v", info)
	}
}
