package service

import (
	"path/filepath"
	"strings"

	"logsnap/collector/utils"

	"github.com/sirupsen/logrus"
)

// LogPath 定义日志路径结构
type LogPath struct {
	Name string
	Path string
}

// RootPath 获取日志相对根路径
func (l *LogPath) RootPath(config *Config) string {
	// 使用相对路径的方式获取根路径
	// 例如：如果LogDir是"/custom/log"，Path是"/custom/log/xyz_max_hmi/server"
	// 则提取"xyz_max_hmi/server"作为根路径

	// 首先获取绝对路径
	absPath := utils.ExpandPath(l.Path)

	// 获取LogDir的绝对路径
	absLogDir := utils.ExpandPath(config.LogRootDir)

	// 如果路径不以LogDir开头，尝试使用相对路径
	if !strings.HasPrefix(absPath, absLogDir) {
		// 尝试获取最后两级目录作为根路径
		parts := strings.Split(absPath, string(filepath.Separator))
		if len(parts) >= 2 {
			return filepath.Join(parts[len(parts)-2:]...)
		}
		return filepath.Base(absPath)
	}

	// 提取相对于LogDir的部分作为根路径
	rel, err := filepath.Rel(absLogDir, absPath)
	if err != nil {
		logrus.Warnf("获取相对路径失败: %v，使用完整路径", err)
		return absPath
	}
	return rel
}

// AbsolutePath 获取绝对路径
func (l *LogPath) AbsolutePath() string {
	return utils.ExpandPath(l.Path)
}
