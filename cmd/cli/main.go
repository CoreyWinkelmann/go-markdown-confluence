package main

import (
	"encoding/json"
	"flag"
	"fmt"
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
	dirDryRun := dirCmd.Bool("dry-run", false, "Simulate without uploading")

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
		handleConvert(*convertInput, *convertOutput)
	case "post":
		postCmd.Parse(os.Args[2:])
		handlePost(*postInput, *postURL, *postUsername, *postAPIToken, *postSpaceKey, *postTitle, *postParentID)
	case "directory":
		dirCmd.Parse(os.Args[2:])
		handleDirectory(*dirPath, *dirMapping, *dirURL, *dirUsername, *dirAPIToken, *dirSpaceKey, *dirDryRun)
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

func handleConvert(input, output string) {
	if input == "" {
		fmt.Println("Error: No input specified")
		return
	}

	var markdownContent string
	var err error
	var tempDir string
	var tempFile string

	if _, err := os.Stat(input); err == nil {
		fmt.Printf("Converting file: %s\n", input)

		fileMapping := map[string]string{
			input: input,
		}

		err = markdownconfluence.ConvertDirectory(filepath.Dir(input), fileMapping, nil)
		if err != nil {
			fmt.Printf("Error during conversion: %v\n", err)
			return
		}

		if output != "" {
			content, readErr := os.ReadFile(input)
			if readErr != nil {
				fmt.Printf("Error reading file: %v\n", readErr)
				return
			}
			markdownContent = string(content)
		}
	} else {
		fmt.Println("Converting input string")

		tempDir, err = os.MkdirTemp("", "markdown-convert")
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

		outputCapturer := &OutputCapturer{}

		err = markdownconfluence.ConvertDirectory(tempDir, fileMapping, outputCapturer)
		if err != nil {
			fmt.Printf("Error during conversion: %v\n", err)
			return
		}

		markdownContent = outputCapturer.GetMarkdown()
	}

	if output != "" {
		if markdownContent != "" {
			dummyClient := &OutputCapturer{}

			if tempFile == "" {
				tempDir, err = os.MkdirTemp("", "markdown-convert")
				if err != nil {
					fmt.Printf("Error creating temporary directory: %v\n", err)
					return
				}
				defer os.RemoveAll(tempDir)

				tempFile = filepath.Join(tempDir, "input.md")
				if err := os.WriteFile(tempFile, []byte(markdownContent), 0644); err != nil {
					fmt.Printf("Error writing temporary file: %v\n", err)
					return
				}
			}

			fileMapping := map[string]string{
				tempFile: tempFile,
			}

			err = markdownconfluence.ConvertDirectory(filepath.Dir(tempFile), fileMapping, dummyClient)
			if err != nil {
				fmt.Printf("Error during conversion: %v\n", err)
				return
			}

			result := dummyClient.GetMarkdown()
			if err := os.WriteFile(output, []byte(result), 0644); err != nil {
				fmt.Printf("Error writing to output file: %v\n", err)
				return
			}
			fmt.Printf("Conversion saved to: %s\n", output)
		}
	} else if markdownContent != "" {
		fmt.Println("Converted Markdown:")

		dummyClient := &OutputCapturer{}

		tempDir, err = os.MkdirTemp("", "markdown-convert")
		if err != nil {
			fmt.Printf("Error creating temporary directory: %v\n", err)
			return
		}
		defer os.RemoveAll(tempDir)

		tempFile = filepath.Join(tempDir, "input.md")
		if err := os.WriteFile(tempFile, []byte(markdownContent), 0644); err != nil {
			fmt.Printf("Error writing temporary file: %v\n", err)
			return
		}

		fileMapping := map[string]string{
			tempFile: tempFile,
		}

		err = markdownconfluence.ConvertDirectory(filepath.Dir(tempFile), fileMapping, dummyClient)
		if err != nil {
			fmt.Printf("Error during conversion: %v\n", err)
			return
		}

		result := dummyClient.GetMarkdown()
		fmt.Println(result)
	}
}

func handlePost(input, confluenceURL, username, apiToken, spaceKey, title, parentID string) {
	if input == "" || confluenceURL == "" || username == "" || apiToken == "" || spaceKey == "" || title == "" {
		fmt.Println("Error: Missing required parameters")
		return
	}

	var tempDir string
	var tempFile string
	var err error

	client := NewConfluenceClient(confluenceURL, username, apiToken)

	if _, statErr := os.Stat(input); statErr == nil {
		fmt.Printf("Converting and posting file: %s\n", input)

		fileMapping := map[string]string{
			input: title,
		}

		err = markdownconfluence.ConvertDirectory(filepath.Dir(input), fileMapping, client)
		if err != nil {
			fmt.Printf("Error during conversion or posting: %v\n", err)
			return
		}

		fmt.Println("Successfully created page in Confluence")
	} else {
		fmt.Println("Converting and posting input string")

		tempDir, err = os.MkdirTemp("", "markdown-post")
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

		err = markdownconfluence.ConvertDirectory(tempDir, fileMapping, client)
		if err != nil {
			fmt.Printf("Error during conversion or posting: %v\n", err)
			return
		}

		fmt.Println("Successfully created page in Confluence")
	}
}

func handleDirectory(dirPath, mappingPath, confluenceURL, username, apiToken, spaceKey string, dryRun bool) {
	if dirPath == "" {
		fmt.Println("Error: Directory path is required")
		return
	}

	if spaceKey == "" {
		spaceKey = "DOCS"
	}

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
	}

	var client markdownconfluence.ConfluenceClient
	if !dryRun && confluenceURL != "" && username != "" && apiToken != "" {
		client = NewConfluenceClient(confluenceURL, username, apiToken)
	}

	err := markdownconfluence.ConvertDirectory(dirPath, fileMapping, client)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Conversion completed successfully")
}

func printHelp() {
	fmt.Println("Markdown to Confluence Converter")
	fmt.Println("--------------------------------")
	fmt.Println("Usage:")
	fmt.Println("  convert --input <markdown_or_file> [--output <file>]")
	fmt.Println("  post --input <markdown_or_file> --url <confluence_url> --username <username> --token <api_token> --space <space_key> --title <title> [--parent <parent_id>]")
	fmt.Println("  directory --path <directory_path> [--mapping <mapping_file>] [--url <confluence_url> --username <username> --token <api_token> --space <space_key>] [--dry-run]")
	fmt.Println("  help, -help     Show this help message")
	fmt.Println("  version, -version    Show version information")
}

func printVersion() {
	fmt.Println("Markdown to Confluence Converter v0.1.0")
	fmt.Println("Released: April 2025")
}

type OutputCapturer struct {
	convertedMarkdown string
}

func (c *OutputCapturer) CreatePage(spaceKey, title, content string, parentID string) (string, error) {
	c.convertedMarkdown = content
	return "dummy-page-id", nil
}

func (c *OutputCapturer) UpdatePage(pageID, title, content, spaceKey string, version int) error {
	c.convertedMarkdown = content
	return nil
}

func (c *OutputCapturer) GetMarkdown() string {
	return c.convertedMarkdown
}
