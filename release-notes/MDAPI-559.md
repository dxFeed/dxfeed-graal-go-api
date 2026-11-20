* **\[MDAPI-559]\[GO]** The tools take the options of the tools of the dxFeed C++ API.
  * `-p`/`--properties` for every tool sets the properties of the endpoints, also as system properties.
  * `connect` takes `-f`/`--from-time`, `-s`/`--source`, `-t`/`--tape`, `-q`/`--quiet` and `--force-stream`; `dump`
    takes `-t`, `-q` and subscribes to all events of all symbols by default.
  * `perftest` takes `--force-stream`, `--cpu-usage-by-core` and `--detach-listener` and prints the CPU usage;
    `latencytest` takes `--force-stream`, `--interval` and `--ignore-exchanges`.
  * `tools qds "<qds-args>" [-p <properties>]` runs the QDS tools of the SDK, such as connect or multiplexor, and
    exits with their status; without arguments it lists them.
  * `tools diagnose <address> <types> <symbols>` checks the connection step by step and prints a report with a
    verdict (DATA, PARTIAL DATA, CONNECTED, NO DATA, CANNOT CONNECT, CANNOT RESOLVE), which is also the exit status
    (0, 2 to 5); `--compare` diagnoses a second address, `--duration` sets the wait, `-o` writes the report to a
    file. Logins, passwords and tokens are masked.
  * `tools nettest p|c <address> -S <n>` publishes or receives generated Trade or TimeAndSale events at a controlled
    rate and prints the rate and the lag, as the NetTest specification `nettest/1` of the SDK 3.18.0 defines them;
    the command line, the events and the output lines are those of the NetTest of the other dxFeed APIs.
  * `tools nettest` writes a JSON summary of the run with `--summary <file>` (schema `nettest-summary/1`, with the
    e2e lag percentiles of an HdrHistogram) and serves the metrics in the Prometheus text format with
    `--metrics-port [<host>:]<port>`.
  * `tools --version` prints the versions of the tools and of the SDK.
  * An unknown option, an extra argument or a wrong value is an error with the exit status 2, as before an unknown
    event type; an error of a run gives the exit status 1.
