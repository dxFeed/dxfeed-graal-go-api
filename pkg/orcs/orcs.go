// Package orcs is the client of the dxFeed ORCS (Order Reconstruction Service): it requests the historical price
// levels (orders) and quotes of a candle symbol from the price level service and checks the price levels:
//
//	service, err := orcs.NewPriceLevelService("orcs.example.com:7777[login=entitle:<token>]")
//	...
//	defer service.Close()
//	ntv, _ := order.ValueOfName("NTV")
//	orders, err := service.GetOrders(candle.NewCandleSymbol("AAPL{=d,gr=s}"), ntv, from, to)
package orcs

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/candle"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/order"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/quote"
)

// GranularityAttributeKey is the key of the candle symbol attribute of the granularity, such as "AAPL{=d,gr=s}": the
// period of the price levels that GetOrders and GetQuotes return, an integer number of seconds of at least 1 second.
const GranularityAttributeKey = "gr"

// PriceLevelService is a client of the price level service of ORCS; create it with NewPriceLevelService. It may be
// used from several goroutines.
type PriceLevelService struct {
	// mu is held for reading by the requests and for writing by Close, so the handle is not released under a request.
	mu     sync.RWMutex
	handle *native.PriceLevelServiceHandle
}

// NewPriceLevelService creates a client of the price level service at the address of its RMI endpoint, such as
// "orcs.example.com:7777". To pass an authorization token, add it to the address:
// "orcs.example.com:7777[login=entitle:<token>]". The client connects in the background, so an unreachable address
// is not an error here. Close it when it is no longer needed.
func NewPriceLevelService(address string) (*PriceLevelService, error) {
	handle, err := native.NewPriceLevelServiceHandle(address)
	if err != nil {
		return nil, err
	}
	return &PriceLevelService{handle: handle}, nil
}

// GetOrders returns the price levels (orders) of the candle symbol and the order source from from to to, sorted by
// time; see GetOrdersWithCaller. It passes the default caller identifier of the SDK.
func (s *PriceLevelService) GetOrders(symbol *candle.CandleSymbol, source *order.Source, from,
	to time.Time) ([]*order.Order, error) {
	return s.getOrders(symbol, source, from, to, nil)
}

// GetOrdersWithCaller returns the price levels (orders) of the candle symbol and the order source from from to to,
// sorted by time. The candle symbol shall have a period that the service supports and the granularity attribute
// (GranularityAttributeKey), such as "AAPL{=d,gr=s}"; otherwise the result is empty. The caller identifies the client
// to the service. It blocks while the request runs.
func (s *PriceLevelService) GetOrdersWithCaller(symbol *candle.CandleSymbol, source *order.Source, from,
	to time.Time, caller string) ([]*order.Order, error) {
	return s.getOrders(symbol, source, from, to, &caller)
}

