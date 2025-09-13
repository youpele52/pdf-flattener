package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var gsCommand = []string{
	"gs", "-sDEVICE=pdfwrite",
	"-dPDFSETTINGS=/prepress",
	"-dNOPAUSE",
	"-dBATCH",
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
	return nil
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
