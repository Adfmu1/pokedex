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
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input: "PriCe iS 	TASTELESS ",
			expected: []string{"price", "is", "tasteless"},
		},
		{
			input: "  OvER mY    	PrecioUs little lamb",
			expected: []string{"over", "my", "precious", "little", "lamb"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(c.expected) != len(actual) {
			t.Errorf("The length of expected string: '%s' is different than the actual string: '%s'", c.expected, actual)
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("The word: '%s' is different thant the expected word: '%s'", word, expectedWord)
			}
		}
	}
}