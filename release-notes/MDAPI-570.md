* **\[MDAPI-570]\[GO]** Fixed a data race in the connection state listeners of `DXEndpoint`: `AddListener` and
  `RemoveListener` may now be called from any goroutine, also while the connection state changes or from a listener.
