# Security policy

## Reporting a vulnerability

Please do not open a public issue for a vulnerability or include session
material in any report. Use GitHub's private vulnerability reporting for this
repository. If that option is unavailable, open an issue that contains no
sensitive details and asks the maintainer for a private reporting channel.

## Session-cookie handling

`gh-attach` needs an authenticated GitHub browser session because GitHub's
native attachment upload flow is not exposed by the public REST or GraphQL API.
The extension reads a `user_session` cookie locally and sends it only to
GitHub-owned endpoints involved in the upload flow.

Treat `GH_ATTACH_USER_SESSION` as a secret with the same care as a signed-in
browser session:

- Prefer automatic discovery from a supported local browser
- Never commit, print, paste, or transmit the cookie value
- Avoid shell history when setting an explicit fallback
- Unset the variable as soon as the upload is complete
- Revoke the GitHub browser session if the value may have been exposed

`gh attach doctor` reports discovery status without printing cookie material.
