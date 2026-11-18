# Contributing

Thank you for helping with the dxFeed Graal Go API. Bug reports, questions, fixes and new features are all welcome.

## Issues

- **Bug or feature request:** open an [issue](https://github.com/dxFeed/dxfeed-graal-go-api/issues/new/choose) with
  the template. For a bug, give the version of the module, the platform and the Go version, and a minimal program
  that reproduces it.
- **Question about the API:** open an issue as well. Questions about dxFeed data, subscriptions or accounts go to
  your dxFeed support contact.
- Look for an existing issue first; add to it rather than open a duplicate.

## Pull requests

1. For anything bigger than a small fix, open an issue first, so we can agree on the approach before you write it.
2. Fork the repository and create a branch from `main`.
3. Set up the build and run the checks as [HOW_TO_DEVELOP.md](HOW_TO_DEVELOP.md) describes.
4. Make the change with tests, keep the existing code style (`gofmt`), and document new exported identifiers.
5. Add a release note fragment (see below) when users of the API will notice the change.
6. Open the pull request against `main` and fill in its template. CI builds and tests it on Windows, Linux and macOS.

Keep a pull request to one topic: a bug fix, a feature, or a refactoring. Small pull requests are reviewed sooner.

## Release notes

A change that users notice gets a fragment in [release-notes/](release-notes/README.md): one file per ticket, written
for users of the API, with breaking changes marked `[BREAKING]`. CI checks the format. The files are named after the
maintainers' ticket keys: in an external pull request, describe the change for users in its description, and the
maintainers add the fragment.

Documentation, CI and test-only changes need no fragment.

## Compatibility

The module follows semantic versioning; [HOW_TO_RELEASE.md](HOW_TO_RELEASE.md#version-numbers) lists what counts as a
breaking change. Say in the pull request if your change breaks anything in the packages outside `internal/`, and why
it cannot be avoided.

## Labels

| Label           | Meaning                                                            |
|-----------------|--------------------------------------------------------------------|
| `bug`           | Something does not work as documented                              |
| `enhancement`   | A new feature or an improvement of an existing one                 |
| `question`      | A question about the API                                           |
| `documentation` | README, guides and Go doc comments                                 |
| `breaking`      | Changes the public API or its behavior incompatibly                |
| `platform`      | Specific to one OS, architecture or compiler                       |
| `dependencies`  | Updates of the Graal Native SDK, Go modules or GitHub Actions      |
| `ci`            | GitHub Actions workflows and the release process                   |

## AI-assisted contributions

You may use AI tools to write code, tests or documentation. You are responsible for what you submit as for anything
else you write:

- understand every line and be able to explain it in the review;
- build and test it yourself, as described above;
- make sure that it contains no code that you have no right to contribute.

Say in the pull request description if a substantial part of it was generated.

## Developer Certificate of Origin

By contributing, you certify the [Developer Certificate of Origin](https://developercertificate.org/): you wrote the
change or otherwise have the right to submit it under the [license](LICENSE) of the project. Sign off each commit to
state this:

```bash
git commit -s
```

It adds a `Signed-off-by: Your Name <your@email>` line with the name and email of your git configuration.

## License

Contributions are licensed under the [Mozilla Public License 2.0](LICENSE), as the rest of the project.
