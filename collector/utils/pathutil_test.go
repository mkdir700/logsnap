package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandPath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Skipf("无法获取用户主目录: %v", err)
	}

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "以波浪号开头的路径",
			path:     "~/Documents",
			expected: filepath.Join(homeDir, "Documents"),
		},
		{
			name:     "普通路径",
			path:     "/tmp/test",
			expected: "/tmp/test",
		},
		{
			name:     "相对路径",
			path:     "./relative/path",
			expected: "./relative/path",
		},
		{
			name:     "空路径",
			path:     "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExpandPath(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}
