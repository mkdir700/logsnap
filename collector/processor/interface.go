package processor

import (
	collector "logsnap/collector"
	"time"
)

// FileInfoProvider 定义文件信息提供者接口
type FileInfoProvider interface {
	// ParseFileInfos 从文件路径列表解析文件信息
	ParseFileInfos(files []string) ([]LogFileInfo, error)
	
	// IsMatch 通过文件名判断是否为匹配的日志文件
	IsMatch(fileName string) bool
} 

// FileProcessorProvider 定义文件处理器提供者接口
// 实现此接口的类型可以直接用于创建GenericFileProcessor
type FileProcessorProvider interface {
	// FindFiles 查找文件
	FindFiles(dirPath string, suffixes ...string) ([]string, error)

	// FilterFiles 过滤文件
	FilterFiles(files []string, startTime, endTime time.Time) ([]LogFileInfo, error)

	// ProcessDir 处理目录中的文件
	ProcessDir(dirPath, outputDir string, startTime, endTime time.Time) ([]collector.FileProcessResult, error)

	// ProcessFile 处理单个文件
	ProcessFile(fileInfo LogFileInfo, startTime, endTime time.Time, outputDir string) (collector.FileProcessResult, error)

	// GetFileSuffixes 获取文件后缀列表
	GetFileSuffixes() []string
}

// TimeProvider 是一个接口，定义了获取时间的方法
type TimeProvider interface {
	GetStartTime() time.Time
}

// FileInfoFilter 定义文件信息筛选接口
type FileInfoFilter interface {
	// ParseFileInfos 从文件路径列表解析文件信息
	ParseFileInfos(files []string) ([]LogFileInfo, error)

	// IsMatch 通过文件名判断是否为匹配的日志文件
	IsMatch(fileName string) bool
}