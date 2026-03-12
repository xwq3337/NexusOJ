package formatter

import (
	"strings"

	"github.com/go-enry/go-enry/v2"
)

// LanguageDetector handles programming language detection
type LanguageDetector struct {
	// Could add caching, custom heuristics, etc.
}

// NewLanguageDetector creates a new language detector
func NewLanguageDetector() *LanguageDetector {
	return &LanguageDetector{}
}

// DetectLanguage detects the programming language from code content
func (d *LanguageDetector) DetectLanguage(code string) string {
	// If code is empty, return empty
	if strings.TrimSpace(code) == "" {
		return ""
	}

	// Try shebang detection first
	if lang := d.detectByShebang(code); lang != "" {
		return lang
	}

	// Try content-based detection using enry
	// enry.GetLanguage uses content heuristics
	lang := enry.GetLanguage("", []byte(code))
	if lang != "" {
		return d.normalizeLanguage(lang)
	}

	return ""
}

// detectByShebang detects language from shebang line
func (d *LanguageDetector) detectByShebang(code string) string {
	lines := strings.Split(code, "\n")
	if len(lines) == 0 {
		return ""
	}

	firstLine := strings.TrimSpace(lines[0])
	if !strings.HasPrefix(firstLine, "#!") {
		return ""
	}

	// Common shebang patterns
	shebangMap := map[string]string{
		"python":  "Python",
		"python3": "Python",
		"python2": "Python",
		"bash":    "Shell",
		"sh":      "Shell",
		"node":    "JavaScript",
		"ruby":    "Ruby",
		"perl":    "Perl",
		"php":     "PHP",
	}

	for pattern, lang := range shebangMap {
		if strings.Contains(firstLine, pattern) {
			return lang
		}
	}

	return ""
}

// normalizeLanguage normalizes language names to our supported set
func (d *LanguageDetector) normalizeLanguage(lang string) string {
	// Map enry language names to our supported languages
	normalization := map[string]string{
		"C":           "C",
		"C++":         "C++",
		"Java":        "Java",
		"Python":      "Python",
		"Go":          "Go",
		"JavaScript":  "JavaScript",
		"TypeScript":  "JavaScript", // Use Prettier for both
		"Rust":        "Rust",
		"Objective-C": "C", // Use clang-format
		"C#":          "C", // Use clang-format
	}

	if normalized, ok := normalization[lang]; ok {
		return normalized
	}

	return lang
}

// DetectLanguageWithFilename detects language using both filename and content
func (d *LanguageDetector) DetectLanguageWithFilename(filename string, code string) string {
	// Try with filename first
	lang, safe := enry.GetLanguageByExtension(filename)
	if !safe || lang == "" {
		// Fallback to content-based detection
		lang = enry.GetLanguage(filename, []byte(code))
	}
	if lang != "" {
		return d.normalizeLanguage(lang)
	}
	return d.DetectLanguage(code)
}
