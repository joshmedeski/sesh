package seshcli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloneVcs(t *testing.T) {
	tests := []struct {
		name       string
		flag       string
		configured string
		want       string
	}{
		{"flag overrides config", "git", "jj", "git"},
		{"config used without flag", "", "jj", "jj"},
		{"neither set defaults to git in the cloner", "", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, cloneVcs(tt.flag, tt.configured))
		})
	}
}
