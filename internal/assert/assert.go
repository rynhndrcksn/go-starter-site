package assert

import (
	"strings"
	"testing"
)

// Equal is a test helper that takes two comparable types and checks if they are equal.
func Equal[T comparable](t *testing.T, got, want T) {
	// Tells Go's test runner that the Equal() function is a test helper.
	// This means that if t.Errorf() is called the Go test runner will report the
	// filename and line number of the code that called Equal() in the output.
	t.Helper()

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

// StringContains is a test helper that takes two strings and checks if a string contains a substring.
func StringContains(t *testing.T, actual, expectedSubstring string) {
	t.Helper()

	if !strings.Contains(actual, expectedSubstring) {
		t.Errorf("got: %q; expected to contain: %q", actual, expectedSubstring)
	}
}

// NilError is a test helper that takes an error and checks if it's nil or not.
func NilError(t *testing.T, actual error) {
	t.Helper()

	if actual != nil {
		t.Errorf("got: %v; expected: nil", actual)
	}
}
