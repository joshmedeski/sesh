package ansi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColorCode(t *testing.T) {
	tests := []struct {
		name  string
		color string
		code  int
		ok    bool
	}{
		{"standard color", "blue", Blue, true},
		{"bright color", "bright-blue", BrightBlue, true},
		{"case and whitespace", "  Bright-Magenta  ", BrightMagenta, true},
		{"gray alias", "gray", BrightBlack, true},
		{"grey alias", "grey", BrightBlack, true},
		{"empty color", "", 0, true},
		{"unsupported color", "orange", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, ok := ColorCode(tt.color)
			assert.Equal(t, tt.code, code)
			assert.Equal(t, tt.ok, ok)
		})
	}
}
