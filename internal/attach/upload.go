package attach

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const userAgent = "gh-attach/0.1"

type Asset struct {
	URL         string `json:"url"`
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
}

type policyResponse struct {
	UploadURL string `json:"upload_url"`
	Asset     struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		ContentType string `json:"content_type"`
		Href        string `json:"href"`
	} `json:"asset"`
	Form                         map[string]string `json:"form"`
	AssetUploadAuthenticityToken string            `json:"asset_upload_authenticity_token"`
}

func NewClient(sessionCookie *http.Cookie) *http.Client {
	jar, _ := cookiejar.New(nil)
	ghURL, _ := url.Parse("https://github.com")

	sameSiteCookie := &http.Cookie{
		Name:     "__Host-user_session_same_site",
		Value:    sessionCookie.Value,
		Domain:   "github.com",
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	}

	jar.SetCookies(ghURL, []*http.Cookie{sessionCookie, sameSiteCookie})
	return &http.Client{Jar: jar, Timeout: 30 * time.Second}
}

func Upload(client *http.Client, owner, repo string, repoID int, path string) (*Asset, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("stat file: %w", err)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("%s is a directory", path)
	}

	filename := filepath.Base(path)
	contentType, err := detectContentType(path)
	if err != nil {
		return nil, err
	}

	uploadToken, err := getUploadToken(client, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("get upload token: %w", err)
	}

	policy, err := requestPolicy(client, owner, repo, uploadToken, repoID, filename, info.Size(), contentType)
	if err != nil {
		return nil, fmt.Errorf("request policy: %w", err)
	}

	if err := uploadToS3(policy, path, filename); err != nil {
		return nil, fmt.Errorf("upload to storage: %w", err)
	}

	asset, err := finalizeUpload(client, owner, repo, policy)
	if err != nil {
		return nil, fmt.Errorf("finalize upload: %w", err)
	}

	if asset.ContentType == "" {
		asset.ContentType = contentType
	}
	return asset, nil
}

func requestPolicy(client *http.Client, owner, repo, uploadToken string, repoID int, filename string, fileSize int64, contentType string) (*policyResponse, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fields := []struct {
		key   string
		value string
	}{
		{key: "name", value: filename},
		{key: "size", value: strconv.FormatInt(fileSize, 10)},
		{key: "content_type", value: contentType},
		{key: "authenticity_token", value: uploadToken},
		{key: "repository_id", value: strconv.Itoa(repoID)},
	}

	for _, field := range fields {
		if err := writer.WriteField(field.key, field.value); err != nil {
			return nil, fmt.Errorf("write form field %s: %w", field.key, err)
		}
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close multipart writer: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://github.com/upload/policies/assets", body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Origin", "https://github.com")
	req.Header.Set("Referer", fmt.Sprintf("https://github.com/%s/%s", owner, repo))
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("expected 201, got %d: %s", resp.StatusCode, truncate(string(respBody), 200))
	}

	var policy policyResponse
	if err := json.NewDecoder(resp.Body).Decode(&policy); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if policy.UploadURL == "" {
		return nil, fmt.Errorf("missing upload_url")
	}
	if policy.Asset.ID == 0 {
		return nil, fmt.Errorf("missing asset id")
	}
	if policy.AssetUploadAuthenticityToken == "" {
		return nil, fmt.Errorf("missing asset_upload_authenticity_token")
	}

	return &policy, nil
}

func finalizeUpload(client *http.Client, owner, repo string, policy *policyResponse) (*Asset, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if err := writer.WriteField("authenticity_token", policy.AssetUploadAuthenticityToken); err != nil {
		return nil, fmt.Errorf("write authenticity_token: %w", err)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close multipart writer: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("https://github.com/upload/assets/%d", policy.Asset.ID), body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Origin", "https://github.com")
	req.Header.Set("Referer", fmt.Sprintf("https://github.com/%s/%s", owner, repo))
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("expected 200, got %d: %s", resp.StatusCode, truncate(string(respBody), 200))
	}

	var result struct {
		Href        string `json:"href"`
		Name        string `json:"name"`
		ContentType string `json:"content_type"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &Asset{
		URL:         result.Href,
		Filename:    result.Name,
		ContentType: result.ContentType,
	}, nil
}

func detectContentType(path string) (string, error) {
	if byExt := mime.TypeByExtension(filepath.Ext(path)); byExt != "" {
		return byExt, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("open file for content type detection: %w", err)
	}
	defer file.Close()

	header := make([]byte, 512)
	n, err := file.Read(header)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("read file header: %w", err)
	}

	return http.DetectContentType(header[:n]), nil
}

func truncate(input string, max int) string {
	if len(input) <= max {
		return input
	}
	return input[:max] + "..."
}
