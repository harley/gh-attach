package attach

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
)

var uploadTokenRe = regexp.MustCompile(`"uploadToken":"([^"]+)"`)

func getUploadToken(client *http.Client, owner, repo string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://github.com/%s/%s", owner, repo), nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetch repo page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("repo page returned %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read repo page: %w", err)
	}

	match := uploadTokenRe.FindSubmatch(body)
	if match == nil {
		return "", fmt.Errorf("uploadToken not found; write access may be required")
	}

	return string(match[1]), nil
}
