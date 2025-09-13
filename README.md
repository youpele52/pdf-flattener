# PDF Flattener

A command-line tool to flatten PDF files using Ghostscript. Flattening a PDF combines all layers, annotations, and form fields into a single layer, making the document more compatible across different PDF readers and preventing editing of form fields and redaction bars.

## Features

- Flatten a single PDF file
- Process all PDF files in a directory (including subdirectories)
- Preserves original files by creating new flattened versions
- Enhanced flattening for redaction bars to make them permanently uneditable
- PNG-based flattening for maximum compatibility and security
- Automatic detection and use of ImageMagick if available for better quality
- Clear success/failure indicators with emoji feedback
- Cross-platform support (macOS, Linux, Windows)

## Prerequisites

### Required Software

- **Go** (version 1.16 or higher): [Download Go](https://golang.org/dl/)
- **Ghostscript**: This tool relies on Ghostscript to process PDF files

### Optional Software

- **ImageMagick**: For enhanced PDF flattening quality (especially for redaction bars)
  - The tool will automatically detect and use ImageMagick if available
  - If not available, it will fall back to using Ghostscript

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

### Installing ImageMagick (Optional but Recommended)

#### macOS
```bash
brew install imagemagick
```

#### Linux
Debian/Ubuntu:
```bash
sudo apt-get install imagemagick
```

RHEL/CentOS:
```bash
sudo yum install imagemagick
```

#### Windows
Download and install from the [ImageMagick website](https://imagemagick.org/script/download.php)

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

To flatten a single PDF file (creates a new file with "_flattened" suffix):
```bash
go run . /path/to/your/file.pdf
```

To flatten and replace the original PDF file:
```bash
go run . /path/to/your/file.pdf -replace
```

To process all PDFs in a directory (including subdirectories):
```bash
go run . /path/to/your/folder
```

To process all PDFs in a directory and replace the original files:
```bash
go run . /path/to/your/folder -replace
```

### Running the compiled binary

If you've built the application:

To flatten a single PDF file (creates a new file with "_flattened" suffix):
```bash
./pdf-flattener /path/to/your/file.pdf
```

To flatten and replace the original PDF file:
```bash
./pdf-flattener /path/to/your/file.pdf -replace
```

To process all PDFs in a directory (including subdirectories):
```bash
./pdf-flattener /path/to/your/folder
```

To process all PDFs in a directory and replace the original files:
```bash
./pdf-flattener /path/to/your/folder -replace
```

## Enhanced Flattening Process

This tool uses a multi-stage approach to ensure thorough flattening of PDFs:

1. **Initial Flattening**: Uses Ghostscript to perform basic PDF flattening
2. **Advanced Flattening** (on macOS): 
   - Converts the PDF to high-resolution PNG images (200 DPI)
   - Converts the PNG images back to PDF
   - This process ensures that all elements, including redaction bars, become permanent and uneditable
   - Uses ImageMagick if available (for better quality) or falls back to Ghostscript

## Output

For each processed PDF file, a new file with `_flattened` appended to the filename will be created in the same directory as the original file.

Example:
```
Flattening: /Documents/example.pdf -> /Documents/example_flattened.pdf
✅ Successfully flattened: /Documents/example_flattened.pdf

Performing additional flattening pass for better redaction handling...
Step 1: Converting PDF to high-resolution PNG images...
Step 2: Converting PNG images back to PDF...
✅ PNG-based flattening completed - redaction bars should now be permanent
```

## Error Handling

The tool provides clear error messages with emoji indicators:

- ✅ Success: PDF successfully flattened
- ❌ Error: Failed to flatten PDF
- ⚠️ Warning: PDF was created but might have issues

## License

[MIT License](LICENSE)