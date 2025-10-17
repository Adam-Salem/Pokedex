package main

import (
	"testing"
)

func TestCleanInput(t *testing.T){
	cases := []struct {
	input    string
	expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input: "pikachu SQUIRTLE bulBAsaur CHARMANDER",
			expected: []string{"pikachu", "squirtle", "bulbasaur", "charmander"},
		},
		{
			input: "WHATS UP BRO I REALLY HOPE YOUR STRING FUNCTION WORKS PROPERLY",
			expected: []string{"whats", "up", "bro", "i", "really", "hope", "your", "string", "function", "works", "properly"},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		// Check the length of the actual slice against the expected slice
		// if they don't match, use t.Errorf to print an error message
		// and fail the test
		if len(actual) != len(c.expected){
			t.Errorf("Mismatched lengths of expected vs actual output")
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			// Check each word in the slice
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
			if word != expectedWord{
				t.Errorf("Word not equal to expected word")
			}
		}
	}
}