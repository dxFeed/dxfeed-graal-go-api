// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package osub_test

import (
	"fmt"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
)

// Strings turns plain strings into symbols for one AddSymbols call; a symbol of another kind goes into the same list.
func ExampleStrings() {
	symbols := append(osub.Strings("AAPL", "IBM"), osub.NewWildcardSymbol())
	for _, symbol := range symbols {
		fmt.Printf("%T\n", symbol)
	}
	// Output:
	// osub.String
	// osub.String
	// *osub.WildcardSymbol
}
