package helper

import (
	"os"
	"testing"
)

func LoadFile(t *testing.T, path string) []byte {
	t.Helper()

	bt, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read from %q: %v", path, err)
	}
	return bt
}
