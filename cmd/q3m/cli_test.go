package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// buildBinary builds the CLI binary for testing and returns the path.
func buildBinary(t *testing.T) string {
	t.Helper()
	bin := t.TempDir() + "/q3m"
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = "."
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	return bin
}

func runCLI(t *testing.T, bin string, args ...string) (string, string, int) {
	t.Helper()
	cmd := exec.Command(bin, args...)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	code := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			code = exitErr.ExitCode()
		} else {
			t.Fatalf("run %v: %v", args, err)
		}
	}
	return stdout.String(), stderr.String(), code
}

func TestCLIVersion(t *testing.T) {
	bin := buildBinary(t)
	out, _, code := runCLI(t, bin, "--version")
	if code != 0 {
		t.Fatalf("--version exited %d", code)
	}
	if !strings.Contains(out, "q3m version") {
		t.Errorf("--version output = %q, want 'q3m version ...'", out)
	}
}

func TestCLIEncodeText(t *testing.T) {
	bin := buildBinary(t)
	out, _, code := runCLI(t, bin, "encode", "48.8584", "2.2945")
	if code != 0 {
		t.Fatalf("encode exited %d", code)
	}
	trimmed := strings.TrimSpace(out)
	parts := strings.Split(trimmed, ".")
	if len(parts) != 3 {
		t.Errorf("encode output = %q, want 3 dot-separated words", trimmed)
	}
}

func TestCLIEncodeJSON(t *testing.T) {
	bin := buildBinary(t)
	out, _, code := runCLI(t, bin, "encode", "48.8584", "2.2945", "--json")
	if code != 0 {
		t.Fatalf("encode --json exited %d", code)
	}
	var result map[string]any
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out)
	}
	for _, key := range []string{"address", "w1", "w2", "w3", "lat", "lon"} {
		if _, ok := result[key]; !ok {
			t.Errorf("missing key %q in JSON output", key)
		}
	}
}

func TestCLIDecodeText(t *testing.T) {
	bin := buildBinary(t)
	// First encode to get a valid address.
	encOut, _, _ := runCLI(t, bin, "encode", "48.8584", "2.2945")
	addr := strings.TrimSpace(encOut)

	out, _, code := runCLI(t, bin, "decode", addr)
	if code != 0 {
		t.Fatalf("decode exited %d", code)
	}
	if !strings.Contains(out, ",") {
		t.Errorf("decode output = %q, want 'lat, lon'", out)
	}
}

func TestCLIDecodeJSON(t *testing.T) {
	bin := buildBinary(t)
	encOut, _, _ := runCLI(t, bin, "encode", "48.8584", "2.2945")
	addr := strings.TrimSpace(encOut)

	out, _, code := runCLI(t, bin, "decode", addr, "--json")
	if code != 0 {
		t.Fatalf("decode --json exited %d", code)
	}
	var result map[string]any
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out)
	}
	for _, key := range []string{"lat", "lon", "address", "w1", "w2", "w3"} {
		if _, ok := result[key]; !ok {
			t.Errorf("missing key %q in decode JSON output", key)
		}
	}
}

func TestCLIInfoText(t *testing.T) {
	bin := buildBinary(t)
	out, _, code := runCLI(t, bin, "info")
	if code != 0 {
		t.Fatalf("info exited %d", code)
	}
	if !strings.Contains(out, "Lambert93") {
		t.Errorf("info output should mention Lambert93, got: %q", out)
	}
}

func TestCLIInfoJSON(t *testing.T) {
	bin := buildBinary(t)
	out, _, code := runCLI(t, bin, "info", "--json")
	if code != 0 {
		t.Fatalf("info --json exited %d", code)
	}
	var result map[string]any
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out)
	}
	for _, key := range []string{"projection", "total_cells", "dict_size"} {
		if _, ok := result[key]; !ok {
			t.Errorf("missing key %q in info JSON output", key)
		}
	}
}

func TestCLIEncodeOutOfBounds(t *testing.T) {
	bin := buildBinary(t)
	_, stderr, code := runCLI(t, bin, "encode", "0", "0")
	if code == 0 {
		t.Error("encode with out-of-bounds coords should fail")
	}
	if !strings.Contains(stderr, "erreur") {
		t.Errorf("stderr = %q, want error message", stderr)
	}
}

func TestCLIEncodeInvalidArgs(t *testing.T) {
	bin := buildBinary(t)
	_, _, code := runCLI(t, bin, "encode", "abc", "2.0")
	if code == 0 {
		t.Error("encode with invalid lat should fail")
	}
}

func TestCLIDecodeUnknownWord(t *testing.T) {
	bin := buildBinary(t)
	_, stderr, code := runCLI(t, bin, "decode", "xyzzy.hello.world")
	if code == 0 {
		t.Error("decode with unknown words should fail")
	}
	if !strings.Contains(stderr, "erreur") {
		t.Errorf("stderr = %q, want error message", stderr)
	}
}

func TestCLINoArgs(t *testing.T) {
	bin := buildBinary(t)
	out, _, code := runCLI(t, bin)
	if code != 0 {
		t.Errorf("no args should show help, got exit code %d", code)
	}
	if !strings.Contains(out, "q3m") {
		t.Errorf("help output should mention q3m, got: %q", out)
	}
}

func TestMain(m *testing.M) {
	// Ensure we're in the right directory for go build.
	os.Chdir(".")
	os.Exit(m.Run())
}
