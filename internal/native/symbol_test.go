package native

import (
	"strings"
	"testing"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
)

// stubSource is an indexed event source with any type and name.
type stubSource struct {
	sourceType events.EventSourceType
	name       *string
}

func (s stubSource) Name() *string                { return s.name }
func (s stubSource) Id() int64                    { return 7 }
func (s stubSource) Type() events.EventSourceType { return s.sourceType }

func TestNewCSymbolRejectsInvalidSources(t *testing.T) {
	name := "STUB"
	cases := map[string]struct {
		source events.IndexedEventSourceInterface
		want   string
	}{
		"unknown type": {stubSource{sourceType: 42, name: &name}, "unsupported indexed event source type 42"},
		"no name":      {stubSource{sourceType: events.IndexedEventSourceType}, "without a name"},
		"nil source":   {nil, "without a source"},
	}
	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			symbol, err := newCSymbol(osub.NewIndexedEventSubscriptionSymbol("AAPL", c.source))
			if err == nil {
				freeCSymbol(symbol)
				t.Fatal("no error")
			}
			if !strings.Contains(err.Error(), c.want) {
				t.Fatalf("error %q does not contain %q", err, c.want)
			}
		})
	}
}

func TestNewCSymbolListFreesOnError(t *testing.T) {
	if _, err := newCSymbolList([]any{"AAPL", osub.NewWildcardSymbol(), 3.5}); err == nil {
		t.Fatal("a list with a float symbol built without an error")
	}
	list, err := newCSymbolList([]any{"AAPL", osub.NewTimeSeriesSubscriptionSymbol("IBM", 5)})
	if err != nil {
		t.Fatal(err)
	}
	freeCSymbolList(list)
}
