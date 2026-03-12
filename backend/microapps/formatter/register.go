package formatter

// RegisterAllFormatters registers all built-in formatters to the given registry
func RegisterAllFormatters(registry *Registry) {
	registry.Register(NewClangFormatFormatter())
	registry.Register(NewBlackFormatter())
	registry.Register(NewGofmtFormatter())
	registry.Register(NewPrettierFormatter())
	registry.Register(NewRustfmtFormatter())
}
