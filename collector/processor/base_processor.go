package processor

import (
	"fmt"
	collector "logsnap/collector"
	"path/filepath"
	"time"
)

// BaseProcessor 定义基础处理器结构
type BaseProcessor struct {
	Name      string   // 处理器名称
	LogDir    Path   // 基础目录	
	OutputDir string // 输出目录
}

// NewBaseProcessor 创建基础处理器
func NewBaseProcessor(name, logDir, outputDir string) *BaseProcessor {
	baseDirPath := NewPath(logDir)
	return &BaseProcessor{
		Name:      name,
		LogDir:    baseDirPath,
		OutputDir: outputDir,
	}
}

// GetName 返回处理器名称
func (p *BaseProcessor) GetName() string {
	return p.Name
}

// GetLogPath 返回日志文件的路径
func (p *BaseProcessor) GetLogPath() (string, error) {
	return p.LogDir.GetAbsolutePath()
}

// GetOutputDir 返回日志文件的输出目录
func (p *BaseProcessor) GetOutputDir() string {
	return p.OutputDir
}

// Collect 处理日志文件的通用方法，子类可以覆盖
func (p *BaseProcessor) Collect(startTime, endTime time.Time, rootOutputDir string) (string, []collector.FileProcessResult, error) {
	// 创建文件处理器
	fileProcessors := p.CreateFileProcessor()

	if len(fileProcessors) == 0 {
		return "", nil, fmt.Errorf("没有找到文件处理器")
	}
	
	// 拼接成完整的输出目录
	outputDir := filepath.Join(rootOutputDir, p.OutputDir)

	// 使用通用的收集方法
	return CollectWithProcessor(p, fileProcessors, startTime, endTime, outputDir)
}

// CreateFileProcessor 创建文件处理器，子类应该覆盖此方法
func (p *BaseProcessor) CreateFileProcessor() []FileProcessorProvider {
	// 默认实现，子类应该覆盖
	return nil
}

type BaseProcessorProvider struct {
	FileInfoFilter FileInfoProvider
	Suffixes []string
}

// NewBaseProcessorProvider 创建基础处理器提供者
// 参数:
//  fileInfoFilter: 文件信息过滤器
//  suffixes: 文件后缀
//  processStrategy: 处理策略
func NewBaseProcessorProvider(
	fileInfoFilter FileInfoProvider,
	suffixes []string,
) *BaseProcessorProvider {
	return &BaseProcessorProvider{	
		FileInfoFilter: fileInfoFilter,
		Suffixes: suffixes,
	}
}

// FindFiles 查找日志文件
// 参数:
//  dirPath: 目录路径
//  suffixes: 文件后缀
// 返回:
//  日志文件路径列表
//  错误信息
func (p *BaseProcessorProvider) FindFiles(dirPath string, suffixes ...string) ([]string, error) {
	return DefaultFindLogFiles(dirPath, suffixes...)
}

// FilterFiles 过滤日志文件
// 参数:
//  files: 日志文件路径列表
//  startTime: 开始时间
//  endTime: 结束时间
// 返回:
//  日志文件信息列表
//  错误信息
func (p *BaseProcessorProvider) FilterFiles(files []string, startTime, endTime time.Time) ([]LogFileInfo, error) {
	return FilterFiles(files, p.FileInfoFilter, startTime, endTime, nil)
}

// ProcessDir 处理目录
// 参数:
//  dirPath: 目录路径
//  outputDir: 输出目录
//  startTime: 开始时间
//  endTime: 结束时间
// 返回:
func (p *BaseProcessorProvider) ProcessDir(dirPath, outputDir string, startTime, endTime time.Time) ([]collector.FileProcessResult, error) {
	return DefaultProcessDir(p, dirPath, outputDir, startTime, endTime)
}

// ProcessFile 处理日志文件, 默认是直接复制
// 参数:
//  fileInfo: 日志文件信息
//  startTime: 开始时间
//  endTime: 结束时间
//  outputDir: 输出目录
// 返回:
func (p *BaseProcessorProvider) ProcessFile(fileInfo LogFileInfo, startTime, endTime time.Time, outputDir string) (collector.FileProcessResult, error) {
	return ProcessLogWithStrategy(fileInfo, outputDir, &CopyLogProcessor{})
}

// GetFileSuffixes 返回文件后缀
// 返回:
//  文件后缀列表
func (p *BaseProcessorProvider) GetFileSuffixes() []string {
	return p.Suffixes
}