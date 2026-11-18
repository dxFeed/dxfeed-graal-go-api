* **\[MDAPI-567]\[GO]** Fixed `Profile.IsShortSaleRestricted`: it checks the short sale restriction, as in Java and C++.
  Before, it compared the trading status and returned true for a halted instrument.
