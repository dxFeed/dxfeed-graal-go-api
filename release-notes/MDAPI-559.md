* **\[MDAPI-559]\[GO]** The tools take the options of the tools of the dxFeed C++ API.
  * `-p`/`--properties` for every tool sets the properties of the endpoints, also as system properties.
  * `connect` takes `-f`/`--from-time`, `-s`/`--source`, `-t`/`--tape`, `-q`/`--quiet` and `--force-stream`; `dump`
    takes `-t`, `-q` and subscribes to all events of all symbols by default.
  * `perftest` takes `--force-stream`, `--cpu-usage-by-core` and `--detach-listener` and prints the CPU usage;
    `latencytest` takes `--force-stream`, `--interval` and `--ignore-exchanges`.
  * `tools qds "<qds-args>" [-p <properties>]` runs the QDS tools of the SDK, such as connect or multiplexor, and
    exits with their status; without arguments it lists them.
  * `tools --version` prints the versions of the tools and of the SDK.
  * An unknown option, an extra argument or a wrong value is an error with the exit status 2, as before an unknown
    event type; an error of a run gives the exit status 1.
