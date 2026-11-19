* **\[MDAPI-554]\[GO]** Added candle symbols with attributes.
  * `candle.CandleSymbol` parses and formats the symbols of the Java and C++ APIs: `candle.ParseCandleSymbol` (an error
    for an invalid attribute), `candle.NewCandleSymbol` and `candle.NewCandleSymbolWithAttributes`; the string form is
    normalized. `Symbol` and `SetSymbol` are deprecated.
  * The attributes: `candle.CandlePeriod` and `candle.CandleType`, `candle.CandlePrice`, `candle.CandleSession`,
    `candle.CandleAlignment`, `candle.CandleExchange`, `candle.CandlePriceLevel` and `candle.CandleDataType`.
