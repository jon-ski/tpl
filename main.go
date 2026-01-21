package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	icsv "github.com/jon-ski/tpl/internal/csv"
	ijson "github.com/jon-ski/tpl/internal/json"
	ixml "github.com/jon-ski/tpl/internal/xml"
)

type Flags struct {
	inputFile    *string
	templateFile *string

	debug *bool
}

func flagParse() *Flags {
	flags := &Flags{
		inputFile:    flag.String("i", "", "input file path"),
		templateFile: flag.String("t", "main.tmpl", "template file path (default: 'main.tmpl')"),
		debug:        flag.Bool("d", false, "sets log level to debug"),
	}
	flag.Parse()
	return flags
}

var logLevel = new(slog.LevelVar)

func main() {
	// Parse flags
	flags := flagParse()

	// Set up log
	logLevel.Set(slog.LevelInfo)
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	if *flags.debug {
		logLevel.Set(slog.LevelDebug)
	}

	// Open input file
	inputFile, err := os.Open(*flags.inputFile)
	if err != nil {
		slog.Error("failed to open input file", slog.String("err", err.Error()))
		os.Exit(1)
	}

	// Parse Input File to a map based on the extension
	data, err := createMap(inputFile)
	if err != nil {
		slog.Error("failed to create data map", slog.String("err", err.Error()))
		os.Exit(1)
	}

	// Parse Template file
	tmpl, err := template.New("root").Funcs(templateFuncs()).ParseFiles(*flags.templateFile)
	// tmpl, err := template.ParseFiles(*flags.templateFile)
	if err != nil {
		slog.Error("failed to parse template", slog.String("err", err.Error()))
		os.Exit(1)
	}
	slog.Debug("template parsed", slog.String("name", tmpl.Name()), slog.String("templates", tmpl.DefinedTemplates()))

	// Execute template
	err = tmpl.ExecuteTemplate(os.Stdout, filepath.Base(*flags.templateFile), data)
	// tmpl.Templates()[0].Execute(os.Stdout, data)
	// err = tmpl.ExecuteTemplate(os.Stdout, *flags.templateFile, data)
	if err != nil {
		slog.Error("failed to execute template", slog.String("err", err.Error()))
	}

	slog.Info("template executed successfully")
}

func createMap(inputFile *os.File) (any, error) {
	// Determine file extension
	ext := getFileType(inputFile.Name())

	// Create map based on file type
	switch ext {
	case "csv":
		return icsv.Parse(inputFile)
	case "xml":
		return ixml.Parse(inputFile)
	case "json":
		return ijson.Parse(inputFile)
	default:
		return nil, fmt.Errorf("file type '%s' not supported", ext)
	}
}

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
