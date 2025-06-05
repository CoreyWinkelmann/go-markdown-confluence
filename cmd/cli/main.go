package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go-markdown-confluence/internal/confluence"
	"go-markdown-confluence/pkg/markdownconfluence"
	"os"
	"path/filepath"
)

func main() {
	helpFlag := flag.Bool("help", false, "Show usage information")
	versionFlag := flag.Bool("version", false, "Show version information")

	convertCmd := flag.NewFlagSet("convert", flag.ExitOnError)
	convertInput := convertCmd.String("input", "", "Markdown input (file or string)")
	convertOutput := convertCmd.String("output", "", "Output file (optional)")
	convertDryRun := convertCmd.Bool("dry-run", false, "Skip Confluence upload and output JSON")

	postCmd := flag.NewFlagSet("post", flag.ExitOnError)
	postInput := postCmd.String("input", "", "Markdown input (file or string)")
	postURL := postCmd.String("url", "", "Confluence URL")
	postUsername := postCmd.String("username", "", "Confluence username")
	postAPIToken := postCmd.String("token", "", "Confluence API token")
	postSpaceKey := postCmd.String("space", "", "Confluence space key")
	postTitle := postCmd.String("title", "", "Page title")
	postParentID := postCmd.String("parent", "", "Parent page ID (optional)")

	dirCmd := flag.NewFlagSet("directory", flag.ExitOnError)
	dirPath := dirCmd.String("path", "", "Path to the directory containing Markdown files")
	dirMapping := dirCmd.String("mapping", "", "Path to JSON file with file mappings")
	dirURL := dirCmd.String("url", "", "Confluence URL")
	dirUsername := dirCmd.String("username", "", "Confluence username")
	dirAPIToken := dirCmd.String("token", "", "Confluence API token")
	dirSpaceKey := dirCmd.String("space", "", "Confluence space key (default: DOCS)")
	dirDryRun := dirCmd.Bool("dry-run", false, "Skip uploading to Confluence")
	dirOutputDir := dirCmd.String("output-directory", "", "Directory to save converted JSON files (when using --dry-run)")

	flag.Parse()

	if *helpFlag {
		printHelp()
		return
	}

	if *versionFlag {
		printVersion()
		return
	}

	if len(os.Args) < 2 {
		printHelp()
		return
	}

	switch os.Args[1] {
	case "convert":
		convertCmd.Parse(os.Args[2:])
		handleConvert(*convertInput, *convertOutput, *convertDryRun)
	case "post":
		postCmd.Parse(os.Args[2:])
		handlePost(*postInput, *postURL, *postUsername, *postAPIToken, *postSpaceKey, *postTitle, *postParentID)
	case "directory":
		dirCmd.Parse(os.Args[2:])
		handleDirectory(*dirPath, *dirMapping, *dirURL, *dirUsername, *dirAPIToken, *dirSpaceKey, *dirDryRun, *dirOutputDir)
	case "help":
		printHelp()
	case "version":
		printVersion()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printHelp()
		os.Exit(1)
	}
}

