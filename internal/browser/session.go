package browser

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/browserutils/kooky"
	_ "github.com/browserutils/kooky/browser/brave"
	_ "github.com/browserutils/kooky/browser/chrome"
	_ "github.com/browserutils/kooky/browser/chromium"
	_ "github.com/browserutils/kooky/browser/edge"
)

func GetGitHubSession() (*http.Cookie, error) {
	if value := envSession(); value != "" {
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
		return nil, fmt.Errorf("read browser cookies: %w", err)
	}

	return nil, fmt.Errorf("no github.com user_session cookie found in supported browsers; set GH_ATTACH_USER_SESSION to override")
}

func envSession() string {
	if value := os.Getenv("GH_ATTACH_USER_SESSION"); value != "" {
		return value
	}
	return os.Getenv("GITHUB_USER_SESSION")
}
