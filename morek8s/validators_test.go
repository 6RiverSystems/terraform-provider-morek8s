package morek8s

import "testing"

func TestValidateK8sFile(t *testing.T) {
	invalid := `{
"x": 1,
"y": 2
}`
	ws, es := validateK8sFile(invalid, "data")

	if len(ws) > 0 {
		t.Errorf("Unexpected number of warnings. Expected 0, got %d\nWrnings: %#v", len(ws), ws)
	}

	if len(es) == 0 {
		t.Errorf("Expected validation error but got none")
	}
}
