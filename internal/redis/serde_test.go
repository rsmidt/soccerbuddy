package redis

import "testing"

func TestFlattenToString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Empty string returns empty string.
		{"", ""},
		// Non-array string returns itself.
		{"something", "something"},
		// Valid JSON array returns the first element.
		{`["wrapped","in","array"]`, "wrapped"},
		{`["onlyone"]`, "onlyone"},
		// Valid JSON array but empty array should return empty string.
		{`[]`, ""},
		// Malformed JSON starting with '[' should be treated as a normal string.
		{"[not, valid, json]", "[not, valid, json]"},
	}

	for _, tt := range tests {
		got := FlattenToString(tt.input)
		if got != tt.expected {
			t.Errorf("FlattenToString(%q) = %q; want %q", tt.input, got, tt.expected)
		}
	}
}
