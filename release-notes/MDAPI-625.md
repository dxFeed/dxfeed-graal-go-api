* **\[MDAPI-625]\[GO]** `tools nettest` follows the NetTest specification `nettest/1` more closely.
  * A consumer that listens on a port no longer stops with `no_data` while no peer is connected: the time without
    data counts only with open connections, as the `connections` field of the statistics lines shows them.
  * Without `-r`, the metrics have no sample of `nettest_schedule_slip_seconds`, only its family.
  * In the summary file, `schedule_slip.max_ns` leaves out the warm-up when a measured interval exists, and
    `growth_ns_per_s` is truncated as in the other NetTest tools; `runtime.gc_count` and `runtime.gc_time_ms` are the
    growth over the run, and `runtime.cpu_percent_mean` is the mean from the start of the instances.
  * `-r` below 1 or above `1000G` is an error with the text of the other NetTest tools.
  * The event times and the e2e lag have the nanosecond resolution on Windows too: the event time is the wall-clock
    time at the start plus the monotonic time since then.
  * `-Ddxscheme.nanoTime` is on for `true` in any case, as in Java.
  * An option that takes a value takes the next argument whatever it is, as in the other NetTest tools: `-n --help`
    names the run `--help`; a help option anywhere else prints the usage, even after an invalid argument.
  * The `no_data` error gives its limit in milliseconds, as the other NetTest tools do.
  * A signal while the symbols are read ends the run with the summary `stopped_by_signal` and without the endpoints;
    a second signal ends the process at once, with the exit code 128 plus the signal number.
  * A run shorter than the period prints its last block, counted from the start, when it lasts at least one second.
