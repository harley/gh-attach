package app

import (
	"strings"
	"testing"
)

func TestDoctorIsRecognizedAsACommand(t *testing.T) {
	if !isDoctorCommand([]string{"doctor"}) {
		t.Fatal("expected doctor command to be recognized")
	}
	if !isDoctorCommand([]string{"doctor", "unexpected"}) {
		t.Fatal("expected doctor invocation with extra arguments to reach doctor validation")
	}
	if isDoctorCommand([]string{"screenshot.png"}) {
		t.Fatal("expected upload shortcut not to be recognized as doctor")
	}
}

func TestDoctorRejectsExtraArguments(t *testing.T) {
	err := Run([]string{"doctor", "unexpected"})
	if err == nil || !strings.Contains(err.Error(), "does not accept arguments") {
		t.Fatalf("expected a doctor argument error, got %v", err)
	}
}
