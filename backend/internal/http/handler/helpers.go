package handler

import "fmt"

// isValidationError checks if an error originated from input validation.
// Shared utility for all handlers in this package.
func isValidationError(err error) bool {
	msg := fmt.Sprintf("%s", err)
	if len(msg) < 16 {
		return false
	}
	return msg[:16] == "validation error"
}
