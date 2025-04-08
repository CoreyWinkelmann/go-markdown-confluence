// Package converter provides functionality to convert parsed AST to Confluence ADF format.
package converter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"

	"go-markdown-confluence/internal/confluence"
)

// ConvertToADF converts a parsed AST node to an ADFDocument.
func ConvertToADF(n ast.Node, source []byte) (*confluence.ADFDocument, error) {
	if n == nil {
		return nil, fmt.Errorf("Invalid Markdown: AST node is nil")
	}

	doc := &confluence.ADFDocument{
		Type:    "doc",
		Content: []interface{}{},
	}

	err := ast.Walk(n, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			switch n.Kind() {
			case ast.KindDocument:
			default:
				if n.Kind() == ast.KindText && len(n.Text(source)) == 0 {
					return ast.WalkContinue, nil
				}
			}

			switch n.Kind() {
			case ast.KindDocument:

			case ast.KindHeading:
				v := n.(*ast.Heading)
				heading := &confluence.ADFHeading{
					Type: "heading",
					Attrs: confluence.HeadingAttrs{
						Level: v.Level,
					},
					Content: []interface{}{},
				}
				doc.Content = append(doc.Content, heading)

			case ast.KindParagraph:
				paragraph := &confluence.ADFParagraph{
					Type:    "paragraph",
					Content: []interface{}{},
				}
				doc.Content = append(doc.Content, paragraph)

			case ast.KindText:
				v := n.(*ast.Text)
				text := &confluence.ADFText{
					Type: "text",
					Text: string(v.Segment.Value(source)),
				}
				addToParent(doc, text)

			case ast.KindEmphasis:
				v := n.(*ast.Emphasis)
				markType := "strong"
				if v.Level == 1 {
					markType = "em"
				}

				parent := getCurrentParent(doc)
				if parent != nil {
					if textParent, ok := parent.(*confluence.ADFText); ok {
						textParent.Marks = append(textParent.Marks, confluence.Mark{Type: markType})
					}
				}

			case ast.KindLink:
				v := n.(*ast.Link)
				link := &confluence.ADFLink{
					Type: "link",
					Attrs: confluence.LinkAttrs{
						Href: string(v.Destination),
					},
					Content: []interface{}{},
				}
				addToParent(doc, link)

			case ast.KindImage:
				v := n.(*ast.Image)
				image := &confluence.ADFImage{
					Type: "image",
					Attrs: confluence.ImageAttrs{
						Src: string(v.Destination),
						Alt: string(v.Title),
					},
				}
				addToParent(doc, image)
				return ast.WalkSkipChildren, nil

			case ast.KindCodeBlock, ast.KindFencedCodeBlock:
				var language string
				if fenced, ok := n.(*ast.FencedCodeBlock); ok {
					language = string(fenced.Language(source))
				}

				codeBlock := &confluence.ADFCodeBlock{
					Type: "codeBlock",
					Attrs: confluence.CodeBlockAttrs{
						Language: language,
					},
					Content: []interface{}{},
				}
				doc.Content = append(doc.Content, codeBlock)

				if fenced, ok := n.(*ast.FencedCodeBlock); ok {
					lines := fenced.Lines()
					var codeText strings.Builder
					for i := 0; i < lines.Len(); i++ {
						line := lines.At(i)
						codeText.Write(line.Value(source))
					}

					codeContent := &confluence.ADFText{
						Type: "text",
						Text: codeText.String(),
					}
					codeBlock.Content = append(codeBlock.Content, codeContent)
					return ast.WalkSkipChildren, nil
				}

			case ast.KindCodeSpan:
				parent := getCurrentParent(doc)
				if parent != nil {
					if textParent, ok := parent.(*confluence.ADFText); ok {
						textParent.Marks = append(textParent.Marks, confluence.Mark{Type: "code"})
					}
				}

			case ast.KindList:
				v := n.(*ast.List)
				listType := "bulletList"
				if v.IsOrdered() {
					listType = "orderedList"
				}

				list := &confluence.ADFList{
					Type:    listType,
					Content: []interface{}{},
				}
				doc.Content = append(doc.Content, list)

			case ast.KindListItem:
				listItem := &confluence.ADFListItem{
					Type:    "listItem",
					Content: []interface{}{},
				}

				if len(doc.Content) > 0 {
					if list, ok := doc.Content[len(doc.Content)-1].(*confluence.ADFList); ok {
						list.Content = append(list.Content, listItem)
					}
				}

			case ast.KindThematicBreak:
				rule := &confluence.ADFRule{
					Type: "rule",
				}
				doc.Content = append(doc.Content, rule)

			case ast.KindBlockquote:
				v := n.(*ast.Blockquote)
				if child := v.FirstChild(); child != nil {
					if textNode, ok := child.(*ast.Text); ok {
						text := string(textNode.Segment.Value(source))
						panelType := ""
						if strings.HasPrefix(strings.ToLower(text), "**info:**") {
							panelType = "info"
						} else if strings.HasPrefix(strings.ToLower(text), "**warning:**") {
							panelType = "warning"
						} else if strings.HasPrefix(strings.ToLower(text), "**error:**") {
							panelType = "error"
						}

						if panelType != "" {
							panel := &confluence.ADFPanel{
								Type: "panel",
								Attrs: confluence.PanelAttrs{
									PanelType: panelType,
								},
								Content: []interface{}{},
							}
							doc.Content = append(doc.Content, panel)
							return ast.WalkSkipChildren, nil
						}
					}
				}

				blockquote := &confluence.ADFBlockquote{
					Type:    "blockquote",
					Content: []interface{}{},
				}
				doc.Content = append(doc.Content, blockquote)
			}
		}
		return ast.WalkContinue, nil
	})

	if err != nil {
		return nil, err
	}

	return doc, nil
}

