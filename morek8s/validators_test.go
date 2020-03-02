package morek8s

import "testing"

func TestValidateK8sFile(t *testing.T) {
	invalid := `{
    "apiVersion": "v1",
    "data": {
        "password": "MWYyZDFlMmU2N2Rm",
    },
    "username": "YWRtaW4="
    "kind": "Secret",
    "metadata": {
        "name": "my_embedded_secret",
        "namespace": "default"
    },
    "type": "Opaque"
}`
	ws, es := validateK8sFile(invalid, "data")

	if len(ws) > 0 {
		t.Errorf("Unexpected number of warnings. Expected 0, got %d\nWarnings: %#v", len(ws), ws)
	}

	if len(es) == 0 {
		t.Errorf("Expected validation error but got none")
	}
}
