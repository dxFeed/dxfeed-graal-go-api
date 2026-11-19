* **\[MDAPI-557]\[GO]** Added the `schedule` package: trading schedules, days and sessions.
  * `schedule.FromProfile`, `FromProfileAndVenue` and `FromDefinition` return a `Schedule`; `TradingVenues`,
    `DownloadDefaults` and `SetDefaults` manage the venues and the default schedule data.
  * `Schedule`, `Day` and `Session` find the day or the session at a time, the sessions of a day, and the previous
    or next day or session that a `DayFilter` or `SessionFilter` accepts, such as `TradingDays` or `RegularSessions`;
    a filter may test the days or sessions with a Go function.
  * The new `schedule` sample prints the schedule of an instrument at a time.
