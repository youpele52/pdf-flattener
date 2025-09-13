package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

var gsCommand = []string{
	"gs", "-sDEVICE=pdfwrite",
	"-dPDFSETTINGS=/prepress",
	"-dNOPAUSE",
	"-dBATCH",
	"-dPrinted=true",
}

// flattenPDF runs Ghostscript to flatten a single PDF.
func flattenPDF(inputPDF string, replace bool) error {
	dir := filepath.Dir(inputPDF)
	base := filepath.Base(inputPDF)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)

	// Determine the output path
	var outputPDF string
	var tempOutput string

	if replace {
		// Create a temporary file path for processing
		tempOutput = filepath.Join(dir, name+"_temp_"+ext)
		outputPDF = inputPDF // For display purposes
	} else {
		outputPDF = filepath.Join(dir, name+"_flattened.pdf")
		tempOutput = outputPDF // No temp needed when not replacing
	}
	// Use the temporary output path for Ghostscript
	args := append(gsCommand[1:], "-o", tempOutput, inputPDF)
	cmd := exec.Command(gsCommand[0], args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Printf("\nFlattening: %s -> %s\n", inputPDF, outputPDF)

	err := cmd.Run()
	if err != nil {
		// Clean up temp file if it exists
		os.Remove(tempOutput)
		return fmt.Errorf("\n❌ failed to flatten PDF: %w", err)
	}

	// Check if the output file was created and has content
	if info, err := os.Stat(tempOutput); err != nil || info.Size() == 0 {
		// Clean up temp file if it exists
		os.Remove(tempOutput)
		if err != nil {
			return fmt.Errorf("\n❌ flattened PDF not created: %w", err)
		}
		return fmt.Errorf("\n⚠️ flattened PDF was created but is empty")
	}

	// If we're replacing the original file, do the replacement now
	if replace {
		// Remove the original file
		if err := os.Remove(inputPDF); err != nil {
			// Clean up temp file
			os.Remove(tempOutput)
			return fmt.Errorf("\n❌ could not remove original file for replacement: %w", err)
		}

		// Rename the temp file to the original file name
		if err := os.Rename(tempOutput, inputPDF); err != nil {
			return fmt.Errorf("\n❌ could not rename temp file to original: %w", err)
		}
	}

	fmt.Printf("\n✅ Successfully flattened: %s\n", outputPDF)

	// On Mac, perform an additional pass for better handling of redaction bars
	if runtime.GOOS == "darwin" {
		finalOutput := outputPDF
		if replace {
			finalOutput = inputPDF
		}

		// Run the additional flattening pass
		performAdditionalFlatteningPass(finalOutput, dir, name, ext)
	}

	return nil
}

