# TPL

tpl is a simple command-line tool that processes structured input
files (CSV, JSON, XML, etc.) using the Go template engine and outputs
formatted content. This allows users to easily generate structured
files such as configuration files, import files, or even custom code
snippets from simple data sources.

## Goals

- **Flexible Input Formats**: Supports CSV, JSON, XML, and other text-based formats.
- **Go Template Engine**: Leverages Go’s powerful templating system for dynamic output.

## Installation

```sh
git clone https://github.com/jon-ski/tpl.git
cd tpl
go build -o tpl
```

Alternatively, install via `go install`:

```sh
go install github.com/jon-ski/tpl@latest
```

## Usage

Basic usage involves piping an input file into `tpl`, specifying a Go template file as an argument:

```sh
./tpl -i input.csv -t main.tmpl > output.txt
```

### Example: Generating Configuration Files

#### Input (`config.csv`):
```csv
name,ip,port
server1,192.168.1.10,8080
server2,192.168.1.11,9090
```

#### Template (`config.tmpl`):
```
server_name: {{ .name }}
address: {{ .ip }}:{{ .port }}
```

#### Command:
```sh
config.csv | ./tpl -i config.csv -t config.tmpl > config_output.txt
```

#### Output (`config_output.txt`):
```
server_name: server1
address: 192.168.1.10:8080

server_name: server2
address: 192.168.1.11:9090
```

### Example: Generating Random Code

Use JSON input to generate random struct definitions:

#### Input (`structs.json`):
```json
[
  {"name": "User", "fields": [{"type": "string", "name": "Username"}, {"type": "int", "name": "Age"}]},
  {"name": "Product", "fields": [{"type": "string", "name": "Name"}, {"type": "float64", "name": "Price"}]}
]
```

#### Template (`structs.tmpl`):
```go
{{ range . }}
type {{ .name }} struct {
  {{ range .fields }}{{ .name }} {{ .type }}
  {{ end }}
}
{{ end }}
```

#### Command:
```sh
./tpl -i structs.json -t structs.tmpl > models.go
```

#### Output (`models.go`):
```go
type User struct {
  Username string
  Age int
}

type Product struct {
  Name string
  Price float64
}
```

## License

MIT License.
