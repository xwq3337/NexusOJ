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
		[]string{"c", "cpp", "java", "javascript", "objc"},
		[]string{}, // Args will be set in Format method
	)

	return &ClangFormatFormatter{CommandFormatter: base}
}

// Format implements the Formatter interface
func (f *ClangFormatFormatter) Format(ctx context.Context, code string, options FormatOptions) (string, error) {
	// Determine file extension based on language
	extensions := map[string]string{
		"c":          "c",
		"cpp":        "cpp",
		"java":       "java",
		"javascript": "js",
		"objc":       "m",
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
	// Google LLVM Chromium Mozilla WebKit Microsoft GNU
	f.args = []string{
		"-style=Microsoft",
		fmt.Sprintf("-assume-filename=%s", filename),
	}

	return f.CommandFormatter.Format(ctx, code, options)
}
