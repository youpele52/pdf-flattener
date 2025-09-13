package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Check if we have at least one argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . /path/to/your/file.pdf [-replace] or go run . /path/to/your/folder [-replace]")
		os.Exit(1)
	}
	
	// Check if the replace flag is present
	replace := false
	filePath := os.Args[1]
	
	if len(os.Args) >= 3 && os.Args[2] == "-replace" {
		replace = true
	}

	if err := checkGhostscriptInstalled(); err != nil {
		log.Fatal(err)
	}

	info, err := os.Stat(filePath)
	if err != nil {
		log.Fatalf("Cannot access argument: %v", err)
	}
	
	if replace {
		fmt.Println("⚠️  Replace mode enabled: original files will be overwritten")
	}
	
	if info.IsDir() {
		if err := flattenFolder(filePath, replace); err != nil {
			log.Fatalf("Error while flattening folder: %v", err)
		}
	} else if strings.ToLower(filepath.Ext(info.Name())) == ".pdf" {
		if err := flattenPDF(filePath, replace); err != nil {
			log.Fatalf("Error while flattening file: %v", err)
		}
	} else {
		fmt.Println("Argument must be a PDF file or folder")
	}
}
