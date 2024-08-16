package utils

import "fmt"

// WrapError formats and wraps an error message with additional context.
func WrapError(message string, err error) error {
	return fmt.Errorf("%s: %w", message, err)
}
