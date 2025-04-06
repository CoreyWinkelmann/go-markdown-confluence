package main

import (
	"fmt"
	"go-markdown-confluence/pkg/markdownconfluence"
	"os"
	"strings"
)

func main() {
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "help":
			printHelp()
		case "convert":
			handleConvert()
		case "post":
			handlePost()
		case "version":
			printVersion()
		default:
			fmt.Println("Unknown command. Use 'help' for usage information.")
		}
	} else {
		printHelp()
	}
}

func handleConvert() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: convert <markdown_string_or_file>")
		return
	}

	input := os.Args[2]
	var markdown string
	var err error

	// Check if the input is a file path
	if _, err := os.Stat(input); err == nil {
		// It's a file, read its content
		content, readErr := os.ReadFile(input)
		if readErr != nil {
			fmt.Printf("Error reading file: %v\n", readErr)
			return
		}
		markdown = string(content)
		fmt.Printf("Converting file: %s\n", input)
	} else {
		// It's a direct markdown string
		markdown = input
		fmt.Println("Converting input string")
	}

	converted, err := markdownconfluence.Convert(markdown)
	if err != nil {
		fmt.Printf("Error during conversion: %v\n", err)
		return
	}

	// Check if output should go to a file
	if len(os.Args) > 3 && strings.HasPrefix(os.Args[3], "--output=") {
		outputPath := strings.TrimPrefix(os.Args[3], "--output=")
		if err := os.WriteFile(outputPath, []byte(converted), 0644); err != nil {
			fmt.Printf("Error writing to output file: %v\n", err)
			return
		}
		fmt.Printf("Conversion saved to: %s\n", outputPath)
	} else {
		fmt.Println("Converted Markdown:")
		fmt.Println(converted)
	}
}

func handlePost() {
	if len(os.Args) < 6 {
		fmt.Println("Usage: post <markdown_string_or_file> <confluence_url> <username> <api_token> <space_key> <title> [parent_id]")
		return
	}

	input := os.Args[2]
	confluenceURL := os.Args[3]
	username := os.Args[4]
	apiToken := os.Args[5]
	spaceKey := os.Args[6]
	title := os.Args[7]

	var parentID string
	if len(os.Args) > 8 {
		parentID = os.Args[8]
	}

	var markdown string

	// Check if the input is a file path
	if _, err := os.Stat(input); err == nil {
		// It's a file, read its content
		content, readErr := os.ReadFile(input)
		if readErr != nil {
			fmt.Printf("Error reading file: %v\n", readErr)
			return
		}
		markdown = string(content)
		fmt.Printf("Converting and posting file: %s\n", input)
	} else {
		// It's a direct markdown string
		markdown = input
		fmt.Println("Converting and posting input string")
	}

	// Step 1: Convert markdown to Confluence format (pure function, part of functional core)
	converted, err := markdownconfluence.Convert(markdown)
	if err != nil {
		fmt.Printf("Error during conversion: %v\n", err)
		return
	}

	// Step 2: Post to Confluence (side effect, part of imperative shell)
	client := NewConfluenceClient(confluenceURL, username, apiToken)
	pageID, err := client.CreatePage(spaceKey, title, converted, parentID)
	if err != nil {
		fmt.Printf("Error posting to Confluence: %v\n", err)
		return
	}

	fmt.Printf("Successfully created page in Confluence with ID: %s\n", pageID)
}

func printHelp() {
	fmt.Println("Markdown to Confluence Converter")
	fmt.Println("--------------------------------")
	fmt.Println("Usage:")
	fmt.Println("  help                         Show this help message")
	fmt.Println("  version                      Show version information")
	fmt.Println("  convert <markdown_or_file>   Convert Markdown to Confluence format")
	fmt.Println("    Options:")
	fmt.Println("    --output=<file>            Save output to specified file")
	fmt.Println("  post <markdown_or_file> <confluence_url> <username> <api_token> <space_key> <title> [parent_id]")
	fmt.Println("                                Convert and post to Confluence")
}

func printVersion() {
	fmt.Println("Markdown to Confluence Converter v0.1.0")
	fmt.Println("Released: April 2025")
}
