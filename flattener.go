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
func flattenPDF(inputPDF string) error {
	dir := filepath.Dir(inputPDF)
	base := filepath.Base(inputPDF)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	outputPDF := filepath.Join(dir, name+"_flattened.pdf")
	args := append(gsCommand[1:], "-o", outputPDF, inputPDF)
	cmd := exec.Command(gsCommand[0], args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Printf("\nFlattening: %s -> %s\n", inputPDF, outputPDF)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("\n❌ failed to flatten PDF: %w", err)
	}

	// Check if the output file was created and has content
	if info, err := os.Stat(outputPDF); err != nil || info.Size() == 0 {
		if err != nil {
			return fmt.Errorf("\n❌ flattened PDF not created: %w", err)
		}
		return fmt.Errorf("\n⚠️ flattened PDF was created but is empty")
	}

	fmt.Printf("\n✅ Successfully flattened: %s\n", outputPDF)
	return nil
}

// flattenFolder walks the folder recursively and flattens each PDF.
func flattenFolder(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.ToLower(filepath.Ext(info.Name())) == ".pdf" {
			return flattenPDF(path)
		}
		return nil
	})
}
