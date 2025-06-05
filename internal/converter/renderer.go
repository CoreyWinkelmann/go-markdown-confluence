// Package converter provides functionality to convert parsed AST to Confluence ADF format.
package converter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"

	"go-markdown-confluence/internal/confluence"
	"go-markdown-confluence/internal/mermaid"
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
				text := string(v.Segment.Value(source))
				if strings.HasPrefix(text, "[[") && strings.HasSuffix(text, "]]") {
					target := strings.TrimSuffix(strings.TrimPrefix(text, "[["), "]]")
					linkText := target
					if parts := strings.SplitN(target, "|", 2); len(parts) == 2 {
						target = parts[0]
						linkText = parts[1]
					}
					link := &confluence.ADFLink{
						Type:    "link",
						Attrs:   confluence.LinkAttrs{Href: target},
						Content: []interface{}{&confluence.ADFText{Type: "text", Text: linkText}},
					}
					addToParent(doc, link)
				} else if strings.HasPrefix(text, ":") && strings.HasSuffix(text, ":") {
					emoji := &confluence.ADFEmoji{
						Type: "emoji",
						Attrs: confluence.EmojiAttrs{
							ShortName: text,
						},
					}
					addToParent(doc, emoji)
				} else {
					textNode := &confluence.ADFText{
						Type: "text",
						Text: text,
					}
					addToParent(doc, textNode)
				}

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

				lines := n.Lines()
				var codeText strings.Builder
				for i := 0; i < lines.Len(); i++ {
					line := lines.At(i)
					codeText.Write(line.Value(source))
				}

				codeStr := codeText.String()

				if language == "adf" || language == "adf-json" {
					var raw interface{}
					if err := json.Unmarshal([]byte(codeStr), &raw); err == nil {
						doc.Content = append(doc.Content, raw)
						return ast.WalkSkipChildren, nil
					}
				}

				if language == "mermaid" {
					imgPath, _ := mermaid.RenderDiagram(codeStr)
					image := &confluence.ADFImage{Type: "image", Attrs: confluence.ImageAttrs{Src: imgPath}}
					addToParent(doc, image)
					return ast.WalkSkipChildren, nil
				}

				codeBlock := &confluence.ADFCodeBlock{
					Type:    "codeBlock",
					Attrs:   confluence.CodeBlockAttrs{Language: language},
					Content: []interface{}{&confluence.ADFText{Type: "text", Text: codeStr}},
				}
				doc.Content = append(doc.Content, codeBlock)
				return ast.WalkSkipChildren, nil

			case ast.KindCodeSpan:
				parent := getCurrentParent(doc)
				if parent != nil {
					if textParent, ok := parent.(*confluence.ADFText); ok {
						textParent.Marks = append(textParent.Marks, confluence.Mark{Type: "code"})
					}
				}

			// Ensure no additional paragraph content is added for task lists and decision items
			case ast.KindList:
				v := n.(*ast.List)
				isTaskList := false
				for child := v.FirstChild(); child != nil; child = child.NextSibling() {
					if item, ok := child.(*ast.ListItem); ok {
						if item.FirstChild() != nil {
							text := string(item.FirstChild().Text(source))
							if len(text) > 0 && (text[0] == '[' && (text[1] == ' ' || text[1] == 'x') && text[2] == ']') {
								isTaskList = true
								break
							}
						}
					}
				}

				if isTaskList {
					taskList := &confluence.ADFTaskList{
						Type:    "taskList",
						Content: []interface{}{},
					}
					doc.Content = append(doc.Content, taskList)
					return ast.WalkSkipChildren, nil // Ensure proper return values
				}

				listType := "bulletList"
				if v.IsOrdered() {
					listType = "orderedList"
				}

				list := &confluence.ADFList{
					Type:    listType,
					Content: []interface{}{},
				}
				doc.Content = append(doc.Content, list)
				return ast.WalkContinue, nil // Ensure proper return values

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
				textContent := strings.TrimSpace(string(v.Text(source)))
				if strings.HasPrefix(strings.ToLower(textContent), "[!") {
					end := strings.Index(textContent, "]")
					if end != -1 {
						callout := strings.ToLower(textContent[2:end])
						content := strings.TrimSpace(textContent[end+1:])
						panelType := "info"
						if callout == "warning" || callout == "caution" || callout == "danger" {
							panelType = "warning"
						}
						panel := &confluence.ADFPanel{
							Type:    "panel",
							Attrs:   confluence.PanelAttrs{PanelType: panelType},
							Content: []interface{}{&confluence.ADFParagraph{Type: "paragraph", Content: []interface{}{&confluence.ADFText{Type: "text", Text: content}}}},
						}
						doc.Content = append(doc.Content, panel)
						return ast.WalkSkipChildren, nil
					}
				}
				if strings.HasPrefix(strings.ToLower(textContent), "decision:") {
					decision := &confluence.ADFDecisionItem{
						Type: "decisionItem",
						Attrs: confluence.DecisionItemAttrs{
							State: "DECIDED",
						},
					}
					doc.Content = append(doc.Content, decision)
					return ast.WalkSkipChildren, nil // Ensure proper return values
				}

				blockquote := &confluence.ADFBlockquote{
					Type:    "blockquote",
					Content: []interface{}{},
				}
				doc.Content = append(doc.Content, blockquote)
				return ast.WalkContinue, nil // Ensure proper return values

			case ast.KindHTMLBlock:
				v := n.(*ast.HTMLBlock)
				if strings.Contains(string(v.Text(source)), "placeholder") {
					placeholder := &confluence.ADFPlaceholder{
						Type: "placeholder",
						Attrs: confluence.PlaceholderAttrs{
							Text: "Add your content here",
						},
					}
					addToParent(doc, placeholder)
				}
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

// Prevent adding empty paragraphs and redundant content
func addToParent(doc *confluence.ADFDocument, node interface{}) {
	// Directly add standalone elements without wrapping
	switch node.(type) {
	case *confluence.ADFEmoji, *confluence.ADFPlaceholder, *confluence.ADFTaskList, *confluence.ADFDecisionItem:
		if len(doc.Content) > 0 {
			// Remove the last element if it's an empty paragraph
			if lastElem, ok := doc.Content[len(doc.Content)-1].(*confluence.ADFParagraph); ok && len(lastElem.Content) == 0 {
				doc.Content = doc.Content[:len(doc.Content)-1]
			}
		}
		doc.Content = append(doc.Content, node)
		return
	}

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
