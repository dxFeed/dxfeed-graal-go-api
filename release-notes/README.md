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
- Mark a breaking change with `**BREAKING:**` at the start of the top-level item or the nested item it concerns.
- No blank lines: they would split the item.

CI checks the fragments on every pull request:

```bash
go run ./internal/cmd/release check
```
