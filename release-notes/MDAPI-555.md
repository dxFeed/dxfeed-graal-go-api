* **\[MDAPI-555]\[GO]** Instrument profiles: token, all fields by name, writer.
  * `ipf.InstrumentProfileReader.ReadFromFileWithToken` reads profiles with an `auth.Token`.
  * `ipf.InstrumentProfileField`, `ipf.GetField`, `SetField`, `GetNumericField`, `SetNumericField`, `GetDateField` and
    `SetDateField` read and write the standard and the custom fields of a profile by name; `ipf.FormatNumber`,
    `ParseNumber`, `FormatDate`, `ParseDate` and the `ipf.Type*` instrument types.
  * `ipf.WriteToFile` writes profiles to a file in the IPF format, plain or compressed.