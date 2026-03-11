// package data gets data from various datasources.
package data

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/jon-ski/tpl/internal/csv"
	"github.com/jon-ski/tpl/internal/env"
)

type Format string

const (
	FormatAuto Format = "auto"
	FormatJSON Format = "json"
	FormatYAML Format = "yaml"
	FormatTOML Format = "toml"
	FormatCSV  Format = "csv"
	FormatTSV  Format = "tsv"
)

type DataSource struct {
	Name   string
	Scheme string
	Path   string
	URL    *url.URL // normalized URL (file/http/https)
	Format Format
}

func (d DataSource) String() string {
	return fmt.Sprintf(
		"Name:'%s' | Scheme:'%s' | Path:'%s' | URI:'%+v'",
		d.Name,
		d.Scheme,
		d.Path,
		d.URL,
	)
}

func (d DataSource) HasValidScheme() bool {
	switch d.Scheme {
	case "": // same as file
		return true
	case "file":
		return true
	}

	return false
}

func (d DataSource) GetData() (any, error) {
	switch d.Scheme {
	case "", "file":
		return d.getFileData()
	}
	return nil, fmt.Errorf("no data implementation for [%s]", d.Scheme)
}

func (d DataSource) getFileData() (any, error) {
	ext := d.Ext()
	switch ext {
	case "csv":
		return d.getCsvData()
	case "json":
		return d.getJsonData()
	}

	return nil, fmt.Errorf("format not supported: %s", ext)
}

func (d DataSource) getCsvData() ([]map[string]string, error) {
	f, err := os.Open(d.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	data, err := csv.Parse(f)
	if err != nil {
		return data, fmt.Errorf("failed to parse csv: %w", err)
	}
	return data, nil
}

func (d DataSource) getJsonData() (map[string]any, error) {
	f, err := os.Open(d.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	b, err := io.ReadAll(f)
	var data map[string]any
	err = json.Unmarshal(b, &data)
	return data, nil
}

func (d DataSource) Ext() string {
	q := d.URL.Query().Get("format")
	if q != "" {
		return q
	}
	filepathExt := strings.TrimPrefix(filepath.Ext(d.Path), ".")
	if filepathExt != "" {
		return filepathExt
	}
	return ""
}

type DataSourceList []DataSource

func (d *DataSourceList) String() string {
	return fmt.Sprintf("%v", []DataSource(*d))
}

func (d *DataSourceList) Set(v string) error {
	name, spec, ok := strings.Cut(v, "=")
	if !ok || name == "" || spec == "" {
		return fmt.Errorf("invalid -d %q (expected name=spec)", v)
	}
	ds, err := parseDataSourceValue(name, spec)
	if err != nil {
		return fmt.Errorf("failed to parse datasource: %w", err)
	}

	if !ds.HasValidScheme() {
		return fmt.Errorf("invalid scheme [%s]", ds.Scheme)
	}

	*d = append(*d, ds)
	return nil
}

func parseDataSourceValue(key, value string) (DataSource, error) {
	var ds DataSource
	u, err := parseUrl(value)
	if err != nil {
		return ds, fmt.Errorf("failed to parse URL: %w", err)
	}

	ds.Name = key
	ds.Scheme = u.Scheme

	// Prefer path, fall back to opaque
	pathRaw := u.Opaque

	if u.Path != "" {
		pathRaw = u.Path
	}

	ds.Path = filepath.FromSlash(pathRaw)
	ds.URL = u
	return ds, nil
}

// parseUrl is just a wrapper for url.Parse() in case I need to inject some
// logic at some point
func parseUrl(s string) (*url.URL, error) {
	u, err := url.Parse(s)
	if err != nil {
		return u, err
	}

	return u, nil
}

func GetAll(sources DataSourceList) (map[string]any, error) {
	data := make(map[string]any)
	var err error

	// Attach environment data
	data["env"], err = env.GetEnv()
	if err != nil {
		return data, err
	}

	// Loop over sources and attempt to fetch data
	for _, source := range sources {
		sourceData, err := source.GetData()
		if err != nil {
			return data, fmt.Errorf("failed to get data for source [%s]: %w", source.Name, err)
		}
		data[source.Name] = sourceData
	}

	return data, nil
}
