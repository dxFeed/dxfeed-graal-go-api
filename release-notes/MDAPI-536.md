* **\[MDAPI-536]\[GO]** Releases are tagged `vX.Y.Z`, so `go get` and the Go module proxy resolve them; `v0.4.0` is
  the first. Pushing a version tag publishes a GitHub Release with the notes of that version from `ReleaseNotes.md`.
  * Added `ReleaseNotes.md` with the history of earlier versions. Changes since the last release are kept as one file
    per Jira ticket in `release-notes/` and checked in CI.
  * `HOW_TO_RELEASE.md` describes the release steps and the versioning policy.
