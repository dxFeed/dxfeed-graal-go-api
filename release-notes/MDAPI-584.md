* **\[MDAPI-584]\[GO]** Fixed a panic of the `tools` program on `--help`, `-h`, `help` or an unknown tool name
  * Fixed `--help`, `-h` and `help`: they print the usage with the list of tools, as `tools` without arguments does,
    and exit with status 0 without starting the SDK.
  * Fixed an unknown tool name: it prints an error naming the tool and the usage, and exits with status 2.
  * Fixed the order of the list of tools: it is alphabetical instead of random.
  * Fixed the help of `latencytest`: it described the performance counters of `perftest` instead of the latency test.