// getCurrentParent returns the last content element that can contain child nodes.
func getCurrentParent(doc *confluence.ADFDocument) interface{} {
	if len(doc.Content) == 0 {
		return nil
	}

	lastElem := doc.Content[len(doc.Content)-1]

	switch v := lastElem.(type) {
	case *confluence.ADFParagraph:
		if len(v.Content) > 0 {
			return v.Content[len(v.Content)-1]
		}
		return v
	case *confluence.ADFHeading:
		if len(v.Content) > 0 {
			return v.Content[len(v.Content)-1]
		}
		return v
	case *confluence.ADFBlockquote:
		if len(v.Content) > 0 {
			return v.Content[len(v.Content)-1]
		}
		return v
	case *confluence.ADFCodeBlock:
		if len(v.Content) > 0 {
			return v.Content[len(v.Content)-1]
		}
		return v
	case *confluence.ADFLink:
		if len(v.Content) > 0 {
			return v.Content[len(v.Content)-1]
		}
		return v
	case *confluence.ADFList:
		if len(v.Content) > 0 {
			listItem := v.Content[len(v.Content)-1]
			if li, ok := listItem.(*confluence.ADFListItem); ok {
				if len(li.Content) > 0 {
					return li.Content[len(li.Content)-1]
				}
				return li
			}
		}
		return v
	}

	return lastElem
}

// addToParent adds a node to the appropriate parent in the document.
func addToParent(doc *confluence.ADFDocument, node interface{}) {
	if len(doc.Content) == 0 {
		paragraph := &confluence.ADFParagraph{
			Type:    "paragraph",
			Content: []interface{}{},
		}
		doc.Content = append(doc.Content, paragraph)
	}

	lastElem := doc.Content[len(doc.Content)-1]

	switch v := lastElem.(type) {
	case *confluence.ADFParagraph:
		v.Content = append(v.Content, node)
	case *confluence.ADFHeading:
		v.Content = append(v.Content, node)
	case *confluence.ADFBlockquote:
		v.Content = append(v.Content, node)
	case *confluence.ADFCodeBlock:
		v.Content = append(v.Content, node)
	case *confluence.ADFLink:
		v.Content = append(v.Content, node)
	case *confluence.ADFList:
		if len(v.Content) > 0 {
			listItem := v.Content[len(v.Content)-1]
			if li, ok := listItem.(*confluence.ADFListItem); ok {
				li.Content = append(li.Content, node)
			}
		}
	default:
		paragraph := &confluence.ADFParagraph{
			Type:    "paragraph",
			Content: []interface{}{node},
		}
		doc.Content = append(doc.Content, paragraph)
	}
}

// SerializeToJSON converts an ADFDocument to its JSON representation.
func SerializeToJSON(doc *confluence.ADFDocument) (string, error) {
	bytes, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error serializing to JSON: %w", err)
	}
	return string(bytes), nil
}
