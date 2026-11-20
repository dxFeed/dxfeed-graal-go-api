* **\[MDAPI-598]\[GO]\[BREAKING]** dxFeed Graal Native SDK updated to 3.18.0, the minimum version of the native
  library now.
  * **\[BREAKING]** A library of SDK 3.17 is refused at start-up: install the new one with
    `go run github.com/dxfeed/dxfeed-graal-go-api/cmd/install-sdk` and point `CGO_LDFLAGS` and the run-time library
    path to its directory (the directory name holds the SDK version).
  * `DXEndpoint.ConnectorIOStatistics()` returns the I/O statistics of each message connector as
    `api.MessageConnectorIOStatistics`: the data and subscription records, the bytes, the sums of the data lags and of
    the round-trip times in microseconds, read and written, and the open and closed connections. It returns an error
    that wraps `api.ErrIOStatisticsUnsupported` when the library does not support them.
  * The SDK deserializes only the dxFeed, QD and standard JDK classes: an attachment of `Message` or `Configuration`
    of another class comes as a nil attachment, and an ORCS `PriceLevelService` result of another class fails the
    request with an error.
