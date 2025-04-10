package main

import (
	"fmt"
	"go-markdown-confluence/pkg/markdownconfluence"
	"os"
	"path/filepath"
	"testing"
)

func main() {
	// Use absolute paths directly
	dirPath := "/home/cagenix/source/go-markdown-confluence/examples"
	outputDir := "/home/cagenix/source/go-markdown-confluence/output"

	// Debugging: Print paths
	fmt.Printf("Examples directory: %s\n", dirPath)
	fmt.Printf("Output directory: %s\n", outputDir)

	// Check if examples directory exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Examples directory does not exist: %s\n", dirPath)
		os.Exit(1)
	}

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// Walk through the examples directory
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", path, err)
		}

		if info.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		fmt.Printf("Processing file: %s\n", path)

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		convertedContent, err := markdownconfluence.Convert(string(content))
		if err != nil {
			return fmt.Errorf("failed to convert file %s: %w", path, err)
		}

		outputPath := filepath.Join(outputDir, filepath.Base(path[:len(path)-3])+".json")
		if err := os.WriteFile(outputPath, []byte(convertedContent), 0644); err != nil {
			return fmt.Errorf("failed to write converted file %s: %w", outputPath, err)
		}

		fmt.Printf("Converted and saved: %s\n", outputPath)
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during conversion: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Println("All files processed successfully.")
	}
}

// Add integration tests for new ADF features
func TestIntegrationWithNewADFFeatures(t *testing.T) {
	markdown := `# Test Document\n\n:smile:\n\n<div>placeholder</div>\n\n- [ ] Task 1\n- [x] Task 2\n\n> Decision: Approve the proposal`
	expected := `{
		"type": "doc",
		"content": [
			{"type": "heading", "attrs": {"level": 1}, "content": [{"type": "text", "text": "Test Document"}]},
			{"type": "emoji", "attrs": {"shortName": ":smile:"}},
			{"type": "placeholder", "attrs": {"text": "Add your content here"}},
			{"type": "taskList", "content": []},
			{"type": "decisionItem", "attrs": {"state": "DECIDED"}}
		]
	}`

	result, err := markdownconfluence.Convert(markdown)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}
