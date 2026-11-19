// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

// Package candlewebservice downloads the history of candles and other time series events from the dxFeed candle web
// service with HistoryEndpoint, without a feed connection:
//
//	endpoint, err := candlewebservice.NewHistoryEndpointBuilder().
//		WithAddress("https://tools.dxfeed.com/candledata").
//		WithUserName(user).
//		WithPassword(password).
//		Build()
//	...
//	events, err := endpoint.GetTimeSeries(eventcodes.Candle, osub.Candle(candle.NewCandleSymbol("AAPL{=d}")), from, to)
package candlewebservice

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/auth"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

// Compression is the compression of the data that the service sends.
type Compression int32

const (
	// CompressionNone requests the data without compression.
	CompressionNone Compression = iota
	// CompressionZip requests the data in a zip archive; the default.
	CompressionZip
	// CompressionGzip requests the data compressed with gzip.
	CompressionGzip
)

// String returns the name of the compression, such as "Zip".
func (c Compression) String() string {
	switch c {
	case CompressionNone:
		return "None"
	case CompressionZip:
		return "Zip"
	case CompressionGzip:
		return "Gzip"
	default:
		return fmt.Sprintf("Compression(%d)", int32(c))
	}
}

// Format is the format of the data that the service sends.
type Format int32

const (
	// FormatText requests the data in the text format of the dxFeed tape files.
	FormatText Format = iota
	// FormatCsv requests the data as comma-separated values.
	FormatCsv
	// FormatBinary requests the data in the binary format of the dxFeed tape files; the default.
	FormatBinary
)

// String returns the name of the format, such as "Binary".
func (f Format) String() string {
	switch f {
	case FormatText:
		return "Text"
	case FormatCsv:
		return "Csv"
	case FormatBinary:
		return "Binary"
	default:
		return fmt.Sprintf("Format(%d)", int32(f))
	}
}

// HistoryEndpointBuilder collects the settings of a HistoryEndpoint; create it with NewHistoryEndpointBuilder. Each
// With method changes the builder and returns it, so the calls chain. A builder is not safe for concurrent use.
type HistoryEndpointBuilder struct {
	address, userName, password, authToken *string
	compression                            Compression
	format                                 Format
}

// NewHistoryEndpointBuilder returns a builder with the zip compression and the binary format; the address, the user
// name, the password and the authentication token are not set.
func NewHistoryEndpointBuilder() *HistoryEndpointBuilder {
	return &HistoryEndpointBuilder{compression: CompressionZip, format: FormatBinary}
}

// WithAddress sets the address of the service, such as "https://tools.dxfeed.com/candledata"; it is required.
func (b *HistoryEndpointBuilder) WithAddress(address string) *HistoryEndpointBuilder {
	b.address = &address
	return b
}

// WithUserName sets the user name of the basic authentication.
func (b *HistoryEndpointBuilder) WithUserName(userName string) *HistoryEndpointBuilder {
	b.userName = &userName
	return b
}

// WithPassword sets the password of the basic authentication.
func (b *HistoryEndpointBuilder) WithPassword(password string) *HistoryEndpointBuilder {
	b.password = &password
	return b
}

// WithAuthToken sets the authentication token, which the requests send as a bearer token.
func (b *HistoryEndpointBuilder) WithAuthToken(authToken string) *HistoryEndpointBuilder {
	b.authToken = &authToken
	return b
}

// WithToken sets the credentials from an auth.Token: the user name and the password of a Basic token, otherwise the
// value of the token, which the requests send as a bearer token. A nil token clears the authentication token.
func (b *HistoryEndpointBuilder) WithToken(token *auth.Token) *HistoryEndpointBuilder {
	if token == nil {
		b.authToken = nil
		return b
	}
	if user, basic := token.User(); basic {
		password, _ := token.Password()
		return b.WithUserName(user).WithPassword(password)
	}
	return b.WithAuthToken(token.Value())
}

// WithCompression sets the compression of the data.
func (b *HistoryEndpointBuilder) WithCompression(compression Compression) *HistoryEndpointBuilder {
	b.compression = compression
	return b
}

// WithFormat sets the format of the data.
func (b *HistoryEndpointBuilder) WithFormat(format Format) *HistoryEndpointBuilder {
	b.format = format
	return b
}

// Build creates an endpoint with the current settings; later changes of the builder do not affect it. It fails when
// the address is not set or empty, or the compression or the format is unknown. Close the endpoint when it is no
// longer needed.
func (b *HistoryEndpointBuilder) Build() (*HistoryEndpoint, error) {
	if b.address == nil || *b.address == "" {
		return nil, errors.New("candlewebservice: the address of the history endpoint is not set")
	}
	if b.compression < CompressionNone || b.compression > CompressionGzip {
		return nil, fmt.Errorf("candlewebservice: unknown compression %v", b.compression)
	}
	if b.format < FormatText || b.format > FormatBinary {
		return nil, fmt.Errorf("candlewebservice: unknown format %v", b.format)
	}
	handle, err := native.NewHistoryEndpointHandle(native.HistoryEndpointSettings{
		Address:     b.address,
		UserName:    b.userName,
		Password:    b.password,
		AuthToken:   b.authToken,
		Compression: int32(b.compression),
		Format:      int32(b.format),
	})
	if err != nil {
		return nil, err
	}
	return &HistoryEndpoint{handle: handle, address: *b.address, compression: b.compression, format: b.format}, nil
}

// HistoryEndpoint downloads time series events from the dxFeed candle web service; create it with
// NewHistoryEndpointBuilder. It may be used from several goroutines; Close waits for the running requests.
type HistoryEndpoint struct {
	// mu is held for reading by the requests and for writing by Close, so the handle is not released under a request.
	mu          sync.RWMutex
	handle      *native.HistoryEndpointHandle
	address     string
	compression Compression
	format      Format
}

// Address returns the address of the service.
func (e *HistoryEndpoint) Address() string { return e.address }

// Compression returns the compression of the data.
func (e *HistoryEndpoint) Compression() Compression { return e.compression }

// Format returns the format of the data.
func (e *HistoryEndpoint) Format() Format { return e.format }

// GetTimeSeries downloads the events of the type and the symbol within the time range from from to to, such as the
// candles of osub.Candle(candle.NewCandleSymbol("AAPL{=d}")) or the TimeAndSale events of osub.String("AAPL"). The
// type must be a time series event, such as Candle or TimeAndSale. It blocks while the request runs and returns the
// events in the order the service sends them, *candle.Candle values for Candle; an event of a type that the Go API
// does not support yet is skipped. The result may be empty.
func (e *HistoryEndpoint) GetTimeSeries(code eventcodes.EventCode, symbol osub.Symbol, from, to time.Time) ([]any,
	error) {
	if symbol == nil {
		return nil, errors.New("candlewebservice: the symbol is nil")
	}
	if to.Before(from) {
		return nil, fmt.Errorf("candlewebservice: the time range ends (%v) before it starts (%v)", to, from)
	}
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.handle.TimeSeries(code.NativeCode(), symbol, from.UnixMilli(), to.UnixMilli())
}

// Close waits for the running GetTimeSeries calls and releases the endpoint; later GetTimeSeries calls fail. Close may
// be called more than once.
func (e *HistoryEndpoint) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.handle.Free()
}
