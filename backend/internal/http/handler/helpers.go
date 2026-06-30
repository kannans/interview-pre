package handler

import (
	"errors"
	"fmt"

	"github.com/kannans/interview-prep/backend/internal/service"
)

// isValidationError checks if an error originated from input validation.
func isValidationError(err error) bool {
	msg := fmt.Sprintf("%s", err)
	if len(msg) < 16 {
		return false
	}
	return msg[:16] == "validation error"
}

// isAIError checks if an error originated from the AI service layer.
func isAIError(err error) bool {
	return errors.Is(err, service.ErrAIFailed)
}
