// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package candle

import (
	"fmt"
	"strconv"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/eventsymbols"
)

// CandleSessionAttributeKey is the key of the session in a candle symbol string: "tho", which stands for "trading
// hours only".
const CandleSessionAttributeKey = "tho"

// CandleSession is the session attribute of a candle symbol: which trading sessions the candles are built from.
// The zero value is CandleSessionAny, the default.
type CandleSession uint8

const (
	// CandleSessionAny uses all trading sessions in the candle aggregation.
	CandleSessionAny CandleSession = iota
	// CandleSessionRegular uses only the regular trading session data in the candle aggregation.
	CandleSessionRegular

	// CandleSessionDefault is the default session: CandleSessionAny.
	CandleSessionDefault = CandleSessionAny
)

var candleSessions = [...]string{
	CandleSessionAny:     "false",
	CandleSessionRegular: "true",
}

// String returns the string form of the session: "false" for any session, "true" for the regular one.
func (s CandleSession) String() string {
	if int(s) >= len(candleSessions) {
		return "CandleSession(" + strconv.Itoa(int(s)) + ")"
	}
	return candleSessions[s]
}

// IsRegular reports whether only the regular trading session is used.
func (s CandleSession) IsRegular() bool {
	return s == CandleSessionRegular
}

// ChangeAttributeForSymbol returns the candle symbol string with this session set; the default one removes it.
func (s CandleSession) ChangeAttributeForSymbol(symbol string) string {
	if s == CandleSessionDefault {
		return eventsymbols.RemoveAttributeStringByKey(symbol, CandleSessionAttributeKey)
	}
	return eventsymbols.ChangeAttributeStringByKey(symbol, CandleSessionAttributeKey, s.String())
}

// ParseCandleSession parses the string form of a session: the exact string or its prefix in any case ("t", "FALSE").
// It is an error for an empty or an unknown string.
func ParseCandleSession(s string) (CandleSession, error) {
	if s == "" {
		return 0, fmt.Errorf("missing candle session")
	}
	for i, str := range candleSessions {
		if str == s {
			return CandleSession(i), nil
		}
	}
	for i, str := range candleSessions {
		if hasPrefixFoldASCII(str, s) {
			return CandleSession(i), nil
		}
	}
	return 0, fmt.Errorf("unknown candle session: %s", s)
}

// CandleSessionFromSymbol returns the session of the candle symbol string. As Java Boolean.parseBoolean does, only
// the value "true" in any case is the regular session, any other value or no attribute is the default one.
func CandleSessionFromSymbol(symbol string) CandleSession {
	value, ok := eventsymbols.GetAttributeStringByKey(symbol, CandleSessionAttributeKey)
	if ok && equalFoldASCII(value, "true") {
		return CandleSessionRegular
	}
	return CandleSessionDefault
}

// NormalizeCandleSessionAttribute returns the candle symbol string with the session in its canonical form: only
// "true" in any case stays (as "true"), any other value is the default one and is removed.
func NormalizeCandleSessionAttribute(symbol string) string {
	return normalizeAttribute(symbol, CandleSessionAttributeKey, CandleSessionDefault.String(),
		func(value string) (string, bool, error) {
			if !equalFoldASCII(value, "true") {
				return "", true, nil
			}
			return CandleSessionRegular.String(), false, nil
		})
}
