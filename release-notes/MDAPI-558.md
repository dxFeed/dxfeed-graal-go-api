* **\[MDAPI-558]\[GO]** Added the event models of the new package `model`.
  * `model.NewIndexedTxModelBuilder` builds an `IndexedTxModel` of the indexed events of a symbol, such as the
    orders of an order book: it passes the complete transactions of its sources to a `TxModelListener`, with batch
    and snapshot processing; `Attach`, `Detach` and `SetSources` change its subscription.
  * `model.NewTimeSeriesTxModelBuilder` builds a `TimeSeriesTxModel` of the time series events of a symbol from a
    time, such as candles: it passes their complete transactions to a `TimeSeriesTxModelListener`; `SetFromTime`
    changes the time of its subscription. The new interface `events.TimeSeriesEvent` describes these events.
  * The new `inc-order-snapshot` sample prints the snapshot and the incremental updates of the orders of a source.
