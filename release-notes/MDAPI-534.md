* **\[MDAPI-534]\[GO]** Fixed event codes after `AnalyticOrder`, which were shifted by one against the native SDK:
  `SpreadOrder` subscribed to OTC Markets orders and `Series` to spread orders, and publishing a `SpreadOrder` panicked.
  * **BREAKING:** the values of `eventcodes.SpreadOrder`, `Series` and `OptionSale` change; code that stores the
    numbers must be updated.
  * Added the codes `OtcMarketsOrder`, `TextMessage`, `MarketMaker`, `TextConfiguration`, `NuamOrder`,
    `NuamTimeAndSale`, `NuamTrade` and `OrderImbalance`; their Go event types are not implemented yet.
  * `CreateSubscription` and `Publish` return an error for an event type without a Go implementation, and received
    events of such a type are skipped instead of panicking.