func handleConvert(input, output string, dryRun bool) {
	if input == "" {
		fmt.Println("Error: No input specified")
		return
	}

	var markdownContent string
	var tempFile string

	if _, err := os.Stat(input); err == nil {
		fmt.Printf("Converting file: %s\n", input)

		fileMapping := map[string]string{
			input: input,
		}

		options := markdownconfluence.DefaultConvertOptions()
		options.DryRun = true // Always dry run for convert command

		// If output is specified, use it as the output directory
		if output != "" {
			options.OutputDirectory = filepath.Dir(output)
		}

		err = markdownconfluence.ConvertDirectoryWithOptions(filepath.Dir(input), fileMapping, nil, options, "DOCS")
		if err != nil {
			fmt.Printf("Error during conversion: %v\n", err)
			return
		}

		// If no output was specified but we want to display the result
		if output == "" {
			content, readErr := os.ReadFile(input)
			if readErr != nil {
				fmt.Printf("Error reading file: %v\n", readErr)
				return
			}
			markdownContent = string(content)
		}
	} else {
		fmt.Println("Converting input string")

		tempDir, err := os.MkdirTemp("", "markdown-convert")
		if err != nil {
			fmt.Printf("Error creating temporary directory: %v\n", err)
			return
		}
		defer os.RemoveAll(tempDir)

		tempFile = filepath.Join(tempDir, "input.md")
		if err := os.WriteFile(tempFile, []byte(input), 0644); err != nil {
			fmt.Printf("Error writing temporary file: %v\n", err)
			return
		}

		fileMapping := map[string]string{
			tempFile: tempFile,
		}

		options := markdownconfluence.DefaultConvertOptions()
		options.DryRun = true

		if output != "" {
			options.OutputDirectory = filepath.Dir(output)
		}

		outputCapturer := &OutputCapturer{}

		err = markdownconfluence.ConvertDirectoryWithOptions(tempDir, fileMapping, outputCapturer, options, "DOCS")
		if err != nil {
			fmt.Printf("Error during conversion: %v\n", err)
			return
		}

		markdownContent = outputCapturer.GetMarkdown()
	}

	// If output is specified but we haven't written to it yet
	if output != "" && !dryRun {
		result, err := convert(markdownContent)
		if err != nil {
			fmt.Printf("Error during conversion: %v\n", err)
			return
		}

		// If the output doesn't have a .json extension, add it
		outputPath := output
		if filepath.Ext(outputPath) != ".json" {
			outputPath += ".json"
		}

		if err := os.WriteFile(outputPath, []byte(result), 0644); err != nil {
			fmt.Printf("Error writing to output file: %v\n", err)
			return
		}
		fmt.Printf("Conversion saved to: %s\n", outputPath)
	} else if markdownContent != "" && !dryRun {
		// Display the result if we're not using dry-run
		result, err := convert(markdownContent)
		if err != nil {
			fmt.Printf("Error during conversion: %v\n", err)
			return
		}
		fmt.Println("Converted Markdown to ADF JSON:")
		fmt.Println(result)
	}
}

