* **\[MDAPI-548]\[GO]** The API calls the current (`dxfg_v4`) functions of the dxFeed Graal Native SDK instead of
  the ones that SDK 3.16 deprecates and SDK 4.0.0 removes.
  * Failures of adding or removing symbols, adding a listener and publishing events are returned as errors; before,
    they were lost.
  * Publishing a `Profile` without a status reason, or another event with an unset string field, no longer crashes
    the process.
