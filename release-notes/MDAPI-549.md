* **\[MDAPI-549]\[GO]\[BREAKING]** The API returns errors instead of panicking.
  * A failure to create the isolate, or a native library of another version, is returned by the first call and every
    later one; before, every call panicked.
  * Java exceptions are returned as `common.JavaError` with the class name, the message and the stack trace; read it
    with `errors.As`. The error text is the class name and the message; the stack trace is in `JavaError.StackTrace`.
  * A panic of an event or connection state listener is logged and no longer ends the process.
  * Events of a type that the API does not support yet are skipped, and each such type is logged once.
  * Subscription symbols are passed to the SDK in C memory: time series and indexed event symbols follow the cgo
    pointer rules, and an unsupported symbol in `AddSymbols` or `RemoveSymbols` is an error.
  * **\[BREAKING]** `DXFeedSubscription.Clear` and `DXFeedSubscription.Close` return an `error`.
  * **\[BREAKING]** `api.SetSystemProperty` returns an `error`, and `api.GetSystemProperty` returns `(string, error)`.
  * **\[BREAKING]** `parser.ParseSymbols` returns `([]any, error)`; before, a parse failure gave an empty list.
