// Package markdownconfluence provides functionality to convert Markdown to Confluence-compatible format.
package markdownconfluence

import (
	"fmt"
	"go-markdown-confluence/internal/converter"
	"go-markdown-confluence/internal/parser"
	"os"
	"path/filepath"
)

// convert takes a Markdown string and converts it to a Confluence-compatible format.
// If the input is an empty string, it returns an empty string without error.
func convert(markdown string) (string, error) {
	for _, r := range markdown {
		if r == '\x00' {
			return "", fmt.Errorf("Invalid Markdown: Contains null character")
		}
	}

	mdParser := parser.NewMarkdownParser()
	document := mdParser.Parse(markdown)

	if document == nil {
		return "", nil
	}

	if !document.HasChildren() {
		return "", fmt.Errorf("Invalid Markdown: Parsed AST has no children")
	}

	markdownBytes := []byte(markdown)
	adfDocument, err := converter.ConvertToADF(document, markdownBytes)
	if err != nil {
		return "", fmt.Errorf("Failed to convert Markdown to Confluence format: %w", err)
	}

	jsonContent, err := converter.SerializeToJSON(adfDocument)
	if err != nil {
		return "", fmt.Errorf("Failed to serialize Confluence document to JSON: %w", err)
	}

	return jsonContent, nil
}

// ConvertDirectory takes a directory path, processes all Markdown files within it,
// and converts them to Confluence-compatible format. It uses a file mapping to handle
// renamed or moved files for upserts.
func ConvertDirectory(dirPath string, fileMapping map[string]string, confluenceClient ConfluenceClient) error {
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", path, err)
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".md" {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		confluenceContent, err := convert(string(content))
		if err != nil {
			return fmt.Errorf("failed to convert file %s: %w", path, err)
		}

		targetPath, exists := fileMapping[path]
		if !exists {
			targetPath = path
		}

		title := filepath.Base(targetPath)
		title = title[:len(title)-len(filepath.Ext(title))]

		spaceKey := "DOCS"
		var parentID string

		if confluenceClient != nil {
			pageID, err := confluenceClient.CreatePage(spaceKey, title, confluenceContent, parentID)
			if err != nil {
				return fmt.Errorf("failed to upload file %s to Confluence: %w", path, err)
			}
			fmt.Printf("Successfully uploaded to Confluence: %s (Page ID: %s)\n", targetPath, pageID)
		} else {
			fmt.Printf("Would upload to Confluence (dry run): %s\n", targetPath)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("error processing directory %s: %w", dirPath, err)
	}

	return nil
}

// ConfluenceClient interface defines the operations needed for Confluence
type ConfluenceClient interface {
	CreatePage(spaceKey, title, content string, parentID string) (string, error)
	UpdatePage(pageID, title, content, spaceKey string, version int) error
}
