* **\[MDAPI-561]\[GO]** Added the client of the ORCS price level service.
  * `orcs.NewPriceLevelService` connects to the service: `GetOrders`, `GetQuotes` and their `WithCaller` forms
    request the price levels (orders) and the quotes of a candle symbol with the granularity attribute
    (`orcs.GranularityAttributeKey`) within a time range, `GetAuthOrderSource` returns the order sources and the
    symbols available to the client.
  * `orcs.NewAuthOrderSource` creates an `AuthOrderSource` from the symbols by the order source identifier.
  * `orcs.ValidatePriceLevels` checks the price levels as `PriceLevelChecker` does; the new
    `price-level-service-client` sample requests and checks them.
