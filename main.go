package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run . /path/to/your/file.pdf or go run . /path/to/your/folder")
		os.Exit(1)
	}

	if err := checkGhostscriptInstalled(); err != nil {
		log.Fatal(err)
	}

	arg := os.Args[1]
	info, err := os.Stat(arg)
	if err != nil {
		log.Fatalf("Cannot access argument: %v", err)
	}
	if info.IsDir() {
		if err := flattenFolder(arg); err != nil {
			log.Fatalf("Error while flattening folder: %v", err)
		}
	} else if strings.ToLower(filepath.Ext(info.Name())) == ".pdf" {
		if err := flattenPDF(arg); err != nil {
			log.Fatalf("Error while flattening file: %v", err)
		}
	} else {
		fmt.Println("Argument must be a PDF file or folder")
	}
}
