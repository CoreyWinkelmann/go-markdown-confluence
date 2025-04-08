# go-markdown-confluence

`go-markdown-confluence` is a Go-based tool designed to convert Markdown files into Confluence-compatible formats. It simplifies the process of migrating or sharing Markdown content within Confluence by providing a seamless conversion pipeline.

## Features

- **Markdown Parsing**: Supports GitHub Flavored Markdown (GFM), including tables, task lists, and code blocks.
- **Confluence Integration**: Converts Markdown into Confluence-compatible JSON or other formats.
- **Customizable Rendering**: Extend or modify the rendering logic to suit specific Confluence requirements.
- **Nested File Support**: Handles nested Markdown files and outputs corresponding structured JSON files.

## Project Structure

```
cmd/
  cli/                # Command-line interface for the tool
examples/            # Example Markdown files demonstrating features
internal/
  confluence/        # Confluence client and related functions
  converter/         # Markdown to Confluence renderer
  parser/            # Markdown parsing logic
output/              # Generated output files (e.g., JSON)
pkg/
  markdownconfluence/ # Core library for Markdown to Confluence conversion
test/               # Test cases for conversion logic
```

## Getting Started

### Prerequisites

- Go 1.20 or later
- A Confluence account (if integrating with Confluence)

### Installation

Clone the repository:

```bash
git clone https://github.com/your-repo/go-markdown-confluence.git
cd go-markdown-confluence
```

Build the CLI:

```bash
cd cmd/cli
go build -o markdown-confluence
```

### Usage

Convert a Markdown file to Confluence JSON:

```bash
./markdown-confluence -input examples/basic.md -output output/basic.json
```

### Process an Entire Directory

The tool also supports processing an entire directory of Markdown files, including nested directories. This is useful for converting multiple Markdown files into Confluence-compatible formats in one go.

#### Usage

To process a directory, use the `-input-dir` flag to specify the directory containing Markdown files and the `-output-dir` flag to specify the directory where the converted files will be saved:

```bash
./markdown-confluence -input-dir examples/ -output-dir output/
```

#### Example

Input Directory (`examples/`):
```
examples/
  basic.md
  advanced-features.md
  nested/
    folder/
      nested1.md
      nested2.md
```

Output Directory (`output/`):
```
output/
  basic.json
  advanced-features.json
  nested/
    folder/
      nested1.json
      nested2.json
```

The tool will recursively process all Markdown files in the specified input directory and maintain the directory structure in the output directory.

### Example

Input (`examples/basic.md`):

```markdown
# Example

This is a basic example.
```

Output (`output/basic.json`):

```json
{
  "type": "doc",
  "content": [
    {
      "type": "heading",
      "content": "Example"
    },
    {
      "type": "paragraph",
      "content": "This is a basic example."
    }
  ]
}
```

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
