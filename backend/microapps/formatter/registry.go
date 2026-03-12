package formatter

import (
	"sort"
	"sync"
)

// Registry manages available formatters
type Registry struct {
	sync.RWMutex
	formatters map[string]Formatter   // by name
	byLanguage map[string][]Formatter // by language
}

// NewRegistry creates a new formatter registry
func NewRegistry() *Registry {
	r := &Registry{
		formatters: make(map[string]Formatter),
		byLanguage: make(map[string][]Formatter),
	}

	// Register all built-in formatters
	RegisterAllFormatters(r)

	return r
}

// Register adds a formatter to the registry
func (r *Registry) Register(f Formatter) {
	r.Lock()
	defer r.Unlock()

	r.formatters[f.Name()] = f

	for _, lang := range f.Languages() {
		r.byLanguage[lang] = append(r.byLanguage[lang], f)
	}
}

// GetByName retrieves a formatter by name
func (r *Registry) GetByName(name string) (Formatter, bool) {
	r.RLock()
	defer r.RUnlock()

	f, ok := r.formatters[name]
	return f, ok
}

// GetByLanguage retrieves a formatter for a specific language
func (r *Registry) GetByLanguage(language string) (Formatter, bool) {
	r.RLock()
	defer r.RUnlock()

	formatters, ok := r.byLanguage[language]
	if !ok || len(formatters) == 0 {
		return nil, false
	}

	// Return the first available formatter
	return formatters[0], true
}

// SupportedLanguages returns a sorted list of all supported languages
func (r *Registry) SupportedLanguages() []string {
	r.RLock()
	defer r.RUnlock()

	languages := make([]string, 0, len(r.byLanguage))
	for lang := range r.byLanguage {
		languages = append(languages, lang)
	}
	sort.Strings(languages)
	return languages
}

// AllFormatters returns all registered formatters
func (r *Registry) AllFormatters() []Formatter {
	r.RLock()
	defer r.RUnlock()

	formatters := make([]Formatter, 0, len(r.formatters))
	for _, f := range r.formatters {
		formatters = append(formatters, f)
	}
	return formatters
}
