package formatter

// GofmtFormatter formats Go code using gofmt
type GofmtFormatter struct {
	*CommandFormatter
}

// NewGofmtFormatter creates a new gofmt formatter
func NewGofmtFormatter() *GofmtFormatter {
	base := NewCommandFormatter(
		"gofmt",
		"gofmt",
		[]string{"go"},
		[]string{}, // gofmt reads from stdin by default
	)

	return &GofmtFormatter{CommandFormatter: base}
}
