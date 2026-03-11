# TPL

tpl is a simple command-line tool that processes structured input
files (CSV, JSON, XML, etc.) using the Go template engine and outputs
formatted content. This allows users to easily generate structured
files such as configuration files, import files, or even custom code
snippets from simple data sources.

## Goals

- **Flexible Input Formats**: Supports CSV, JSON, XML, and other text-based formats.
- **Go Template Engine**: Leverages Go’s powerful templating system for dynamic output.

## Usage

```sh
tpl -i "{{ . }}"
```

## License

MIT License.
