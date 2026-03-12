package formatter

// RustfmtFormatter formats Rust code using rustfmt
type RustfmtFormatter struct {
	*CommandFormatter
}

// NewRustfmtFormatter creates a new rustfmt formatter
func NewRustfmtFormatter() *RustfmtFormatter {
	base := NewCommandFormatter(
		"rustfmt",
		"rustfmt",
		[]string{"rust"},
		[]string{"--emit", "stdout"}, // Emit to stdout instead of modifying in-place
	)

	return &RustfmtFormatter{CommandFormatter: base}
}
