# How to release

A release is a tag `vX.Y.Z` on `main`. Pushing the tag runs the [Release](.github/workflows/release.yml) workflow:
it builds and tests on all platforms, checks that the tag is the latest version in [ReleaseNotes.md](ReleaseNotes.md),
publishes a GitHub Release with that version's notes and the tools archives, and asks the Go module proxy for the
version, so `go get` sees it at once.

## Version numbers

The module follows [semantic versioning](https://semver.org) and stays at v0 until the API is complete.

- **v0:** a breaking change raises the minor version (`v0.4.0` → `v0.5.0`); other changes raise the minor or the patch
  version. A breaking change is marked `[BREAKING]` in its release note
  ([format](release-notes/README.md#format)).
- **v1 and later:** a breaking change raises the major version. Go requires a module path ending in `/vN` from v2 on
  (`github.com/dxfeed/dxfeed-graal-go-api/v2`), so a new major version changes the import path.

Breaking, for packages outside `internal/`:

- removing or renaming an exported identifier, changing a function signature or the type of an exported field;
- adding a method to an exported interface;
- changing documented behavior or which errors a function returns;
- dropping a supported platform or requiring users to change their build steps.

Not breaking: anything under `internal/`, a newer Graal Native SDK with the same API, a higher minimum Go version
(still mentioned in the release notes).

Tags `0.0.1`…`0.3.0` were made without the `v` prefix; they stay, and Go tools see them only as pseudo-versions.

## Steps

1. Make sure every merged change has its fragment in [release-notes/](release-notes/README.md) and that CI on `main`
   is green.
2. Update a clean checkout of `main` and preview the release:

   ```bash
   git switch main
   git pull --ff-only
   go run ./internal/cmd/release prepare v0.5.0 --dry-run
   ```

3. Prepare it: the command moves the fragments into `ReleaseNotes.md` under `## v0.5.0`, commits `Release v0.5.0` and
   creates the tag `v0.5.0` locally.

   ```bash
   go run ./internal/cmd/release prepare v0.5.0
   ```

4. Check the commit, then push it and the tag together; `--atomic` pushes neither if `main` has moved on:

   ```bash
   git push --atomic origin HEAD v0.5.0
   ```

5. Watch the Release workflow and check the GitHub Release and `go list -m github.com/dxfeed/dxfeed-graal-go-api@v0.5.0`.

To check a version's notes and archives without publishing, run the Release workflow by hand (`workflow_dispatch`)
with the version; it skips the publish job.

If the workflow fails, re-run the failed jobs when the cause is transient. A pushed tag is never moved or reused, since
the Go proxy may already have cached it: if the code needs a fix, merge it with its fragment and release the next patch
version.
