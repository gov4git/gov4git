## Integration tests

This folder contains integration tests that perform network access to a test repo on GitHub.

Run tests using:

```go test -v -tags=integration ./...```

Tests do not require an auth token. If in the future tests need an auth token, use:

```GITHUB_AUTH_TOKEN=YOUR_TOKEN go test -v -tags=integration ./...```
