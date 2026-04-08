# gh-attach

`gh-attach` is a GitHub CLI extension for native GitHub attachments from the terminal.

It targets the gap in `gh` where PRs and issues accept markdown bodies, but there is no first-class CLI flow for uploading screenshots or other supported assets and embedding them before create.

## Install

```bash
gh extension install harley/gh-attach
```

For local development:

```bash
go build -o bin/gh-attach ./cmd/gh-attach
gh attach --help
```

## Usage

```bash
# Upload files and print markdown
gh attach upload screenshot.png
gh attach screenshot.png

# Create a PR with uploaded attachments appended to the body
gh attach pr create \
  --attach screenshot.png \
  --title "Fix layout regression" \
  --body "This fixes the mobile spacing issue."
```

## Why Native Attachments

- Uses GitHub's native attachment flow rather than release assets
- Preserves private-repo privacy semantics
- Keeps attachments on GitHub's own `user-attachments` infrastructure
- Renders images as markdown images and other files as markdown links

## Current Scope

- Upload native attachments and print markdown
- Wrap `gh pr create` with `--attach`
- Support `--body` and `--body-file` when composing the final PR body

## Constraints

- Native upload requires an authenticated GitHub browser session
- If browser cookie discovery does not work on your machine, set `GH_ATTACH_USER_SESSION`
- On macOS, `gh-attach` also checks known Chromium-clone stores such as `Comet`
- `--fill`, `--fill-first`, and `--fill-verbose` are not supported with `--attach` yet
- The first cut is optimized for private-repo screenshots and attachments, not CI automation
