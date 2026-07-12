package app

import "testing"

func TestDoctorIsRecognizedAsACommand(t *testing.T) {
	if !isDoctorCommand([]string{"doctor"}) {
		t.Fatal("expected doctor command to be recognized")
	}
	if isDoctorCommand([]string{"screenshot.png"}) {
		t.Fatal("expected upload shortcut not to be recognized as doctor")
	}
}
