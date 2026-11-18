* **\[MDAPI-553]\[GO]** Added the remaining event types of Graal Native SDK 3.16 and the event flags API.
  * `events.EventFlag` (`TxPending`, `RemoveEvent`, `SnapshotBegin`, `SnapshotEnd`, `SnapshotSnip` and others), the
    `events.IndexedEvent` interface and the `IsSnapshotBegin`, `IsSnapshotEnd`, `IsSnapshotSnip`,
    `IsSnapshotEndOrSnip`, `IsPending` and `IsRemove` helpers.
  * New events: `series.Series`, `theoprice.TheoPrice`, `candle.DailyCandle`, and `misc.Message`,
    `misc.TextMessage`, `misc.Configuration`, `misc.TextConfiguration`.
  * New market events: `optionsale.OptionSale`, `marketmaker.MarketMaker`, `order.OtcMarketsOrder` (with
    `order.OtcMarketsPriceType`) and `orderimbalance.OrderImbalance` (with `orderimbalance.ImbalanceSide` and
    `orderimbalance.AuctionType`).
  * New Nuam Exchange events: `order.NuamOrder` (with `order.NuamOrderType` and `order.NuamTimeInForceType`),
    `timeandsale.NuamTimeAndSale` and `trade.NuamTrade`.
  * `TimeAndSale.String` prints the buyer and the seller instead of the addresses of their strings.
