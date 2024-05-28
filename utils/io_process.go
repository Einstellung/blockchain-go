package utils

import (
	"fmt"
	"os"
)

// printErr is like fmt.Printf, but writes to stderr.
func PrintErr(m string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, m, args...)
}

// withColor wraps a string with color tags for display in the messages text box.
func WithColor(color, msg string) string {
	return fmt.Sprintf("[%s]%s[-]", color, msg)
}