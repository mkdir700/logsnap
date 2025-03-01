package utils

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

// ExpandPath 展开路径中的波浪号
func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			logrus.Warnf("无法获取用户主目录: %v", err)
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}
