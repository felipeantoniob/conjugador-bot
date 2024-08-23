package db

import (
	"database/sql"
	"testing"
)

func TestNullStringToString(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		input    sql.NullString
	}{
		{
			name:     "Valid NullString with value",
			expected: "test",
			input:    sql.NullString{String: "test", Valid: true},
		},
		{
			name:     "Valid NullString with empty value",
			expected: "",
			input:    sql.NullString{String: "", Valid: true},
		},
		{
			name:     "Invalid NullString",
			expected: "",
			input:    sql.NullString{String: "test", Valid: false},
		},
		{
			name:     "Invalid NullString with empty value",
			expected: "",
			input:    sql.NullString{String: "", Valid: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NullStringToString(tt.input)
			if got != tt.expected {
				t.Errorf("NullStringToString() = %v, want %v", got, tt.expected)
			}
		})
	}
}
