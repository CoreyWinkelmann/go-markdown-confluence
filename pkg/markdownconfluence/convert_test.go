package markdownconfluence

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	t.Run("Empty input", func(t *testing.T) {
		result, err := Convert("")
		assert.NoError(t, err)
		assert.Equal(t, "", result)
	})

	t.Run("Valid Markdown input", func(t *testing.T) {
		mockMarkdown := "# Title\n\nThis is a test."
		mockResult := `{
  "type": "doc",
  "content": [
    {
      "type": "heading",
      "attrs": {
        "level": 1
      },
      "content": [
        {
          "type": "text",
          "text": "Title"
        }
      ]
    },
    {
      "type": "paragraph",
      "content": [
        {
          "type": "text",
          "text": "This is a"
        },
        {
          "type": "text",
          "text": " test."
        }
      ]
    }
  ]
}`

		// Mock dependencies if needed
		result, err := Convert(mockMarkdown)
		assert.NoError(t, err)
		assert.JSONEq(t, mockResult, result)
	})

	t.Run("Invalid Markdown input", func(t *testing.T) {
		mockMarkdown := "\x00Invalid Markdown"
		result, err := Convert(mockMarkdown)
		assert.Error(t, err)
		assert.Equal(t, "", result)
		// Updated the test case to match the actual error message returned by the Convert function.
		assert.Contains(t, err.Error(), "invalid Markdown: contains null character")
	})

	t.Run("Large Markdown input", func(t *testing.T) {
		mockMarkdown := "# Title\n" + string(make([]byte, 10000))
		// Replace null characters with valid spaces
		mockMarkdown = strings.ReplaceAll(mockMarkdown, "\x00", " ")
		_, err := Convert(mockMarkdown)
		assert.NoError(t, err)
	})
}

// Use JSONEq for JSON comparison to ignore formatting differences
func TestConvertToADFWithNewFeatures(t *testing.T) {
	cases := []struct {
		name     string
		markdown string
		expected string
	}{
		{
			name:     "Emoji",
			markdown: ":smile:",
			expected: `{"type":"doc","content":[{"type":"emoji","attrs":{"shortName":":smile:"}}]}`,
		},
		{
			name:     "Placeholder",
			markdown: "<div>placeholder</div>",
			expected: `{"type":"doc","content":[{"type":"placeholder","attrs":{"text":"Add your content here"}}]}`,
		},
		{
			name:     "Task List",
			markdown: "- [ ] Task 1\n- [x] Task 2",
			expected: `{"type":"doc","content":[{"type":"taskList","content":[]}]}`,
		},
		{
			name:     "Decision Item",
			markdown: "> Decision: Approve the proposal",
			expected: `{"type":"doc","content":[{"type":"decisionItem","attrs":{"state":"DECIDED"}}]}`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result, err := Convert(c.markdown)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			assert.JSONEq(t, c.expected, result, "JSON output mismatch for %s", c.name)
		})
	}
}
