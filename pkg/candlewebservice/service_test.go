package candlewebservice

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/auth"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/candle"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

// candleTape is a response of the service in the text format without compression: one daily candle of AAPL.
const candleTape = "==DXP3\ttype=tape\n" +
	"=Candle\tEventSymbol\tTime\tSequence\tCount\tOpen\tHigh\tLow\tClose\tVolume\n" +
	"Candle\tAAPL{=d}\t20260101-000000-0000\t0\t3\t40\t43\t39.5\t42.5\t1000\n"

// fakeService serves candleTape and records the requests.
type fakeService struct {
	mu       sync.Mutex
	requests []*http.Request
}

func (s *fakeService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	s.requests = append(s.requests, r)
	s.mu.Unlock()
	_, _ = fmt.Fprint(w, candleTape)
}

func (s *fakeService) last(t *testing.T) *http.Request {
	t.Helper()
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.requests) == 0 {
		t.Fatal("the service got no request")
	}
	return s.requests[len(s.requests)-1]
}

// TestGetTimeSeriesFromService downloads candles from a local fake of the service and checks the request and the
// parsed events.
func TestGetTimeSeriesFromService(t *testing.T) {
	service := &fakeService{}
	server := httptest.NewServer(service)
	defer server.Close()
	from := time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC)
	to := from.AddDate(0, 2, 0)

	cases := []struct {
		name          string
		builder       *HistoryEndpointBuilder
		authorization string
	}{
		{"basic", NewHistoryEndpointBuilder().WithUserName("user").WithPassword("secret"), "Basic dXNlcjpzZWNyZXQ="},
		{"bearer", NewHistoryEndpointBuilder().WithAuthToken("token"), "Bearer token"},
		{"basic token", NewHistoryEndpointBuilder().WithToken(mustToken(auth.NewBasicUser("user", "secret"))),
			"Basic dXNlcjpzZWNyZXQ="},
		{"bearer token", NewHistoryEndpointBuilder().WithToken(mustToken(auth.NewBearer("token"))), "Bearer token"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			endpoint, err := c.builder.WithAddress(server.URL + "/candledata").WithCompression(CompressionNone).
				WithFormat(FormatText).Build()
			if err != nil {
				t.Fatalf("build: %v", err)
			}
			defer func() { _ = endpoint.Close() }()
			events, err := endpoint.GetTimeSeries(eventcodes.Candle, osub.String("AAPL{=d}"), from, to)
			if err != nil {
				t.Fatalf("get time series: %v", err)
			}

			request := service.last(t)
			query := request.URL.Query()
			for key, want := range map[string]string{"records": "Candle", "symbols": "AAPL{=d}", "format": "text",
				"compression": "none"} {
				if got := query.Get(key); got != want {
					t.Errorf("query %s=%q, want %q (request %s)", key, got, want, request.URL)
				}
			}
			if got := request.Header.Get("Authorization"); got != c.authorization {
				t.Errorf("Authorization %q, want %q", got, c.authorization)
			}

			if len(events) != 1 {
				t.Fatalf("got %d events, want 1: %v", len(events), events)
			}
			got, ok := events[0].(*candle.Candle)
			if !ok {
				t.Fatalf("got %T, want *candle.Candle", events[0])
			}
			want := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()
			if got.EventSymbol().String() != "AAPL{=d}" || got.Time() != want || got.Count() != 3 ||
				got.Open() != 40 || got.High() != 43 || got.Low() != 39.5 || got.Close() != 42.5 ||
				got.Volume() != 1000 {
				t.Fatalf("candle %v differs from the served one", got)
			}
		})
	}
}

func mustToken(token *auth.Token, err error) *auth.Token {
	if err != nil {
		panic(err)
	}
	return token
}