func convert(markdownContent string) (string, error) {
	if markdownContent == "" {
		return "", nil
	}

	tempDir, err := os.MkdirTemp("", "markdown-convert-temp")
	if err != nil {
		return "", fmt.Errorf("error creating temporary directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	tempFile := filepath.Join(tempDir, "input.md")
	if err := os.WriteFile(tempFile, []byte(markdownContent), 0644); err != nil {
		return "", fmt.Errorf("error writing temporary file: %w", err)
	}

	fileMapping := map[string]string{
		tempFile: tempFile,
	}

	dummyClient := &OutputCapturer{}
	options := markdownconfluence.DefaultConvertOptions()
	options.DryRun = true

	err = markdownconfluence.ConvertDirectoryWithOptions(tempDir, fileMapping, dummyClient, options, "DOCS")
	if err != nil {
		return "", fmt.Errorf("error during conversion: %w", err)
	}

	return dummyClient.GetMarkdown(), nil
}

func handlePost(input, confluenceURL, username, apiToken, spaceKey, title, parentID string) {
	if input == "" || confluenceURL == "" || username == "" || apiToken == "" || spaceKey == "" || title == "" {
		fmt.Println("Error: Missing required parameters")
		return
	}

	var tempFile string

	client := confluence.NewConfluenceClient(confluenceURL, username, apiToken)

	if _, statErr := os.Stat(input); statErr == nil {
		fmt.Printf("Converting and posting file: %s\n", input)

		fileMapping := map[string]string{
			input: title,
		}

		options := markdownconfluence.DefaultConvertOptions()
		options.DefaultSpaceKey = spaceKey

		err := markdownconfluence.ConvertDirectoryWithOptions(filepath.Dir(input), fileMapping, client, options, "DOCS")
		if err != nil {
			fmt.Printf("Error during conversion or posting: %v\n", err)
			return
		}

		fmt.Println("Successfully created page in Confluence")
	} else {
		fmt.Println("Converting and posting input string")

		tempDir, err := os.MkdirTemp("", "markdown-post")
		if err != nil {
			fmt.Printf("Error creating temporary directory: %v\n", err)
			return
		}
		defer os.RemoveAll(tempDir)

		tempFile = filepath.Join(tempDir, title+".md")
		if err := os.WriteFile(tempFile, []byte(input), 0644); err != nil {
			fmt.Printf("Error writing temporary file: %v\n", err)
			return
		}

		fileMapping := map[string]string{
			tempFile: title,
		}

		options := markdownconfluence.DefaultConvertOptions()
		options.DefaultSpaceKey = spaceKey

		err = markdownconfluence.ConvertDirectoryWithOptions(tempDir, fileMapping, client, options, "DOCS")
		if err != nil {
			fmt.Printf("Error during conversion or posting: %v\n", err)
			return
		}

		fmt.Println("Successfully created page in Confluence")
	}
}

func handleDirectory(dirPath, mappingPath, confluenceURL, username, apiToken, spaceKey string, dryRun bool, outputDir string) {
	fmt.Println("Starting directory conversion process...")

	if dirPath == "" {
		fmt.Println("Error: Directory path is required")
		return
	}

	// Get and display absolute path for the directory
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		fmt.Printf("Warning: Failed to get absolute path: %v\n", err)
	} else {
		fmt.Printf("Processing directory: %s (resolved to: %s)\n", dirPath, absPath)
	}

	// Check if directory exists
	dirInfo, err := os.Stat(dirPath)
	if err != nil {
		fmt.Printf("Error: Directory not accessible: %v\n", err)
		return
	}

	if !dirInfo.IsDir() {
		fmt.Printf("Error: Path is not a directory: %s\n", dirPath)
		return
	}

	// List files in the directory
	fmt.Println("Files in directory:")
	files, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Printf("Error: Failed to read directory: %v\n", err)
	} else {
		for _, file := range files {
			fmt.Printf("  - %s (isDir=%v)\n", file.Name(), file.IsDir())
		}
	}

	// Check output directory
	if dryRun && outputDir != "" {
		if _, err := os.Stat(outputDir); os.IsNotExist(err) {
			fmt.Printf("Creating output directory: %s\n", outputDir)
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				fmt.Printf("Error: Failed to create output directory: %v\n", err)
				return
			}
		}
	}

	if spaceKey == "" {
		spaceKey = "DOCS"
	}

	// Create an explicit file mapping for all markdown files if no mapping provided
	fileMapping := make(map[string]string)
	if mappingPath != "" {
		mappingFile, err := os.ReadFile(mappingPath)
		if err != nil {
			fmt.Printf("Error: Failed to read mapping file: %v\n", err)
			return
		}

		err = json.Unmarshal(mappingFile, &fileMapping)
		if err != nil {
			fmt.Printf("Error: Failed to parse mapping file: %v\n", err)
			return
		}
	} else {
		// No mapping file provided, create mapping for all markdown files
		err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("error accessing path %s: %w", path, err)
			}

			if info.IsDir() {
				return nil
			}

			if filepath.Ext(path) == ".md" {
				absPath, _ := filepath.Abs(path)
				fileMapping[absPath] = absPath
				fmt.Printf("Added to file mapping: %s\n", absPath)
			}
			return nil
		})

		if err != nil {
			fmt.Printf("Error: Failed to walk directory: %v\n", err)
			return
		}
	}

	fmt.Printf("File mapping contains %d entries\n", len(fileMapping))

	// Ensure consistent usage of the ConfluenceAPI interface
	var client confluence.ConfluenceAPI = confluence.NewConfluenceClient("", "", "")
	if !dryRun && confluenceURL != "" && username != "" && apiToken != "" {
		client = confluence.NewConfluenceClient(confluenceURL, username, apiToken)
	}

	options := markdownconfluence.DefaultConvertOptions()
	options.DryRun = dryRun
	options.OutputDirectory = outputDir
	options.DefaultSpaceKey = spaceKey

	fmt.Printf("Converting with options: DryRun=%v, OutputDirectory=%s\n", options.DryRun, options.OutputDirectory)

	for oldPath, newTitle := range fileMapping {
		if !dryRun {
			// Check if a page with the old title exists
			oldTitle := filepath.Base(oldPath)
			page, err := client.GetPageByTitle(spaceKey, oldTitle)
			if err != nil {
				fmt.Printf("Error checking for existing page: %v\n", err)
				continue
			}

			if page != nil {
				// Update the page title and parent if it exists
				fmt.Printf("Updating page '%s' to new title '%s'\n", oldTitle, newTitle)
				err = client.UpdatePage(page.ID, newTitle, page.Body.Storage.Value, spaceKey, page.Version.Number+1)
				if err != nil {
					fmt.Printf("Error updating page: %v\n", err)
					continue
				}
			} else {
				// Create a new page if it doesn't exist
				fmt.Printf("Creating new page '%s'\n", newTitle)
				_, err = client.CreatePage(spaceKey, newTitle, "", "")
				if err != nil {
					fmt.Printf("Error creating page: %v\n", err)
					continue
				}
			}
		}
	}

	err = markdownconfluence.ConvertDirectoryWithOptions(dirPath, fileMapping, client, options, "DOCS")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// List files in output directory after conversion
	if dryRun && outputDir != "" {
		fmt.Println("Files in output directory after conversion:")
		outFiles, err := os.ReadDir(outputDir)
		if err != nil {
			fmt.Printf("Error: Failed to read output directory: %v\n", err)
		} else {
			for _, file := range outFiles {
				fmt.Printf("  - %s\n", file.Name())
			}
		}
	}

	fmt.Println("Conversion completed successfully")
}

