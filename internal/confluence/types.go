// Package confluence provides types and functionality to interact with Confluence API.
package confluence

// Page represents a Confluence page.
type Page struct {
	ID        string     `json:"id,omitempty"`
	Type      string     `json:"type"`
	Title     string     `json:"title"`
	Space     Space      `json:"space"`
	Body      Body       `json:"body"`
	Version   *Version   `json:"version,omitempty"`
	Ancestors []Ancestor `json:"ancestors,omitempty"`
}

// Space represents a Confluence space.
type Space struct {
	Key string `json:"key"`
}

// Body represents the body of a Confluence page.
type Body struct {
	Storage Storage `json:"storage"`
}

// Storage represents the storage format of a Confluence page body.
type Storage struct {
	Value          string `json:"value"`
	Representation string `json:"representation"`
}

// Version represents the version of a Confluence page.
type Version struct {
	Number int `json:"number"`
}

// Ancestor represents an ancestor of a Confluence page.
type Ancestor struct {
	ID string `json:"id"`
}

// ADFDocument represents the root of an Atlassian Document Format (ADF) structure.
type ADFDocument struct {
	Type    string        `json:"type"`
	Content []interface{} `json:"content"`
}

// ADFParagraph represents a paragraph in ADF.
type ADFParagraph struct {
	Type    string        `json:"type"`
	Content []interface{} `json:"content"`
}

// ADFHeading represents a heading in ADF.
type ADFHeading struct {
	Type    string        `json:"type"`
	Attrs   HeadingAttrs  `json:"attrs"`
	Content []interface{} `json:"content"`
}

// HeadingAttrs represents attributes for a heading.
type HeadingAttrs struct {
	Level int `json:"level"`
}

// ADFText represents a text node in ADF.
type ADFText struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Marks []Mark `json:"marks,omitempty"`
}

// Mark represents formatting (e.g., bold, italic) in ADF.
type Mark struct {
	Type string `json:"type"`
}

// ADFEmphasis represents emphasized text (bold, italic) in ADF.
// Note: In ADF, this is represented as text with marks.
type ADFEmphasis struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Marks []Mark `json:"marks"`
}

// ADFLink represents a hyperlink in ADF.
type ADFLink struct {
	Type    string        `json:"type"`
	Attrs   LinkAttrs     `json:"attrs"`
	Content []interface{} `json:"content"`
}

// LinkAttrs represents attributes for a link.
type LinkAttrs struct {
	Href string `json:"href"`
}

// ADFImage represents an image in ADF.
type ADFImage struct {
	Type  string     `json:"type"`
	Attrs ImageAttrs `json:"attrs"`
}

// ImageAttrs represents attributes for an image.
type ImageAttrs struct {
	Src    string `json:"src"`
	Alt    string `json:"alt,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

// ADFCodeBlock represents a code block in ADF.
type ADFCodeBlock struct {
	Type    string         `json:"type"`
	Attrs   CodeBlockAttrs `json:"attrs,omitempty"`
	Content []interface{}  `json:"content"`
}

// CodeBlockAttrs represents attributes for a code block.
type CodeBlockAttrs struct {
	Language string `json:"language,omitempty"`
}

// ADFCodeSpan represents inline code in ADF.
type ADFCodeSpan struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Marks []Mark `json:"marks"`
}

// ADFList represents a list (ordered or unordered) in ADF.
type ADFList struct {
	Type    string        `json:"type"`
	Attrs   ListAttrs     `json:"attrs,omitempty"`
	Content []interface{} `json:"content"`
}

// ListAttrs represents attributes for a list.
type ListAttrs struct {
	Order string `json:"order,omitempty"` // "ordered" or "bullet"
}

// ADFListItem represents an item in a list in ADF.
type ADFListItem struct {
	Type    string        `json:"type"`
	Content []interface{} `json:"content"`
}

// ADFRule represents a horizontal rule in ADF.
type ADFRule struct {
	Type string `json:"type"`
}

// ADFBlockquote represents a blockquote in ADF.
type ADFBlockquote struct {
	Type    string        `json:"type"`
	Content []interface{} `json:"content"`
}

// ADFTable represents a table in ADF.
type ADFTable struct {
	Type    string        `json:"type"`
	Content []ADFTableRow `json:"content"`
}

// ADFTableRow represents a table row in ADF.
type ADFTableRow struct {
	Type    string         `json:"type"`
	Content []ADFTableCell `json:"content"`
}

// ADFTableCell represents a table cell in ADF.
type ADFTableCell struct {
	Type    string        `json:"type"`
	Content []interface{} `json:"content"`
}

// ADFPanel represents a panel in ADF.
type ADFPanel struct {
	Type    string        `json:"type"`
	Attrs   PanelAttrs    `json:"attrs"`
	Content []interface{} `json:"content"`
}

// PanelAttrs represents the attributes of a panel in ADF.
type PanelAttrs struct {
	PanelType string `json:"panelType"`
}
