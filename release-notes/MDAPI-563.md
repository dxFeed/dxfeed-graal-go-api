* **\[MDAPI-563]\[GO]** Go samples for the remaining samples of the C++ API.
  * Tape file samples: `write-tape-file` publishes quotes into a tape file, `file-parser` prints the events of a
    tape file, `convert-tape-file` converts one tape file into another. `publish-profiles` serves Profile events of
    `:TEST` symbols, and `request-profile` requests one of them. The tests run these samples, without network, on
    every CI platform.
  * API samples: `connect` prints events of any types and symbols, optionally from a time; `quote-and-trade` shows
    listeners, a promise and snapshots of quotes and trades; `fetch-daily-candles` fetches the daily candles of the
    last 20 days; `last-event-console` keeps the last events of many symbols in memory and prints those of the symbols
    typed in the console.
