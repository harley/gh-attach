package browser

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/browserutils/kooky"
	_ "github.com/browserutils/kooky/browser/brave"
	_ "github.com/browserutils/kooky/browser/chrome"
	_ "github.com/browserutils/kooky/browser/chromium"
	_ "github.com/browserutils/kooky/browser/edge"
)

var sessionValuePattern = regexp.MustCompile(`(?i)user_session=[^;\s]*`)

func GetGitHubSession() (*http.Cookie, error) {
	if value, _ := envSession(); value != "" {
		return &http.Cookie{
			Name:     "user_session",
			Value:    value,
			Domain:   "github.com",
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
		}, nil
	}

	cookies, err := kooky.ReadCookies(
		context.Background(),
		kooky.Valid,
		kooky.DomainHasSuffix("github.com"),
		kooky.Name("user_session"),
	)

	if len(cookies) > 0 {
		return &cookies[0].Cookie, nil
	}
	if err != nil {
		return nil, errors.New(explainSessionError(err))
	}

	return nil, errors.New(explainSessionError(errors.New("no github.com user_session cookie found")))
}

func envSession() (string, string) {
	if value := os.Getenv("GH_ATTACH_USER_SESSION"); value != "" {
		return value, "GH_ATTACH_USER_SESSION"
	}
	if value := os.Getenv("GITHUB_USER_SESSION"); value != "" {
		return value, "GITHUB_USER_SESSION (legacy alias)"
	}
	return "", ""
}

// Diagnostic reports session discovery without exposing cookie material.
type Diagnostic struct {
	EnvironmentOverride bool
	SessionFound        bool
	Source              string
	Problem             string
}

func Diagnose() Diagnostic {
	if _, source := envSession(); source != "" {
		return Diagnostic{EnvironmentOverride: true, SessionFound: true, Source: source}
	}

	_, err := GetGitHubSession()
	if err != nil {
		return Diagnostic{Problem: err.Error()}
	}
	return Diagnostic{SessionFound: true, Source: "supported browser cookie store"}
}

func (d Diagnostic) String() string {
	envStatus := "not configured"
	if d.EnvironmentOverride {
		envStatus = "configured"
	}
	sessionStatus := "not found"
	if d.SessionFound {
		sessionStatus = "available"
	}

	var out strings.Builder
	fmt.Fprintln(&out, "gh-attach session diagnostics")
	fmt.Fprintf(&out, "Environment override: %s\n", envStatus)
	fmt.Fprintf(&out, "GitHub browser session: %s\n", sessionStatus)
	if d.Source != "" {
		fmt.Fprintf(&out, "Session source: %s\n", d.Source)
	}
	if d.Problem != "" {
		fmt.Fprintf(&out, "\nProblem:\n%s\n", d.Problem)
	}
	return strings.TrimRight(out.String(), "\n")
}

func explainSessionError(err error) string {
	if err == nil {
		return ""
	}
	detail := sessionValuePattern.ReplaceAllString(err.Error(), "user_session=[REDACTED]")
	lower := strings.ToLower(detail)

	if strings.Contains(lower, "keychain") ||
		strings.Contains(lower, "safe storage") ||
		strings.Contains(lower, "-25293") {
		return "GitHub browser cookie found, but macOS Keychain could not unlock the browser's Safe Storage key. Unlock your login keychain and retry, sign in with another supported browser profile, or set GH_ATTACH_USER_SESSION explicitly. Run `gh attach doctor` to verify the session setup.\n\nUnderlying browser error: " + detail
	}

	if strings.Contains(lower, "no github.com user_session") {
		return "No authenticated GitHub browser session was found. Sign in to github.com using Chrome, Chromium, Brave, or Edge, then retry. For headless or unsupported environments, set GH_ATTACH_USER_SESSION explicitly. Run `gh attach doctor` to verify the session setup."
	}

	return "Could not read an authenticated GitHub browser session. Retry after signing in to github.com, or set GH_ATTACH_USER_SESSION explicitly. Run `gh attach doctor` for details.\n\nUnderlying browser error: " + detail
}
