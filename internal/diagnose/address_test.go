// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package diagnose

import (
	"reflect"
	"testing"
)

func TestSplitAddress(t *testing.T) {
	for address, want := range map[string][]Target{
		"tls+gzip+user:pass@demo.example.com:7300[login=secret]": {
			{Part: "tls+gzip+user:pass@demo.example.com:7300[login=secret]", Host: "demo.example.com", Port: 7300}},
		"a.example.com:1, (b.example.com:2,c.example.com:3)[user=x,password=y]": {
			{Part: "a.example.com:1", Host: "a.example.com", Port: 1},
			{Part: "b.example.com:2", Host: "b.example.com", Port: 2},
			{Part: "c.example.com:3", Host: "c.example.com", Port: 3}},
		"ssl[trustStore=a.jks]+host.example.com:443": {
			{Part: "ssl[trustStore=a.jks]+host.example.com:443", Host: "host.example.com", Port: 443}},
		":7300":         {{Part: ":7300", Skip: SkipServer}},
		"tape:file.txt": {{Part: "tape:file.txt", Skip: SkipNotSocket}},
		"dxlink:wss://demo.example.com/dxlink-ws": {{Part: "dxlink:wss://demo.example.com/dxlink-ws", Skip: SkipNotSocket}},
		"host:7300-7310": {{Part: "host:7300-7310", Skip: SkipPortRange}},
		"host":           {{Part: "host", Skip: SkipNoPort}},
		"host:70000":     {{Part: "host:70000", Skip: SkipOutOfRange}},
		"host:0":         {{Part: "host:0", Skip: SkipOutOfRange}},
		"[::1]:7300":     {{Part: "[::1]:7300", Skip: SkipIPv6}},
		"(a:1":           {{Part: "(a:1", Host: "(a", Port: 1}},
		" , ":            nil,
	} {
		if got := SplitAddress(address); !reflect.DeepEqual(got, want) {
			t.Errorf("SplitAddress(%q) = %+v, want %+v", address, got, want)
		}
	}
}
