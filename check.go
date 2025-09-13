package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// checkGhostscriptInstalled checks if Ghostscript is properly installed and available
func checkGhostscriptInstalled() error {
	// Try to run a simple Ghostscript command to verify it's the real Ghostscript
	cmd := exec.Command("gs", "--version")
	output, err := cmd.CombinedOutput()

	if err != nil || !strings.Contains(string(output), ".") {
		// Provide detailed installation instructions based on OS
		var installInstructions string
		if runtime.GOOS == "darwin" {
			installInstructions = "On macOS, install Ghostscript with: 'brew install ghostscript'"
		} else if runtime.GOOS == "linux" {
			installInstructions = "On Linux, install Ghostscript with: 'sudo apt-get install ghostscript' (Debian/Ubuntu) or 'sudo yum install ghostscript' (RHEL/CentOS)"
		} else if runtime.GOOS == "windows" {
			installInstructions = "On Windows, download and install Ghostscript from https://www.ghostscript.com/download/gsdnld.html"
		}

		return fmt.Errorf("Ghostscript (gs) is not installed or not working properly.\n%s\nNote: If 'gs' is aliased to another command, please unalias it or use the full path to Ghostscript", installInstructions)
	}

	return nil
}
