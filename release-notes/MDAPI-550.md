* **\[MDAPI-550]\[GO]\[BREAKING]** The public API is cleaned up before v1.
  * **\[BREAKING]** The package `pkg/api/Osub` is renamed to `pkg/api/osub`.
  * **\[BREAKING]** The helper packages `pkg/mathutil`, `pkg/timeutil`, `pkg/formatutil` and `pkg/parser` are moved to
    `internal/` and are no longer part of the public API.
