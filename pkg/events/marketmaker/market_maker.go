// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

// Package marketmaker provides the MarketMaker event: the top bid and ask of one market participant.
package marketmaker

import (
	"fmt"
	"math"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/formatutil"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

// The index holds the exchange code in bits 32 to 47 and the market maker ID (MMID) in bits 0 to 31.
const (
	exchangeMask  = 0xFFFF
	exchangeShift = 32
	mmidMask      = 0xFFFFFFFF
)

// MarketMaker is the top bid and ask of one market participant (a market maker) on one exchange. It is an indexed
// event with the default source: the index holds the exchange code and the market maker ID, see SetExchangeCode and
// SetMarketMaker.
type MarketMaker struct {
	eventSymbol *string
	eventTime   int64
	eventFlags  int32
	index       int64
	bidTime     int64
	bidPrice    float64
	bidSize     float64
	bidCount    int64
	askTime     int64
	askPrice    float64
	askSize     float64
	askCount    int64
}

// NewMarketMaker creates a MarketMaker for the given event symbol.
func NewMarketMaker(eventSymbol string) *MarketMaker {
	return &MarketMaker{
		eventSymbol: &eventSymbol,
		bidPrice:    math.NaN(),
		bidSize:     math.NaN(),
		askPrice:    math.NaN(),
		askSize:     math.NaN(),
	}
}

// Type returns the event code of this event type.
func (m *MarketMaker) Type() eventcodes.EventCode {
	return eventcodes.MarketMaker
}

// Source is always DEFAULT for MarketMaker.
func (m *MarketMaker) Source() *events.IndexedEventSource {
	return events.DefaultIndexedEventSource()
}

// EventSymbol returns the symbol of this event.
func (m *MarketMaker) EventSymbol() *string {
	return m.eventSymbol
}

// SetEventSymbol sets the symbol of this event.
func (m *MarketMaker) SetEventSymbol(eventSymbol string) {
	m.eventSymbol = &eventSymbol
}

// EventTime returns the time, in milliseconds, when the event was received or last modified by the source.
func (m *MarketMaker) EventTime() int64 {
	return m.eventTime
}

// SetEventTime sets the event time in milliseconds.
func (m *MarketMaker) SetEventTime(eventTime int64) {
	m.eventTime = eventTime
}

// EventFlags returns the transactional event flags; see events.EventFlag.
func (m *MarketMaker) EventFlags() int32 {
	return m.eventFlags
}

// SetEventFlags sets the transactional event flags.
func (m *MarketMaker) SetEventFlags(eventFlags int32) {
	m.eventFlags = eventFlags
}

// Index returns the unique per-symbol index of this event: the exchange code and the market maker ID.
func (m *MarketMaker) Index() int64 {
	return m.index
}

// SetIndex sets the index of this event, which changes both the exchange code and the market maker ID; it returns an
// error if the index is negative.
func (m *MarketMaker) SetIndex(index int64) error {
	if index < 0 {
		return fmt.Errorf("negative index: %d", index)
	}
	m.index = index
	return nil
}

// Time returns the time of the last bid or ask change in milliseconds: the later of BidTime and AskTime.
func (m *MarketMaker) Time() int64 {
	return max(m.bidTime, m.askTime)
}

// ExchangeCode returns the exchange code of this event, a UTF-16 code unit.
func (m *MarketMaker) ExchangeCode() int16 {
	return int16(m.index >> exchangeShift & exchangeMask)
}

// SetExchangeCode sets the exchange code of this event; the market maker ID is kept.
func (m *MarketMaker) SetExchangeCode(exchangeCode int16) {
	m.index = m.index&^(exchangeMask<<exchangeShift) | int64(uint16(exchangeCode))<<exchangeShift
}

// MarketMaker returns the market maker ID (MMID) of this event, or nil when it is not set.
func (m *MarketMaker) MarketMaker() *string {
	return decodeShortString(m.index & mmidMask)
}

// SetMarketMaker sets the market maker ID (MMID) of this event; the exchange code is kept. The ID has up to 8
// characters from U+0001 to U+00FF, and only the last 4 of them are stored, as in the Java API. Nil or an empty string
// clears the ID. It returns an error for a longer string or an illegal character.
func (m *MarketMaker) SetMarketMaker(marketMaker *string) error {
	var code int64
	if marketMaker != nil {
		var err error
		if code, err = encodeShortString(*marketMaker); err != nil {
			return err
		}
	}
	m.index = m.index&^mmidMask | code&mmidMask
	return nil
}

// BidTime returns the time of the last bid change in milliseconds since the Unix epoch.
func (m *MarketMaker) BidTime() int64 {
	return m.bidTime
}

// SetBidTime sets the time of the last bid change in milliseconds since the Unix epoch.
func (m *MarketMaker) SetBidTime(bidTime int64) {
	m.bidTime = bidTime
}

// BidPrice returns the bid price.
func (m *MarketMaker) BidPrice() float64 {
	return m.bidPrice
}

// SetBidPrice sets the bid price.
func (m *MarketMaker) SetBidPrice(bidPrice float64) {
	m.bidPrice = bidPrice
}

// BidSize returns the bid size.
func (m *MarketMaker) BidSize() float64 {
	return m.bidSize
}

// SetBidSize sets the bid size.
func (m *MarketMaker) SetBidSize(bidSize float64) {
	m.bidSize = bidSize
}

// BidCount returns the number of the individual bid orders.
func (m *MarketMaker) BidCount() int64 {
	return m.bidCount
}

// SetBidCount sets the number of the individual bid orders.
func (m *MarketMaker) SetBidCount(bidCount int64) {
	m.bidCount = bidCount
}

// AskTime returns the time of the last ask change in milliseconds since the Unix epoch.
func (m *MarketMaker) AskTime() int64 {
	return m.askTime
}

// SetAskTime sets the time of the last ask change in milliseconds since the Unix epoch.
func (m *MarketMaker) SetAskTime(askTime int64) {
	m.askTime = askTime
}

// AskPrice returns the ask price.
func (m *MarketMaker) AskPrice() float64 {
	return m.askPrice
}

// SetAskPrice sets the ask price.
func (m *MarketMaker) SetAskPrice(askPrice float64) {
	m.askPrice = askPrice
}

// AskSize returns the ask size.
func (m *MarketMaker) AskSize() float64 {
	return m.askSize
}

// SetAskSize sets the ask size.
func (m *MarketMaker) SetAskSize(askSize float64) {
	m.askSize = askSize
}

// AskCount returns the number of the individual ask orders.
func (m *MarketMaker) AskCount() int64 {
	return m.askCount
}

// SetAskCount sets the number of the individual ask orders.
func (m *MarketMaker) SetAskCount(askCount int64) {
	m.askCount = askCount
}

// String returns a text representation of the event.
func (m *MarketMaker) String() string {
	return "MarketMaker{" + formatutil.FormatString(m.EventSymbol()) +
		", eventTime=" + formatutil.FormatTime(m.EventTime()) +
		", eventFlags=" + formatutil.HexFormat(int64(m.EventFlags())) +
		", index=" + formatutil.HexFormat(m.Index()) +
		", time=" + formatutil.FormatTime(m.Time()) +
		", exchange=" + formatutil.FormatChar(rune(uint16(m.ExchangeCode()))) +
		", marketMaker=" + formatutil.FormatString(m.MarketMaker()) +
		", bidTime=" + formatutil.FormatTime(m.bidTime) +
		", bidPrice=" + formatutil.FormatFloat64(m.bidPrice) +
		", bidSize=" + formatutil.FormatFloat64(m.bidSize) +
		", bidCount=" + formatutil.FormatInt64(m.bidCount) +
		", askTime=" + formatutil.FormatTime(m.askTime) +
		", askPrice=" + formatutil.FormatFloat64(m.askPrice) +
		", askSize=" + formatutil.FormatFloat64(m.askSize) +
		", askCount=" + formatutil.FormatInt64(m.askCount) +
		"}"
}

// encodeShortString packs up to 8 characters up to U+00FF, one byte each, the last character in the lowest byte; zero
// characters are skipped. It is the ShortString encoding of the Java API.
func encodeShortString(s string) (int64, error) {
	runes := []rune(s)
	if len(runes) > 8 {
		return 0, fmt.Errorf("invalid market maker: the string is too long: %q", s)
	}
	var code uint64
	for _, r := range runes {
		if r > 0xFF {
			return 0, fmt.Errorf("invalid market maker: the string contains an illegal character: %q", s)
		}
		if r != 0 {
			code = code<<8 | uint64(r)
		}
	}
	return int64(code), nil
}

// decodeShortString is the reverse of encodeShortString: a zero code is nil, zero bytes are skipped.
func decodeShortString(code int64) *string {
	if code == 0 {
		return nil
	}
	var runes []rune
	for shift := 56; shift >= 0; shift -= 8 {
		if c := rune(uint64(code) >> shift & 0xFF); c != 0 {
			runes = append(runes, c)
		}
	}
	s := string(runes)
	return &s
}
