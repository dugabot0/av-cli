package output

import (
	"encoding/json"
	"fmt"
	"os"
)

// Print writes v as JSON to stdout. If pretty is true, output is indented.
// On marshal error it writes to stderr and exits with code 1.
func Print(v any, pretty bool) {
	var b []byte
	var err error
	if pretty {
		b, err = json.MarshalIndent(v, "", "  ")
	} else {
		b, err = json.Marshal(v)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: marshal JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(b))
}

// Logf writes a formatted message to stderr (used for status/debug output).
func Logf(quiet bool, format string, args ...any) {
	if !quiet {
		fmt.Fprintf(os.Stderr, format+"\n", args...)
	}
}
