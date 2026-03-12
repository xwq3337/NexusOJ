package formatter

import (
	"os/exec"
)

// BlackFormatter formats Python code using black
type BlackFormatter struct {
	*CommandFormatter
}

// NewBlackFormatter creates a new black formatter
func NewBlackFormatter() *BlackFormatter {
	base := NewCommandFormatter(
		"black",
		"black",
		[]string{"python"},
		[]string{"-", "-"}, // Read from stdin, output to stdout
	)

	return &BlackFormatter{CommandFormatter: base}
}

// Available checks if black is available
func (f *BlackFormatter) Available() bool {
	// Black might be installed as 'black' or 'black3'
	if f.CommandFormatter.Available() {
		return true
	}

	// Try alternative names
	for _, alt := range []string{"black3"} {
		_, err := exec.LookPath(alt)
		if err == nil {
			f.command = alt
			return true
		}
	}

	// Reset to default
	f.command = "black"
	return false
}
