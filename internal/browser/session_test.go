package browser

import (
	"errors"
	"strings"
	"testing"
)

func TestExplainSessionErrorForMacOSKeychainFailure(t *testing.T) {
	err := errors.New("decrypting cookie user_session=: keyring password retrieval failed: error reading 'Chrome Safe Storage' keychain password: The user name or passphrase you entered is not correct. (-25293)")

	message := explainSessionError(err)

	for _, want := range []string{
		"macOS Keychain",
		"GH_ATTACH_USER_SESSION",
		"gh attach doctor",
	} {
		if !strings.Contains(message, want) {
			t.Fatalf("expected %q in message:\n%s", want, message)
		}
	}
}

func TestExplainSessionErrorForMissingBrowserSession(t *testing.T) {
	message := explainSessionError(errors.New("no github.com user_session cookie found"))

	if !strings.Contains(message, "Sign in to github.com") {
		t.Fatalf("expected sign-in guidance, got:\n%s", message)
	}
	if !strings.Contains(message, "GH_ATTACH_USER_SESSION") {
		t.Fatalf("expected environment fallback, got:\n%s", message)
	}
}

func TestFormatDiagnosticDoesNotRevealSessionValue(t *testing.T) {
	diagnostic := Diagnostic{
		EnvironmentOverride: true,
		SessionFound:        true,
		Source:              "GH_ATTACH_USER_SESSION",
	}

	output := diagnostic.String()

	if !strings.Contains(output, "Environment override: configured") {
		t.Fatalf("expected configured status, got:\n%s", output)
	}
	if strings.Contains(output, "user_session=") {
		t.Fatalf("diagnostic must not reveal cookie material:\n%s", output)
	}
}

func TestExplainSessionErrorRedactsCookieMaterial(t *testing.T) {
	err := errors.New("decrypting cookie user_session=super-secret-value; Path=/; Domain=github.com: keychain failed")

	message := explainSessionError(err)

	if strings.Contains(message, "super-secret-value") {
		t.Fatalf("error guidance must redact cookie material:\n%s", message)
	}
	if !strings.Contains(message, "user_session=[REDACTED]") {
		t.Fatalf("expected explicit redaction marker:\n%s", message)
	}
}