func printHelp() {
	fmt.Println("Markdown to Confluence Converter")
	fmt.Println("--------------------------------")
	fmt.Println("Usage:")
	fmt.Println("  convert --input <markdown_or_file> [--output <file>] [--dry-run]")
	fmt.Println("  post --input <markdown_or_file> --url <confluence_url> --username <username> --token <api_token> --space <space_key> --title <title> [--parent <parent_id>]")
	fmt.Println("  directory --path <directory_path> [--mapping <mapping_file>] [--url <confluence_url> --username <username> --token <api_token> --space <space_key>] [--dry-run] [--output-directory <directory>]")
	fmt.Println("  help, -help     Show this help message")
	fmt.Println("  version, -version    Show version information")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --dry-run             Skip uploading to Confluence")
	fmt.Println("  --output-directory    Directory to save converted JSON files when using --dry-run")
	fmt.Println("                        Files will be saved in a structure mirroring the original paths")
}

func printVersion() {
	fmt.Println("Markdown to Confluence Converter v0.1.0")
	fmt.Println("Released: April 2025")
}

type OutputCapturer struct {
	convertedMarkdown string
	Output            []string // added field to capture output logs
}

// Add the missing CreateParentPage method to make OutputCapturer implement markdownconfluence.ConfluenceClient
func (o *OutputCapturer) CreateParentPage(spaceKey, title, parentID string) (string, error) {
	o.Output = append(o.Output, fmt.Sprintf("Would create parent page '%s' in space '%s' with parent ID: %s",
		title, spaceKey, parentID))
	return fmt.Sprintf("dummy-parent-page-id-%s", title), nil
}

// Make sure the CreatePage method signature matches the interface requirement
func (o *OutputCapturer) CreatePage(spaceKey, title, content, parentID string) (string, error) {
	// If this method already exists, ensure its signature matches exactly
	o.Output = append(o.Output, fmt.Sprintf("Would create page '%s' in space '%s' with parent ID: %s",
		title, spaceKey, parentID))
	return fmt.Sprintf("dummy-page-id-%s", title), nil
}

func (c *OutputCapturer) UpdatePage(pageID, title, content, spaceKey string, version int) error {
	c.convertedMarkdown = content
	return nil
}

func (c *OutputCapturer) GetPageByTitle(spaceKey, title string) (*confluence.Page, error) {
	return nil, nil
}

func (c *OutputCapturer) GetPageByID(pageID string) (*confluence.Page, error) {
	return nil, nil
}

func (c *OutputCapturer) DeletePage(pageID string) error {
	return nil
}

func (c *OutputCapturer) UploadAttachment(pageID, filePath string) error {
	c.Output = append(c.Output, fmt.Sprintf("Would upload attachment %s to page %s", filePath, pageID))
	return nil
}

func (c *OutputCapturer) GetMarkdown() string {
	return c.convertedMarkdown
}
