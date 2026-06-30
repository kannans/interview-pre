package service

import (
	"fmt"
	"strings"
)

// Supported practice submission languages.
var supportedLanguages = map[string]struct{}{
	"javascript": {},
	"typescript": {},
	"python":     {},
	"go":         {},
	"ruby":       {},
}

const defaultPracticeLanguage = "javascript"

// NormalizePracticeLanguage returns a supported language id or an error.
func NormalizePracticeLanguage(raw string) (string, error) {
	lang := strings.ToLower(strings.TrimSpace(raw))
	if lang == "" {
		return defaultPracticeLanguage, nil
	}
	if _, ok := supportedLanguages[lang]; !ok {
		return "", fmt.Errorf("validation error: language must be one of javascript, typescript, python, go, ruby")
	}
	return lang, nil
}
