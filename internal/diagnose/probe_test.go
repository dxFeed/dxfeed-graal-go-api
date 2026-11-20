// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package diagnose

import (
	"context"
	"net"
	"testing"
	"time"
)

func TestProbeHostNotResolved(t *testing.T) {
	// The .invalid domain never resolves (RFC 2606).
	probe := ProbeHost(context.Background(), "diagnose-test.invalid", 7300, time.Second)
	if probe.Err == "" || len(probe.IPs) != 0 {
		t.Errorf("%+v, want an error and no IP addresses", probe)
	}
}

func TestProbeHostConnects(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	probe := ProbeHost(context.Background(), "127.0.0.1", port, time.Second)
	if probe.Err != "" || len(probe.IPs) != 1 || probe.IPs[0].IP != "127.0.0.1" || probe.IPs[0].Err != "" {
		t.Errorf("open port: %+v", probe)
	}
	_ = listener.Close()
	probe = ProbeHost(context.Background(), "127.0.0.1", port, time.Second)
	if probe.Err != "" || len(probe.IPs) != 1 || probe.IPs[0].Err == "" {
		t.Errorf("closed port: %+v, want a TCP error", probe)
	}
}