func (s *PriceLevelService) getOrders(symbol *candle.CandleSymbol, source *order.Source, from, to time.Time,
	caller *string) ([]*order.Order, error) {
	if symbol == nil {
		return nil, errNilSymbol
	}
	if source == nil {
		return nil, errors.New("orcs: the order source is nil")
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	events, err := s.handle.Orders(symbol, source, from.UnixMilli(), to.UnixMilli(), caller)
	if err != nil {
		return nil, err
	}
	return typed[*order.Order](events)
}

// GetQuotes returns the quotes of the candle symbol from from to to, sorted by time; see GetQuotesWithCaller. It
// passes the default caller identifier of the SDK.
func (s *PriceLevelService) GetQuotes(symbol *candle.CandleSymbol, from, to time.Time) ([]*quote.Quote, error) {
	return s.getQuotes(symbol, from, to, nil)
}

// GetQuotesWithCaller returns the quotes of the candle symbol from from to to, sorted by time. The candle symbol shall
// have the granularity attribute (GranularityAttributeKey), such as "AAPL{gr=s}". The caller identifies the client to
// the service. It blocks while the request runs.
func (s *PriceLevelService) GetQuotesWithCaller(symbol *candle.CandleSymbol, from, to time.Time,
	caller string) ([]*quote.Quote, error) {
	return s.getQuotes(symbol, from, to, &caller)
}

func (s *PriceLevelService) getQuotes(symbol *candle.CandleSymbol, from, to time.Time,
	caller *string) ([]*quote.Quote, error) {
	if symbol == nil {
		return nil, errNilSymbol
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	events, err := s.handle.Quotes(symbol, from.UnixMilli(), to.UnixMilli(), caller)
	if err != nil {
		return nil, err
	}
	return typed[*quote.Quote](events)
}

// GetAuthOrderSource returns the order sources and their symbols that are available to the client, as of now.
func (s *PriceLevelService) GetAuthOrderSource() (*AuthOrderSource, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, err := s.handle.AuthOrderSource()
	if err != nil {
		return nil, err
	}
	return newAuthOrderSource(data), nil
}

// Close closes the connection to the service, waits for the running requests and releases the client; later requests
// fail. Closing the connection does not end a running request: Close returns when the request does, and a request to
// an unreachable service fails after the request timeout of the SDK (a minute by default). Close may be called more
// than once.
func (s *PriceLevelService) Close() error {
	s.mu.RLock()
	err := s.handle.Close()
	s.mu.RUnlock()
	s.mu.Lock()
	defer s.mu.Unlock()
	return errors.Join(err, s.handle.Free())
}

var errNilSymbol = errors.New("orcs: the candle symbol is nil")

// typed converts the events that the SDK returned to the expected type.
func typed[T any](events []interface{}) ([]T, error) {
	result := make([]T, 0, len(events))
	for _, e := range events {
		t, ok := e.(T)
		if !ok {
			var want T
			return nil, fmt.Errorf("orcs: got %T, want %T", e, want)
		}
		result = append(result, t)
	}
	return result, nil
}

// SourceSymbols are the symbols available to the client for an order source.
type SourceSymbols struct {
	// Source is the order source.
	Source *order.Source
	// Symbols are the symbols, sorted.
	Symbols []string
}

// AuthOrderSource holds the order sources and their symbols that are available to the client; get it with
// PriceLevelService.GetAuthOrderSource or create it with NewAuthOrderSource. It is a copy that does not change.
type AuthOrderSource struct {
	byIDs     map[int32][]string
	bySources []SourceSymbols
}

// NewAuthOrderSource creates an AuthOrderSource from the symbols by the order source identifier. Two order sources
// with the same identifier cannot be given, and a nil symbol list is an empty one.
func NewAuthOrderSource(symbolsByID map[int32][]string) (*AuthOrderSource, error) {
	data, err := native.NewAuthOrderSource(symbolsByID)
	if err != nil {
		return nil, err
	}
	return newAuthOrderSource(data), nil
}

func newAuthOrderSource(data native.AuthOrderSourceData) *AuthOrderSource {
	a := &AuthOrderSource{byIDs: data.ByIDs, bySources: make([]SourceSymbols, len(data.BySources))}
	for i, e := range data.BySources {
		a.bySources[i] = SourceSymbols{Source: e.Source, Symbols: e.Symbols}
	}
	return a
}

// ByIDs returns the sorted symbols by the identifier of the order source. The result is a copy.
func (a *AuthOrderSource) ByIDs() map[int32][]string {
	result := make(map[int32][]string, len(a.byIDs))
	for id, symbols := range a.byIDs {
		result[id] = append([]string(nil), symbols...)
	}
	return result
}

// ByOrderSources returns the order sources with their sorted symbols, ordered by the identifier of the order source.
// The result is a copy.
func (a *AuthOrderSource) ByOrderSources() []SourceSymbols {
	result := make([]SourceSymbols, len(a.bySources))
	for i, e := range a.bySources {
		result[i] = SourceSymbols{Source: e.Source, Symbols: append([]string(nil), e.Symbols...)}
	}
	return result
}

// ValidatePriceLevels checks the price levels (orders) that GetOrders returned, in the order of the list: at each
// moment of time the bid price shall be lower than the ask price, and no gap between the times of two subsequent
// orders shall be greater than timeGapBound. Note that in some cases it is normal for the bid to be greater than or
// equal to the ask. The details, and the quotes when printQuotes is set, go to the log of the SDK. An empty list is
// valid; a nil order is an error.
func ValidatePriceLevels(orders []*order.Order, timeGapBound time.Duration, printQuotes bool) (bool, error) {
	events := make([]interface{}, len(orders))
	for i, o := range orders {
		if o == nil {
			return false, fmt.Errorf("orcs: the order %d is nil", i)
		}
		events[i] = o
	}
	return native.ValidatePriceLevels(events, timeGapBound.Milliseconds(), printQuotes)
}
