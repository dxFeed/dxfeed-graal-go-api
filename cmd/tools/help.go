// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package main

// The help of the tools. A tool prints it when it gets no arguments or -h, --help or -help, and with an error of its
// arguments. The options may stand anywhere; an option with a value takes the next argument.

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
  -s, --source <source>   Subscribe to the indexed events of the order source, such as NTV; not with -f.
  -p, --properties <name=value,...>
                          Set the properties of the endpoints, also as system properties, such as
                          dxfeed.aggregationPeriod=1s.
  -t, --tape <file>       Write the received events to the tape file, such as out.txt[format=text].
  -q, --quiet             Do not print the events.
  --force-stream          Enforce a streaming contract for the subscription: the StreamFeed role instead of Feed.

Samples:
  tools connect demo.dxfeed.com:7300 Quote AAPL
  tools connect demo.dxfeed.com:7300 Candle "AAPL{=d}" -f 20230901Z
  tools connect demo.dxfeed.com:7300 Order AAPL -s NTV
  tools connect dxlink:wss://demo.dxfeed.com/market-data/dxlink-ws Quote AAPL
`

const diagnoseHelp = `Diagnose: checks the connection to the address step by step and prints a report with a verdict: the DNS
lookup and a TCP connect to each IP address of each host:port, the state changes and the connection errors of the
endpoint, the events of each type with the first ones, their flags and their lag, and the last event per symbol.
Logins, passwords and tokens are masked in the report.

Usage:
  tools diagnose <address> <types> <symbols> [<options>]

Where:
` + addressArgument + allTypesArgument + symbolsArgument + `
Options:
  --compare <address>     Diagnose this address too, with the same types, symbols and properties, and compare the
                          two.
  --duration <seconds>    The time to wait for the events of an address, from 1 to 3600 seconds, 10 by default.
  -o, --report <file>     Also write the report to the file.
  -p, --properties <name=value,...>
                          Set the properties of the endpoint, also as system properties, such as
                          dxfeed.aggregationPeriod=1s.

Verdicts and exit statuses (the worst of the two addresses with --compare):
  DATA                 0  Events came for every type and symbol.
  PARTIAL DATA         2  Events came for some types and symbols only.
  CONNECTED, NO DATA   3  The endpoint connected, but no events came.
  CANNOT CONNECT       4  The endpoint did not connect.
  CANNOT RESOLVE       5  The host name does not resolve.
An error of the arguments gives the exit status 2 and another error 1.

Samples:
  tools diagnose demo.dxfeed.com:7300 Quote,Trade AAPL,IBM
  tools diagnose demo.dxfeed.com:7300 Quote AAPL --compare dxlink:wss://demo.dxfeed.com/dxlink-ws --duration 20
  tools diagnose demo.dxfeed.com:7300 Quote AAPL -o report.txt
`

const qdsHelp = `Qds: runs the QDS tools of the dxFeed Graal Native SDK, such as connect, dump, multiplexor or feed, with their
command line in one argument. Without arguments it prints the list of the QDS tools.

Usage:
  tools qds "<qds-args>" [<options>]

Where:
  <qds-args>  The command line of the QDS tools: the name of a tool and its arguments, such as
              "connect demo.dxfeed.com:7300 Quote AAPL"; several tools run together when separated by +. The other
              arguments, except the first -p, are added to it. "help <tool>" prints the help of a QDS tool.

Options:
  -p, --properties <name=value,...>
                          Set system properties before the QDS tools run, such as dxfeed.aggregationPeriod=1s.

-h, --help and -help print this help only as the first argument; later they go to the QDS tools. The exit status is
that of the QDS tools, 0 if they succeed and 1 if they fail; a wrong option gives 2.

Samples:
  tools qds "connect demo.dxfeed.com:7300 Quote AAPL -stat 5s"
  tools qds "help connect"
`

const dumpHelp = `Dump: prints the events received from the address and, with -t, writes them to a tape file. It uses the StreamFeed
role, so the events of a tape file come as they are, without conflation.

Usage:
  tools dump <address> [<types> [<symbols>]] [<options>]

Where:
` + addressArgument + allTypesArgument + `             The default is all.
` + symbolsArgument + `             The default is all.

Options:
  -t, --tape <file>       Write the events to the tape file, such as out.txt[format=text].
  -p, --properties <name=value,...>
                          Set the properties of the endpoints, also as system properties, such as
                          dxfeed.aggregationPeriod=1s.
  -q, --quiet             Do not print the events.

Samples:
  tools dump demo.dxfeed.com:7300 Quote AAPL,IBM -t out.txt[format=text] -q
  tools dump in.tape[speed=max] Quote,Profile,TimeAndSale all -t out.tape -q
  tools dump in.tape[speed=max]
`

const perfTestHelp = `PerfTest: connects to the address and prints every 2 seconds the rate of events, the rate of listener calls, the
number of events in a call and the CPU usage of the process.

Usage:
  tools perftest <address> <types> <symbols> [<options>]

Where:
` + addressArgument + allTypesArgument + symbolsArgument + `
Options:
  -p, --properties <name=value,...>
                          Set the properties of the endpoints, also as system properties, such as
                          dxfeed.aggregationPeriod=1s.
  --force-stream          Enforce a streaming contract for the subscription: the StreamFeed role instead of Feed.
  --cpu-usage-by-core     Show the CPU usage relative to one core (a fully loaded core is 100%), not to all of them.
  --detach-listener       Do not add a listener to the subscription: the events are received but not counted.

Sample:
  tools perftest demo.dxfeed.com:7300 TimeAndSale AAPL,IBM,MSFT
`

const latencyTestHelp = `LatencyTest: connects to the address, subscribes to TimeAndSale events and calculates their latency: the delay
between the time of an event and its receipt, in milliseconds (min, max, mean, 99th percentile, etc).

Usage:
  tools latencytest <address> [<types> [<symbols>]] [<options>]

Where:
` + addressArgument + allTypesArgument + `             The default is all. Only TimeAndSale events are measured; the others are counted.
` + symbolsArgument + `             The default is all.

Options:
  -p, --properties <name=value,...>
                          Set the properties of the endpoints, also as system properties, such as
                          dxfeed.aggregationPeriod=1s.
  --interval <seconds>    The measurement interval in seconds, 2 by default.
  --ignore-exchanges <codes>
                          Leave out the events of the exchange codes, a comma-separated list without spaces, such
                          as Q,N.
  --force-stream          Enforce a streaming contract for the subscription: the StreamFeed role instead of Feed.

Sample:
  tools latencytest demo.dxfeed.com:7300 TimeAndSale AAPL,IBM,MSFT
`
