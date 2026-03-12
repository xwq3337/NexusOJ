package formatter

import "fmt"

// Base error types
var (
	// ErrLanguageNotDetected is returned when language cannot be detected
	ErrLanguageNotDetected = fmt.Errorf("language could not be detected")
	// ErrEmptyCode is returned when input code is empty
	ErrEmptyCode = fmt.Errorf("code is empty")
	// ErrFormatterNotFound is returned when no formatter is found for a language
	ErrFormatterNotFound = fmt.Errorf("formatter not found")
)

// FormatterNotAvailableError is returned when no formatter is available for a language
type FormatterNotAvailableError struct {
	Language string
	Hint     string
}

func (e *FormatterNotAvailableError) Error() string {
	msg := fmt.Sprintf("no formatter available for language: %s", e.Language)
	if e.Hint != "" {
		msg += fmt.Sprintf(" (supported languages: %s)", e.Hint)
	}
	return msg
}

// FormatterBinaryNotFoundError is returned when the formatter binary is not found
type FormatterBinaryNotFoundError struct {
	Formatter   string
	Language    string
	InstallHint string
}

func (e *FormatterBinaryNotFoundError) Error() string {
	msg := fmt.Sprintf("formatter binary '%s' not found for language: %s", e.Formatter, e.Language)
	if e.InstallHint != "" {
		msg += fmt.Sprintf("\nInstallation: %s", e.InstallHint)
	}
	return msg
}

// FormatError wraps formatting errors with context
type FormatError struct {
	Formatter string
	Language  string
	Code      string
	Err       error
}

func (e *FormatError) Error() string {
	msg := fmt.Sprintf("format error with %s for %s", e.Formatter, e.Language)
	if e.Code != "" {
		// Show first 100 characters of code for debugging
		codePreview := e.Code
		if len(codePreview) > 100 {
			codePreview = codePreview[:100] + "..."
		}
		msg += fmt.Sprintf("\nCode preview: %s", codePreview)
	}
	msg += fmt.Sprintf("\nError: %v", e.Err)
	return msg
}

func (e *FormatError) Unwrap() error {
	return e.Err
}
