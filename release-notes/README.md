# Release note fragments

Every pull request with a Jira ticket adds or extends `release-notes/MDAPI-NNN.md`, where `MDAPI-NNN` is the ticket
key. A release moves all fragments into [ReleaseNotes.md](../ReleaseNotes.md) under the new version heading.

## Format

One top-level item per file, for the ticket in the file name; further changes of the same ticket go into nested items:

```markdown
* **\[MDAPI-NNN]\[GO]** What changed for users of the API.
  A long item continues on lines indented by two spaces.
  * A further change of the same ticket.
```

- Write for API users: what changed and what they have to do, not how it was implemented.
- Mark a ticket with a breaking change as `**\[MDAPI-NNN]\[GO]\[BREAKING]**`, and start each breaking nested item with
  `**\[BREAKING]**`:

  ```markdown
  * **\[MDAPI-NNN]\[GO]\[BREAKING]** What changed for users of the API.
    * **\[BREAKING]** What users have to change in their code.
  ```

- No blank lines: they would split the item.

CI checks the fragments on every pull request:

```bash
go run ./internal/cmd/release check
```

CI also compares the exported API with the latest release tag ([gorelease](https://pkg.go.dev/golang.org/x/exp/cmd/gorelease)).
An incompatible change, such as a removed or changed function, type, method or field, fails the check unless a
fragment that the pull request adds or changes is marked `[BREAKING]`; the incompatible changes already on `main` do
not count. On a clean working tree with the tags and `origin/main` fetched and cgo set up as for a build:

```bash
go run ./internal/cmd/release api
```
