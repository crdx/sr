package main

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubstitution_Replace_Fixed(t *testing.T) {
	tests := []struct {
		name        string
		fixedString string
		replacement string
		input       string
		expected    string
	}{
		{
			name:        "simple replacement",
			fixedString: "foo",
			replacement: "bar",
			input:       "foo baz foo",
			expected:    "bar baz bar",
		},
		{
			name:        "no match",
			fixedString: "foo",
			replacement: "bar",
			input:       "baz qux",
			expected:    "baz qux",
		},
		{
			name:        "regex metacharacters are literal",
			fixedString: "foo.*bar",
			replacement: "replaced",
			input:       "foo.*bar baz",
			expected:    "replaced baz",
		},
		{
			name:        "dollar sign in pattern is literal",
			fixedString: "$100",
			replacement: "£100",
			input:       "costs $100",
			expected:    "costs £100",
		},
		{
			name:        "replacement is literal (no capture groups)",
			fixedString: "foo",
			replacement: "$1",
			input:       "foo bar",
			expected:    "$1 bar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub := &substitution{
				fixedString: tt.fixedString,
				replacement: tt.replacement,
				isFixed:     true,
			}
			assert.Equal(t, tt.expected, sub.replace(tt.input))
		})
	}
}

func TestSubstitution_Replace_Regex(t *testing.T) {
	tests := []struct {
		name        string
		pattern     string
		replacement string
		input       string
		expected    string
	}{
		{
			name:        "simple regex",
			pattern:     "foo",
			replacement: "bar",
			input:       "foo baz foo",
			expected:    "bar baz bar",
		},
		{
			name:        "capture group",
			pattern:     `(\w+): (\d+)`,
			replacement: "$2 $1",
			input:       "count: 42",
			expected:    "42 count",
		},
		{
			name:        "wildcard",
			pattern:     "foo.*bar",
			replacement: "replaced",
			input:       "foo123bar baz",
			expected:    "replaced baz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub := &substitution{
				pattern:     regexp.MustCompile(tt.pattern),
				replacement: tt.replacement,
				isFixed:     false,
			}
			assert.Equal(t, tt.expected, sub.replace(tt.input))
		})
	}
}
