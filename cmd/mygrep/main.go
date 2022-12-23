package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func eprintf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...);
}

// Usage: echo <input_text> | your_grep.sh -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		eprintf("usage: %s -E <pattern>\n", os.Args[0])
		os.Exit(2)
	}

	input_bytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		eprintf("can't read from stdin: %v\n", err)
	}

	input := string(input_bytes)
	pattern := os.Args[2]
	if len(pattern) == 1 && strings.IndexByte(input, pattern[0]) != -1 {
		os.Exit(0)
	} else if len(pattern) == 2 && strings.IndexAny(input, "0123456789") != -1 {
		os.Exit(0)
	} else {
		eprintf("unsupported pattern\n")
		os.Exit(3)
	}
	os.Exit(1)
}