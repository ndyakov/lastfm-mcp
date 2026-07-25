# Contributing

Thank you for improving `lastfm-mcp`.

## Before opening a change

- Use an issue for substantial API or security design changes.
- Never include Last.fm credentials, session keys, bearer tokens, passwords, or private listening data.
- Keep read-only behavior as the default; new mutation or authentication features must be explicitly gated.

## Development

Fork the repository, create a focused branch, and run:

```sh
go mod tidy
gofmt -w .
go vet ./...
go test -race ./...
```

Add deterministic tests for behavior changes. Tests must not depend on the live Last.fm service. Use an `httptest` server and `lastfm.WithBaseURL` for API fixtures.

Commit messages should be concise and imperative. Pull requests should explain behavior, risks, validation, and any user-facing configuration changes. By contributing, you agree that your contribution is licensed under the project's MIT License.
