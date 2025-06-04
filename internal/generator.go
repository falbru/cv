package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/goccy/go-yaml"
)

type Education struct {
	Name      string `yaml:"institution"`
	Location  string `yaml:"location"`
	Note      string `yaml:"degree"`
	StartDate string `yaml:"start_date"`
	EndDate   string `yaml:"end_date"`
}

type Experience struct {
	Name      string   `yaml:"company"`
	Location  string   `yaml:"location"`
	Title     string   `yaml:"position"`
	StartDate string   `yaml:"start_date"`
	EndDate   string   `yaml:"end_date"`
	Sentences []string `yaml:"description"`
}

type Project struct {
	Name      string   `yaml:"name"`
	URL       string   `yaml:"github_url"`
	Sentences []string `yaml:"description"`
}

type PersonalInfo struct {
	Name         string `yaml:"name"`
	Email        string `yaml:"email"`
	PhoneNumber  string `yaml:"phone"`
	DOB          string `yaml:"date_of_birth"`
	LinkedInUser string `yaml:"linkedin"`
	GitHubUser   string `yaml:"github"`
	Website      string `yaml:"website"`
}

type Skillset struct {
	Name     string   `yaml:"name"`
	Keywords []string `yaml:"keywords"`
}

type CV struct {
	PersonalInfo `yaml:"personal_info"`
	Educations   []Education       `yaml:"education"`
	Experiences  []Experience      `yaml:"experiences"`
	Projects     []Project         `yaml:"projects"`
	Skills       []Skillset        `yaml:"skills"`
	Translations map[string]string `yaml:"translations"`
}

func Generate(dataFilePath string, templatePath string, outputName string) error {
	var me CV

	fileContent, err := os.ReadFile(dataFilePath)
	if err != nil {
		return fmt.Errorf("failed to read YAML file: %w", err)
	}

	err = yaml.Unmarshal(fileContent, &me)
	if err != nil {
		return fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	funcMap := template.FuncMap{
		"join": strings.Join,
	}

	tmpl, err := template.New(filepath.Base(templatePath)).Funcs(funcMap).ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	texFileName := outputName + ".tex"
	texFile, err := os.Create(texFileName)
	if err != nil {
		return fmt.Errorf("failed to create .tex file: %w", err)
	}
	defer texFile.Close()

	err = tmpl.Execute(texFile, me)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	err = compilePDF(texFileName)
	if err != nil {
		return fmt.Errorf("failed to compile PDF: %w", err)
	}

	return nil
}

func compilePDF(texFileName string) error {
	dir := filepath.Dir(texFileName)
	if dir == "" {
		dir = "."
	}

	cmd := exec.Command("pdflatex", "-interaction=nonstopmode", "-output-directory="+dir, texFileName)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "pdflatex output:\n%s\n", string(output))
		return fmt.Errorf("pdflatex compilation failed: %w", err)
	}

	outputStr := string(output)
	if strings.Contains(outputStr, "Error") { // TODO Possibly not safe
		fmt.Printf("LaTeX compilation warnings/errors:\n%s\n", outputStr)
	}

	return nil
}

func CleanupLatexFiles(dir string) error {
	extensions := []string{".aux", ".log", ".tex"}

	for _, ext := range extensions {
		files, err := filepath.Glob(filepath.Join(dir, "*"+ext))
		if err != nil {
			return err
		}

		for _, file := range files {
			if _, err := os.Stat(file); err == nil {
				os.Remove(file)
			}

		}
	}
	return nil
}
