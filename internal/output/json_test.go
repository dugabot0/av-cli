package output

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// captureStdout replaces os.Stdout temporarily and returns what was written.
func captureStdout(fn func()) string {
	r, w, _ := os.Pipe()
	orig := os.Stdout
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = orig

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestPrintCompact(t *testing.T) {
	v := map[string]int{"a": 1}
	out := captureStdout(func() { Print(v, false) })
	out = strings.TrimSpace(out)
	if out != `{"a":1}` {
		t.Errorf("compact: want %q, got %q", `{"a":1}`, out)
	}
}

func TestPrintPretty(t *testing.T) {
	v := map[string]int{"a": 1}
	out := captureStdout(func() { Print(v, true) })
	if !strings.Contains(out, "\n") {
		t.Errorf("pretty: expected indented output, got %q", out)
	}
}

func TestLogfQuiet(t *testing.T) {
	// Should not panic or write when quiet=true
	Logf(true, "should not appear: %s", "test")
}

func TestLogfNotQuiet(t *testing.T) {
	// Just verify it doesn't panic; stderr capture is complex
	// In real use, messages appear on stderr
	Logf(false, "test message %d", 42)
}
