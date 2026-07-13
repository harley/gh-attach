# Contributing

Thank you for helping improve `gh-attach`.

## Development setup

You need the Go version declared in `go.mod` (or newer) and an authenticated
GitHub CLI installation.

```bash
git clone https://github.com/harley/gh-attach.git
cd gh-attach
go test ./...
go build -o bin/gh-attach ./cmd/gh-attach
gh extension install .
gh attach --help
```

Local extension installation uses the executable `gh-attach` script in the
repository root. It runs `bin/gh-attach` when present and otherwise falls back
to `go run`.

## Before opening a pull request

```bash
git ls-files -z '*.go' | xargs -0 gofmt -w
go test ./...
go vet ./...
go run golang.org/x/vuln/cmd/govulncheck@v1.6.0 ./...
./script/generate-third-party-notices.sh /tmp/THIRD_PARTY_NOTICES.txt
go build -o bin/gh-attach ./cmd/gh-attach
```

Keep pull requests focused. Add tests for behavior changes and never include
browser cookies, GitHub session values, or attachment URLs from private
repositories in fixtures, logs, screenshots, or issue reports.

## Releases

Maintainers publish a semantic version tag such as `v0.1.0`. The release
workflow cross-compiles the extension, creates a GitHub release, publishes the
platform-specific assets expected by `gh extension install`, and generates
build provenance attestations. Each release also includes generated third-party
license notices for every module linked into a published platform binary.

Tags containing a hyphen, such as `v0.2.0-rc.1`, produce prereleases and are
appropriate for testing the release path.
