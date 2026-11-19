// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package profile_test

import (
	"fmt"
	"testing"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/profile"
)

// IsShortSaleRestricted depends on the short sale restriction only, IsTradingHalted on the trading status only.
func TestShortSaleRestrictedAndTradingHalted(t *testing.T) {
	for _, ssr := range []profile.ShortSaleRestriction{
		profile.UndefinedShortSaleRestriction, profile.ActiveShortSaleRestriction, profile.InactiveShortSaleRestriction,
	} {
		for _, status := range []profile.TradingStatus{
			profile.UndefinedTradingStatus, profile.HaltedTradingStatus, profile.ActiveTradingStatus,
		} {
			t.Run(fmt.Sprintf("ssr=%d/status=%d", ssr, status), func(t *testing.T) {
				p := profile.NewProfile("AAPL")
				p.SetShortSaleRestriction(ssr)
				p.SetTradingStatus(status)
				if want := ssr == profile.ActiveShortSaleRestriction; p.IsShortSaleRestricted() != want {
					t.Errorf("IsShortSaleRestricted = %v, want %v", p.IsShortSaleRestricted(), want)
				}
				if want := status == profile.HaltedTradingStatus; p.IsTradingHalted() != want {
					t.Errorf("IsTradingHalted = %v, want %v", p.IsTradingHalted(), want)
				}
			})
		}
	}
}
