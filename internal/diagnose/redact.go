// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

// Package diagnose holds the parts of the Diagnose tool of cmd/tools that need no native library: the redaction of
// secrets, the splitting of an address into the sockets to probe, the network probe and the verdict.
package diagnose

import "strings"

// Masked replaces a secret.
const Masked = "***"

// exactSecretKeys are the keys whose values are secret; secretKeySuffixes are the ends of such keys. Both are lower
// case and compared with the last dot-separated segment of a key.
var (
	exactSecretKeys    = []string{"login", "user", "username", "auth", "authorization"}
	secretKeySuffixes  = []string{"password", "passwd", "pwd", "token", "secret", "credentials", "apikey"}
	userInfoBoundaries = "()=,[]\"'"
)

// Redact masks the secrets of a text: the value of a key=value pair whose key names a login, a password, a token or
// another secret, and the user:password part of an address, such as "tls+me:pass@host:7300", which becomes
// "tls+***@host:7300".
func Redact(text string) string {
	out := make([]byte, 0, len(text))
	for i := 0; i < len(text); {
		switch text[i] {
		case '=':
			if isSecretKey(keyBefore(text, i)) {
				out = append(out, '=')
				out = append(out, Masked...)
				i = valueEnd(text, i+1)
				continue
			}
		case '@':
			if start, ok := userInfoStart(out); ok {
				last := lastAtBeforeBoundary(text, i)
				out = append(out[:start], Masked...)
				out = append(out, '@')
				i = last + 1
				continue
			}
		}
		out = append(out, text[i])
		i++
	}
	return string(out)
}

func isKeyChar(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '.' || c == '_' || c == '-'
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f' || c == '\v'
}

// keyBefore returns the key that ends at the "=" at index i.
func keyBefore(text string, i int) string {
	k := i
	for k > 0 && isKeyChar(text[k-1]) {
		k--
	}
	return text[k:i]
}

func isSecretKey(key string) bool {
	if key == "" {
		return false
	}
	key = strings.ToLower(key[strings.LastIndexByte(key, '.')+1:])
	for _, exact := range exactSecretKeys {
		if key == exact {
			return true
		}
	}
	for _, suffix := range secretKeySuffixes {
		if strings.HasSuffix(key, suffix) {
			return true
		}
	}
	return false
}

// valueEnd returns the index after the value that starts at index j: a quoted value ends after its closing quote,
// another one before a comma, a closing bracket, a quote or a space.
func valueEnd(text string, j int) int {
	if j < len(text) && (text[j] == '"' || text[j] == '\'') {
		if end := strings.IndexByte(text[j+1:], text[j]); end >= 0 {
			return j + 1 + end + 1
		}
		return len(text)
	}
	for j < len(text) && !strings.ContainsRune(",]\"'", rune(text[j])) && !isSpace(text[j]) {
		j++
	}
	return j
}

func isUserInfoBoundary(c byte) bool {
	return strings.IndexByte(userInfoBoundaries, c) >= 0 || isSpace(c)
}

// userInfoStart returns where the user information of an address starts in the output written before an "@": after
// a scheme ("//"), a codec ("+") or a path ("/") within the fragment that follows the last boundary character. There
// is user information only when it has a colon, that is a password.
func userInfoStart(out []byte) (int, bool) {
	f := len(out)
	for f > 0 && !isUserInfoBoundary(out[f-1]) {
		f--
	}
	fragment := string(out[f:])
	start := 0
	if slashes := strings.LastIndex(fragment, "//"); slashes >= 0 {
		start = slashes + 2
	}
	colon := strings.IndexByte(fragment[start:], ':')
	if colon < 0 {
		return 0, false
	}
	if cut := strings.LastIndexAny(fragment[start:start+colon], "+/"); cut >= 0 {
		start += cut + 1
	}
	return f + start, true
}

// lastAtBeforeBoundary returns the index of the last "@" from index i up to the next boundary character: a password
// may contain "@".
func lastAtBeforeBoundary(text string, i int) int {
	last := i
	for j := i; j < len(text) && !isUserInfoBoundary(text[j]); j++ {
		if text[j] == '@' {
			last = j
		}
	}
	return last
}
