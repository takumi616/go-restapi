package helper

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func AssertJson(t *testing.T, expected, actual []byte) {
	t.Helper()

	var jsonExpected, jsonActual any
	if err := json.Unmarshal(expected, &jsonExpected); err != nil {
		t.Fatalf("Failed to unmarshal expected %q: %v", expected, err)
	}
	if err := json.Unmarshal(actual, &jsonActual); err != nil {
		t.Fatalf("Failed to unmarshal actual %q: %v", actual, err)
	}

	if diff := cmp.Diff(jsonActual, jsonExpected); diff != "" {
		t.Errorf("found differs: (-actual +expected)\n%s", diff)
	}
}

func AssertResponse(t *testing.T, actual *http.Response, expectedStatus int, expectedBody []byte) {
	t.Helper()

	t.Cleanup(func() { _ = actual.Body.Close() })
	actualBody, err := io.ReadAll(actual.Body)
	if err != nil {
		t.Fatal(err)
	}

	if actual.StatusCode != expectedStatus {
		t.Fatalf("expected status %d, but actual %d, body: %q", expectedStatus, actual.StatusCode, actualBody)
	}

	if len(actualBody) == 0 && len(expectedBody) == 0 {
		return
	}

	AssertJson(t, expectedBody, actualBody)
}
