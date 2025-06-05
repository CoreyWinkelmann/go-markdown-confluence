package mermaid

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

// GeneratePlaceholder creates a placeholder image file for a mermaid diagram.
// In a real implementation this would render the diagram to an image.
func GeneratePlaceholder(diagram string) (string, error) {
	tmpDir, err := os.MkdirTemp("", "mermaid")
	if err != nil {
		return "", err
	}
	file := filepath.Join(tmpDir, "diagram.png")
	_ = ioutil.WriteFile(file, []byte(diagram), 0644)
	return file, nil
}

// RenderDiagram attempts to render a mermaid diagram to an image using the mmdc CLI.
// If mmdc is not available, it falls back to GeneratePlaceholder.
func RenderDiagram(diagram string) (string, error) {
	if _, err := exec.LookPath("mmdc"); err == nil {
		tmpDir, err := os.MkdirTemp("", "mermaid")
		if err != nil {
			return "", err
		}
		src := filepath.Join(tmpDir, "diagram.mmd")
		if err := ioutil.WriteFile(src, []byte(diagram), 0644); err != nil {
			return "", err
		}
		out := filepath.Join(tmpDir, "diagram.png")
		cmd := exec.Command("mmdc", "-i", src, "-o", out)
		if err := cmd.Run(); err == nil {
			return out, nil
		}
	}
	return GeneratePlaceholder(diagram)
}
