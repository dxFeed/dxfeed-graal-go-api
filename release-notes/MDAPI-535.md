* **\[MDAPI-535]\[GO]** Fixed `NewEndpointWithProperties`: connection state listeners added with `AddListener` are now
  called, as for `NewEndpoint`.
  * Fixed a memory leak: the property keys and values passed to `NewEndpointWithProperties` were never freed
    ([#22](https://github.com/dxFeed/dxfeed-graal-go-api/issues/22), reported by @OvOhao).
