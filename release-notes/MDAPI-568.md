* **\[MDAPI-568]\[GO]** Fixed `Time`, `Sequence`, `SetTime` and `SetSequence` of `Order`, `AnalyticOrder` and
  `SpreadOrder`: they now use the time-sequence field, as in Java and C++. Before, received orders reported time 0,
  and `SetTime` or `SetSequence` overwrote the index with the order source.