// performAdditionalFlatteningPass performs a conversion to high-resolution PNG and back to PDF
// to completely flatten all elements including redaction bars.
func performAdditionalFlatteningPass(inputPDF, dir, name, ext string) {
	// Create temporary files for the PNG conversion process
	tempPNGDir := filepath.Join(dir, "temp_png_"+name)
	secondPassOutput := filepath.Join(dir, name+"_second_pass_"+ext)
	
	// Create temp directory for PNG files
	err := os.MkdirAll(tempPNGDir, 0755)
	if err != nil {
		fmt.Printf("\n⚠️ Could not create temporary directory: %v\n", err)
		return
	}
	defer os.RemoveAll(tempPNGDir) // Clean up temp directory when done
	
	fmt.Printf("\nPerforming PNG conversion for complete redaction flattening...\n")
	
	// Step 1: Convert PDF to high-resolution PNG images
	pngPrefix := filepath.Join(tempPNGDir, "page")
	pdfToPngArgs := []string{
		"gs",
		"-dSAFER",
		"-dBATCH",
		"-dNOPAUSE",
		"-sDEVICE=png16m", // 24-bit RGB color PNG
		"-r200", // 300 DPI for high resolution
		"-dTextAlphaBits=4",
		"-dGraphicsAlphaBits=4",
		"-dMaxBitmap=500000000", // Allow large bitmaps
		"-sOutputFile=" + pngPrefix + "-%03d.png",
		inputPDF,
	}
	
	cmd := exec.Command(pdfToPngArgs[0], pdfToPngArgs[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	fmt.Printf("\nStep 1: Converting PDF to high-resolution PNG images...\n")
	if err := cmd.Run(); err != nil {
		fmt.Printf("\n⚠️ PNG conversion failed, but original flattening succeeded: %v\n", err)
		return
	}
	
	// Step 2: Convert PNG images back to a single PDF
	pngFiles, err := filepath.Glob(filepath.Join(tempPNGDir, "*.png"))
	if err != nil || len(pngFiles) == 0 {
		fmt.Printf("\n⚠️ No PNG files found: %v\n", err)
		return
	}
	
	// Sort PNG files to ensure correct page order
	sort.Strings(pngFiles)
	
	// Check if ImageMagick is installed
	isMagickInstalled := checkCommandExists("magick")
	if !isMagickInstalled {
		// Also check for older 'convert' command
		isMagickInstalled = checkCommandExists("convert")
	}
	
	var pngToPdfArgs []string
	
	if isMagickInstalled {
		// Use ImageMagick's magick for reliable PNG->PDF assembly
		magickCmd := "magick"
		if !checkCommandExists("magick") && checkCommandExists("convert") {
			// Use older 'convert' command if 'magick' is not available
			magickCmd = "convert"
		}
		
		pngToPdfArgs = []string{magickCmd}
		// Append all PNG files first (already sorted)
		pngToPdfArgs = append(pngToPdfArgs, pngFiles...)
		// Add processing options and output file
		pngToPdfArgs = append(pngToPdfArgs,
			"-units", "PixelsPerInch",
			"-density", "300",
			"-strip",
			"-compress", "zip",
			secondPassOutput)
		
		fmt.Printf("\nUsing ImageMagick for PDF creation...\n")
	} else {
		// Fall back to Ghostscript if ImageMagick is not available
		fmt.Printf("\nImageMagick not found, falling back to Ghostscript...\n")
		
		// Create a text file with the list of PNG files for Ghostscript
		pngListFile := filepath.Join(tempPNGDir, "pnglist.txt")
		pngListContent := ""
		for _, pngFile := range pngFiles {
			pngListContent += pngFile + "\n"
		}
		
		if err := os.WriteFile(pngListFile, []byte(pngListContent), 0644); err != nil {
			fmt.Printf("\n⚠️ Failed to create PNG list file: %v\n", err)
			return
		}
		
		pngToPdfArgs = []string{
			"gs",
			"-dSAFER",
			"-dBATCH",
			"-dNOPAUSE",
			"-sDEVICE=pdfwrite",
			"-dPDFSETTINGS=/prepress",
			"-dCompatibilityLevel=1.4",
			"-sOutputFile=" + secondPassOutput,
			"-c", "30000000 setvmthreshold",
			"@" + pngListFile,
		}
	}
	
	cmd = exec.Command(pngToPdfArgs[0], pngToPdfArgs[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	fmt.Printf("\nStep 2: Converting PNG images back to PDF...\n")
	if err := cmd.Run(); err != nil {
		os.Remove(secondPassOutput)
		fmt.Printf("\n⚠️ PDF creation from PNGs failed: %v\n", err)
		return
	}
	
	// Check if the output PDF was created and has content
	if info, err := os.Stat(secondPassOutput); err != nil || info.Size() == 0 {
		os.Remove(secondPassOutput)
		fmt.Printf("\n⚠️ PDF creation failed, but original flattening succeeded\n")
		return
	}
	
	// Replace the input file with the completely flattened version
	os.Remove(inputPDF)
	os.Rename(secondPassOutput, inputPDF)
	
	fmt.Printf("\n✅ PNG-based flattening completed - redaction bars should now be permanent\n")
}

// flattenFolder walks the folder recursively and flattens each PDF.
func flattenFolder(root string, replace bool) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.ToLower(filepath.Ext(info.Name())) == ".pdf" {
			return flattenPDF(path, replace)
		}
		return nil
	})
}
