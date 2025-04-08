// Package markdownconfluence provides functionality to convert Markdown to Confluence-compatible format.
package markdownconfluence

import (
	"fmt"
	"go-markdown-confluence/internal/converter"
	"go-markdown-confluence/internal/parser"
	"os"
	"path/filepath"
	"strings"
)

// ConfluenceClient defines the interface that any Confluence client must implement.
type ConfluenceClient interface {
	// CreateParentPage creates a parent page in Confluence.
	CreateParentPage(spaceKey, title, parentID string) (string, error)
	// CreatePage creates a page in Confluence.
	CreatePage(spaceKey, title, content, parentID string) (string, error)
}

// ConversionResult holds the result of a Markdown file conversion.
type ConversionResult struct {
	FilePath         string // Original Markdown file path
	Title            string // Page title derived from filename
	ConvertedContent string // Converted content in ADF JSON format
	TargetPath       string // Target path after applying mapping
}

// Convert takes a Markdown string and converts it to a Confluence-compatible format.
// If the input is an empty string, it returns an empty string without error.
func Convert(markdown string) (string, error) {
	for _, r := range markdown {
		if r == '\x00' {
			return "", fmt.Errorf("invalid Markdown: contains null character")
		}
	}

	mdParser := parser.NewMarkdownParser()
	document := mdParser.Parse(markdown)

	if document == nil {
		return "", nil
	}

	if !document.HasChildren() {
		return "", fmt.Errorf("invalid Markdown: parsed AST has no children")
	}

	markdownBytes := []byte(markdown)
	adfDocument, err := converter.ConvertToADF(document, markdownBytes)
	if err != nil {
		return "", fmt.Errorf("failed to convert Markdown to Confluence format: %w", err)
	}

	jsonContent, err := converter.SerializeToJSON(adfDocument)
	if err != nil {
		return "", fmt.Errorf("failed to serialize Confluence document to JSON: %w", err)
	}

	return jsonContent, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ConvertDirectoryOptions holds options for the ConvertDirectory function.
type ConvertDirectoryOptions struct {
	DryRun          bool   // If true, skip uploading to Confluence
	OutputDirectory string // Directory to save converted files (only used when DryRun is true)
	DefaultSpaceKey string // Default space key to use for Confluence
}

// DefaultConvertOptions returns the default options for ConvertDirectory.
func DefaultConvertOptions() *ConvertDirectoryOptions {
	return &ConvertDirectoryOptions{
		DryRun:          false,
		OutputDirectory: "",
		DefaultSpaceKey: "DOCS",
	}
}

// ConvertDirectoryWithResults takes a directory path, processes all Markdown files within it,
// converts them to Confluence-compatible format, and returns the conversion results.
// This is useful for testing the conversion without uploading to Confluence.
func ConvertDirectoryWithResults(dirPath string, fileMapping map[string]string, options *ConvertDirectoryOptions) ([]ConversionResult, error) {
	if options == nil {
		options = DefaultConvertOptions()
	}

	var results []ConversionResult

	var markdownFiles []string
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

		markdownFiles = append(markdownFiles, path)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error finding markdown files: %w", err)
	}

	// Process each markdown file
	for _, path := range markdownFiles {
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", path, err)
		}

		confluenceContent, err := Convert(string(content))
		if err != nil {
			return nil, fmt.Errorf("failed to convert file %s: %w", path, err)
		}

		targetPath, exists := fileMapping[path]
		if !exists {
			targetPath = path
		}

		title := filepath.Base(targetPath)
		title = title[:len(title)-len(filepath.Ext(title))]

		// Add result to the collection
		results = append(results, ConversionResult{
			FilePath:         path,
			Title:            title,
			ConvertedContent: confluenceContent,
			TargetPath:       targetPath,
		})

		// Save converted content to file if in dry run mode with output directory specified
		if options.DryRun && options.OutputDirectory != "" {
			outputDir := options.OutputDirectory

			// Create output directory if it doesn't exist
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create output directory %s: %w", outputDir, err)
			}

			// Create a subfolder structure mirroring the original path if needed
			relPath, err := filepath.Rel(dirPath, filepath.Dir(path))
			if err != nil {
				relPath = "" // If we can't get a relative path, use the root output directory
			}

			// Update ConvertDirectoryWithResults to handle nested folder structures
			outputSubdir := filepath.Join(outputDir, relPath)
			if err := os.MkdirAll(outputSubdir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create output subdirectory %s: %w", outputSubdir, err)
			}

			// Ensure the nested folder structure is preserved in the output directory
			outputPath := filepath.Join(outputSubdir, title+".json")
			if err := os.WriteFile(outputPath, []byte(confluenceContent), 0644); err != nil {
				return nil, fmt.Errorf("failed to write converted file %s: %w", outputPath, err)
			}
		}
	}

	return results, nil
}

// ConvertDirectory takes a directory path, processes all Markdown files within it,
// and converts them to Confluence-compatible format. It uses a file mapping to handle
// renamed or moved files for upserts.
func ConvertDirectory(dirPath string, fileMapping map[string]string, confluenceClient ConfluenceClient) error {
	return ConvertDirectoryWithOptions(dirPath, fileMapping, confluenceClient, nil, "")
}

// ConvertDirectoryWithOptions is like ConvertDirectory but allows specifying options.
func ConvertDirectoryWithOptions(dirPath string, fileMapping map[string]string, confluenceClient ConfluenceClient, options *ConvertDirectoryOptions, spaceKey string) error {
	if options == nil {
		options = DefaultConvertOptions()
	}

	if options.DryRun {
		_, err := ConvertDirectoryWithResults(dirPath, fileMapping, options)
		if err != nil {
			return fmt.Errorf("error during dry run conversion: %w", err)
		}
		return nil
	}

	if len(fileMapping) == 0 {
		fileMapping = make(map[string]string)
		filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && filepath.Ext(path) == ".md" {
				absPath, _ := filepath.Abs(path)
				fileMapping[absPath] = absPath
			}
			return nil
		})
	}

	results, err := ConvertDirectoryWithResults(dirPath, fileMapping, options)
	if err != nil {
		return err
	}

	parentPageIDs := make(map[string]string)

	for _, result := range results {
		relPath, err := filepath.Rel(dirPath, filepath.Dir(result.FilePath))
		if err != nil {
			return fmt.Errorf("failed to determine relative path for %s: %w", result.FilePath, err)
		}

		pathParts := strings.Split(filepath.ToSlash(filepath.Clean(relPath)), "/")
		currentParentID := ""
		for _, part := range pathParts {
			if part == "." || part == "" {
				continue
			}

			if parentPageID, exists := parentPageIDs[part]; exists {
				currentParentID = parentPageID
				continue
			}

			pageID, err := confluenceClient.CreateParentPage(spaceKey, part, currentParentID)
			if err != nil {
				return fmt.Errorf("failed to create parent page %s: %w", part, err)
			}
			parentPageIDs[part] = pageID
			currentParentID = pageID
		}

		_, err = confluenceClient.CreatePage(spaceKey, result.Title, result.ConvertedContent, currentParentID)
		if err != nil {
			return fmt.Errorf("failed to upload file %s to Confluence: %w", result.FilePath, err)
		}
	}
	return nil
}
