package browser

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/browserutils/kooky"
	_ "github.com/browserutils/kooky/browser/brave"
	_ "github.com/browserutils/kooky/browser/chrome"
	_ "github.com/browserutils/kooky/browser/chromium"
	kbchromium "github.com/browserutils/kooky/browser/chromium"
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
		if cookie, fallbackErr := readKnownChromiumCloneSession(context.Background()); fallbackErr == nil {
			return cookie, nil
		}
		return nil, fmt.Errorf("read browser cookies: %w; also checked clone stores: %s; or set GH_ATTACH_USER_SESSION", err, cloneStoreList())
	}

	if cookie, fallbackErr := readKnownChromiumCloneSession(context.Background()); fallbackErr == nil {
		return cookie, nil
	}

	return nil, fmt.Errorf("no github.com user_session cookie found in supported browsers or clone stores (%s); set GH_ATTACH_USER_SESSION to override", cloneStoreList())
}

func envSession() string {
	if value := os.Getenv("GH_ATTACH_USER_SESSION"); value != "" {
		return value
	}
	return os.Getenv("GITHUB_USER_SESSION")
}

func readKnownChromiumCloneSession(ctx context.Context) (*http.Cookie, error) {
	for _, storePath := range knownChromiumCloneStores() {
		cookies, err := kbchromium.ReadCookies(
			ctx,
			storePath,
			kooky.Valid,
			kooky.DomainHasSuffix("github.com"),
			kooky.Name("user_session"),
		)
		if err != nil || len(cookies) == 0 {
			continue
		}
		return &cookies[0].Cookie, nil
	}

	return nil, fmt.Errorf("no github.com user_session cookie found in known Chromium clone stores")
}

func knownChromiumCloneStores() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	return []string{
		filepath.Join(home, "Library", "Application Support", "Comet", "Default", "Cookies"),
		filepath.Join(home, "Library", "Application Support", "Comet", "OpenClaw", "Cookies"),
	}
}

func cloneStoreList() string {
	return filepath.Join("Library", "Application Support", "Comet", "Default", "Cookies") + ", " +
		filepath.Join("Library", "Application Support", "Comet", "OpenClaw", "Cookies")
}
