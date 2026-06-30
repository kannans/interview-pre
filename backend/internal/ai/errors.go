package ai

import "errors"

var (
	ErrNotConfigured   = errors.New("openai api key is not configured")
	ErrEmptyResponse   = errors.New("openai returned empty response")
	ErrInvalidJSON     = errors.New("openai response is not valid json")
	ErrValidation      = errors.New("openai response failed schema validation")
	ErrNotImplemented  = errors.New("ai capability not implemented yet")
)
