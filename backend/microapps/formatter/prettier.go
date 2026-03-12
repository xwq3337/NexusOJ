package formatter

import (
	"context"
	"fmt"
)

// PrettierFormatter formats JavaScript/TypeScript using prettier
type PrettierFormatter struct {
	*CommandFormatter
}

// NewPrettierFormatter creates a new prettier formatter
func NewPrettierFormatter() *PrettierFormatter {
	base := NewCommandFormatter(
		"prettier",
		"prettier",
		[]string{"JavaScript", "TypeScript", "JSON", "CSS", "HTML"},
		[]string{}, // Args will be set in Format method
	)

	return &PrettierFormatter{CommandFormatter: base}
}

// Format implements the Formatter interface
func (f *PrettierFormatter) Format(ctx context.Context, code string, options FormatOptions) (string, error) {
	// Determine file extension based on language
	extensions := map[string]string{
		"JavaScript": "js",
		"TypeScript": "ts",
		"JSON":       "json",
		"CSS":        "css",
		"HTML":       "html",
	}

	ext := "js" // default
	if options.Language != "" {
		if e, ok := extensions[options.Language]; ok {
			ext = e
		}
	}

	// Set stdin filepath for prettier (required for stdin mode)
	stdinFilepath := fmt.Sprintf("stdin.%s", ext)
	if options.Filename != "" {
		stdinFilepath = options.Filename
	}

	// Set args for this format call
	f.args = []string{
		"--stdin-filepath", stdinFilepath,
	}

	return f.CommandFormatter.Format(ctx, code, options)
}
