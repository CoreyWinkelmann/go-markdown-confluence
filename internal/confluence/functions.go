// Package confluence provides types and data structures for Confluence pages.
package confluence

// NewPage creates a new page with the given parameters
func NewPage(title, spaceKey, content, parentID string) Page {
	page := Page{
		Type:  "page",
		Title: title,
		Space: Space{
			Key: spaceKey,
		},
		Body: Body{
			Storage: Storage{
				Value:          content,
				Representation: "atlas_doc_format", // Changed from "wiki" to the ADF format
			},
		},
	}

	// Add parent page as ancestor if provided
	if parentID != "" {
		page.Ancestors = []Ancestor{
			{ID: parentID},
		}
	}

	return page
}

// NewPageWithVersion creates a new page with version information for updates
func NewPageWithVersion(title, spaceKey, content string, version int) Page {
	page := NewPage(title, spaceKey, content, "")
	page.Version = &Version{
		Number: version,
	}
	return page
}
