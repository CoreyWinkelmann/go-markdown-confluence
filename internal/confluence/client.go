// Package confluence provides types for Confluence API entities.
package confluence

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Ensure ConfluenceClient is defined as part of the package
var _ ConfluenceAPI = (*ConfluenceClient)(nil)

// Define ConfluenceClient as a struct implementing the ConfluenceClient interface
type ConfluenceClient struct {
	BaseURL    string
	Username   string
	APIToken   string
	HTTPClient *http.Client
}

// Implement the NewConfluenceClient function to create a new ConfluenceClient instance
func NewConfluenceClient(baseURL, username, apiToken string) *ConfluenceClient {
	return &ConfluenceClient{
		BaseURL:    baseURL,
		Username:   username,
		APIToken:   apiToken,
		HTTPClient: &http.Client{},
	}
}

// Implement the CreateParentPage method for the ConfluenceClient struct
func (c *ConfluenceClient) CreateParentPage(spaceKey, title, parentID string) (string, error) {
	// Placeholder implementation for creating a parent page
	fmt.Printf("Creating parent page: SpaceKey=%s, Title=%s, ParentID=%s\n", spaceKey, title, parentID)
	return "mock-parent-page-id", nil
}

// Implement the CreatePage method for the ConfluenceClient struct
func (c *ConfluenceClient) CreatePage(spaceKey, title, content string, parentID string) (string, error) {
	// Placeholder implementation for creating a page
	fmt.Printf("Creating page: SpaceKey=%s, Title=%s, ParentID=%s\n", spaceKey, title, parentID)
	return "mock-page-id", nil
}

// Implement the UpdatePage method for the ConfluenceClient struct
func (c *ConfluenceClient) UpdatePage(pageID, title, content, spaceKey string, version int) error {
	// Placeholder implementation for updating a page
	fmt.Printf("Updating page: PageID=%s, Title=%s, Version=%d\n", pageID, title, version)
	return nil
}

// UploadAttachment uploads a file as an attachment to the specified page.
func (c *ConfluenceClient) UploadAttachment(pageID, filePath string) error {
	fmt.Printf("Uploading attachment %s to page %s\n", filePath, pageID)
	return nil
}

// Implement GetPageByTitle in the Confluence client
func (c *ConfluenceClient) GetPageByTitle(spaceKey, title string) (*Page, error) {
	url := fmt.Sprintf("%s/rest/api/content?spaceKey=%s&title=%s", c.BaseURL, spaceKey, title)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	request.SetBasicAuth(c.Username, c.APIToken)
	response, err := c.HTTPClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil, nil // Page not found
	} else if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	var result struct {
		Results []Page `json:"results"`
	}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result.Results) == 0 {
		return nil, nil // No pages found
	}

	return &result.Results[0], nil
}
