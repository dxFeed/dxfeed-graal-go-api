// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package diagnose

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

// MaxProbedIPs is the greatest number of the IP addresses of a host that are probed.
const MaxProbedIPs = 16

// ProbeTimeout is the timeout of the DNS lookup and of each TCP connect of the probe.
const ProbeTimeout = 3 * time.Second

// IPProbe is the result of a TCP connect to an IP address.
type IPProbe struct {
	IP string
	// Elapsed is how long the connect took.
	Elapsed time.Duration
	// Err is why the connect failed; empty when it succeeded.
	Err string
}

// HostProbe is the result of the probe of a host: its IP addresses and a TCP connect to each.
type HostProbe struct {
	// Err is why the host did not resolve; empty when it did.
	Err string
	// Resolved is the number of the IP addresses of the host; only the first MaxProbedIPs are in IPs.
	Resolved int
	IPs      []IPProbe
}

// ProbeHost resolves the host and opens a TCP connection to the port of each of its IP addresses, at most
// MaxProbedIPs of them, one at a time with the timeout each; a connection is closed at once and nothing is sent.
func ProbeHost(ctx context.Context, host string, port int, timeout time.Duration) HostProbe {
	lookup, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	addresses, err := net.DefaultResolver.LookupIPAddr(lookup, host)
	if err != nil {
		return HostProbe{Err: errorText(err, timeout)}
	}
	if len(addresses) == 0 {
		return HostProbe{Err: "no addresses"}
	}
	result := HostProbe{Resolved: len(addresses)}
	if len(addresses) > MaxProbedIPs {
		addresses = addresses[:MaxProbedIPs]
	}
	dialer := net.Dialer{Timeout: timeout}
	for _, address := range addresses {
		start := time.Now()
		conn, err := dialer.DialContext(ctx, "tcp", net.JoinHostPort(address.String(), strconv.Itoa(port)))
		probe := IPProbe{IP: address.String(), Elapsed: time.Since(start)}
		if err != nil {
			probe.Err = errorText(err, timeout)
		} else {
			_ = conn.Close()
		}
		result.IPs = append(result.IPs, probe)
	}
	return result
}

// errorText returns the text of a network error without the operation and the address, which the report has already.
func errorText(err error, timeout time.Duration) string {
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, os.ErrDeadlineExceeded) {
		return fmt.Sprintf("timed out after %d ms", timeout.Milliseconds())
	}
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		if dnsErr.IsTimeout {
			return fmt.Sprintf("timed out after %d ms", timeout.Milliseconds())
		}
		return dnsErr.Err
	}
	var opErr *net.OpError
	if errors.As(err, &opErr) && opErr.Err != nil {
		var sysErr *os.SyscallError
		if errors.As(opErr.Err, &sysErr) {
			return sysErr.Err.Error()
		}
		return opErr.Err.Error()
	}
	return err.Error()
}
