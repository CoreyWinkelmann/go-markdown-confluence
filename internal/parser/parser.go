// Package parser provides functionality to parse Markdown content into AST.
package parser

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// MarkdownParser handles parsing of markdown content
type MarkdownParser struct {
	parser   parser.Parser
	markdown goldmark.Markdown
}

// NewMarkdownParser creates a new MarkdownParser instance
func NewMarkdownParser() *MarkdownParser {
	p := parser.NewParser(
		parser.WithBlockParsers(parser.DefaultBlockParsers()...),
		parser.WithInlineParsers(parser.DefaultInlineParsers()...),
	)

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
		),
		goldmark.WithParser(p),
	)

	return &MarkdownParser{
		parser:   p,
		markdown: md,
	}
}

// Parse parses the given markdown string into an AST
func (mp *MarkdownParser) Parse(markdown string) ast.Node {
	if len(markdown) == 0 {
		return nil
	}

	reader := text.NewReader([]byte(markdown))
	astRoot := mp.markdown.Parser().Parse(reader)

	return astRoot
}
