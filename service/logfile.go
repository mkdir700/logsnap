package service

import (
	"os"
	"path/filepath"
	"time"
)

// LogFile 定义日志文件结构
type LogFile struct {
	Name      string    // 文件名
	Path      string    // 文件路径
	Size      int64     // 文件大小
	ModTime   time.Time // 修改时间
	IsDir     bool      // 是否是目录
	LogPath   *LogPath  // 所属日志路径
	RootPath  string    // 相对根路径
	Selected  bool      // 是否被选中
	Uploading bool      // 是否正在上传
}

// NewLogFile 创建新的日志文件对象
func NewLogFile(path string, info os.FileInfo, logPath *LogPath, config *Config) *LogFile {
	rootPath := ""
	if logPath != nil {
		rootPath = logPath.RootPath(config)
	}

	return &LogFile{
		Name:     info.Name(),
		Path:     path,
		Size:     info.Size(),
		ModTime:  info.ModTime(),
		IsDir:    info.IsDir(),
		LogPath:  logPath,
		RootPath: rootPath,
	}
}

// RelativePath 获取相对于根路径的路径
func (l *LogFile) RelativePath() string {
	if l.LogPath == nil || l.RootPath == "" {
		return l.Name
	}

	// 获取相对于LogPath的路径
	rel, err := filepath.Rel(l.LogPath.Path, l.Path)
	if err != nil {
		return l.Name
	}

	// 如果是当前目录，则返回文件名
	if rel == "." {
		return l.Name
	}

	return filepath.Join(l.RootPath, rel)
}
