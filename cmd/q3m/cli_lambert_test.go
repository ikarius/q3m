package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestCLITolamText(t *testing.T) {
	bin := buildBinary(t)
	out, _, code := runCLI(t, bin, "tolam", "48.8584", "2.2945")
	if code != 0 {
		t.Fatalf("tolam exited %d", code)
	}
	if !strings.Contains(out, ",") {
		t.Errorf("tolam output = %q, want 'E, N'", out)
	}
}

func TestCLITolamJSON(t *testing.T) {
	bin := buildBinary(t)
	out, _, code := runCLI(t, bin, "tolam", "48.8584", "2.2945", "--json")
	if code != 0 {
		t.Fatalf("tolam --json exited %d", code)
	}
	var result map[string]any
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out)
	}
	for _, key := range []string{"e", "n", "lat", "lon"} {
		if _, ok := result[key]; !ok {
			t.Errorf("missing key %q in tolam JSON output", key)
		}
	}
}

func TestCLIFromlamText(t *testing.T) {
	bin := buildBinary(t)
	// Encode Tour Eiffel to Lambert93 first.
	tolamOut, _, _ := runCLI(t, bin, "tolam", "48.8584", "2.2945")
	parts := strings.Split(strings.TrimSpace(tolamOut), ", ")
	if len(parts) != 2 {
		t.Fatalf("tolam output = %q, want 'E, N'", tolamOut)
	}

	out, _, code := runCLI(t, bin, "fromlam", parts[0], parts[1])
	if code != 0 {
		t.Fatalf("fromlam exited %d", code)
	}
	if !strings.Contains(out, "48.8584") {
		t.Errorf("fromlam output = %q, want latitude ~48.8584", out)
	}
}

func TestCLIFromlamJSON(t *testing.T) {
	bin := buildBinary(t)
	// Encode Tour Eiffel to Lambert93 first.
	tolamOut, _, _ := runCLI(t, bin, "tolam", "48.8584", "2.2945")
	parts := strings.Split(strings.TrimSpace(tolamOut), ", ")
	if len(parts) != 2 {
		t.Fatalf("tolam output = %q, want 'E, N'", tolamOut)
	}

	out, _, code := runCLI(t, bin, "fromlam", parts[0], parts[1], "--json")
	if code != 0 {
		t.Fatalf("fromlam --json exited %d", code)
	}
	var result map[string]any
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out)
	}
	for _, key := range []string{"lat", "lon", "e", "n"} {
		if _, ok := result[key]; !ok {
			t.Errorf("missing key %q in fromlam JSON output", key)
		}
	}
}

func TestCLITolamInvalidArgs(t *testing.T) {
	bin := buildBinary(t)
	_, _, code := runCLI(t, bin, "tolam", "abc", "2.0")
	if code == 0 {
		t.Error("tolam with invalid lat should fail")
	}
}

func TestCLIFromlamInvalidArgs(t *testing.T) {
	bin := buildBinary(t)
	_, _, code := runCLI(t, bin, "fromlam", "abc", "6862047")
	if code == 0 {
		t.Error("fromlam with invalid E should fail")
	}
}

func TestCLITolamMissingArgs(t *testing.T) {
	bin := buildBinary(t)
	_, _, code := runCLI(t, bin, "tolam", "48.8584")
	if code == 0 {
		t.Error("tolam with missing arg should fail")
	}
}

func TestCLIFromlamMissingArgs(t *testing.T) {
	bin := buildBinary(t)
	_, _, code := runCLI(t, bin, "fromlam")
	if code == 0 {
		t.Error("fromlam with no args should fail")
	}
}
