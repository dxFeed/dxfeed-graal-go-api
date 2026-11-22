// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package nettest

// HelpText is the help of the NetTest tool.
const HelpText = `NetTest: publishes (mode p) or receives (mode c) generated Trade or TimeAndSale events at a controlled rate and
prints every period the rate and the lag, as the NetTest specification nettest/1 of the dxFeed Graal Native SDK
defines them (docs/wrapper-guide/nettest.md of its sources).
spec_version nettest/1

Usage:
  tools nettest p|c <address> -S <n>[@<file>] [<options>]

Where:
  p|c        The mode: p publishes the events, c receives them.
  <address>  host:port connects to the address, :port listens on the port.

Options:
  -S <n>, -S <n>@<file>   The symbols: n generated ones, or the first n of an IPF file; n from 1 to 100000000.
                          Required.
  -C <n>                  The instances, each with its own endpoint, from 1 to 1000; 1 by default and with :port.
  -r, --rate-limit */<rate>
                          p only: the events per second of all instances together, such as */500k or */1.5M;
                          as fast as possible by default.
  -s, --stat <time>       The statistics period, such as 500ms, 10s or 1m; 10s by default.
  -n, --name <name>       The name in the output lines; the endpoints are named <name>.<i>. NetTest by default.
  -c, --collector stream|ticker
                          The contract of the consumer, which selects the roles of the endpoints.
  --event Trade|TimeAndSale
                          The event type, Trade by default.
  --batch <n>             p only: the events per publishing call, from 1 to 4194304; 100 by default.
  --reuse                 p only: republish the events of the pool instead of new events per batch.
  --role <role>           The role: PUBLISHER or STREAM_PUBLISHER on p, FEED or STREAM_FEED on c; by default from
                          -c: STREAM_PUBLISHER with stream, PUBLISHER otherwise on p; FEED with ticker, STREAM_FEED
                          otherwise on c.
  --duration <time>       End the run after the time with the exit status 0.
  --warmup <time>         The warm-up time: on c from the first event, on p from the first batch.
  --summary <file>        Write a JSON summary of the run (schema nettest-summary/1) to the file at the end.
  --metrics-port [<host>:]<port>
                          Serve the metrics in the Prometheus text format at http://<host>:<port>/metrics;
                          the host is 127.0.0.1 by default, the port 0 takes a free port.
  -D<key>=<value>         Set a system property before the first endpoint, also on the endpoints;
                          -Ddxscheme.nanoTime=false turns off the nanosecond time for peers without it.
  -h, -help, --help       Print this help.

Not supported: -X and -XX (the options of a Graal isolate), --legacy.

The output lines, once per period:
  {*NetTestProducerSide-<name>*} data <rps> rps, lag <us> us, instances <C>, connections <M>
      <name>.<i> data <rps> rps, lag <us> us
  {*NetTestRuntime-<name>*} heap <used> MB of <committed> MB, gc <count> (<time> ms), cpu <percent>%
On p the lag is the time of filling and publishing a batch, on c the time from the publishing to the receipt.

Exit status: 0 at the end of --duration or on Ctrl+C, 1 on an error of the run, 2 on an error of the command line
and without arguments.

Samples:
  tools nettest p :7001 -S 1000 -r */100k -c stream
  tools nettest c 127.0.0.1:7001 -S 1000 -C 4 -c stream --duration 1m
`
