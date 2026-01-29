package utils

import "testing"

func TestGreet(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Test with name",
			input:    "Alice",
			expected: "Hello, Alice! Welcome to home-ctrl!",
		},
		{
			name:     "Test with empty name",
			input:    "",
			expected: "Hello, ! Welcome to home-ctrl!",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Greet(tc.input)
			if result != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}