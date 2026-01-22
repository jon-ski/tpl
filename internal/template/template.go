package template

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

// RunTemplate parses and executes a Go text/template text using the provided data.
func RunTemplate(t string, data any) error {
	slog.Debug("starting")
	// Parse template
	tmpl, err := template.New("root").Funcs(templateFuncs()).Parse(t)
	// tmpl, err := template.ParseFiles(*flags.templateFile)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}
	slog.Debug("template parsed", slog.String("name", tmpl.Name()), slog.String("templates", tmpl.DefinedTemplates()))

	// Execute template
	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	slog.Debug("template executed successfully")

	return nil
}

// RunTemplatePath parses and executes a Go text/template file using the provided data.
// The template is loaded from templateFilePath, augmented with custom template
// functions, and executed with its base filename as the template name.
// Output is written directly to standard output.
// An error is returned if parsing or execution fails.
func RunTemplatePath(templateFilePath string, data any) error {
	slog.Debug("starting")
	// Parse Template file
	tmpl, err := template.New("root").Funcs(templateFuncs()).ParseFiles(templateFilePath)
	// tmpl, err := template.ParseFiles(*flags.templateFile)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}
	slog.Debug("template parsed", slog.String("name", tmpl.Name()), slog.String("templates", tmpl.DefinedTemplates()))

	// Execute template
	err = tmpl.ExecuteTemplate(os.Stdout, filepath.Base(templateFilePath), data)
	// tmpl.Templates()[0].Execute(os.Stdout, data)
	// err = tmpl.ExecuteTemplate(os.Stdout, *flags.templateFile, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	slog.Info("template executed successfully")

	return nil
}

// templateFuncs returns a template.FuncMap containing utility functions
// intended for use inside Go templates.
// The function map includes helpers for:
//   - Integer and floating-point math operations.
//   - String manipulation and inspection.
//   - String-to-number conversions with fallback values on error.
//   - List and sequence construction.
//   - Map construction and key extraction.
//   - JSON and XML marshaling for embedding structured data in templates.
//
// These functions are registered with templates at parse time and are
// designed to favor convenience and expressiveness within template logic
// over strict error handling.
func templateFuncs() template.FuncMap {
	return template.FuncMap{
		// math
		"mod": func(i, j int) int {
			return i % j
		},
		"add": func(i, j int) int {
			return i + j
		},
		"sub": func(i, j int) int {
			return i - j
		},
		"div": func(i, j int) int {
			return i / j
		},
		"abs": func(i int) int {
			if i < 0 {
				return -i
			}
			return i
		},
		"min": func(a, b int) int {
			if a < b {
				return a
			}
			return b
		},
		"max": func(a, b int) int {
			if a > b {
				return a
			}
			return b
		},
		"addf": func(i, j float64) float64 {
			return i + j
		},
		"subf": func(i, j float64) float64 {
			return i - j
		},
		"divf": func(i, j float64) float64 {
			return i / j
		},
		"absf": func(i float64) float64 {
			if i < 0 {
				return -i
			}
			return i
		},
		"minf": func(a, b float64) float64 {
			if a < b {
				return a
			}
			return b
		},
		"maxf": func(a, b float64) float64 {
			if a > b {
				return a
			}
			return b
		},

		// strings
		"toUpper":    strings.ToUpper,
		"toLower":    strings.ToLower,
		"trim":       strings.TrimSpace,
		"trimLeft":   strings.TrimLeft,
		"trimRight":  strings.TrimRight,
		"trimPrefix": strings.TrimPrefix,
		"trimSuffix": strings.TrimSuffix,
		"hasPrefix":  strings.HasPrefix,
		"hasSuffix":  strings.HasSuffix,
		"replace":    strings.Replace,
		"split":      strings.Split,
		"join":       strings.Join,
		"joinEmpty":  joinEmpty,
		"contains":   strings.Contains,
		"count":      strings.Count,
		"lastIndex":  strings.LastIndex,
		"repeat":     strings.Repeat,

		// Conversions
		"atoi": func(s string) int {
			x, err := strconv.Atoi(s)
			if err != nil {
				return math.MaxInt // TODO: Figure out what value to use if err
			}
			return x
		},

		"atof": func(s string) float64 {
			x, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return math.MaxFloat64
			}
			return x
		},

		// Lists
		"list": func(v ...any) []any {
			return v
		},
		"seq": func(start, end int) []int {
			if start > end {
				return []int{}
			}
			s := make([]int, end-start+1)
			for i := range s {
				s[i] = start + i
			}
			return s
		},
		"len": func(l []any) int {
			return len(l)
		},

		// Maps
		"dict": func(kv ...any) (map[string]any, error) {
			m := make(map[string]any)
			for i := 0; i < len(kv); i += 2 {
				key, ok := kv[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict keys must be strings")
				}
				m[key] = kv[i+1]
			}
			return m, nil
		},

		"keys": func(m map[string]any) []string {
			list := make([]string, 0, len(m))
			for k := range m {
				list = append(list, k)
			}
			return list
		},

		// Marshallers
		"json": func(v any) string {
			b, err := json.Marshal(v)
			if err != nil {
				return fmt.Sprintf("error marshalling json: %s", err)
			}
			return string(b)
		},

		"xml": func(v any) string {
			b, err := xml.Marshal(v)
			if err != nil {
				return fmt.Sprintf("error marshalling xml: %s", err)
			}
			return string(b)
		},
	}
}

// joinEmpty joins a list of strings, but ignores empty strings (removing whitespace counts as empty)
func joinEmpty(sep string, s ...string) string {
	if len(s) == 0 {
		return ""
	}
	if len(s) == 1 {
		return s[0]
	}

	var buf strings.Builder
	for i := range s {
		if strings.TrimSpace(s[i]) != "" {
			buf.WriteString(s[i])
			if i < len(s)-1 {
				buf.WriteString(sep)
			}
		}
	}
	return buf.String()
}
