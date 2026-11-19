// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package main

// The help of the tools. A tool prints it when it gets fewer than three arguments. The arguments come first and the
// options after them, since the value of an option is not told apart from an argument.

const addressArgument = `  <address>  The address to connect to: host:port, a dxLink URL or a tape file, such as demo.dxfeed.com:7300.
             To pass an authorization token, add it to the address: demo.dxfeed.com:7300[login=entitle:<token>].
`

const allTypesArgument = `  <types>    A comma-separated list of event types, case-insensitive, such as Quote,Trade,TimeAndSale,Candle,Order;
             all is every event type.
`

const symbolsArgument = `  <symbols>  A comma-separated list of symbols, such as IBM,AAPL,MSFT; a candle symbol has its aggregation, such as
             AAPL{=d}; all is the wildcard symbol.
`

const connectHelp = `Connect: connects to the address and prints the events of the symbols.

Usage:
  tools connect <address> <types> <symbols> [<options>]

Where:
` + addressArgument + allTypesArgument + symbolsArgument + `
Options:
  -f, --from-time <time>  Subscribe to the time series from the time, such as 20230901Z or 2023-09-01T10:00:00Z.
  -p, --properties <name=value,...>
                          Set the properties of the endpoint, such as dxfeed.experimental.dxlink.enable=true.
  -q, --quite             Do not print the events.
  --force-stream          Enforce a streaming contract for the subscription: the StreamFeed role instead of Feed.

Samples:
  tools connect demo.dxfeed.com:7300 Quote AAPL
  tools connect demo.dxfeed.com:7300 Candle "AAPL{=d}" -f 20230901Z
  tools connect dxlink:wss://demo.dxfeed.com/market-data/dxlink-ws Quote AAPL -p dxfeed.experimental.dxlink.enable=true
`

const dumpHelp = `Dump: prints the events received from the address and, with -t, writes them to a tape file. It uses the StreamFeed
role, so the events of a tape file come as they are, without conflation.

Usage:
  tools dump <address> <types> <symbols> [<options>]

Where:
` + addressArgument + allTypesArgument + symbolsArgument + `
Options:
  -t, --tape <file>       Write the events to the tape file, such as out.txt[format=text].
  -p, --properties <name=value,...>
                          Set the properties of the endpoints, such as dxfeed.wildcard.enable=true for the symbol all.
  -q, --quite             Do not print the events.

Samples:
  tools dump demo.dxfeed.com:7300 Quote AAPL,IBM -t out.txt[format=text] -q
  tools dump in.tape[speed=max] Quote,Profile,TimeAndSale all -t out.tape -q -p dxfeed.wildcard.enable=true
`

const perfTestHelp = `PerfTest: connects to the address and prints every 2 seconds the rate of events, the rate of listener calls and the
number of events in a call.

Usage:
  tools perftest <address> <types> <symbols> [<options>]

Where:
` + addressArgument + `  <types>    Quote, TimeAndSale or Quote,TimeAndSale, case-insensitive; the Quote events are also printed.
` + symbolsArgument + `
Options:
  --force-stream          Enforce a streaming contract for the subscription: the StreamFeed role instead of Feed.

Sample:
  tools perftest demo.dxfeed.com:7300 TimeAndSale AAPL,IBM,MSFT
`

const latencyTestHelp = `LatencyTest: connects to the address, subscribes to TimeAndSale events and calculates their latency: the delay
between the time of an event and its receipt, in milliseconds (min, max, mean, 99th percentile, etc).

Usage:
  tools latencytest <address> <types> <symbols> [<options>]

Where:
` + addressArgument + `  <types>    TimeAndSale, case-insensitive; only these events are measured.
` + symbolsArgument + `
Options:
  --ignore-exchanges <codes>
                          Leave out the events of the exchange codes, a comma-separated list without spaces, such
                          as Q,N.
  --force-stream          Enforce a streaming contract for the subscription: the StreamFeed role instead of Feed.

Sample:
  tools latencytest demo.dxfeed.com:7300 TimeAndSale AAPL,IBM,MSFT
`
