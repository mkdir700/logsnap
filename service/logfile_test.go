package service

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 模拟FileInfo接口
type mockFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
	sys     interface{}
}

func (m mockFileInfo) Name() string       { return m.name }
func (m mockFileInfo) Size() int64        { return m.size }
func (m mockFileInfo) Mode() os.FileMode  { return m.mode }
func (m mockFileInfo) ModTime() time.Time { return m.modTime }
func (m mockFileInfo) IsDir() bool        { return m.isDir }
func (m mockFileInfo) Sys() interface{}   { return m.sys }

func TestNewLogFile(t *testing.T) {
	// 创建模拟的FileInfo
	now := time.Now()
	info := mockFileInfo{
		name:    "test.log",
		size:    1024,
		mode:    0644,
		modTime: now,
		isDir:   false,
	}
	
	// 创建LogPath
	logPath := &LogPath{
		Name: "测试日志",
		Path: "/path/to/logs",
	}
	
	// 创建配置
	config := &Config{
		LogRootDir: "/path/to",
	}
	
	// 创建LogFile
	logFile := NewLogFile("/path/to/logs/test.log", info, logPath, config)
	
	// 验证LogFile字段
	assert.Equal(t, "test.log", logFile.Name, "文件名应该正确")
	assert.Equal(t, "/path/to/logs/test.log", logFile.Path, "路径应该正确")
	assert.Equal(t, int64(1024), logFile.Size, "大小应该正确")
	assert.Equal(t, now, logFile.ModTime, "修改时间应该正确")
	assert.False(t, logFile.IsDir, "应该不是目录")
	assert.Equal(t, logPath, logFile.LogPath, "LogPath应该正确")
	assert.Equal(t, "logs", logFile.RootPath, "RootPath应该正确")
	
	// 测试没有LogPath的情况
	logFile = NewLogFile("/path/to/test.log", info, nil, config)
	assert.Equal(t, "", logFile.RootPath, "没有LogPath时RootPath应该为空")
}

func TestRelativePath(t *testing.T) {
	// 创建临时目录结构
	tempDir, err := os.MkdirTemp("", "logsnap-test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// 创建子目录
	subDir := filepath.Join(tempDir, "logs")
	err = os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatalf("创建子目录失败: %v", err)
	}
	
	// 创建测试文件
	testFilePath := filepath.Join(subDir, "test.log")
	testFile, err := os.Create(testFilePath)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}
	testFile.Close()
	
	// 获取文件信息
	fileInfo, err := os.Stat(testFilePath)
	if err != nil {
		t.Fatalf("获取文件信息失败: %v", err)
	}
	
	// 创建LogPath
	logPath := &LogPath{
		Name: "测试日志",
		Path: subDir,
	}
	
	// 创建LogFile
	logFile := &LogFile{
		Name:     fileInfo.Name(),
		Path:     testFilePath,
		Size:     fileInfo.Size(),
		ModTime:  fileInfo.ModTime(),
		IsDir:    fileInfo.IsDir(),
		LogPath:  logPath,
		RootPath: "logs",
	}
	
	// 测试RelativePath
	relPath := logFile.RelativePath()
	assert.Equal(t, "logs/test.log", relPath, "相对路径应该正确")
	
	// 测试没有LogPath的情况
	logFile.LogPath = nil
	relPath = logFile.RelativePath()
	assert.Equal(t, "test.log", relPath, "没有LogPath时应该返回文件名")
	
	// 测试有LogPath但没有RootPath的情况
	logFile.LogPath = logPath
	logFile.RootPath = ""
	relPath = logFile.RelativePath()
	assert.Equal(t, "test.log", relPath, "没有RootPath时应该返回文件名")
	
	// 测试文件直接在LogPath目录下的情况
	logFile.RootPath = "logs"
	relPath = logFile.RelativePath()
	assert.Equal(t, "logs/test.log", relPath, "直接在LogPath下的文件应该返回正确的相对路径")
}
