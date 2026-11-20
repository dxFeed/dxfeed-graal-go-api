// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package diagnose

import (
	"strconv"
	"strings"
)

// Target is a part of an address: a host and a port to probe, or a part that is skipped with the reason.
type Target struct {
	// Part is the part of the address as written, trimmed, with its properties.
	Part string
	Host string
	Port int
	// Skip is why the part is not probed; empty when it is.
	Skip string
}

// The reasons of a skipped part.
const (
	SkipIPv6       = "an IPv6 literal is not probed"
	SkipNotSocket  = "not a host:port socket address"
	SkipNoPort     = "no port"
	SkipServer     = "a server address (no host)"
	SkipPortRange  = "a range of ports"
	SkipOutOfRange = "the port is out of range"
)

// SplitAddress splits an address into the parts to probe: the addresses separated by commas, also within a group in
// parentheses, without the properties in brackets, the codecs before "+" and the user information before "@". Only a
// host:port socket address is probed: a server address, a tape file, a dxLink URL, a range of ports and an IPv6
// literal are skipped.
func SplitAddress(address string) []Target {
	var targets []Target
	for _, part := range splitTopLevel(address) {
		if part[0] == '(' {
			if end := matchingParen(part); end > 0 {
				targets = append(targets, SplitAddress(part[1:end])...)
				continue
			}
		}
		targets = append(targets, splitPart(part))
	}
	return targets
}

// splitTopLevel splits the text by the commas outside parentheses and brackets and drops the empty parts.
func splitTopLevel(text string) []string {
	var parts []string
	parens, brackets, start := 0, 0, 0
	add := func(end int) {
		if part := strings.TrimSpace(text[start:end]); part != "" {
			parts = append(parts, part)
		}
		start = end + 1
	}
	for i := 0; i < len(text); i++ {
		switch text[i] {
		case '(':
			parens++
		case ')':
			parens--
		case '[':
			brackets++
		case ']':
			brackets--
		case ',':
			if parens <= 0 && brackets <= 0 {
				add(i)
			}
		}
	}
	add(len(text))
	return parts
}

// matchingParen returns the index of the parenthesis that closes the one at index 0, or -1.
func matchingParen(text string) int {
	depth := 0
	for i := 0; i < len(text); i++ {
		switch text[i] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

func splitPart(part string) Target {
	t := Target{Part: part}
	var s strings.Builder
	ipv6 := false
	for i := 0; i < len(part); i++ {
		if part[i] != '[' {
			s.WriteByte(part[i])
			continue
		}
		end := strings.IndexByte(part[i:], ']')
		if end < 0 {
			end = len(part) - i
		}
		if !strings.Contains(part[i:i+end], "=") {
			ipv6 = true
		}
		i += end
	}
	socket := s.String()
	switch {
	case ipv6:
		t.Skip = SkipIPv6
		return t
	case strings.Contains(socket, "://"):
		t.Skip = SkipNotSocket
		return t
	}
	socket = socket[strings.LastIndexByte(socket, '+')+1:]
	socket = strings.TrimSpace(socket[strings.LastIndexByte(socket, '@')+1:])
	colon := strings.LastIndexByte(socket, ':')
	if colon < 0 {
		t.Skip = SkipNoPort
		return t
	}
	if colon == 0 {
		t.Skip = SkipServer
		return t
	}
	host, port := socket[:colon], socket[colon+1:]
	switch {
	case strings.ContainsAny(host, ":/"):
		t.Skip = SkipNotSocket
	case strings.Contains(port, "-"):
		t.Skip = SkipPortRange
	case port == "" || len(port) > 5 || strings.Trim(port, "0123456789") != "":
		t.Skip = SkipNotSocket
	default:
		n, _ := strconv.Atoi(port)
		if n < 1 || n > 65535 {
			t.Skip = SkipOutOfRange
		} else {
			t.Host, t.Port = host, n
		}
	}
	return t
}
