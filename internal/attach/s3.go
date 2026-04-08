package attach

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

func uploadToS3(policy *policyResponse, path, filename string) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	orderedKeys := []string{
		"key",
		"acl",
		"policy",
		"X-Amz-Algorithm",
		"X-Amz-Credential",
		"X-Amz-Date",
		"X-Amz-Signature",
		"Content-Type",
		"Cache-Control",
		"x-amz-meta-Surrogate-Control",
	}

	written := make(map[string]bool, len(orderedKeys))
	for _, key := range orderedKeys {
		value, ok := policy.Form[key]
		if !ok {
			continue
		}
		if err := writer.WriteField(key, value); err != nil {
			return fmt.Errorf("write form field %s: %w", key, err)
		}
		written[key] = true
	}

	for key, value := range policy.Form {
		if written[key] {
			continue
		}
		if err := writer.WriteField(key, value); err != nil {
			return fmt.Errorf("write form field %s: %w", key, err)
		}
	}

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return fmt.Errorf("create file form field: %w", err)
	}

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("copy file into multipart body: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("close multipart writer: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, policy.UploadURL, body)
	if err != nil {
		return fmt.Errorf("create upload request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Origin", "https://github.com")
	req.Header.Set("User-Agent", userAgent)

	resp, err := (&http.Client{Timeout: 120 * time.Second}).Do(req)
	if err != nil {
		return fmt.Errorf("execute upload request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("storage returned %d: %s", resp.StatusCode, truncate(string(respBody), 300))
	}

	return nil
}
