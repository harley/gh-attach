# gh-attach

`gh-attach` is a GitHub CLI extension for native GitHub attachments from the terminal.

It targets the gap in `gh` where PRs and issues accept markdown bodies, but there is no first-class CLI flow for uploading screenshots or other supported assets and embedding them before create.

> [!IMPORTANT]
> This extension uses GitHub's browser-backed native upload flow, which is not a
> documented public API. GitHub may change it without notice. `gh-attach` never
> prints session-cookie values, but an explicit session override must be treated
> as a secret.

## Install

```bash
gh extension install harley/gh-attach
```

Upgrade later with:

```bash
gh extension upgrade attach
```

For local development:

```bash
go build -o bin/gh-attach ./cmd/gh-attach
gh attach --help
```

## Quick start

```bash
# Upload files and print markdown
gh attach upload screenshot.png
gh attach screenshot.png

# Diagnose browser-session discovery without printing cookie values
gh attach doctor

# Create a PR with uploaded attachments appended to the body
gh attach pr create \
  --attach screenshot.png \
  --title "Fix layout regression" \
  --body "This fixes the mobile spacing issue."
```

Use `--repo OWNER/REPO` when running outside the target repository:

```bash
gh attach upload --repo owner/repo --json screenshot.png
```

Image assets are rendered as Markdown images. Other supported assets are
rendered as Markdown links. `--json` returns the uploaded URL, filename, and
content type for each file.

## Why Native Attachments

- Uses GitHub's native attachment flow rather than release assets
- Preserves private-repo privacy semantics
- Keeps attachments on GitHub's own `user-attachments` infrastructure
- Renders images as markdown images and other files as markdown links

## Current Scope

- Upload native attachments and print markdown
- Wrap `gh pr create` with `--attach`
- Support `--body` and `--body-file` when composing the final PR body

## Authentication

Native uploads require both:

- An authenticated GitHub CLI session for repository lookup
- An authenticated `github.com` browser session for the native upload flow

Cookie discovery supports Chrome, Chromium, Brave, and Edge through the public
[`kooky`](https://github.com/browserutils/kooky) library. Run the redacted
diagnostic first when discovery fails:

```bash
gh attach doctor
```

For headless or unsupported environments, `GH_ATTACH_USER_SESSION` is the
explicit fallback. Treat it as a secret, avoid shell history, and unset it after
use. `GITHUB_USER_SESSION` is accepted only as a backwards-compatible alias.

See [SECURITY.md](SECURITY.md) before using an explicit session override.

## Constraints

- Native upload requires an authenticated GitHub browser session
- Cookie discovery supports Chrome, Chromium, Brave, and Edge through the public `kooky` library
- If browser cookie discovery does not work on your machine, run `gh attach doctor`, then sign in with a supported browser or set `GH_ATTACH_USER_SESSION`
- `GITHUB_USER_SESSION` remains accepted as a backwards-compatible alias and is identified as legacy by `gh attach doctor`
- On macOS, cookie decryption depends on the browser's Safe Storage key in the login Keychain. `gh attach doctor` reports Keychain access failures with recovery guidance
- `--fill`, `--fill-first`, and `--fill-verbose` are not supported with `--attach` yet
- The first cut is optimized for private-repo screenshots and attachments, not CI automation

## Supported platforms

Versioned releases provide precompiled binaries for the platforms supported by
GitHub CLI's official extension precompile workflow. Install with
`gh extension install harley/gh-attach`; no local Go toolchain is required for a
released version.

## Project status

`gh-attach` is an early-stage community project. It is not verified or endorsed
by GitHub. See [CONTRIBUTING.md](CONTRIBUTING.md) for local development and
release guidance.

## License

MIT
