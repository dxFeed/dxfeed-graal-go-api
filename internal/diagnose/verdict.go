// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package diagnose

// Kind is the kind of a verdict; a greater kind is a worse one.
type Kind int

// The kinds of a verdict, from the best to the worst.
const (
	Data Kind = iota
	PartialData
	ConnectedNoData
	CannotConnect
	CannotResolve
)

// String returns the name of the kind as the report prints it.
func (k Kind) String() string {
	switch k {
	case Data:
		return "DATA"
	case PartialData:
		return "PARTIAL DATA"
	case ConnectedNoData:
		return "CONNECTED, NO DATA"
	case CannotConnect:
		return "CANNOT CONNECT"
	default:
		return "CANNOT RESOLVE"
	}
}

// ExitCode returns the exit status of the tool for the kind: 0 for data, 2 to 5 for the others.
func (k Kind) ExitCode() int {
	if k == Data {
		return 0
	}
	return int(k) + 1
}

// Wildcard is the symbol that stands for every symbol.
const Wildcard = "*"

// Pair is an event type and a symbol.
type Pair struct{ Type, Symbol string }

// Input is what a verdict is computed from.
type Input struct {
	// Types and Symbols are the requested ones; the symbol Wildcard matches any symbol.
	Types, Symbols []string
	// Received are the type and symbol pairs of the received events.
	Received map[Pair]bool
	// WasConnected is true if the endpoint was connected at some time.
	WasConnected bool
	// LastError is the text of the last connection error of the endpoint, LastErrorIsDNS is true when it is a host
	// name that does not resolve; HasLastError is false when there is none.
	LastError      string
	LastErrorIsDNS bool
	HasLastError   bool
	// NoHostResolved is true when the probe had hosts and none of them resolved.
	NoHostResolved bool
}

// Verdict is the result of a diagnosis.
type Verdict struct {
	Kind Kind
	Hint string
	// Missing are the requested type and symbol pairs without events, as "<type> <symbol>", for PartialData only.
	Missing []string
}

// The hints of the verdicts.
const (
	HintPartial = "No events for some symbols: check the spelling of the symbols and the entitlements for them; a " +
		"symbol that does not trade now is silent too."
	HintNoData = "The endpoint connected but no events came: likely not entitled or not in the endpoint filter; " +
		"check the endpoint IPF. A symbol that does not trade now is silent too."
	HintResolve   = "The host name does not resolve: check the address and the DNS of this machine."
	HintNoError   = "The endpoint did not connect and reported no connection error."
	lastErrorHint = "Last connection error: "
)

// ComputeVerdict computes the verdict of a diagnosis.
func ComputeVerdict(in Input) Verdict {
	present := 0
	var missing []string
	for _, eventType := range in.Types {
		for _, symbol := range in.Symbols {
			if hasData(in.Received, eventType, symbol) {
				present++
			} else {
				missing = append(missing, eventType+" "+symbol)
			}
		}
	}
	switch {
	case present > 0 && len(missing) == 0:
		return Verdict{Kind: Data}
	case present > 0:
		return Verdict{Kind: PartialData, Hint: HintPartial, Missing: missing}
	case in.WasConnected:
		hint := HintNoData
		if in.HasLastError {
			hint += " " + lastErrorHint + in.LastError
		}
		return Verdict{Kind: ConnectedNoData, Hint: hint}
	case in.HasLastError && in.LastErrorIsDNS:
		return Verdict{Kind: CannotResolve, Hint: HintResolve + " Last error: " + in.LastError}
	case !in.HasLastError && in.NoHostResolved:
		return Verdict{Kind: CannotResolve, Hint: HintResolve}
	case in.HasLastError:
		return Verdict{Kind: CannotConnect, Hint: lastErrorHint + in.LastError}
	default:
		return Verdict{Kind: CannotConnect, Hint: HintNoError}
	}
}

func hasData(received map[Pair]bool, eventType, symbol string) bool {
	if symbol != Wildcard {
		return received[Pair{eventType, symbol}]
	}
	for pair := range received {
		if pair.Type == eventType {
			return true
		}
	}
	return false
}
