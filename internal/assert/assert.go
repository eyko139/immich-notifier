package assert

import (
    "testing"
    "strings"
)

func AssertEqual[T comparable](t *testing.T, actual, expected T) {
    t.Helper()

    if actual != expected {
        t.Errorf("Expected %v, got %v", expected, actual)
    }
}

func StringContains(t *testing.T, actual, expected string) {
    t.Helper()

    if !strings.Contains(actual, expected) {
        t.Errorf("Expected %v, got %v", expected, actual)
    }
}
