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
		assert.Contains(t, err.Error(), "Invalid Markdown")
	})

	t.Run("Large Markdown input", func(t *testing.T) {
		mockMarkdown := "# Title\n" + string(make([]byte, 10000))
		// Replace null characters with valid spaces
		mockMarkdown = strings.ReplaceAll(mockMarkdown, "\x00", " ")
		_, err := Convert(mockMarkdown)
		assert.NoError(t, err)
	})
}
