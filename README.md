# PDF Flattener

A simple command-line tool to flatten PDF files using Ghostscript. Flattening a PDF combines all layers, annotations, and form fields into a single layer, making the document more compatible across different PDF readers and preventing editing of form fields.

## Features

- Flatten a single PDF file
- Process all PDF files in a directory (including subdirectories)
- Preserves original files by creating new flattened versions
- Clear success/failure indicators with emoji feedback
- Cross-platform support (macOS, Linux, Windows)

## Prerequisites

### Required Software

- **Go** (version 1.16 or higher): [Download Go](https://golang.org/dl/)
- **Ghostscript**: This tool relies on Ghostscript to process PDF files

### Installing Ghostscript

#### macOS
```bash
brew install ghostscript
```

#### Linux
Debian/Ubuntu:
```bash
sudo apt-get install ghostscript
```

RHEL/CentOS:
```bash
sudo yum install ghostscript
```

#### Windows
Download and install from the [Ghostscript website](https://www.ghostscript.com/download/gsdnld.html)

## Installation

1. Clone this repository:
```bash
git clone https://github.com/yourusername/pdf-flattener.git
cd pdf-flattener
```

2. Build the application (optional):
```bash
go build
```

## Usage

### Running with Go

To flatten a single PDF file:
```bash
go run . /path/to/your/file.pdf
```

To process all PDFs in a directory (including subdirectories):
```bash
go run . /path/to/your/folder
```

### Running the compiled binary

If you've built the application:

To flatten a single PDF file:
```bash
./pdf-flattener /path/to/your/file.pdf
```

To process all PDFs in a directory (including subdirectories):
```bash
./pdf-flattener /path/to/your/folder
```

## Output

For each processed PDF file, a new file with `_flattened` appended to the filename will be created in the same directory as the original file.

Example:
```
Flattening: /Documents/example.pdf -> /Documents/example_flattened.pdf
✅ Successfully flattened: /Documents/example_flattened.pdf
```

## Error Handling

The tool provides clear error messages with emoji indicators:

- ✅ Success: PDF successfully flattened
- ❌ Error: Failed to flatten PDF
- ⚠️ Warning: PDF was created but might have issues

## License

[MIT License](LICENSE)