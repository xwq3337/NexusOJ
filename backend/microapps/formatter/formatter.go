package formatter

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// Formatter is the core interface for code formatters
type Formatter interface {
	// Format formats the given code and returns the result
	Format(ctx context.Context, code string, options FormatOptions) (string, error)
	// Name returns the formatter name
	Name() string
	// Languages returns the list of supported languages
	Languages() []string
	// Available checks if the formatter binary is available
	Available() bool
}

// FormatOptions contains formatting options
type FormatOptions struct {
	// Filename is used for language detection and formatter-specific config
	Filename string
	// Language overrides auto-detection if specified
	Language string
	// FormatterArgs are additional arguments to pass to the formatter
	FormatterArgs []string
}

// FormatResult contains the result of formatting
type FormatResult struct {
	Formatted string
	Language  string
	Formatter string
}

// Service is the main formatting service
type Service struct {
	registry *Registry
	detector *LanguageDetector
	timeout  time.Duration
}

// NewService creates a new formatting service
func NewService() *Service {
	return &Service{
		registry: NewRegistry(),
		detector: NewLanguageDetector(),
		timeout:  30 * time.Second,
	}
}

// SetTimeout sets the timeout for formatter execution
func (s *Service) SetTimeout(timeout time.Duration) {
	s.timeout = timeout
}

// FormatCode is the main entry point - matches the required signature
// This function formats code and returns the formatted result.
// If language is empty, it will auto-detect the language.
func (s *Service) FormatCode(code string, language string) (string, error) {
	return s.FormatCodeWithContext(context.Background(), code, language)
}

// FormatCodeWithContext formats code with context support
func (s *Service) FormatCodeWithContext(ctx context.Context, code string, language string) (string, error) {
	// 1. Input validation
	if strings.TrimSpace(code) == "" {
		return "", fmt.Errorf("%w: code is empty", ErrEmptyCode)
	}

	// 2. Language detection
	detectedLang := language
	if detectedLang == "" {
		detectedLang = s.detector.DetectLanguage(code)
		if detectedLang == "" {
			return "", fmt.Errorf("%w: could not detect language from code", ErrLanguageNotDetected)
		}
	}

	// 3. Formatter selection
	formatter, ok := s.registry.GetByLanguage(detectedLang)
	if !ok {
		supportedLangs := s.registry.SupportedLanguages()
		return "", &FormatterNotAvailableError{
			Language: detectedLang,
			Hint:     fmt.Sprintf("%v", supportedLangs),
		}
	}

	// 4. Availability check
	if !formatter.Available() {
		return "", &FormatterBinaryNotFoundError{
			Formatter:   formatter.Name(),
			Language:    detectedLang,
			InstallHint: getInstallHint(formatter.Name()),
		}
	}

	// 5. Execution with timeout
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	result, err := formatter.Format(ctx, code, FormatOptions{
		Language: detectedLang,
	})

	// 6. Error wrapping with context
	if err != nil {
		return "", &FormatError{
			Formatter: formatter.Name(),
			Language:  detectedLang,
			Code:      truncateCode(code, 100),
			Err:       err,
		}
	}

	// 7. Output validation
	if strings.TrimSpace(result) == "" {
		return "", fmt.Errorf("formatter returned empty output")
	}

	return result, nil
}

// getInstallHint returns installation instructions for formatters
func getInstallHint(formatter string) string {
	hints := map[string]string{
		"clang-format": "Install via: brew install clang-format (macOS) or apt install clang-format (Ubuntu)",
		"black":        "Install via: pip install black",
		"gofmt":        "Install Go: https://golang.org/dl/",
		"prettier":     "Install via: npm install -g prettier",
		"rustfmt":      "Install via: rustup component add rustfmt",
	}

	if hint, ok := hints[formatter]; ok {
		return hint
	}
	return fmt.Sprintf("Install %s", formatter)
}

// truncateCode truncates code to a specified length for error messages
func truncateCode(code string, maxLen int) string {
	if len(code) <= maxLen {
		return code
	}
	return code[:maxLen] + "..."
}
