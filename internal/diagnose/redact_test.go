// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package diagnose

import "testing"

func TestRedact(t *testing.T) {
	for text, want := range map[string]string{
		"host:7300[login=abc,password=s3cr3t]":       "host:7300[login=***,password=***]",
		"host:7300[Password=abc]":                    "host:7300[Password=***]",
		"x[dxfeed.password=abc,dxfeed.user=me]":      "x[dxfeed.password=***,dxfeed.user=***]",
		"token=abc def":                              "token=*** def",
		`password="a b" next`:                        "password=*** next",
		"password='a,b'":                             "password=***",
		"host:7300[password=a;b)c,login=d]":          "host:7300[password=***,login=***]",
		"Quote{AAPL, bid=1.5@Q}":                     "Quote{AAPL, bid=1.5@Q}",
		"keyStore=a.jks":                             "keyStore=a.jks",
		"trustStorePassword=x1 authToken=y apiKey=z": "trustStorePassword=*** authToken=*** apiKey=***",
		"clientSecret=s keyStorePwd=p":               "clientSecret=*** keyStorePwd=***",
		"tls+me:s3cr3t@host:7300":                    "tls+***@host:7300",
		"me:p@ss@host:7300":                          "***@host:7300",
		"dxlink:wss://me:p:ss@host/dxlink-ws":        "dxlink:wss://***@host/dxlink-ws",
		"ssl[trustStore=a.jks]+me:x@host:7300":       "ssl[trustStore=a.jks]+***@host:7300",
		"(me:x@a:1,b:2)":                             "(***@a:1,b:2)",
		"connect to me:p/+@x@host:7300 failed":       "connect to ***@host:7300 failed",
		"user@host:7300":                             "user@host:7300",
		"no secrets here":                            "no secrets here",
		"=x":                                         "=x",
		"password=":                                  "password=***",
	} {
		if got := Redact(text); got != want {
			t.Errorf("Redact(%q) = %q, want %q", text, got, want)
		}
	}
}
