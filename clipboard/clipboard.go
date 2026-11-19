package clipboard

import (
	"context"
	"fmt"

	"golang.design/x/clipboard"
)

// Clipboard provides a cross-platform interface for accessing the system's clipboard.
type Clipboard interface {
	// Read returns the clipboard content as a string.
	Read() (string, error)
}

type RealClipboard struct{}

// NewClipboard creates a new instance of RealClipboard.
func NewClipboard() Clipboard {
	return &RealClipboard{}
}

// Read reads the clipboard content and returns it as a string.
// If no clipboard is available, or it doesn't hold text data, an error is returned.
func (c *RealClipboard) Read() (string, error) {
	if err := clipboard.Init(); err != nil {
		return "", fmt.Errorf("clipboard unavailable: %w", err)
	}
	bytes, err := clipboard.Read(context.Background(), clipboard.FmtText)
	if err != nil {
		return "", fmt.Errorf("cannot read text from clipboard: %w", err)
	}
	return string(bytes), nil
}
