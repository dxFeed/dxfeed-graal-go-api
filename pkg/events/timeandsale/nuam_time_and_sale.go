package timeandsale

import (
	"github.com/dxfeed/dxfeed-graal-go-api/internal/formatutil"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

// NuamTimeAndSale is a TimeAndSale of the Nuam Exchange with the match ID of the trade.
type NuamTimeAndSale struct {
	TimeAndSale
	matchId int64
}

// NewNuamTimeAndSale creates a NuamTimeAndSale for the given event symbol.
func NewNuamTimeAndSale(eventSymbol string) *NuamTimeAndSale {
	return &NuamTimeAndSale{TimeAndSale: *NewTimeAndSale(eventSymbol)}
}

// Type returns the event code of this event.
func (t *NuamTimeAndSale) Type() eventcodes.EventCode {
	return eventcodes.NuamTimeAndSale
}

// MatchId returns the match ID of the trade.
func (t *NuamTimeAndSale) MatchId() int64 {
	return t.matchId
}

// SetMatchId sets the match ID of the trade.
func (t *NuamTimeAndSale) SetMatchId(matchId int64) {
	t.matchId = matchId
}

// String returns a string representation of this event.
func (t *NuamTimeAndSale) String() string {
	return "NuamTimeAndSale{" + t.fieldsString() +
		", matchId=" + formatutil.FormatInt64(t.matchId) +
		"}"
}
