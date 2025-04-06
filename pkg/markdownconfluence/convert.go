// Package markdownconfluence provides functionality to convert Markdown to Confluence-compatible format.
package markdownconfluence

import (
	"fmt"
	"go-markdown-confluence/internal/converter"
	"go-markdown-confluence/internal/parser"
)

// Convert takes a Markdown string and converts it to a Confluence-compatible format.
// If the input is an empty string, it returns an empty string without error.
func Convert(markdown string) (string, error) {
	if markdown == "" {
		fmt.Println("Debug: Input Markdown is empty")
		return "", nil
	}

	fmt.Printf("Debug: Input Markdown size: %d bytes\n", len(markdown))

	// Debugging large input
	if len(markdown) > 1000 {
		fmt.Println("Debug: Large input detected")
		for i, r := range markdown {
			if r == '\x00' {
				fmt.Printf("Debug: Null character found at position %d\n", i)
				break
			}
		}
	}

	// Pre-validation for invalid characters
	for _, r := range markdown {
		if r == '\x00' {
			return "", fmt.Errorf("Invalid Markdown: Contains null character")
		}
	}

	// Step 1: Parse the markdown content
	mdParser := parser.NewMarkdownParser()
	document := mdParser.Parse(markdown)

	if document == nil {
		fmt.Println("Debug: Parsed AST is nil")
		return "", fmt.Errorf("Invalid Markdown: Parsed AST is nil")
	}

	fmt.Printf("Debug: Parsed AST node type: %T\n", document)
	fmt.Printf("Debug: Parsed AST has children: %v\n", document.HasChildren())

	if !document.HasChildren() {
		return "", fmt.Errorf("Invalid Markdown: Parsed AST has no children")
	}

	markdownBytes := []byte(markdown) // Convert string to bytes for processing
	adfDocument, err := converter.ConvertToADF(document, markdownBytes)
	if err != nil {
		return "", fmt.Errorf("Failed to convert Markdown to Confluence format: %w", err)
	}

	// Step 3: Serialize the ADF document to JSON
	jsonContent, err := converter.SerializeToJSON(adfDocument)
	if err != nil {
		return "", fmt.Errorf("Failed to serialize Confluence document to JSON: %w", err)
	}

	return jsonContent, nil
}
