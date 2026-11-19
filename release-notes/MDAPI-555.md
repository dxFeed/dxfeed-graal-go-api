* **\[MDAPI-555]\[GO]** Instrument profiles: token, all fields by name, writer, collector and live connection.
  * `ipf.InstrumentProfileReader.ReadFromFileWithToken` reads profiles with an `auth.Token`.
  * `ipf.InstrumentProfileField`, `ipf.GetField`, `SetField`, `GetNumericField`, `SetNumericField`, `GetDateField` and
    `SetDateField` read and write the standard and the custom fields of a profile by name; `ipf.FormatNumber`,
    `ParseNumber`, `FormatDate`, `ParseDate` and the `ipf.Type*` instrument types.
  * `ipf.WriteToFile` writes profiles to a file in the IPF format, plain or compressed.  * `live.InstrumentProfileCollector` collects profiles and notifies `InstrumentProfileUpdateListener`s;
    `live.InstrumentProfileConnection` keeps it up to date from a file, a URL or a server (`CreateConnection`, `Start`,
    `WaitUntilCompleted`, `State`, state change listeners, update period); the new `live-ipf` sample prints the updates.