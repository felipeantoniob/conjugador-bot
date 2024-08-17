package utils

import (
	"errors"
	"testing"
)

func TestWrapError(t *testing.T) {
	// Define the base error to wrap
	baseErr := errors.New("base error")

	tests := []struct {
		name    string
		message string
		err     error
		want    string
	}{
		{
			name:    "Wrap base error with message",
			message: "an error occurred",
			err:     baseErr,
			want:    "an error occurred: base error",
		},
		{
			name:    "Wrap nil error",
			message: "a nil error",
			err:     nil,
			want:    "a nil error: %!w(<nil>)",
		},
		{
			name:    "Wrap empty message",
			message: "",
			err:     baseErr,
			want:    ": base error",
		},
		{
			name:    "Wrap error with special characters in message",
			message: "error: !@#$%^&*()",
			err:     baseErr,
			want:    "error: !@#$%^&*(): base error",
		},
		{
			name:    "Wrap error with empty base error message",
			message: "error",
			err:     errors.New(""),
			want:    "error: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WrapError(tt.message, tt.err)
			if got.Error() != tt.want {
				t.Errorf("got %v, want %v", got.Error(), tt.want)
			}
		})
	}
}
