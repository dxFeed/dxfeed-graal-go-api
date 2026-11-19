* **\[MDAPI-593]\[GO]** Names of the event types; fixed the unknown event types in the tools
  * `eventcodes.EventCode.String` returns the name of the event type, such as `Quote`, so `%v` and `%s` print the
    name; `eventcodes.ByName` finds the code by the name and `eventcodes.All` returns every code.
  * Fixed: the tools dropped an unknown event type without an error. Now it prints the error and the help of the tool
    and gives status 2; the names are case-insensitive, `all` is matched in any case and includes `DailyCandle`, and
    `DailyCandle` is accepted by name.
  * The connect and file-parser samples take the names of the event types from `eventcodes` and accept `DailyCandle`
    too.