// Package parser provides functionality to parse Markdown content into AST.
package parser

import (
	"fmt"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// MarkdownParser handles parsing of markdown content
type MarkdownParser struct {
	parser parser.Parser
}

// NewMarkdownParser creates a new MarkdownParser instance
func NewMarkdownParser() *MarkdownParser {
	p := parser.NewParser(
		parser.WithBlockParsers(parser.DefaultBlockParsers()...),
		parser.WithInlineParsers(parser.DefaultInlineParsers()...),
	)

	return &MarkdownParser{
		parser: p,
	}
}

// Parse takes a markdown string and parses it into an AST
func (mp *MarkdownParser) Parse(markdown string) ast.Node {
	if len(markdown) == 0 {
		fmt.Println("Debug: Markdown input is empty")
		return nil
	}

	fmt.Printf("Debug: Markdown input size: %d bytes\n", len(markdown))

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM, // GitHub Flavored Markdown
		),
		goldmark.WithParser(mp.parser),
	)

	reader := text.NewReader([]byte(markdown))
	parsedNode := md.Parser().Parse(reader)

	if parsedNode == nil {
		fmt.Println("Debug: goldmark parser returned nil")
	} else {
		fmt.Println("Debug: goldmark parser successfully returned a node")
	}

	return parsedNode
}
