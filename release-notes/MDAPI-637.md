* **\[MDAPI-637]\[GO]** NetTest on Windows reads the event time from `GetSystemTimePreciseAsFileTime`, as the C++ and
  Rust tools do: the start time plus the monotonic time drifted from the system clock by hundreds of microseconds, so
  the e2e lags against the other APIs were wrong by up to about 1 ms
