// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package candlewebservice

import (
	"net"
	"strings"
	"testing"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/auth"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

func TestBuilderDefaults(t *testing.T) {
	endpoint, err := NewHistoryEndpointBuilder().WithAddress("http://127.0.0.1:1/candledata").Build()
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	defer func() { _ = endpoint.Close() }()
	if endpoint.Address() != "http://127.0.0.1:1/candledata" || endpoint.Compression() != CompressionZip ||
		endpoint.Format() != FormatBinary {
		t.Fatalf("endpoint %q, %v, %v; want the address, Zip and Binary", endpoint.Address(), endpoint.Compression(),
			endpoint.Format())
	}
}

func TestBuilderValidation(t *testing.T) {
	cases := []struct {
		name    string
		builder *HistoryEndpointBuilder
		want    string
	}{
		{"no address", NewHistoryEndpointBuilder(), "address"},
		{"empty address", NewHistoryEndpointBuilder().WithAddress(""), "address"},
		{"unknown compression", NewHistoryEndpointBuilder().WithAddress("http://x").WithCompression(3),
			"Compression(3)"},
		{"negative compression", NewHistoryEndpointBuilder().WithAddress("http://x").WithCompression(-1),
			"Compression(-1)"},
		{"unknown format", NewHistoryEndpointBuilder().WithAddress("http://x").WithFormat(7), "Format(7)"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			endpoint, err := c.builder.Build()
			if err == nil {
				_ = endpoint.Close()
				t.Fatal("Build returned no error")
			}
			if !strings.Contains(err.Error(), c.want) {
				t.Fatalf("error %q does not mention %q", err, c.want)
			}
		})
	}
}

// TestBuilderSettingsAtBuild builds with every setting and checks that the endpoint keeps the settings of the
// moment of Build.
func TestBuilderSettingsAtBuild(t *testing.T) {
	token, err := auth.NewBearer("token")
	if err != nil {
		t.Fatalf("token: %v", err)
	}
	builder := NewHistoryEndpointBuilder().WithAddress("http://127.0.0.1:1/a").WithUserName("user").
		WithPassword("password").WithToken(token).WithCompression(CompressionGzip).WithFormat(FormatCsv)
	endpoint, err := builder.Build()
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	defer func() { _ = endpoint.Close() }()
	builder.WithAddress("http://127.0.0.1:1/b").WithCompression(CompressionNone).WithFormat(FormatText).WithToken(nil)
	if endpoint.Address() != "http://127.0.0.1:1/a" || endpoint.Compression() != CompressionGzip ||
		endpoint.Format() != FormatCsv {
		t.Fatalf("endpoint %q, %v, %v changed with the builder", endpoint.Address(), endpoint.Compression(),
			endpoint.Format())
	}
}

func TestEnumStrings(t *testing.T) {
	for value, want := range map[Compression]string{CompressionNone: "None", CompressionZip: "Zip",
		CompressionGzip: "Gzip", 9: "Compression(9)"} {
		if got := value.String(); got != want {
			t.Errorf("Compression %d: %q, want %q", int32(value), got, want)
		}
	}
	for value, want := range map[Format]string{FormatText: "Text", FormatCsv: "Csv", FormatBinary: "Binary",
		9: "Format(9)"} {
		if got := value.String(); got != want {
			t.Errorf("Format %d: %q, want %q", int32(value), got, want)
		}
	}
}

// TestGetTimeSeriesErrors checks the argument checks and that a request to a closed port fails instead of hanging.
func TestGetTimeSeriesErrors(t *testing.T) {
	// A port that was free a moment ago: nothing listens on it.
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	address := "http://" + listener.Addr().String() + "/candledata"
	_ = listener.Close()

	endpoint, err := NewHistoryEndpointBuilder().WithAddress(address).Build()
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	symbol := osub.String("AAPL{=d}")
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(24 * time.Hour)

	if _, err := endpoint.GetTimeSeries(eventcodes.Candle, nil, from, to); err == nil {
		t.Error("GetTimeSeries with a nil symbol returned no error")
	}
	if _, err := endpoint.GetTimeSeries(eventcodes.Candle, symbol, to, from); err == nil {
		t.Error("GetTimeSeries with a reversed range returned no error")
	}

	done := make(chan error, 1)
	go func() {
		_, err := endpoint.GetTimeSeries(eventcodes.Candle, symbol, from, to)
		done <- err
	}()
	select {
	case err := <-done:
		if err == nil {
			t.Error("GetTimeSeries from a closed port returned no error")
		}
	case <-time.After(30 * time.Second):
		t.Fatal("GetTimeSeries from a closed port did not return within 30 s")
	}

	if err := endpoint.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	if err := endpoint.Close(); err != nil {
		t.Fatalf("second close: %v", err)
	}
	if _, err := endpoint.GetTimeSeries(eventcodes.Candle, symbol, from, to); err == nil {
		t.Error("GetTimeSeries after Close returned no error")
	}
}
