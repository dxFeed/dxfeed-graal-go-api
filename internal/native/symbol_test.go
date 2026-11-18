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
	invalid := osub.NewTimeSeriesSubscriptionSymbol(3.5, 0)
	if _, err := newCSymbolList([]osub.Symbol{osub.String("AAPL"), osub.NewWildcardSymbol(), invalid}); err == nil {
		t.Fatal("a list with a float inner symbol built without an error")
	}
	list, err := newCSymbolList([]osub.Symbol{osub.String("AAPL"), osub.NewTimeSeriesSubscriptionSymbol(osub.String("IBM"), 5)})
	if err != nil {
		t.Fatal(err)
	}
	freeCSymbolList(list)
}

func TestNewCSymbolOfString(t *testing.T) {
	for _, symbol := range []any{"AAPL{=d}", osub.String("AAPL{=d}")} {
		cSymbol, err := newCSymbol(symbol)
		if err != nil {
			t.Fatalf("%T: %v", symbol, err)
		}
		if got, ok := stringSymbolText(cSymbol); !ok || got != "AAPL{=d}" {
			t.Errorf("%T: got %q, string symbol %v", symbol, got, ok)
		}
		freeCSymbol(cSymbol)
	}
}
