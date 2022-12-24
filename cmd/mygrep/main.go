package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func eprintf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}

type Group struct {
	chars string
	inverted bool
}

func groupMatches(g Group, ch byte) bool {
	idx := strings.IndexByte(g.chars, ch)
	if g.inverted {
		return idx == -1
	} else {
		return idx != -1
	}
}

func parseGroups(pattern string) ([]Group, string) {
	groups := make([]Group, 0, 16)
	cur_group := Group{"", true}
	in_backslash := false
	in_group_start := false
	in_group := false
	for i := range pattern {
		ch := pattern[i]
		if in_backslash {
			if ch == 'd' {
				cur_group.chars += "0123456789"
			} else if ch == 'w' {
				cur_group.chars += "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_"
			} else {
				cur_group.chars += string(ch);
			}
			if !in_group {
				groups = append(groups, cur_group)
				cur_group = Group{"", true}
			} else {
				in_group_start = false
			}
		} else {
			if ch == '[' {
				if in_group {
					return nil, "can't have nested character groups"
				}
				in_group = true
				in_group_start = true
			} else if in_group && ch == ']' {
				if in_group_start {
					return nil, "empty character group is invalid"
				}
				in_group = false
				groups = append(groups, cur_group)
				cur_group = Group{"", true}
			} else if in_group_start && ch == '^' {
				cur_group.inverted = true
				in_group_start = false
			} else {
				cur_group.chars += string(ch);
				if !in_group {
					groups = append(groups, cur_group)
					cur_group = Group{"", true}
				} else {
					in_group_start = false
				}
			}
		}
	}
	if in_group {
		return nil, "unterminated group"
	}
	return groups, ""
}

// Usage: echo <input_text> | your_grep.sh -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		eprintf("usage: %s -E <pattern>\n", os.Args[0])
		os.Exit(2)
	}

	pattern := os.Args[2]
	groups, errstr := parseGroups(pattern)
	if errstr != "" {
		eprintf("can't parse pattern: %s\n", errstr)
		os.Exit(3)
	}

	input_bytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		eprintf("can't read from stdin: %v\n", err)
		os.Exit(4)
	}

	input := string(input_bytes)
	if len(input) < len(groups) {
		os.Exit(1)
	}
	input_max := len(input) - len(groups)
	outer: for i := 0; i < input_max; i++ {
		for j := range groups {
			if !groupMatches(groups[j], input[i+j]) {
				continue outer;
			}
		}
		os.Exit(0)
	}
	os.Exit(1)
}
