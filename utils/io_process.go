package utils

import (
	"fmt"
	"os"
)

// printErr is like fmt.Printf, but writes to stderr.
func PrintErr(m string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, m, args...)
}

// colorPrint prints the first part of the message in the specified color
// and the rest of the parts in the default terminal color.
func ColorPrint(color string, parts ...string) {
    // Define ANSI escape codes for colors
    var colorCode string
    reset := "\033[0m"

    // Select color code based on the input color
    switch color {
    case "black":
        colorCode = "\033[30m"
    case "red":
        colorCode = "\033[31m"
    case "green":
        colorCode = "\033[32m"
    case "yellow":
        colorCode = "\033[33m"
    case "blue":
        colorCode = "\033[34m"
    case "magenta":
        colorCode = "\033[35m"
    case "cyan":
        colorCode = "\033[36m"
    case "white":
        colorCode = "\033[37m"
    default:
        // Default to terminal's default color if the color is not recognized
        colorCode = reset
    }

    if len(parts) > 0 {
        // Print the first part with the selected color
        fmt.Fprint(os.Stdout, colorCode+parts[0]+reset)
    }

    // Print the rest of the parts with the default terminal color
    for i := 1; i < len(parts); i++ {
        fmt.Fprint(os.Stdout, parts[i])
    }

    // Add a new line at the end
    fmt.Fprintln(os.Stdout)
}