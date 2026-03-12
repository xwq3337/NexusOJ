package formatter

import (
	"context"
	"fmt"
)

// ClangFormatFormatter formats C/C++/Java/JavaScript using clang-format
type ClangFormatFormatter struct {
	*CommandFormatter
}

// NewClangFormatFormatter creates a new clang-format formatter
func NewClangFormatFormatter() *ClangFormatFormatter {
	base := NewCommandFormatter(
		"clang-format",
		"clang-format",
		[]string{"C", "C++", "Java", "JavaScript", "Objective-C"},
		[]string{}, // Args will be set in Format method
	)

	return &ClangFormatFormatter{CommandFormatter: base}
}

// Format implements the Formatter interface
func (f *ClangFormatFormatter) Format(ctx context.Context, code string, options FormatOptions) (string, error) {
	// Determine file extension based on language
	extensions := map[string]string{
		"C":           "c",
		"C++":         "cpp",
		"Java":        "java",
		"JavaScript":  "js",
		"Objective-C": "m",
	}

	ext := "cpp" // default
	if options.Language != "" {
		if e, ok := extensions[options.Language]; ok {
			ext = e
		}
	}

	// Set filename for clang-format
	filename := fmt.Sprintf("stdin.%s", ext)
	if options.Filename != "" {
		filename = options.Filename
	}

	// Set args for this format call
	f.args = []string{
		"-style=Google",
		fmt.Sprintf("-assume-filename=%s", filename),
	}

	return f.CommandFormatter.Format(ctx, code, options)
}
