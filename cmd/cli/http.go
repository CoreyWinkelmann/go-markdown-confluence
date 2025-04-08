package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-markdown-confluence/internal/confluence"
	"go-markdown-confluence/pkg/markdownconfluence"
	"net/http"
	"time"
)

// ConfluenceClient handles communication with the Confluence API
type ConfluenceClient struct {
	baseURL    string
	username   string
	apiToken   string
	httpClient *http.Client
}

// NewConfluenceClient creates a new Confluence API client
func NewConfluenceClient(baseURL, username, apiToken string) *ConfluenceClient {
	return &ConfluenceClient{
		baseURL:  baseURL,
		username: username,
		apiToken: apiToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CreatePage creates a new page in Confluence
func (c *ConfluenceClient) CreatePage(spaceKey, title, content string, parentID string) (string, error) {
	// Create page request payload
	page := confluence.NewPage(title, spaceKey, content, parentID)

	// Prepare the HTTP request
	reqBody, err := json.Marshal(page)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	url := fmt.Sprintf("%s/rest/api/content", c.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.username, c.apiToken)

	// Send the HTTP request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("received error status code: %d", resp.StatusCode)
	}

	// Parse the response
	var response struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	return response.ID, nil
}

// UpdatePage updates an existing page in Confluence
func (c *ConfluenceClient) UpdatePage(pageID, title, content, spaceKey string, version int) error {
	// Create page update payload
	page := confluence.NewPageWithVersion(title, spaceKey, content, version)

	// Prepare the HTTP request
	reqBody, err := json.Marshal(page)
	if err != nil {
		return fmt.Errorf("error marshaling request: %w", err)
	}

	url := fmt.Sprintf("%s/rest/api/content/%s", c.baseURL, pageID)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.username, c.apiToken)

	// Send the HTTP request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("received error status code: %d", resp.StatusCode)
	}

	return nil
}

// Ensure ConfluenceClient implements the markdownconfluence.ConfluenceClient interface
var dummyClient markdownconfluence.ConfluenceClient = nil
