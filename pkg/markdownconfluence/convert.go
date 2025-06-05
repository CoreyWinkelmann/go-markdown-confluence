// Package markdownconfluence provides functionality to convert Markdown to Confluence-compatible format.
package markdownconfluence

import (
	"fmt"
	"regexp"
	"strings"

	"go-markdown-confluence/internal/confluence"
	"go-markdown-confluence/internal/converter"
	"go-markdown-confluence/internal/parser"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

// ConfluenceClient defines the interface that any Confluence client must implement.
type ConfluenceClient interface {
	// CreateParentPage creates a parent page in Confluence.
	CreateParentPage(spaceKey, title, parentID string) (string, error)
	// CreatePage creates a page in Confluence.
	CreatePage(spaceKey, title, content, parentID string) (string, error)
	// UpdatePage updates an existing page in Confluence.
	UpdatePage(pageID, title, content, spaceKey string, version int) error
	// GetPageByTitle retrieves a page by its title.
	GetPageByTitle(spaceKey, title string) (*confluence.Page, error)
	// UploadAttachment uploads an attachment to the given page and returns the attachment ID or URL.
	UploadAttachment(pageID, filePath string) error
}

// ConversionResult holds the result of a Markdown file conversion.
type ConversionResult struct {
	FilePath         string   // Original Markdown file path
	Title            string   // Page title derived from filename
	ConvertedContent string   // Converted content in ADF JSON format
	TargetPath       string   // Target path after applying mapping
	ImagePaths       []string // Paths to image files referenced in the Markdown
	PageID           string   // Existing Confluence page ID for updates
}

// Convert takes a Markdown string and converts it to a Confluence-compatible format.
// If the input is an empty string, it returns an empty string without error.
func stripObsidianComments(markdown string) string {
	re := regexp.MustCompile(`%%.*?%%`)
	return re.ReplaceAllString(markdown, "")
}

func replaceWikiLinks(markdown string) string {
	re := regexp.MustCompile(`\[\[(.*?)\]\]`)
	return re.ReplaceAllStringFunc(markdown, func(m string) string {
		match := re.FindStringSubmatch(m)
		if len(match) < 2 {
			return m
		}
		inner := match[1]
		escaped := strings.ReplaceAll(inner, " ", "%20")
		return "[" + inner + "](" + escaped + ")"
	})
}

func extractImagePaths(markdown string) []string {
	re := regexp.MustCompile(`!\[[^\]]*\]\(([^)]+)\)`)
	matches := re.FindAllStringSubmatch(markdown, -1)
	var paths []string
	for _, m := range matches {
		if len(m) > 1 {
			paths = append(paths, m[1])
		}
	}
	return paths
}

func extractFrontmatter(markdown string) (map[string]interface{}, string) {
	if !strings.HasPrefix(markdown, "---") {
		return nil, markdown
	}

	end := strings.Index(markdown[3:], "---")
	if end == -1 {
		return nil, markdown
	}

	front := markdown[3 : 3+end]
	body := markdown[3+end+3:]

	m := make(map[string]interface{})
	if err := yaml.Unmarshal([]byte(front), &m); err != nil {
		return nil, markdown
	}

	return m, strings.TrimLeft(body, "\n")
}

func Convert(markdown string) (string, error) {
	markdown = stripObsidianComments(markdown)
	markdown = replaceWikiLinks(markdown)

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
		contentBytes, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", path, err)
		}

		fm, body := extractFrontmatter(string(contentBytes))
		imagePaths := extractImagePaths(body)
		pageID, _ := fm["connie-page-id"].(string)

		confluenceContent, err := Convert(body)
		if err != nil {
			return nil, fmt.Errorf("failed to convert file %s: %w", path, err)
		}

		targetPath, exists := fileMapping[path]
		if !exists {
			targetPath = path
		}

		title := filepath.Base(targetPath)
		title = title[:len(title)-len(filepath.Ext(title))]
		if v, ok := fm["connie-title"].(string); ok && v != "" {
			title = v
		}

		// Add result to the collection
		results = append(results, ConversionResult{
			FilePath:         path,
			Title:            title,
			ConvertedContent: confluenceContent,
			TargetPath:       targetPath,
			ImagePaths:       imagePaths,
			PageID:           pageID,
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

		pageID := result.PageID
		if pageID != "" {
			err = confluenceClient.UpdatePage(pageID, result.Title, result.ConvertedContent, spaceKey, 1)
			if err != nil {
				return fmt.Errorf("failed to update page %s: %w", pageID, err)
			}
		} else {
			pageID, err = confluenceClient.CreatePage(spaceKey, result.Title, result.ConvertedContent, currentParentID)
			if err != nil {
				return fmt.Errorf("failed to upload file %s to Confluence: %w", result.FilePath, err)
			}
		}

		for _, img := range result.ImagePaths {
			absPath := filepath.Join(filepath.Dir(result.FilePath), img)
			if err := confluenceClient.UploadAttachment(pageID, absPath); err != nil {
				return fmt.Errorf("failed to upload attachment %s: %w", absPath, err)
			}
		}
	}
	return nil
}
