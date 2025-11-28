package md

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/StevenACoffman/j2m"
	cf "github.com/kentaro-m/blackfriday-confluence"
	bf "github.com/russross/blackfriday/v2"
)

// ANSI color codes for terminal output.
const (
	ansiReset = "\033[0m"
)

// Named colors mapped to ANSI codes.
var namedColors = map[string]string{
	"black":   "\033[30m",
	"red":     "\033[31m",
	"green":   "\033[32m",
	"yellow":  "\033[33m",
	"blue":    "\033[34m",
	"magenta": "\033[35m",
	"purple":  "\033[35m", // alias for magenta
	"cyan":    "\033[36m",
	"white":   "\033[37m",
	"orange":  "\033[38;5;208m",
	"pink":    "\033[38;5;213m",
	"brown":   "\033[38;5;130m",
	"gray":    "\033[90m",
	"grey":    "\033[90m",
}

// colorTagRegex matches {color:xxx}...{color} patterns.
var colorTagRegex = regexp.MustCompile(`\{color(?::([^}]+))?\}([\s\S]*?)\{color\}`)

// hexColorRegex matches hex color codes like #de350b or #f00.
var hexColorRegex = regexp.MustCompile(`^#?([0-9a-fA-F]{6}|[0-9a-fA-F]{3})$`)

// ToJiraMD translates CommonMark to Jira flavored markdown.
func ToJiraMD(md string) string {
	if md == "" {
		return md
	}

	renderer := &cf.Renderer{Flags: cf.IgnoreMacroEscaping}
	r := bf.New(bf.WithRenderer(renderer), bf.WithExtensions(bf.CommonExtensions))

	return string(renderer.Render(r.Parse([]byte(md))))
}

// stripColorTags removes {color:xxx}...{color} tags, keeping only the content.
func stripColorTags(input string) string {
	return colorTagRegex.ReplaceAllString(input, "$2")
}

// normalizeLineEndings converts Windows-style \r\n to Unix-style \n.
// This is needed because j2m doesn't handle \r\n properly for tables.
func normalizeLineEndings(input string) string {
	return strings.ReplaceAll(input, "\r\n", "\n")
}

// fixEscapedMarkup handles escaped/alternative Jira markup syntax.
// {*}text{*} -> *text* (bold)
// {_}text{_} -> _text_ (italic).
func fixEscapedMarkup(input string) string {
	// Replace {*}...{*} with *...*
	input = strings.ReplaceAll(input, "{*}", "*")
	// Replace {_}...{_} with _..._
	input = strings.ReplaceAll(input, "{_}", "_")
	return input
}

// FromJiraMD translates Jira flavored markdown to CommonMark.
// Color tags are stripped (content preserved).
func FromJiraMD(jfm string) string {
	// Normalize line endings, fix escaped markup, strip color tags, then convert
	jfm = normalizeLineEndings(jfm)
	jfm = fixEscapedMarkup(jfm)
	return j2m.JiraToMD(stripColorTags(jfm))
}

// ColorPlaceholders stores the mapping from placeholders to ANSI codes.
// This is used to preserve colors through glamour rendering.
type ColorPlaceholders struct {
	mu           sync.Mutex
	placeholders map[string]string // placeholder -> ANSI code
	resetMarker  string            // unique marker for reset
}

// NewColorPlaceholders creates a new ColorPlaceholders instance.
func NewColorPlaceholders() *ColorPlaceholders {
	return &ColorPlaceholders{
		placeholders: make(map[string]string),
		resetMarker:  generateMarker(),
	}
}

const (
	markerByteLen       = 8 // Length of random bytes for marker
	hexColorLen         = 6 // Length of full hex color (e.g., "ff0000")
	shortHexColorLen    = 3 // Length of short hex color (e.g., "f00")
	colorTagSubmatchLen = 3 // Expected submatches: full match, color spec, content
)

// generateMarker creates a unique marker that won't appear in normal text.
func generateMarker() string {
	b := make([]byte, markerByteLen)
	_, _ = rand.Read(b)
	return "CLRM" + hex.EncodeToString(b)
}

// FromJiraMDWithColors translates Jira flavored markdown to CommonMark,
// preserving colors as placeholders that can be replaced after glamour rendering.
// Returns the converted text and a ColorPlaceholders instance for later replacement.
func FromJiraMDWithColors(jfm string) (string, *ColorPlaceholders) {
	cp := NewColorPlaceholders()
	// First, convert color tags to placeholders
	result := cp.processColorTags(jfm)
	// Then convert the rest with j2m
	return j2m.JiraToMD(result), cp
}

// processColorTags converts {color:xxx}...{color} to unique placeholders.
func (cp *ColorPlaceholders) processColorTags(input string) string {
	return colorTagRegex.ReplaceAllStringFunc(input, func(match string) string {
		submatches := colorTagRegex.FindStringSubmatch(match)
		if len(submatches) < colorTagSubmatchLen {
			return match
		}

		colorSpec := strings.ToLower(strings.TrimSpace(submatches[1]))
		content := submatches[2]

		ansiCode := colorToANSI(colorSpec)
		if ansiCode == "" {
			return content // No valid color, just return content
		}

		// Create unique placeholder for this color
		placeholder := generateMarker()
		cp.mu.Lock()
		cp.placeholders[placeholder] = ansiCode
		cp.mu.Unlock()

		return placeholder + content + cp.resetMarker
	})
}

// ReplaceInRendered replaces placeholders with actual ANSI codes in rendered output.
func (cp *ColorPlaceholders) ReplaceInRendered(rendered string) string {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	result := rendered
	for placeholder, ansiCode := range cp.placeholders {
		result = strings.ReplaceAll(result, placeholder, ansiCode)
	}
	result = strings.ReplaceAll(result, cp.resetMarker, ansiReset)
	return result
}

// colorToANSI converts a color specification to an ANSI escape code.
func colorToANSI(colorSpec string) string {
	// Check named colors first
	if code, ok := namedColors[colorSpec]; ok {
		return code
	}

	// Check for hex color
	if hexColorRegex.MatchString(colorSpec) {
		return hexToANSI(colorSpec)
	}

	return ""
}

// hexToANSI converts a hex color code to ANSI 24-bit color escape sequence.
func hexToANSI(hex string) string {
	hex = strings.TrimPrefix(hex, "#")

	// Expand 3-char hex to 6-char
	if len(hex) == shortHexColorLen {
		hex = string(hex[0]) + string(hex[0]) +
			string(hex[1]) + string(hex[1]) +
			string(hex[2]) + string(hex[2])
	}

	if len(hex) != hexColorLen {
		return ""
	}

	r, err := strconv.ParseInt(hex[0:2], 16, 64)
	if err != nil {
		return ""
	}
	g, err := strconv.ParseInt(hex[2:4], 16, 64)
	if err != nil {
		return ""
	}
	b, err := strconv.ParseInt(hex[4:6], 16, 64)
	if err != nil {
		return ""
	}

	// Use 24-bit true color ANSI escape sequence
	return "\033[38;2;" + strconv.FormatInt(r, 10) + ";" +
		strconv.FormatInt(g, 10) + ";" +
		strconv.FormatInt(b, 10) + "m"
}
