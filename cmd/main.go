package main

import (
	"cv/internal"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	DataDir = "data"
	TemplateFile = "template.tex"
	OutputDir = "output"
)

func main() {
	if _, err := os.Stat(TemplateFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Template file '%s' does not exist\n", TemplateFile)
		os.Exit(1)
	}

	files, err := filepath.Glob(filepath.Join(DataDir, "*.yml"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading data directory: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Printf("No .yml files found in %s directory\n", DataDir)
		return
	}

	fmt.Printf("Found %d YAML file(s) in %s/\n", len(files), DataDir)
	fmt.Printf("Using template: %s\n", TemplateFile)
	fmt.Printf("Output directory: %s/\n\n", OutputDir)

	for _, yamlFile := range files {
		baseName := strings.TrimSuffix(filepath.Base(yamlFile), ".yml")
		outputPath := filepath.Join(OutputDir, baseName)

		fmt.Printf("Generating %s.pdf from %s...\n", baseName, yamlFile)

		err := internal.Generate(yamlFile, TemplateFile, outputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error during generation: %v\n", err)
			os.Exit(1)
		}
	}

	err = internal.CleanupLatexFiles(OutputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error cleaning up LaTeX files: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Cleaned up auxiliary LaTeX files\n")

	fmt.Printf("\nAll CVs generated successfully!\n")
}
