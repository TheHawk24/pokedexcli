package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " Hello world ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "HELLO THERE MY GUY ",
			expected: []string{"hello", "there", "my", "guy"},
		},
		{
			input:    "   Check ",
			expected: []string{"check"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Expected Length: %v\nActual Length: %v", len(c.expected), len(actual))
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Expected: %v\nActual: %v", expectedWord, word)
			}
		}
	}
}
