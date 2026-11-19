* **\[MDAPI-566]\[GO]** The packages have overviews on pkg.go.dev: package `api` introduces the module (endpoint roles,
  subscriptions and symbols, event types, time series, publishing, instrument profiles, threads and errors), the other
  packages describe their concepts and link the related types.
  * Examples on pkg.go.dev: time series and indexed subscriptions, the connection state listener, candle symbols, auth
    tokens, the time format and periods of `dxtime`, and a `log/slog` listener of `logging`.
  * The guides `docs/quick-start.md`, `docs/configuration.md` and `docs/troubleshooting.md` describe the installation,
    the configuration and the usual problems; a test checks the links of the Markdown files.
